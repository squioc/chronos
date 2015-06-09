# LibChronos

    import "github.com/squioc/chronos"

The chronos package provides a job scheduler implementation

## Documentation:

The documentation is available at [https://godoc.org/github.com/squioc/chronos](https://godoc.org/github.com/squioc/chronos), thanks to [GoDoc](http://godoc.org)

## Usage

### func NewChronos
```go
func NewChronos(provider axis.Provider, queue Queue) *Chronos
```
Creates a new Scheduler with an axis.Provider and a chronos.Queue.

### func Run
```go
func (c *Chronos) Run(pushChan, workerChan chan Entry, stopChan chan bool)
```
Starts the scheduler (we recommend to use it in a goroutine).

Parameters:

 - `pushChan` is the input channel, which delivers jobs to schedule.
 - `workerChan` is the output channel, which receives jobs whose the position was reached.
 - `stopChan` allow to stop gracefully the scheduler.


