package chronos

import (
	"container/heap"
	"github.com/squioc/axis"
)

// Entry is the interface that wraps elements
// to send in the scheduler
type Entry interface {
	// Return the position, on the axis, of the entry
	Position() axis.Position
}

// Queue is the interface that holds elements
// as an ordered collection
type Queue interface {
	heap.Interface
	// Return the first element of the queue without remove it
	Peek() interface{}
}
