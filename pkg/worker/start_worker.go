package worker

import (
	// "encoding/json"
	// "fmt"
	"log"

	"github.com/b0nbon1/temporal-lite/pkg/queue"
	"github.com/b0nbon1/temporal-lite/utils"
)

func StartWorker(q queue.TaskQueue) {
	go func() {
		for {
			_, err := q.XReadGeneric("0", 5000, utils.TaskMapper)
			if err != nil {
				log.Println("Error reading from queue:", err)
				continue
			}
			// fmt.Println("Length Tasks:", tasks)
		// 	for _, item := range tasks {
		// 		jobMap, ok := item.(map[string]any)
		// 		if !ok {
		// 			fmt.Println("Unexpected type", item)
		// 			continue
		// 		}
		// 		fmt.Println("Job Map:", jobMap)
		// 		fmt.Println("ID:", jobMap["id"])
		// 		fmt.Println("Name:", jobMap["name"])
		// 		fmt.Println("Retries:", jobMap["retries"])
		// 	}
		}
	}()
	
}
