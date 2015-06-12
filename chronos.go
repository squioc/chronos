package chronos

import (
	"container/heap"
	"github.com/squioc/axis"
)

// Chronos is structure that defines the scheduler
type Chronos struct {
	provider axis.Provider
	queue    Queue
	watcher  axis.Watcher
}

// NewChronos creates a new Scheduler from an axis.Provider and a chronos.Queue
func NewChronos(provider axis.Provider, queue Queue) *Chronos {
	return &Chronos{
		provider: provider,
		queue:    queue,
	}
}

// Run starts the scheduler
func (c *Chronos) Run(pushChan, workerChan chan Entry, stopChan chan bool) {
	loopChan := make(chan axis.Position, 1)
	for {
		select {
		case entry := <-pushChan:
			delay := c.delay(entry.Position())

			// Send the element on the worker channel if the delay is negative
			if delay < 0 {
				workerChan <- entry
			} else {
				// Push the element in the queue
				heap.Push(c.queue, entry)

				// If the watcher is undefined
				if c.watcher == nil {
					c.watcher = c.provider.AfterChan(delay, loopChan)
				}

				// Pop to the next delay in the queue
				c.popUntilNextDelay(workerChan)
			}

		case <-loopChan:
			// Send the miminum element on the worker channel
			workerChan <- heap.Pop(c.queue).(Entry)

			// Send the following elements on the worker channel
			// as long as the delay is negative
			for c.queue.Len() > 0 {
				firstEntry := c.queue.Peek()
				delay := c.delay(firstEntry.(Entry).Position())
				if delay <= 0 {
					workerChan <- heap.Pop(c.queue).(Entry)
				} else {
					// Pause the shipment during the delay to elapse
					// Reengage the watcher
					if c.watcher == nil {
						c.watcher = c.provider.AfterChan(delay, loopChan)
					} else {
						c.watcher.Reset(delay)
					}
					break
				}
			}

		case <-stopChan:
			return

		}
	}
}

func (c *Chronos) delay(position axis.Position) axis.Distance {
	return axis.Distance(position - c.provider.Current())
}

func (c *Chronos) popUntilNextDelay(workerChan chan Entry) {
	for c.queue.Len() > 0 {
		firstEntry := c.queue.Peek()

		delay := c.delay(firstEntry.(Entry).Position())
		if delay <= 0 {
			workerChan <- heap.Pop(c.queue).(Entry)
		} else {
			c.watcher.Reset(delay)
			return
		}
	}
}
