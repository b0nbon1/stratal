package utils

import (
	"fmt"

	"github.com/b0nbon1/stratal/db/dto"
)

// sort the tasks inorder of execuations depending on which depends on the other
func TopoSort(tasks []dto.TaskConfig) ([]dto.TaskConfig, error) {
	idToTask := make(map[string]dto.TaskConfig)
	for _, task := range tasks {
		idToTask[task.ID] = task
	}

	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	for _, task := range tasks {
		for _, dep := range task.DependsOn {
			graph[dep] = append(graph[dep], task.ID)
			inDegree[task.ID]++
		}
		if _, ok := inDegree[task.ID]; !ok {
			inDegree[task.ID] = 0
		}
	}

	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	var sorted []dto.TaskConfig
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
