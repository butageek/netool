package digger

import (
	"fmt"
	"net"
	"strconv"

	"github.com/butageek/netool/formatter"
)

// Digger struct of Digger
type Digger struct {
	Domain string
}

// Dig lookup information of host
// Includes records of IP, NS, CNAME, MX
func (d *Digger) Dig() error {
	addrs, err := digHost(d.Domain)
	if err != nil {
		return err
	}
	nss, err := digNS(d.Domain)
	if err != nil {
		return err
	}
	cname, err := digCNAME(d.Domain)
	if err != nil {
		return err
	}
	mxs, err := digMX(d.Domain)
	if err != nil {
		return err
	}

	formatter := &formatter.Formatter{
		Header:          []string{"Domain", "Type", "Value"},
		Data:            assembleDigData(d, addrs, nss, mxs, cname),
		Border:          false,
		Separator:       " ",
		ColumnSeparator: " ",
	}
	formatter.Print()

	return nil
}

func assembleDigData(d *Digger, addrs, nss, mxs []string, cname string) [][]string {
	var data [][]string

	for _, addr := range addrs {
		row := []string{
			d.Domain,
			"A",
			addr,
		}
		data = append(data, row)
	}
	data = append(data, []string{"", "", ""})

	row := []string{
		d.Domain,
		"CNAME",
		cname,
	}
	data = append(data, row)
	data = append(data, []string{"", "", ""})

	for _, ns := range nss {
		row := []string{
			d.Domain,
			"NS",
			ns,
		}
		data = append(data, row)
	}
	data = append(data, []string{"", "", ""})

	for _, mx := range mxs {
		row := []string{
			d.Domain,
			"MX",
			mx,
		}
		data = append(data, row)
	}

	return data
}

func digHost(host string) ([]string, error) {
	var addrs []string

	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

func digNS(domain string) ([]string, error) {
	var nss []string

	ns, err := net.LookupNS(domain)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(ns); i++ {
		nss = append(nss, ns[i].Host)
	}

	return nss, nil
}

func digCNAME(host string) (string, error) {
	cname, err := net.LookupCNAME(host)
	if err != nil {
		return "", err
	}

	return cname, nil
}

func digMX(domain string) ([]string, error) {
	var mxs []string

	mx, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(mx); i++ {
		mxStr := fmt.Sprintf("%s %s", mx[i].Host, strconv.Itoa(int(mx[i].Pref)))
		mxs = append(mxs, mxStr)
	}

	return mxs, nil
}
