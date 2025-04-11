package queue

type TaskQueue interface {
    Enqueue(task string) error
    Dequeue() (string, error)
	XRead(lastID string, block int64) ([]string, error)
}
