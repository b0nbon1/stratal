package processor

import (
	"fmt"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
)

func ProcessJob(job []db.Task, JobId string) {
	fmt.Println("Processing job:", JobId)

	sortedTasks, err := utils.TopoSort(job)
	if err != nil {
		fmt.Println("Error sorting tasks:", err)
		return
	}
	fmt.Println("Sorted tasks:", sortedTasks)

	for _, task := range sortedTasks {
	err := ExecuteTask(task)
	if err != nil {
		fmt.Printf("Task %s failed: %v", task.ID, err)
		// Optionally stop or continue based on failure policy
		break
	}
}
	

}
