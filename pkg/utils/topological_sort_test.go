package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/internal/storage/db/dto"
)

// Helper function to create a UUID from string
func mustUUID(s string) pgtype.UUID {
	u, _ := uuid.Parse(s)
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

// Helper function to create a task
func createTask(id, name string, dependsOn []string) db.Task {
	return db.Task{
		ID:   mustUUID(id),
		Name: name,
		Config: dto.TaskConfig{
			DependsOn: dependsOn,
		},
	}
}

func TestTopoSort_NoDependencies(t *testing.T) {
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000001", "task1", nil),
		createTask("00000000-0000-0000-0000-000000000002", "task2", nil),
		createTask("00000000-0000-0000-0000-000000000003", "task3", nil),
	}

	sorted, err := TopoSort(tasks)

	require.NoError(t, err)
	assert.Len(t, sorted, 3)
	// All tasks should be present in the result
	taskNames := make(map[string]bool)
	for _, task := range sorted {
		taskNames[task.Name] = true
	}
	assert.True(t, taskNames["task1"])
	assert.True(t, taskNames["task2"])
	assert.True(t, taskNames["task3"])
}

func TestTopoSort_LinearDependencies(t *testing.T) {
	// task1 -> task2 -> task3
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000003", "task3", []string{"task2"}),
		createTask("00000000-0000-0000-0000-000000000001", "task1", nil),
		createTask("00000000-0000-0000-0000-000000000002", "task2", []string{"task1"}),
	}

	sorted, err := TopoSort(tasks)

	require.NoError(t, err)
	require.Len(t, sorted, 3)
	// Verify the order: task1 should come before task2, task2 before task3
	assert.Equal(t, "task1", sorted[0].Name)
	assert.Equal(t, "task2", sorted[1].Name)
	assert.Equal(t, "task3", sorted[2].Name)
}

func TestTopoSort_MultipleDependencies(t *testing.T) {
	// task3 depends on both task1 and task2
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000001", "task1", nil),
		createTask("00000000-0000-0000-0000-000000000002", "task2", nil),
		createTask("00000000-0000-0000-0000-000000000003", "task3", []string{"task1", "task2"}),
	}

	sorted, err := TopoSort(tasks)

	require.NoError(t, err)
	require.Len(t, sorted, 3)

	// task3 should come after both task1 and task2
	task3Index := -1
	task1Index := -1
	task2Index := -1
	for i, task := range sorted {
		switch task.Name {
		case "task1":
			task1Index = i
		case "task2":
			task2Index = i
		case "task3":
			task3Index = i
		}
	}

	assert.True(t, task3Index > task1Index, "task3 should come after task1")
	assert.True(t, task3Index > task2Index, "task3 should come after task2")
}

func TestTopoSort_ComplexDependencies(t *testing.T) {
	// Complex dependency graph:
	// task1 (no deps)
	// task2 (no deps)
	// task3 depends on task1
	// task4 depends on task2
	// task5 depends on task3 and task4
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000001", "task1", nil),
		createTask("00000000-0000-0000-0000-000000000002", "task2", nil),
		createTask("00000000-0000-0000-0000-000000000003", "task3", []string{"task1"}),
		createTask("00000000-0000-0000-0000-000000000004", "task4", []string{"task2"}),
		createTask("00000000-0000-0000-0000-000000000005", "task5", []string{"task3", "task4"}),
	}

	sorted, err := TopoSort(tasks)

	require.NoError(t, err)
	require.Len(t, sorted, 5)

	// Create a position map
	position := make(map[string]int)
	for i, task := range sorted {
		position[task.Name] = i
	}

	// Verify dependencies are satisfied
	assert.True(t, position["task3"] > position["task1"], "task3 should come after task1")
	assert.True(t, position["task4"] > position["task2"], "task4 should come after task2")
	assert.True(t, position["task5"] > position["task3"], "task5 should come after task3")
	assert.True(t, position["task5"] > position["task4"], "task5 should come after task4")
}

func TestTopoSort_CyclicDependency(t *testing.T) {
	// task1 -> task2 -> task3 -> task1 (cycle)
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000001", "task1", []string{"task3"}),
		createTask("00000000-0000-0000-0000-000000000002", "task2", []string{"task1"}),
		createTask("00000000-0000-0000-0000-000000000003", "task3", []string{"task2"}),
	}

	sorted, err := TopoSort(tasks)

	require.Error(t, err)
	assert.Nil(t, sorted)
	assert.Contains(t, err.Error(), "cyclic dependency detected")
}

func TestTopoSort_SelfDependency(t *testing.T) {
	// task1 depends on itself
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000001", "task1", []string{"task1"}),
	}

	sorted, err := TopoSort(tasks)

	require.Error(t, err)
	assert.Nil(t, sorted)
	assert.Contains(t, err.Error(), "cyclic dependency detected")
}

func TestTopoSort_MissingDependency(t *testing.T) {
	// task1 depends on non-existent task2
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000001", "task1", []string{"task2"}),
	}

	sorted, err := TopoSort(tasks)

	require.Error(t, err)
	assert.Nil(t, sorted)
	assert.Contains(t, err.Error(), "dependency 'task2' not found")
	assert.Contains(t, err.Error(), "task 'task1'")
}

func TestTopoSort_EmptyTaskList(t *testing.T) {
	tasks := []db.Task{}

	sorted, err := TopoSort(tasks)

	require.NoError(t, err)
	assert.Empty(t, sorted)
}

func TestTopoSort_SingleTask(t *testing.T) {
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000001", "task1", nil),
	}

	sorted, err := TopoSort(tasks)

	require.NoError(t, err)
	require.Len(t, sorted, 1)
	assert.Equal(t, "task1", sorted[0].Name)
}

func TestTopoSort_DiamondDependency(t *testing.T) {
	// Diamond pattern:
	//     task1
	//    /     \
	// task2   task3
	//    \     /
	//     task4
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000001", "task1", nil),
		createTask("00000000-0000-0000-0000-000000000002", "task2", []string{"task1"}),
		createTask("00000000-0000-0000-0000-000000000003", "task3", []string{"task1"}),
		createTask("00000000-0000-0000-0000-000000000004", "task4", []string{"task2", "task3"}),
	}

	sorted, err := TopoSort(tasks)

	require.NoError(t, err)
	require.Len(t, sorted, 4)

	// Create a position map
	position := make(map[string]int)
	for i, task := range sorted {
		position[task.Name] = i
	}

	// Verify dependencies
	assert.True(t, position["task2"] > position["task1"], "task2 should come after task1")
	assert.True(t, position["task3"] > position["task1"], "task3 should come after task1")
	assert.True(t, position["task4"] > position["task2"], "task4 should come after task2")
	assert.True(t, position["task4"] > position["task3"], "task4 should come after task3")

	// task1 should be first
	assert.Equal(t, "task1", sorted[0].Name)
	// task4 should be last
	assert.Equal(t, "task4", sorted[3].Name)
}

func TestTopoSort_PartialCycleBetweenTasks(t *testing.T) {
	// task1 and task2 depend on each other
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000001", "task1", []string{"task2"}),
		createTask("00000000-0000-0000-0000-000000000002", "task2", []string{"task1"}),
		createTask("00000000-0000-0000-0000-000000000003", "task3", nil),
	}

	sorted, err := TopoSort(tasks)

	require.Error(t, err)
	assert.Nil(t, sorted)
	assert.Contains(t, err.Error(), "cyclic dependency detected")
}

func TestTopoSort_LongChain(t *testing.T) {
	// task1 -> task2 -> task3 -> task4 -> task5
	tasks := []db.Task{
		createTask("00000000-0000-0000-0000-000000000005", "task5", []string{"task4"}),
		createTask("00000000-0000-0000-0000-000000000003", "task3", []string{"task2"}),
		createTask("00000000-0000-0000-0000-000000000001", "task1", nil),
		createTask("00000000-0000-0000-0000-000000000004", "task4", []string{"task3"}),
		createTask("00000000-0000-0000-0000-000000000002", "task2", []string{"task1"}),
	}

	sorted, err := TopoSort(tasks)

	require.NoError(t, err)
	require.Len(t, sorted, 5)

	// Verify the exact order
	assert.Equal(t, "task1", sorted[0].Name)
	assert.Equal(t, "task2", sorted[1].Name)
	assert.Equal(t, "task3", sorted[2].Name)
	assert.Equal(t, "task4", sorted[3].Name)
	assert.Equal(t, "task5", sorted[4].Name)
}
