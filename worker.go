package main

import (
	"fmt"
)

func CreateWorker(id int, WorkerPool chan chan WorkRequest) Worker {
	w := Worker{
		ID:         id,
		Work:       make(chan WorkRequest),
		WorkerPool: WorkerPool,
		QuitChan:   make(chan bool),
	}
	return w
}

func (w Worker) Start() {

	go func() {
		for {
			fmt.Println("Worker waiting for work - ", w.ID)
			w.WorkerPool <- w.Work
			select {
			case work := <-w.Work:
				work.Call(work, w)
			case <-w.QuitChan:
				fmt.Println("Worker stopping - ", w.ID)
			}
		}
	}()

}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}