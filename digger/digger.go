package digger

import (
	"fmt"
	"net"
)

// Digger struct of Digger
type Digger struct {
	Host string
}

// DigIP looks up the IP address for the host
func (d *Digger) DigIP() error {
	ip, err := net.LookupIP(d.Host)
	if err != nil {
		return err
	}
	for i := 0; i < len(ip); i++ {
		fmt.Println(ip[i])
	}

	return nil
}

// DigNS looks up the name servers for the host
func (d *Digger) DigNS() error {
	ns, err := net.LookupNS(d.Host)
	if err != nil {
		return err
	}
	for i := 0; i < len(ns); i++ {
		fmt.Println(ns[i].Host)
	}
	return nil
}

// DigCNAME looks up the CNAME for the host
func (d *Digger) DigCNAME() error {
	cname, err := net.LookupCNAME(d.Host)
	if err != nil {
		return err
	}
	fmt.Println(cname)
	return nil
}

// DigMX looks up the MX for the host
func (d *Digger) DigMX() error {
	mx, err := net.LookupMX(d.Host)
	if err != nil {
		return err
	}
	for i := 0; i < len(mx); i++ {
		fmt.Println(mx[i].Host, mx[i].Pref)
	}
	return nil
}
