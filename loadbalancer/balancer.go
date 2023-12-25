package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Balancer struct {
	Config BalancerConfig
	pool   []*Worker

	mu            sync.RWMutex
	nextWorkerIdx int
}

func StartLoadBalancer(opts ...BalancerConfigFn) {
	config := BalancerConfig{
		Workers:             1,
		WorkerQueueSize:     10,
		Strategy:            RoundRobin,
		WorkerPulseInterval: 2 * time.Second,
	}
	for _, optFn := range opts {
		optFn(&config)
	}
	balancer := Balancer{
		Config:        config,
		nextWorkerIdx: 0,
	}
	for id := 1; id <= balancer.Config.Workers; id++ {
		worker := NewWorker(id, balancer.Config.WorkerQueueSize, balancer.Config.WorkerPulseInterval)
		balancer.pool = append(balancer.pool, worker)
		go func() {
			pulse := worker.Start()
			for {
				timeout := time.NewTimer(3 * time.Second)
				select {
				case <-timeout.C:
					fmt.Printf("Worker %v is timedout :: Creating new worker", worker.id)
				case <-pulse:
				}
			}
		}()
	}
	fmt.Println("Load Balancer Started")
	balancer.StartServer()
}

func (b *Balancer) Handler(w http.ResponseWriter, r *http.Request) {
	// 1) Get Worker
	// 2) Worker.requests <- request
	if b.Config.Strategy == RoundRobin {
		b.mu.Lock()
		worker := b.pool[b.nextWorkerIdx]
		b.nextWorkerIdx = (b.nextWorkerIdx + 1) % len(b.pool)
		b.mu.Unlock()
		request := Request{
			w:    w,
			r:    r,
			done: make(chan struct{}),
		}
		worker.requests <- request
		atomic.AddInt32(&worker.pendingRequests, 1)
		<-request.done
	}
}

func (b *Balancer) StartServer() {
	host := "127.0.0.1"
	port := "8050"
	address := host + ":" + port
	server := http.Server{
		Addr:    address,
		Handler: http.HandlerFunc(b.Handler),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start the load balancer")
	}

}
