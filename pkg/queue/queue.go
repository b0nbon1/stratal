package queue

type TaskQueue interface {
    Enqueue(task any)
    Dequeue() (string, error)
	XReadGeneric(lastID string, block int64, mapper func(map[string]interface{}) (any, error)) ([]any, error)
}
