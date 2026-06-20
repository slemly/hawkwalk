package main

import (
	"errors"
	"flag"
	"strconv"
	"strings"
	"time"
)

type Configs struct {
	ports      []int
	host       string
	timeout    time.Duration
	outputFile string
	workers    uint
}

var (
	flgHost    = flag.String("host", "localhost", "Specify host IP address")
	flgPorts   = flag.String("ports", "1-10", "Specify ports or port range")
	flgWorkers = flag.Uint("workers", 100, "Specify number of workers")
	flgTimeout = flag.Duration("timeout", time.Duration(100000), "Specify timeout duration for scanning") // time is in nanoseconds
	flgOutput  = flag.String("oN", "", "Specify output file")
)

func buildConfigs() (*Configs, error) {
	var port_range, err = parse_port_range(*flgPorts)
	if err != nil {
		return nil, err
	}

	var con = &Configs{
		host:       *flgHost,
		ports:      port_range,
		workers:    *flgWorkers,
		timeout:    *flgTimeout,
		outputFile: *flgOutput,
	}
	return con, nil
}

func parse_port_range(s string) ([]int, error) {
	/*
		Parses any of the following strings into a list of ints, or gives an error if the list has an element that could not be parsed.
		For ranges, excludes final value
			- "80"
			- "80,443,22"
			- "1-1024,80,443"
	*/
	var ret []int
	for _, val := range strings.Split(s, ",") {
		entries := strings.Split(val, "-")
		if len(entries) == 2 {
			lower_bound, low_err := strconv.Atoi(entries[0])
			upper_bound, upp_err := strconv.Atoi(entries[1])
			if upp_err != nil || low_err != nil {
				err := errors.New("Could not parse port arguments. Upper or lower bound of specified range cannot be parsed to integer.")
				return nil, err
			}
			for lower_bound <= upper_bound {
				ret = append(ret, lower_bound)
				lower_bound++
			}
		} else if len(entries) == 1 {
			it, it_err := strconv.Atoi(entries[0])
			if it_err != nil {
				return nil, errors.New("Could not parse port arguments. Listed distinct value cannot be cast to integer.")
			} else {
				ret = append(ret, it)
			}
		} else {
			err := errors.New("Could not parse port arguments. Multiple-hyphenated ranges are not supported.")
			return nil, err
		}
	}
	return ret, nil
}
