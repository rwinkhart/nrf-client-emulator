package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	eAPI "github.com/aristanetworks/goeapi"
	eAPIModule "github.com/aristanetworks/goeapi/module"
	probing "github.com/prometheus-community/pro-bing"
)

// isReachable returns a boolean value to indicate if the
// destination provided as the first user argument is reachable.
func isReachable() bool {
	pinger, _ := probing.NewPinger(os.Args[1])
	pinger.Count = 1
	pinger.Timeout = 1 * time.Second
	_ = pinger.Run()
	if pinger.PacketsRecv > 0 {
		return true
	}
	return false
}

// getFailPorts returns a slice of all fail-able/managed ports
// based on the range provided as the third user argument.
func getFailPorts() []string {
	var failRanges = strings.Split(os.Args[3], ",")
	var failPorts []string
	for _, failRange := range failRanges {
		if strings.Contains(failRange, "-") {
			// store each portion of provided range
			splitRange := strings.Split(failRange, "-")
			var rangePrefix string
			var rangeStart string
			var rangeEnd = splitRange[1]
			for i, char := range splitRange[0] {
				if unicode.IsDigit(char) {
					rangePrefix = splitRange[0][:i]
					rangeStart = splitRange[0][i:]
					break
				}
			}

			// iterate through the range and add each port to failPorts
			rangeEndInt, _ := strconv.Atoi(rangeEnd)
			for i, _ := strconv.Atoi(rangeStart); i <= rangeEndInt; i++ {
				failPorts = append(failPorts, rangePrefix+strconv.Itoa(i))
			}

		} else {
			// add directly to failPorts
			failPorts = append(failPorts, failRange)
		}
	}
	return failPorts
}

func main() {
	// ensure /mnt/flash/eapi.conf exists
	_, err := os.Stat("/mnt/flash/eapi.conf")
	if os.IsNotExist(err) {
		// configuration does not exist; generate a default one
		fmt.Println("/mnt/flash/eapi.conf does not exist; generating template and exiting...")
		f, _ := os.Create("/mnt/flash/eapi.conf")
		f.WriteString("[connection:nrfLocalClient]\nhost=127.0.0.1\nusername=admin\npassword=nrf\nenablepwd=\ntransport=http\n")
		f.Close()
		os.Exit(0)
	}

	// interpret command-line arguments
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <nrf server IP> <interval (seconds)> <fail range>\n", os.Args[0])
		os.Exit(1)
	}
	interval, _ := strconv.Atoi(os.Args[2])
	failPorts := getFailPorts()

	// connect to local switch via eAPI
	node, _ := eAPI.ConnectTo("nrfLocalClient")
	intHandler := eAPIModule.Interface(node)

	var portsDown = true // default to true to ensure ports are always set to up on boot
	for {
		if isReachable() {
			// the NRF server is reachable
			// check if ports are down
			// restore them if they are
			fmt.Println("Reachable!")
			if portsDown {
				fmt.Println("Ports down, restoring...")
				// restore ports
				for _, port := range failPorts {
					fmt.Println("Restoring " + port)
					intHandler.SetShutdown(port, false)
				}
				portsDown = false
			}
		} else {
			// the NRF server is not reachable
			// check if ports are up
			// down them if they are
			fmt.Println("Unreachable!")
			if !portsDown {
				fmt.Println("Ports up, failing...")
				// down ports
				for _, port := range failPorts {
					fmt.Println("Downing " + port)
					intHandler.SetShutdown(port, true)
				}
				portsDown = true
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
