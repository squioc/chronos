package chronos

import (
	"github.com/squioc/axis"
)

// Item is the struct that wraps the interface
// manages in the queue
type Item struct {
	Priority axis.Position
	Value    interface{}
	Index    int
}

func (e Item) Position() axis.Position {
	return e.Priority
}

// PriorityQueue holds a set of Items
type PriorityQueue []*Item

func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	ll := *pq
	return ll[i].Position() < ll[j].Position()
}

func (pq *PriorityQueue) Swap(i, j int) {
	ll := *pq
	ll[i], ll[j] = ll[j], ll[i]
	ll[i].Index = i
	ll[j].Index = j
}

func (pq *PriorityQueue) Push(v interface{}) {
	n := len(*pq)
	item := v.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	ll := *pq
	n := len(ll)
	item := ll[n-1]
	item.Index = -1
	*pq = ll[0 : n-1]
	return item
}

func (pq *PriorityQueue) Peek() interface{} {
	ll := *pq
	n := len(ll)
	if n <= 0 {
		return nil
	}
	item := ll[n-1]
	return item
}
