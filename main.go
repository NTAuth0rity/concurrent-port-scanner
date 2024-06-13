package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

func main() {
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
		for i := 1; i <= maxPort; i++ {
			address := fmt.Sprintf("%s:%d", *addressFlag, i)

			conn, err := net.DialTimeout("tcp", address, time.Duration(*timeoutFlag)*time.Second)
			if err != nil {
				if *verboseFlag {
					fmt.Printf("Port %d closed.\n", i)
				}
				continue
			} else {
				fmt.Printf("Port %d open.\n", i)
				conn.Close()
			}
		}
	}
}
