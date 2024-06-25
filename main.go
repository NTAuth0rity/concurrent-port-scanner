package main

import (
	"flag"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	addressFlag := flag.String("a", "", "Provide address to be scanned.")
	rangeFlag := flag.Bool("f", false, "If used all 65535 ports will be scanned, otherwise only the first 1024 will be tested.")
	timeoutFlag := flag.Int("t", 2, "Timeout duration in seconds (default: 2)")
	verboseFlag := flag.Bool("v", false, "Enable verbose output")

	flag.Parse()

	if *addressFlag != "" {
		var maxPort int
		if *rangeFlag {
			maxPort = 65535
		} else {
			maxPort = 1024
		}

		ports := make(chan int, maxPort)
		results := make(chan int, maxPort)
		numWorkers := runtime.NumCPU()

		for i := 0; i < numWorkers; i++ { // using a limited number of goroutines
			wg.Add(1)
			go worker(*addressFlag, *timeoutFlag, *verboseFlag, ports, results, &wg)
		}

		go func() {
			for j := 1; j <= maxPort; j++ {
				ports <- j
			}
			close(ports) // close "ports" channel
		}()

		go func() { // close "results" channel when all workers are done
			wg.Wait()
			close(results)
		}()

		for port := range results {
			if port != 0 {
				fmt.Printf("Port %d open.\n", port)
			}
		}

	} else {
		fmt.Println("Please provide an address to be scanned using the -a flag.")
		return
	}
}

func worker(address string, timeout int, verbose bool, ports, results chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range ports {
		addr := fmt.Sprintf("%s:%d", address, port)
		conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)

		if err != nil {
			if verbose {
				fmt.Printf("Port %d closed.\n", port)
			}
			results <- 0 // Send 0 to indicate the port is closed
			continue
		}

		conn.Close()
		results <- port
	}
}
