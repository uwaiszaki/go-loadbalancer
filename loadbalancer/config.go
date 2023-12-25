package loadbalancer

import "time"

type LoadBalancerStrategy int

const (
	RoundRobin LoadBalancerStrategy = iota
	MinConnections
)

type BalancerConfig struct {
	Workers             int
	WorkerQueueSize     int
	Strategy            LoadBalancerStrategy
	WorkerTaskTimeout   time.Duration
	WorkerPulseInterval time.Duration
}

type BalancerConfigFn func(*BalancerConfig)

func WithWorkerTaskTimeout(t time.Duration) BalancerConfigFn {
	return func(conf *BalancerConfig) {
		conf.WorkerTaskTimeout = t
	}
}

func WithWorkerPulseInterval(t time.Duration) BalancerConfigFn {
	return func(conf *BalancerConfig) {
		conf.WorkerPulseInterval = t
	}
}

func WithWorkerQueueSize(n int) BalancerConfigFn {
	return func(conf *BalancerConfig) {
		conf.WorkerQueueSize = n
	}
}

func WithWorkers(n int) BalancerConfigFn {
	return func(conf *BalancerConfig) {
		conf.Workers = n
	}
}

func WithRoundRobin(conf *BalancerConfig) {
	conf.Strategy = RoundRobin
}

func WithMinConnectionsStrategy(conf *BalancerConfig) {
	conf.Strategy = MinConnections
}
