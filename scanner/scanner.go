package scanner

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Scanner struct of Scanner
type Scanner struct{}

// Scan scans open ports for the host
func (s *Scanner) Scan(host, port string) error {
	ports := parsePorts(port)
	numPorts := len(ports)
	portChan := make(chan int, numPorts)
	wg := sync.WaitGroup{}

	log.Printf("starting host scan %s", host)
	numScanners := 100
	for i := 1; i <= numScanners; i++ {
		wg.Add(1)
		go scanner(host, portChan, &wg)
	}

	for port := range ports {
		portChan <- port
	}
	close(portChan)

	wg.Wait()

	return nil
}

func parsePorts(portString string) []int {
	var ports []int

	portsSplit := strings.Split(portString, ",")

	for _, port := range portsSplit {
		if strings.Contains(port, "-") {
			portBounds := strings.Split(port, "-")
			portStart, err := strconv.Atoi(portBounds[0])
			if err != nil {
				log.Fatal(err)
			}
			portEnd, err := strconv.Atoi(portBounds[1])
			if err != nil {
				log.Fatal(err)
			}
			for i := portStart; i <= portEnd; i++ {
				ports = append(ports, i)
			}
		} else {
			portNum, err := strconv.Atoi(port)
			if err != nil {
				log.Fatal(err)
			}
			ports = append(ports, portNum)
		}
	}

	return ports
}

func scanner(host string, portChan <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range portChan {
		hostIP := fmt.Sprintf("%s:%d", host, port)

		conn, err := net.DialTimeout("tcp", hostIP, 100*time.Millisecond)
		if err != nil {
			continue
		}
		defer conn.Close()

		fmt.Printf("Open port found: %d\n", port)
	}
}
