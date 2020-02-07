package scanner

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/butageek/netool/formatter"
	"github.com/butageek/netool/reference"
)

// Scanner struct of Scanner
type Scanner struct{}

//ScanNet scans network for hosts that are alive
func (s *Scanner) ScanNet(cidr string) error {
	ips, _ := getIPs(cidr)
	jobChan := make(chan string, len(ips))
	resultChan := make(chan string, 10)

	wgs := sync.WaitGroup{}
	wgr := sync.WaitGroup{}

	fmt.Println()
	log.Printf("Scanning net %s\n", cidr)
	fmt.Println()

	numScanners := 100
	wgs.Add(numScanners)
	for i := 1; i <= numScanners; i++ {
		go netScanner(jobChan, resultChan, &wgs)
	}

	hostsAlive := []net.IP{}
	wgr.Add(1)
	go netReceiver(resultChan, &hostsAlive, &wgr)

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
			Header:    []string{"Host", "MAC Address", "Manufacturer"},
			Border:    false,
			Separator: "|",
		}
		formatter.AssembleNetData(hostsAlive)
		formatter.Print()
	} else {
		log.Println("No host alive found!")
	}

	return nil
}

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

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func netScanner(jobChan <-chan string, resultChan chan<- string, wgs *sync.WaitGroup) {
	defer wgs.Done()

	for ip := range jobChan {
		_, err := exec.Command("ping", "-c1", ip).Output()
		if err != nil {
			continue
		} else {
			fmt.Printf("Found host: %s\n", ip)
			resultChan <- ip
		}
	}
}

func netReceiver(resultChan <-chan string, hostsAlive *[]net.IP, wgr *sync.WaitGroup) {
	defer wgr.Done()

	for ip := range resultChan {
		*hostsAlive = append(*hostsAlive, net.ParseIP(ip))
	}
}

func sortIPs(ips *[]net.IP) {
	sort.Slice(*ips, func(i, j int) bool {
		return bytes.Compare((*ips)[i], (*ips)[j]) < 0
	})
}

// ScanPort scans open ports for the host
func (s *Scanner) ScanPort(host, port string) error {
	portRefArray := reference.PortRefArray{}
	portRefArray.Init()

	ports := parsePorts(port)
	numPorts := len(ports)
	jobChan := make(chan int, numPorts)
	resultChan := make(chan int, 10)

	wgs := sync.WaitGroup{}
	wgr := sync.WaitGroup{}

	fmt.Println()
	log.Printf("Scanning host %s\n", host)
	fmt.Println()

	numScanners := 100
	for i := 1; i <= numScanners; i++ {
		wgs.Add(1)
		go portScanner(host, jobChan, resultChan, &wgs)
	}

	openedPorts := []int{}
	wgr.Add(1)
	go portReceiver(resultChan, &openedPorts, &wgr)

	for _, port := range ports {
		jobChan <- port
	}
	close(jobChan)

	wgs.Wait()
	close(resultChan)
	wgr.Wait()

	if len(openedPorts) > 0 {
		formatter := &formatter.Formatter{
			Header:    []string{"Port", "Protocol", "Service Name", "Description"},
			Border:    false,
			Separator: "|",
		}
		formatter.AssemblePortData(openedPorts, &portRefArray)
		formatter.Print()
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

func portReceiver(resultChan <-chan int, openedPorts *[]int, wgr *sync.WaitGroup) {
	defer wgr.Done()

	for port := range resultChan {
		*openedPorts = append(*openedPorts, port)
	}
}
