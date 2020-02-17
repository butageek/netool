package validator

import (
	"regexp"
)

// Validator struct of Validator
type Validator struct {
	Regex map[string]string
}

// InitValidator initialize Validator struct
func InitValidator() *Validator {
	v := Validator{
		Regex: map[string]string{},
	}
	v.Regex["domain"] = `^\w+\.\w{2,4}$`
	v.Regex["port"] = `^\d+([,-]\d+)*$`
	v.Regex["cidr"] = `^\d{1,3}\.\d{1,3}\.\d{1,3}.\d{1,3}\/\d{1,2}$`

	return &v
}

// IsValid validates argument and returns true if valid
func IsValid(pattern, arg string) bool {
	match, _ := regexp.MatchString(pattern, arg)

	return match
}
