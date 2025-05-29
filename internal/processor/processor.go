package processor

import (
	"fmt"

	db "github.com/b0nbon1/stratal/db/sqlc"
	"github.com/b0nbon1/stratal/utils"
)

func ProcessJob(job db.Job) {
	fmt.Println("Processing job:", job.ID)

	sortedTasks, err := utils.TopoSort(job.Config.Tasks)
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