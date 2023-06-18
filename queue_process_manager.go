package queue_process_manager

import "sync"

// defaultMaxAllowance is the default setting for the queue manager
const defaultMaxAllowance = 10

// QueueItem has to properties QueueTaskFunc and QueueFuncArgs the QueueFuncArgs are automatically passed into the QueueFunc
type QueueItem struct {
	QueueFunc     QueueTaskFunc
	QueueFuncArgs []any
}

type QueueTaskFunc func(args []any)
type QueueTaskConfigFunc func(q *Queue)

type Queue struct {
	countRWMutex      sync.RWMutex
	count             int
	queue             []*QueueItem
	maxQueueAllowance int
}

func NewQueue(configs ...QueueTaskConfigFunc) *Queue {
	q := &Queue{
		maxQueueAllowance: defaultMaxAllowance,
	}

	for _, configFunc := range configs {
		configFunc(q)
	}

	return q
}

// WithCustomMaxQueueAllowance takes an int which is used to set the maximum number of items to run in the queue in parallel
func WithCustomMaxQueueAllowance(max int) QueueTaskConfigFunc {
	return func(q *Queue) {
		q.maxQueueAllowance = max
	}
}

func (q *Queue) incrementQueuedCount() {
	q.countRWMutex.Lock()
	q.count += 1
	q.countRWMutex.Unlock()
}

func (q *Queue) decreaseQueuedCount() {
	q.countRWMutex.RLock()
	defer q.countRWMutex.RUnlock()
	q.count -= 1
}

func (q *Queue) getQueuedCount() int {
	q.countRWMutex.RLock()
	defer q.countRWMutex.RUnlock()
	return q.count
}

// AddToQueue adds your queue items
func (q *Queue) AddToQueue(qI *QueueItem) {
	q.queue = append(q.queue, qI)
}

// ProcessQueue runs the queue process
func (q *Queue) ProcessQueue() {
	for _, item := range q.queue {
		q.incrementQueuedCount()

		go func(item *QueueItem) {
			complete := make(chan struct{}, 1)
			go func(item *QueueItem, complete chan<- struct{}) {
				item.QueueFunc(item.QueueFuncArgs)
				complete <- struct{}{}
			}(item, complete)

			select {
			case <-complete:
				q.decreaseQueuedCount()
				break
			}
		}(item)

		for {
			if q.getQueuedCount() < q.maxQueueAllowance {
				break
			}
		}
	}

	for {
		if q.getQueuedCount() == 0 {
			break
		}
	}
}
