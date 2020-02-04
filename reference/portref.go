package reference

import (
	"os"
	"strconv"

	"github.com/gocarina/gocsv"
)

// PortRef struct of PortRef
type PortRef struct {
	Name     string `csv:"Service Name"`
	PortNum  string `csv:"Port Number"`
	Protocol string `csv:"Transport Protocol"`
	Desc     string `csv:"Description"`
}

// PortRefArray array of PortRef
type PortRefArray []PortRef

// Init initialize PortRefArray
func (p *PortRefArray) Init() {
	os.Chdir("reference")
	f, err := os.OpenFile("service-names-port-numbers.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := gocsv.UnmarshalFile(f, p); err != nil {
		panic(err)
	}
}

// Find finds name of port
func (p *PortRefArray) Find(port int) *PortRef {
	for _, portRef := range *p {
		if portRef.Protocol == "tcp" && portRef.PortNum == strconv.Itoa(port) {
			return &portRef
		}
	}

	return &PortRef{}
}
