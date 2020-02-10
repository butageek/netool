package main

import (
	"errors"
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
			Name: "dig",
			Usage: `looks up the information for the domain
			@arguments:
				Domain - domain to lookup, eg. example.com`,
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					myDigger := &digger.Digger{}
					myDigger.Domain = c.Args().Get(0)
					myDigger.Dig()

					return nil
				}

				return errors.New("Missing argument: Domain => eg: example.com")
			},
		},
		{
			Name: "port",
			Usage: `scan open ports for the host
			@arguments:
				Host - host to scan, eg. www.example.com or 10.10.10.10`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "port",
					Aliases: []string{"p"},
					Value:   "1-1023,3389",
					Usage:   "port number to scan, eg. 80,100-200",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					myScanner := &scanner.Scanner{}
					myScanner.ScanPort(c.Args().Get(0), c.String("port"))

					return nil
				}

				return errors.New("Missing argument: Host => eg: www.example.com or 10.10.10.10")
			},
		},
		{
			Name: "net",
			Usage: `scan network for hosts that are alive
			@arguments:
				CIDR - cidr to scan, eg. 192.168.1.1/24`,
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					myScanner := &scanner.Scanner{}
					myScanner.ScanNet(c.Args().Get(0))

					return nil
				}

				return errors.New("Missing argument: CIDR => eg: 192.168.1.1/24")
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
