package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

func scanner(host string, portChan <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range portChan {
		hostIP := fmt.Sprintf("%s:%d", host, port)

		conn, err := net.DialTimeout("tcp", hostIP, 100*time.Millisecond)
		if err != nil {
			continue
		}
		defer conn.Close()

		fmt.Printf("Port %d is open\n", port)
	}
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

func main() {
	app := &cli.App{
		Name:  "network tool bundle",
		Usage: "query IP, NS, CNAME, MX records; scan network and open ports",
	}

	hostFlag := []cli.Flag{
		&cli.StringFlag{
			Name:     "host",
			Required: true,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "ip",
			Usage: "looks up the IP address for the host",
			Flags: hostFlag,
			Action: func(c *cli.Context) error {
				ip, err := net.LookupIP(c.String("host"))
				if err != nil {
					return err
				}
				for i := 0; i < len(ip); i++ {
					fmt.Println(ip[i])
				}
				return nil
			},
		},
		{
			Name:  "ns",
			Usage: "looks up the name servers for the host",
			Flags: hostFlag,
			Action: func(c *cli.Context) error {
				ns, err := net.LookupNS(c.String("host"))
				if err != nil {
					return err
				}
				for i := 0; i < len(ns); i++ {
					fmt.Println(ns[i].Host)
				}
				return nil
			},
		},
		{
			Name:  "cname",
			Usage: "looks up the CNAME for the host",
			Flags: hostFlag,
			Action: func(c *cli.Context) error {
				cname, err := net.LookupCNAME(c.String("host"))
				if err != nil {
					return err
				}
				fmt.Println(cname)
				return nil
			},
		},
		{
			Name:  "mx",
			Usage: "looks up the MX for the host",
			Flags: hostFlag,
			Action: func(c *cli.Context) error {
				mx, err := net.LookupMX(c.String("host"))
				if err != nil {
					return err
				}
				for i := 0; i < len(mx); i++ {
					fmt.Println(mx[i].Host, mx[i].Pref)
				}
				return nil
			},
		},
		{
			Name:  "scan",
			Usage: "scan open ports for the host",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "host",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "port",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				host := c.String("host")
				ports := parsePorts(c.String("port"))
				numPorts := len(ports)
				portChan := make(chan int, numPorts)
				wg := sync.WaitGroup{}

				log.Printf("starting scan host %s", host)
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
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
