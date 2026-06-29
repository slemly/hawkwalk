package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sync"
)

func write_all_to_file(fname string, addr string, results chan ScanResult) {
	cwd, wd_err := os.Getwd()

	if wd_err == nil {
		fmt.Println("Could not determine current working directory. Content will not be written to file.")
		return
	}

	fpath := path.Join(cwd, fname)

	f, fpath_err := os.Create(fpath)

	if fpath_err != nil {
		fmt.Println("Could not create output file.")
		return
	}

	fmt.Fprintf(f, "Results for %s:\n\n", addr)
	for it := range results {
		var st = it.String()
		fmt.Fprintf(f, "\t%s\n", st)
	}

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

	fname := conf.outputFile
	if fname != "" {
		write_all_to_file(fname, conf.host, ch_results)
	} else {
		for it := range ch_results {
			if it.open {
				fmt.Println(it.String())
			}
		}
	}

}
