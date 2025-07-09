package utils

import (
	"fmt"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

// sort the tasks inorder of execuations depending on which depends on the other
func TopoSort(tasks []db.Task) ([]db.Task, error) {
	idToTask := make(map[string]db.Task)
	for _, task := range tasks {
		idToTask[task.ID.String()] = task
	}

	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	for _, task := range tasks {
		for _, dep := range task.Config.DependsOn {
			graph[dep] = append(graph[dep], task.ID.String())
			inDegree[task.ID.String()]++
		}
		if _, ok := inDegree[task.ID.String()]; !ok {
			inDegree[task.ID.String()] = 0
		}
	}

	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

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
