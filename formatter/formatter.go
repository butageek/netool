package formatter

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/butageek/netool/reference"
	"github.com/google/gopacket/macs"
	"github.com/mostlygeek/arp"
	"github.com/olekukonko/tablewriter"
)

// Formatter struct of Formatter
type Formatter struct {
	Header    []string
	Data      [][]string
	Border    bool
	Separator string
}

// Print prints formatted table of data
func (o *Formatter) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(o.Header)
	table.SetBorder(o.Border)
	table.SetCenterSeparator(o.Separator)
	table.AppendBulk(o.Data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

// AssemblePortData assembles output data for port scanner
func (o *Formatter) AssemblePortData(ports []int, pra *reference.PortRefArray) {
	var data [][]string

	for _, port := range ports {
		portRef := pra.Find(port)
		row := []string{
			strconv.Itoa(port),
			"TCP",
			portRef.Name,
			portRef.Desc,
		}
		data = append(data, row)
	}

	o.Data = data
}

// AssembleNetData assembles output data for net scanner
func (o *Formatter) AssembleNetData(ips []net.IP) {
	var data [][]string
	var row []string

	for _, ip := range ips {
		macStr := arp.Search(ip.String())
		mac, err := net.ParseMAC(macStr)
		if err != nil {
			row = []string{
				ip.String(),
				"",
				"",
			}
			data = append(data, row)
			continue
		}
		prefix := [3]byte{
			mac[0],
			mac[1],
			mac[2],
		}
		manufacturer := macs.ValidMACPrefixMap[prefix]
		row = []string{
			ip.String(),
			mac.String(),
			manufacturer,
		}
		data = append(data, row)
	}

	o.Data = data
}
