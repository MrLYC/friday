package sentry

import (
	"container/heap"
	"sync"
	"time"
)

// DelayEvent :
type DelayEvent struct {
	Event *Event
	Time  time.Time
}

// An qItem is something we manage in a priority queue.
type qItem struct {
	value *DelayEvent // The value of the item; arbitrary.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A priorityQueue implements heap.Interface and holds Items.
type priorityQueue []*qItem

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	ei := pq[i].value
	ej := pq[j].value
	return ei.Time.Before(ej.Time)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*qItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Top() interface{} {
	old := *pq
	n := len(old)
	return old[n-1]
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Timer :
type Timer struct {
	BaseTrigger
	eventQueue priorityQueue
	queueMux   sync.Mutex
}

// Init :
func (t *Timer) Init(s *Sentry) {
	t.eventQueue = make(priorityQueue, 0, 10)
	heap.Init(&(t.eventQueue))
	t.BaseTrigger.Init(s)
}

// AddEvent :
func (t *Timer) AddEvent(event *DelayEvent) {
	t.queueMux.Lock()
	defer t.queueMux.Unlock()
	item := &qItem{
		value: event,
	}
	heap.Push(&(t.eventQueue), item)
	heap.Fix(&(t.eventQueue), item.index)
}

// PopEvent :
func (t *Timer) PopEvent() *DelayEvent {
	t.queueMux.Lock()
	defer t.queueMux.Unlock()
	item := heap.Pop(&(t.eventQueue)).(*qItem)
	return item.value
}
