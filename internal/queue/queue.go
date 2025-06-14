package queue

type TaskQueue interface {
    Enqueue(job string)
    Dequeue() (string, error)
}
