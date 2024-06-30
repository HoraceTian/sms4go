package sms4go

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

type Task struct {
	runAt time.Time
	task  func()
}

// PriorityQueue 实现一个优先级队列，最早到期的任务优先级最高
type PriorityQueue []*Task

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].runAt.Before(pq[j].runAt) }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Task)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// DelayQueue 结构体，包含任务队列和协程池
type DelayQueue struct {
	pq          PriorityQueue
	mu          sync.Mutex
	cond        *sync.Cond
	routinePool *ants.Pool
}

// NewDelayQueue 创建一个新的延时队列
func NewDelayQueue(pool *ants.Pool) *DelayQueue {
	dq := &DelayQueue{
		pq:          make(PriorityQueue, 0),
		routinePool: pool,
	}
	dq.cond = sync.NewCond(&dq.mu)
	go dq.run()
	return dq
}

// AddTask 添加任务到延时队列
func (dq *DelayQueue) AddTask(task func(), delay time.Duration) {
	runAt := time.Now().Add(delay)
	dq.mu.Lock()
	heap.Push(&dq.pq, &Task{runAt: runAt, task: task})
	dq.cond.Signal()
	dq.mu.Unlock()
}

// run 处理任务队列中的任务
func (dq *DelayQueue) run() {
	for {
		dq.mu.Lock()
		for dq.pq.Len() == 0 {
			dq.cond.Wait()
		}
		now := time.Now()
		for dq.pq.Len() > 0 {
			task := dq.pq[0]
			if task.runAt.After(now) {
				break
			}
			heap.Pop(&dq.pq)
			dq.mu.Unlock()

			// 提交任务到协程池
			err := dq.routinePool.Submit(task.task)
			if err != nil {
				fmt.Printf("[sms4go] |- Submit task error: %v\n", err)
			}

			dq.mu.Lock()
		}
		dq.mu.Unlock()
	}
}
