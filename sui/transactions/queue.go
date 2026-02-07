package transactions

import "sync"

type SerialQueue struct {
	mu      sync.Mutex
	running bool
	queue   []func()
}

func (q *SerialQueue) RunTask(task func() error) error {
	ch := make(chan error, 1)
	wrapped := func() { ch <- task() }
	q.mu.Lock()
	q.queue = append(q.queue, wrapped)
	if !q.running {
		q.running = true
		go q.drain()
	}
	q.mu.Unlock()
	return <-ch
}

func (q *SerialQueue) drain() {
	for {
		q.mu.Lock()
		if len(q.queue) == 0 {
			q.running = false
			q.mu.Unlock()
			return
		}
		task := q.queue[0]
		q.queue = q.queue[1:]
		q.mu.Unlock()
		task()
	}
}

type ParallelQueue struct {
	sem chan struct{}
}

func NewParallelQueue(maxTasks int) *ParallelQueue {
	if maxTasks <= 0 {
		maxTasks = 1
	}
	return &ParallelQueue{sem: make(chan struct{}, maxTasks)}
}

func (q *ParallelQueue) RunTask(task func() error) error {
	q.sem <- struct{}{}
	defer func() { <-q.sem }()
	return task()
}
