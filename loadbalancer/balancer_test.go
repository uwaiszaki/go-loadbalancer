package loadbalancer

import (
	"testing"
)

func TestLoadBalancer(t *testing.T) {
	StartLoadBalancer(
		WithWorkerQueueSize(10),
		WithWorkers(2),
		WithRoundRobin,
	)
}
