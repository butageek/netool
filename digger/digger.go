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

// Dig looks up information for given domain
// Includes records of IP, NS, CNAME, MX
func (d *Digger) Dig() error {
	// get A records
	addrs, err := digHost(d.Domain)
	if err != nil {
		return err
	}
	// get NS records
	nss, err := digNS(d.Domain)
	if err != nil {
		return err
	}
	// get CNAME records
	cname, err := digCNAME(d.Domain)
	if err != nil {
		return err
	}
	// get MX records
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

// assembleDigData aseembles data for Formatter struct
func assembleDigData(d *Digger, addrs, nss, mxs []string, cname string) [][]string {
	var data [][]string

	// append A records as rows to data
	for _, addr := range addrs {
		row := []string{
			d.Domain,
			"A",
			addr,
		}
		data = append(data, row)
	}
	data = append(data, []string{"", "", ""})

	// append CNAME records as rows to data
	row := []string{
		d.Domain,
		"CNAME",
		cname,
	}
	data = append(data, row)
	data = append(data, []string{"", "", ""})

	// append NS records as rows to data
	for _, ns := range nss {
		row := []string{
			d.Domain,
			"NS",
			ns,
		}
		data = append(data, row)
	}
	data = append(data, []string{"", "", ""})

	// append MX records as rows to data
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

// digHost gets A records
func digHost(host string) ([]string, error) {
	var addrs []string

	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

// digNS gets NS records
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

// digCNAME gets CNAME records
func digCNAME(host string) (string, error) {
	cname, err := net.LookupCNAME(host)
	if err != nil {
		return "", err
	}

	return cname, nil
}

// digMX gets MX records
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
