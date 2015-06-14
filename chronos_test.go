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

func expect(t *testing.T, expected Item, actual Entry, stopChan chan bool) {
	if actual == nil {
		stopChan <- true
		t.Fatalf("The actual element wasn't expected to be nil")
	}
	if actual.(*Item).Value != expected.Value {
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
	firstEntry := &Item{
		Priority: axis.Position(500),
		Value:    "First",
	}
	secondEntry := &Item{
		Priority: axis.Position(1000),
		Value:    "Second",
	}

	// Act
	go chronos.Run(pushChan, workerChan, stopChan)
	pushChan <- firstEntry
	pushChan <- secondEntry
	// Updates the position
	go provider.Update(newPosition)

	// Assert
	firstValue := <-workerChan
	expect(t, *firstEntry, firstValue, stopChan)
	secondValue := <-workerChan
	expect(t, *secondEntry, secondValue, stopChan)
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
	firstEntry := &Item{
		Priority: axis.Position(1000),
		Value:    "First",
	}
	secondEntry := &Item{
		Priority: axis.Position(500),
		Value:    "Second",
	}

	// Act
	go chronos.Run(pushChan, workerChan, stopChan)
	pushChan <- firstEntry
	pushChan <- secondEntry
	// Updates the position
	go provider.Update(newPosition)

	// Assert
	firstValue := <-workerChan
	expect(t, *secondEntry, firstValue, stopChan)
	secondValue := <-workerChan
	expect(t, *firstEntry, secondValue, stopChan)
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
	firstEntry := &Item{
		Priority: axis.Position(1000),
		Value:    "First",
	}

	// Act
	exitChan := make(chan bool, 1)
	// go routine to check that we leave the Run method when we send a booleen on stopChan
	go func(exitChan chan bool) {
		chronos.Run(pushChan, workerChan, stopChan)
		// sends boolean on the channel when the Run method exits
		exitChan <- true
	}(exitChan)
	// stops the goroutine
	stopChan <- true
	// sends an entry to the goroutine
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

func TestRunWithJobsIn2Sets(t *testing.T) {
	fmt.Println("Test Run with jobs in 2 sets")
	// Arrange
	position := axis.Position(0)
	intermediatePosition := axis.Position(1500)
	terminalPosition := axis.Position(3000)
	provider := axis.NewFakeTime(position)
	pq := new(PriorityQueue)
	chronos := NewChronos(provider, pq)
	pushChan := make(chan Entry, 2)
	workerChan := make(chan Entry, 2)
	stopChan := make(chan bool, 1)
	firstEntry := &Item{
		Priority: axis.Position(500),
		Value:    "First",
	}
	secondEntry := &Item{
		Priority: axis.Position(1000),
		Value:    "Second",
	}
	thirdEntry := &Item{
		Priority: axis.Position(2000),
		Value:    "Third",
	}

	// Act
	go chronos.Run(pushChan, workerChan, stopChan)
	// pushes the first set (firstEntry, secondEntry
	pushChan <- firstEntry
	pushChan <- secondEntry
	// lets the goroutine starts then updates the position
	time.Sleep(5 * time.Millisecond)
	provider.Update(intermediatePosition)
	time.Sleep(5 * time.Millisecond)
	// pushes the second set (thridEntry)
	pushChan <- thirdEntry
	time.Sleep(5 * time.Millisecond)
	provider.Update(terminalPosition)

	// Assert
	firstValue := <-workerChan
	expect(t, *firstEntry, firstValue, stopChan)
	secondValue := <-workerChan
	expect(t, *secondEntry, secondValue, stopChan)
	thirdValue := <-workerChan
	expect(t, *thirdEntry, thirdValue, stopChan)
	stopChan <- true
}

func TestRunWithJobInPast(t *testing.T) {
	fmt.Println("Test Run with a job in a past position")
	// Arrange
	position := axis.Position(1000)
	provider := axis.NewFakeTime(position)
	pq := new(PriorityQueue)
	chronos := NewChronos(provider, pq)
	pushChan := make(chan Entry, 2)
	workerChan := make(chan Entry, 2)
	stopChan := make(chan bool, 1)
	firstEntry := &Item{
		Priority: 500,
		Value:    "First",
	}

	// Act
	go chronos.Run(pushChan, workerChan, stopChan)
	// pushes the item
	pushChan <- firstEntry

	// Assert
	firstValue := <-workerChan
	expect(t, *firstEntry, firstValue, stopChan)
	stopChan <- true
}
