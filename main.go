package main

import (
	"flag"
	"fmt"
	"sync"
)

func write_to_file() {

}

func main() {

	var wg sync.WaitGroup
	flag.Parse()
	conf, _ := buildConfigs()

	port_queue := init_port_queue(conf.ports)
	ch_results := make(chan ScanResult, len(conf.ports))

	for range conf.workers {
		wg.Add(1)
		go scan_and_supply_result_to_channel(
			port_queue,
			conf.host,
			conf.timeout,
			scanPort,
			ch_results,
			&wg,
		)
	}
	wg.Wait()
	close(ch_results)
	fmt.Println("Results: ")
	for it := range ch_results {
		fmt.Println(it)
	}
}
