package utils

import (
	"fmt"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

// sort the tasks inorder of execuations depending on which depends on the other
func TopoSort(tasks []db.Task) ([]db.Task, error) {
	idToTask := make(map[string]db.Task)
	nameToID := make(map[string]string)

	for _, task := range tasks {
		taskID := task.ID.String()
		idToTask[taskID] = task
		nameToID[task.Name] = taskID
	}

	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize in-degree for all tasks
	for _, task := range tasks {
		taskID := task.ID.String()
		inDegree[taskID] = 0
	}

	// Build the dependency graph
	for _, task := range tasks {
		taskID := task.ID.String()
		for _, depName := range task.Config.DependsOn {
			depID, exists := nameToID[depName]
			if !exists {
				return nil, fmt.Errorf("dependency '%s' not found for task '%s'", depName, task.Name)
			}
			graph[depID] = append(graph[depID], taskID)
			inDegree[taskID]++
		}
	}

	// Find tasks with no dependencies (in-degree 0)
	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	// Perform topological sort
	var sorted []db.Task
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		sorted = append(sorted, idToTask[curr])

		for _, neighbor := range graph[curr] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(sorted) != len(tasks) {
		return nil, fmt.Errorf("cyclic dependency detected")
	}

	return sorted, nil
}
