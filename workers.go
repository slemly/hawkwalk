package main

import (
	"sync"
	"time"
)

func init_port_queue(ports []int) chan int {
	ch := make(chan int, len(ports))
	defer close(ch)
	for _, val := range ports {
		ch <- val
	}
	return ch
}

type scanfn func(string, int, time.Duration) ScanResult

func scan_and_supply_result_to_channel(
	portsch chan int,
	host string,
	timeout time.Duration,
	fn scanfn,
	resch chan ScanResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for p := range portsch {
		resch <- fn(host, p, timeout)
	}
}
