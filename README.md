# workeramp
Worker pool that ramps us the workers and then ramps down

## Theory
The worker pool runs the callback that is provided to it.
It checks the return of the callback.
If the callback returns > 0, it adds a worker.

This usually works with a Channel where the callback-func returns the lenth of the channel.
If the channel has more work, the return value would increase the workers - rampup.
If the channel has no more work, the worker thread/goroutine would go away - for rampdown.

## Use

// ..
echan := make(chan int, 9)
//..
ws := NewWorkerSet("Events", 6, 1000, work)
//..


func work() (int, error) {
select {
	case a := <-echan:
		do_something()
	case <-time.After(5 * time.Second):
	}
  return len(echan), nil
  }
