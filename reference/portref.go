package reference

import (
	"os"
	"path/filepath"
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
	// switch working directory to the same as executable
	filePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	os.Chdir(filePath)
	f, err := os.OpenFile("service-names-port-numbers.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := gocsv.UnmarshalFile(f, p); err != nil {
		panic(err)
	}
}

// Find finds name of port and returns portRef object
func (p *PortRefArray) Find(port int) *PortRef {
	for _, portRef := range *p {
		if portRef.Protocol == "tcp" && portRef.PortNum == strconv.Itoa(port) {
			return &portRef
		}
	}

	return &PortRef{}
}
