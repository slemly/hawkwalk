package main

import (
	"fmt"
	"net"
	"time"
)

type ScanResult struct {
	port   int
	open   bool
	banner string
	err    error
}

func scanPort(host string, port int, timeout time.Duration) ScanResult {
	fmt.Printf("Scanning %s:%d (timeout %s)\n", host, port, timeout.String())
	isOpen := false
	banner := ""
	var err error
	var conn, conn_err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if conn_err != nil || conn == nil {
		isOpen = false
		err = conn_err
	} else {
		buff := make([]byte, 256)
		defer conn.Close()
		isOpen = true
		conn.SetReadDeadline(time.Now().Add(timeout))
		count, read_err := conn.Read(buff)
		if read_err != nil {
			fmt.Println("Could not read from buffer.")
			err = read_err
		} else {
			banner = string(buff[:count])
			fmt.Printf("Read %d bytes from port\n", count)
			isOpen = true
		}
	}
	return ScanResult{
		open:   isOpen,
		banner: banner,
		port:   port,
		err:    err,
	}
}
