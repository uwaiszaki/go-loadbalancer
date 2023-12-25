package main

import (
	"github.com/uwaiszaki/go-loadbalancer/loadbalancer"
)

// 1
// Load balancer ->
// Server Running on port 8050
// Pool of workers ([]Worker)
// Health Check for the workers
// Random assign load to any worker

// 2
// Min Connection load balancing

func main() {
	loadbalancer.StartLoadBalancer(
		loadbalancer.WithWorkerQueueSize(10),
		loadbalancer.WithWorkers(2),
		loadbalancer.WithRoundRobin,
	)
}
