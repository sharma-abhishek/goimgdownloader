package main

var WorkerPool chan chan WorkRequest

func StartDispatcher(n int) {
	WorkerPool = make(chan chan WorkRequest, n)

	for i := 0; i < n; i++ {
		//create worker object and start
		w := CreateWorker((i + 1), WorkerPool)
		w.Start()
	}

	// function that takes work from worker queue, pull the worker and assign work to worker
	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerPool
					worker <- work
				}()
			}
		}
	}()
}
