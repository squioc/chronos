package chronos

import (
	"fmt"
	"github.com/squioc/axis"
	"testing"
	"time"
)

/****************************************
*                                       *
*                 Tests                 *
*                                       *
****************************************/

func expect(t *testing.T, expected EntryTest, actual Entry, stopChan chan bool) {
	if actual == nil {
		stopChan <- true
		t.Fatalf("The actual element wasn't expected to be nil")
	}
	if actual.(*EntryTest).Element != expected.Element {
		stopChan <- true
		t.Fatalf("The actual element wasn't the expected one")
	}
}

func TestRunWithJobsInOrder(t *testing.T) {
	fmt.Println("Test Run with jobs in order should send in order")
	// Arrange
	position := axis.Position(0)
	newPosition := axis.Position(1500)
	provider := axis.NewFakeTime(position)
	pq := new(PriorityQueue)
	chronos := NewChronos(provider, pq)
	pushChan := make(chan Entry, 2)
	workerChan := make(chan Entry, 2)
	stopChan := make(chan bool, 1)
	firstEntry := &EntryTest{
		position: axis.Position(500),
		Element:  "First",
	}
	secondEntry := &EntryTest{
		position: axis.Position(1000),
		Element:  "Second",
	}

	// Act
	go chronos.Run(pushChan, workerChan, stopChan)
	pushChan <- firstEntry
	pushChan <- secondEntry
	// Updates the position
	go provider.Update(newPosition)

	// Assert
	firstElement := <-workerChan
	expect(t, *firstEntry, firstElement, stopChan)
	secondElement := <-workerChan
	expect(t, *secondEntry, secondElement, stopChan)
	stopChan <- true
}

func TestRunWithJobsInReverseOrder(t *testing.T) {
	fmt.Println("Test Run with jobs in reverse order should send in reverse order")
	// Arrange
	position := axis.Position(0)
	newPosition := axis.Position(1500)
	provider := axis.NewFakeTime(position)
	pq := new(PriorityQueue)
	chronos := NewChronos(provider, pq)
	pushChan := make(chan Entry, 2)
	workerChan := make(chan Entry, 2)
	stopChan := make(chan bool, 1)
	firstEntry := &EntryTest{
		position: axis.Position(1000),
		Element:  "First",
	}
	secondEntry := &EntryTest{
		position: axis.Position(500),
		Element:  "Second",
	}

	// Act
	go chronos.Run(pushChan, workerChan, stopChan)
	pushChan <- firstEntry
	pushChan <- secondEntry
	// Updates the position
	go provider.Update(newPosition)

	// Assert
	firstElement := <-workerChan
	expect(t, *secondEntry, firstElement, stopChan)
	secondElement := <-workerChan
	expect(t, *firstEntry, secondElement, stopChan)
	stopChan <- true
}

func TestRunWithStop(t *testing.T) {
	fmt.Println("Test Run with Stop should exit")
	// Arrange
	position := axis.Position(0)
	provider := axis.NewFakeTime(position)
	pq := new(PriorityQueue)
	chronos := NewChronos(provider, pq)
	pushChan := make(chan Entry, 2)
	workerChan := make(chan Entry, 2)
	stopChan := make(chan bool, 1)
	firstEntry := &EntryTest{
		position: axis.Position(1000),
		Element:  "First",
	}

	// Act
	exitChan := make(chan bool, 1)
	// go routine to check that we leave the Run method when we send a booleen on stopChan
	go func(exitChan chan bool) {
		chronos.Run(pushChan, workerChan, stopChan)
		// Send boolean on the channel when the Run method exits
		exitChan <- true
	}(exitChan)
	// Lets the goroutine starts then stop the goroutine
	time.Sleep(5)
	stopChan <- true
	// Sends an entry to the goroutine
	pushChan <- firstEntry

	// Assert
	select {
	case <-exitChan:
		// PASS
		return
	case <-workerChan:
		// FAIL, we expected to exit, not to receive an item
		t.Fatalf("Expected exit, not item")
	case <-time.After(2 * time.Second):
		// FAIL, timeout
		t.Fatalf("Timeout. the test exceed the expected duration")
	}
}

/****************************************
*                                       *
*        Structs implementations        *
*                                       *
****************************************/

type EntryTest struct {
	position axis.Position
	Element  interface{}
	index    int
}

func (e EntryTest) Position() axis.Position {
	return e.position
}

type PriorityQueue []*EntryTest

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
	ll[i].index = i
	ll[j].index = j
}

func (pq *PriorityQueue) Push(v interface{}) {
	n := len(*pq)
	item := v.(*EntryTest)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	ll := *pq
	n := len(ll)
	item := ll[n-1]
	item.index = -1
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
