package scanner

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Scanner struct of Scanner
type Scanner struct{}

// Scan scans open ports for the host
func (s *Scanner) Scan(host, port string) error {
	portRefArray := PortRefArray{}
	portRefArray.Init()

	ports := parsePorts(port)
	numPorts := len(ports)
	jobChan := make(chan int, numPorts)
	resultChan := make(chan int, numPorts)

	wg := sync.WaitGroup{}

	fmt.Println()
	log.Printf("Scanning host %s\n", host)
	fmt.Println()
	numScanners := 100
	for i := 1; i <= numScanners; i++ {
		wg.Add(1)
		go scanner(host, jobChan, resultChan, &wg)
	}

	for port := range ports {
		jobChan <- port
	}
	close(jobChan)

	wg.Wait()
	close(resultChan)

	openedPorts := receive(resultChan)
	if len(*openedPorts) > 0 {
		outputResults(openedPorts, &portRefArray)
	} else {
		log.Println("No open ports found!")
	}

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

func scanner(host string, jobChan <-chan int, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range jobChan {
		hostIP := fmt.Sprintf("%s:%d", host, port)

		conn, err := net.DialTimeout("tcp", hostIP, 200*time.Millisecond)
		if err != nil {
			continue
		}
		defer conn.Close()

		fmt.Printf("Found open port: %d\n", port)

		resultChan <- port
	}
}

func receive(resultChan <-chan int) *[]int {
	var openedPorts []int

	for port := range resultChan {
		openedPorts = append(openedPorts, port)
	}

	return &openedPorts
}

func outputResults(ports *[]int, pra *PortRefArray) {
	var data [][]string

	for _, port := range *ports {
		portRef := pra.Find(port)
		row := []string{
			strconv.Itoa(port),
			portRef.Name,
			portRef.Desc,
		}
		data = append(data, row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Port", "Service Name", "Description"})
	table.SetBorder(false)
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	fmt.Println()
	log.Println("Printing open ports")
	fmt.Println()
	table.Render()
	fmt.Println()
}
