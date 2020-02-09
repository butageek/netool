package main

import (
	"log"
	"os"

	"github.com/butageek/netool/digger"
	"github.com/butageek/netool/scanner"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "network tool bundle",
		Usage: "query IP, NS, CNAME, MX records; scan network and open ports",
	}

	app.Commands = []*cli.Command{
		{
			Name:  "dig",
			Usage: "looks up the information for the domain",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "domain",
					Usage:    "domain name to lookup",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				myDigger := &digger.Digger{}
				myDigger.Domain = c.String("domain")
				myDigger.Dig()

				return nil
			},
		},
		{
			Name:  "port",
			Usage: "scan open ports for the host",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "host",
					Usage:    "hostname or IP address to scan",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "port",
					Usage:    "port number to scan, eg. 22,80,100-200",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				myScanner := &scanner.Scanner{}
				myScanner.ScanPort(c.String("host"), c.String("port"))

				return nil
			},
		},
		{
			Name:  "net",
			Usage: "scan network for hosts that are alive",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "cidr",
					Usage:    "network range (CIDR) to scan, eg. 192.168.1.0/24",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				myScanner := &scanner.Scanner{}
				myScanner.ScanNet(c.String("cidr"))

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
