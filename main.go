package main

import (
	"errors"
	"log"
	"os"

	"github.com/butageek/netool/digger"
	"github.com/butageek/netool/scanner"
	"github.com/butageek/netool/validator"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "network tool bundle",
		Usage: "query IP, NS, CNAME, MX records; scan network and open ports",
	}

	app.Commands = []*cli.Command{
		{
			// dig looks up for domain informations
			// includes IP, NS, CNAME, MX records
			Name:  "dig",
			Usage: "looks up the information for the domain",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					// validate argument against Domain format
					v := validator.InitValidator()
					if !validator.IsValid(v.Regex["domain"], c.Args().Get(0)) {
						return errors.New("Wrong argument format: Domain. Example: example.com")
					}

					myDigger := &digger.Digger{}
					myDigger.Domain = c.Args().Get(0)
					myDigger.Dig()

					return nil
				}

				return errors.New("Missing argument: Domain. Example: example.com")
			},
		},
		{
			// port command scans ports for given host
			Name:  "port",
			Usage: "scan open ports for the host",
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

				return errors.New("Missing argument: Host. Example: www.example.com or 10.10.10.10")
			},
		},
		{
			// net scans network for given CIDR in format: xxx.xxx.xxx.xx/xx
			Name:  "net",
			Usage: "scan network for hosts that are alive",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					// validate argument against CIDR format
					v := validator.InitValidator()
					if !validator.IsValid(v.Regex["cidr"], c.Args().Get(0)) {
						return errors.New("Wrong argument format: CIDR. Example: 192.168.1.1/24")
					}

					myScanner := &scanner.Scanner{}
					myScanner.ScanNet(c.Args().Get(0))

					return nil
				}

				return errors.New("Missing argument: CIDR. Example: 192.168.1.1/24")
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
