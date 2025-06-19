package queue

type TaskQueue interface {
    Enqueue(job string) error
    Dequeue() (string, error)
}
