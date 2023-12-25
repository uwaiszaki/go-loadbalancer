package loadbalancer

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type Request struct {
	w    http.ResponseWriter
	r    *http.Request
	done chan struct{}
}

type Worker struct {
	id                  int
	requests            chan Request
	pendingRequests     int32
	workerPulseInterval time.Duration
}

func (w *Worker) Start() chan struct{} {
	// Execute the request and return the response
	pulse := make(chan struct{})
	go func() {
		ticker := time.NewTicker(w.workerPulseInterval)
		for {
			select {
			case <-ticker.C:
				pulse <- struct{}{}
			case request, ok := <-w.requests:
				atomic.AddInt32(&w.pendingRequests, -1)
				if !ok {
					continue
				}
				fmt.Println("Got the Request", request)
				time.Sleep(4 * time.Second)
				writer := request.w
				writer.Write([]byte("Hello From the Worker" + " :: " + strconv.Itoa(w.id)))
				close(request.done)
			}
		}
	}()
	return pulse
	// Health Checks
}

func NewWorker(id, workerQueueSize int, workerPulseInterval time.Duration) *Worker {
	return &Worker{
		id:                  id,
		requests:            make(chan Request, workerQueueSize),
		pendingRequests:     0,
		workerPulseInterval: workerPulseInterval,
	}
}

// func (w *Worker)
