package scanner

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/butageek/netool/formatter"
	"github.com/butageek/netool/reference"
	"github.com/butageek/netool/validator"
)

// Scanner struct of Scanner
type Scanner struct{}

//ScanNet scans network for hosts that are alive
func (s *Scanner) ScanNet(cidr string) error {
	// parse IP addresses for given cidr
	ips, _ := getIPs(cidr)
	// init channels
	jobChan := make(chan string, len(ips))
	resultChan := make(chan string, 10)

	// init WaitGroups
	// wgs for Scanner, wgr for Receiver
	wgs := sync.WaitGroup{}
	wgr := sync.WaitGroup{}

	fmt.Println()
	log.Printf("Scanning net %s\n", cidr)
	fmt.Println()

	// set concurrency limit for Scanner
	numScanners := 100
	wgs.Add(numScanners)
	for i := 1; i <= numScanners; i++ {
		go netScanner(jobChan, resultChan, &wgs)
	}

	hostsAlive := []net.IP{}
	// set one Receiver
	wgr.Add(1)
	go netReceiver(resultChan, &hostsAlive, &wgr)

	// init jobChan using parsed IPs
	for _, ip := range ips {
		jobChan <- ip
	}
	close(jobChan)

	wgs.Wait()
	close(resultChan)
	wgr.Wait()

	sortIPs(&hostsAlive)

	if len(hostsAlive) > 0 {
		formatter := &formatter.Formatter{
			Header:          []string{"Host", "Status", "MAC Address", "Manufacturer"},
			Border:          false,
			Separator:       " ",
			ColumnSeparator: " ",
		}
		formatter.AssembleNetData(hostsAlive)
		formatter.Print()
	} else {
		log.Println("No host alive found!")
	}

	return nil
}

// getIPs parses given CIDR and return IPs in that range
func getIPs(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

// inc increases IP address by 1
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// netScanner pings a host and appends it to resultChan if it's alive
func netScanner(jobChan <-chan string, resultChan chan<- string, wgs *sync.WaitGroup) {
	defer wgs.Done()

	switch runtime.GOOS {
	case "windows":
		for ip := range jobChan {
			out, _ := exec.Command("ping", "-n", "1", ip).Output()
			if strings.Contains(string(out), "Destination host unreachable") {
				continue
			} else {
				fmt.Printf("Found host: %s\n", ip)
				resultChan <- ip
			}
		}
	case "linux":
		for ip := range jobChan {
			_, err := exec.Command("ping", "-c", "1", ip).Output()
			if err != nil {
				continue
			} else {
				fmt.Printf("Found host: %s\n", ip)
				resultChan <- ip
			}
		}
	}
}

// netReceiver get IP from resultChan and appends to host IPs that are alive
func netReceiver(resultChan <-chan string, hostsAlive *[]net.IP, wgr *sync.WaitGroup) {
	defer wgr.Done()

	for ip := range resultChan {
		*hostsAlive = append(*hostsAlive, net.ParseIP(ip))
	}
}

// sortIPs sorts IPs
func sortIPs(ips *[]net.IP) {
	sort.Slice(*ips, func(i, j int) bool {
		return bytes.Compare((*ips)[i], (*ips)[j]) < 0
	})
}

// ScanPort scans open ports for the host
func (s *Scanner) ScanPort(host, port string) error {
	// init port reference object
	portRefArray := reference.PortRefArray{}
	portRefArray.Init()

	// parse ports on the given port argument
	ports := parsePorts(port)
	// init channels
	numPorts := len(ports)
	jobChan := make(chan int, numPorts)
	resultChan := make(chan int, 10)

	// init WaitGroups
	// wgs for Scanner, wgr for Receiver
	wgs := sync.WaitGroup{}
	wgr := sync.WaitGroup{}

	fmt.Println()
	log.Printf("Scanning host %s\n", host)
	fmt.Println()

	// set Scanner concurrency limit
	numScanners := 100
	for i := 1; i <= numScanners; i++ {
		wgs.Add(1)
		go portScanner(host, jobChan, resultChan, &wgs)
	}

	openedPorts := []int{}
	// set Receiver concurrency limit to 1
	wgr.Add(1)
	go portReceiver(resultChan, &openedPorts, &wgr)

	// init jobChan using parsed ports
	for _, port := range ports {
		jobChan <- port
	}
	close(jobChan)

	wgs.Wait()
	close(resultChan)
	wgr.Wait()

	if len(openedPorts) > 0 {
		formatter := &formatter.Formatter{
			Header:          []string{"Port", "Protocol", "Service Name", "Description"},
			Border:          false,
			Separator:       " ",
			ColumnSeparator: " ",
		}
		formatter.AssemblePortData(openedPorts, &portRefArray)
		formatter.Print()
	} else {
		log.Println("No open ports found!")
	}

	return nil
}

// parsePorts parses ports on port argument
// supports comma and dash separated ports. eg. 80,100-200
func parsePorts(portString string) []int {
	var ports []int
	v := validator.InitValidator()

	portsSplit := strings.Split(portString, ",")

	for _, port := range portsSplit {
		if strings.Contains(port, "-") {
			portBounds := strings.Split(port, "-")
			portStart, err := strconv.Atoi(portBounds[0])
			if err != nil {
				log.Fatal(err)
			}
			if !validator.IsValid(v.Regex["port"], strconv.Itoa(portStart)) {
				log.Fatal(errors.New("Wrong argument format: Port. Example: 80,100-200"))
			}
			portEnd, err := strconv.Atoi(portBounds[1])
			if err != nil {
				log.Fatal(err)
			}
			if !validator.IsValid(v.Regex["port"], strconv.Itoa(portEnd)) {
				log.Fatal(errors.New("Wrong argument format: Port. Example: 880,100-2000"))
			}
			for i := portStart; i <= portEnd; i++ {
				ports = append(ports, i)
			}
		} else {
			if !validator.IsValid(v.Regex["port"], port) {
				log.Fatal(errors.New("Wrong argument format: Port. Example: 80,100-200"))
			}
			portNum, err := strconv.Atoi(port)
			if err != nil {
				log.Fatal(err)
			}
			ports = append(ports, portNum)
		}
	}

	return ports
}

// portScanner scans a port and push to resultChan if it's open
func portScanner(host string, jobChan <-chan int, resultChan chan<- int, wgs *sync.WaitGroup) {
	defer wgs.Done()

	for port := range jobChan {
		hostIP := fmt.Sprintf("%s:%d", host, port)

		conn, err := net.DialTimeout("tcp", hostIP, 500*time.Millisecond)
		if err != nil {
			continue
		}
		defer conn.Close()

		fmt.Printf("Found open port: %d\n", port)
		resultChan <- port
	}
}

// portReceiver receives ports from resultChan and appends to openedPorts array
func portReceiver(resultChan <-chan int, openedPorts *[]int, wgr *sync.WaitGroup) {
	defer wgr.Done()

	for port := range resultChan {
		*openedPorts = append(*openedPorts, port)
	}
}
