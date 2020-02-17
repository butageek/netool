/*
Copyright © 2020 Hendry Zhou <hendryzhou889@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/butageek/netool/scanner"
	"github.com/butageek/netool/validator"
	"github.com/spf13/cobra"
)

// netCmd represents the net command
var netCmd = &cobra.Command{
	Use:   "net [CIDR]",
	Short: "scan network for hosts that are alive",
	Long: `scan network for hosts that are alive
Arguments:
	CIDR - CIDR notation. eg. 192.168.1.1/24`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// validate argument against CIDR format
		v := validator.InitValidator()
		if !validator.IsValid(v.Regex["cidr"], args[0]) {
			fmt.Println("Invalid CIDR format")
			cmd.Help()
		}

		myScanner := &scanner.Scanner{}
		myScanner.ScanNet(args[0])
	},
}

func init() {
	rootCmd.AddCommand(netCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// netCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// netCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
