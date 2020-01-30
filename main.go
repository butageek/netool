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
				myDigger := &digger.Digger{}
				myDigger.Host = c.String("host")
				myDigger.DigIP()

				return nil
			},
		},
		{
			Name:  "ns",
			Usage: "looks up the name servers for the host",
			Flags: hostFlag,
			Action: func(c *cli.Context) error {
				myDigger := &digger.Digger{}
				myDigger.Host = c.String("host")
				myDigger.DigNS()

				return nil
			},
		},
		{
			Name:  "cname",
			Usage: "looks up the CNAME for the host",
			Flags: hostFlag,
			Action: func(c *cli.Context) error {
				myDigger := &digger.Digger{}
				myDigger.Host = c.String("host")
				myDigger.DigCNAME()

				return nil
			},
		},
		{
			Name:  "mx",
			Usage: "looks up the MX for the host",
			Flags: hostFlag,
			Action: func(c *cli.Context) error {
				myDigger := &digger.Digger{}
				myDigger.Host = c.String("host")
				myDigger.DigMX()

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
				myScanner := &scanner.Scanner{}
				myScanner.Scan(c.String("host"), c.String("port"))

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
