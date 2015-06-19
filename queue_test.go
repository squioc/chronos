package chronos

import (
	"container/heap"
	"fmt"
	"github.com/squioc/axis"
	"testing"
)

func TestSort(t *testing.T) {
	fmt.Println("Test sorting of the queue")
	// Arrange
	pq := new(PriorityQueue)
	heap.Init(pq)
	firstEntry := &Item{
		Priority: axis.Position(1000),
		Value:    "First",
	}
	secondEntry := &Item{
		Priority: axis.Position(200),
		Value:    "Second",
	}

	// Act
	heap.Push(pq, firstEntry)
	heap.Push(pq, secondEntry)

	// Assert
	if !pq.Less(0, 1) {
		t.Fatalf("The queue is not sorted")
	}
}

func TestPeek(t *testing.T) {
	fmt.Println("Test the peek on the first item of the queue")
	// Arrange
	pq := new(PriorityQueue)
	heap.Init(pq)
	firstEntry := &Item{
		Priority: axis.Position(1000),
		Value:    "First",
	}
	secondEntry := &Item{
		Priority: axis.Position(200),
		Value:    "Second",
	}

	// Act
	heap.Push(pq, firstEntry)
	heap.Push(pq, secondEntry)
	firstItem := pq.Peek()

	// Assert
	if firstItem.(*Item).Position() != secondEntry.Position() || firstItem.(*Item).Value != secondEntry.Value {
		t.Fatalf("Unexpected item as the first item of the queue. expected:%v got:%v", secondEntry, firstItem.(*Item))
	}
}

func TestPop(t *testing.T) {
	fmt.Println("Test the pop on the first item of the queue")
	// Arrange
	pq := new(PriorityQueue)
	heap.Init(pq)
	firstEntry := &Item{
		Priority: axis.Position(1000),
		Value:    "First",
	}
	secondEntry := &Item{
		Priority: axis.Position(200),
		Value:    "Second",
	}

	// Act
	heap.Push(pq, firstEntry)
	heap.Push(pq, secondEntry)
	firstItem := heap.Pop(pq)

	// Assert
	if firstItem.(*Item).Position() != secondEntry.Position() || firstItem.(*Item).Value != secondEntry.Value {
		t.Fatalf("Unexpected item as the first item of the queue. expected:%v got:%v", secondEntry, firstItem.(*Item))
	}
}
