package main

import (
	"flag"
	"fmt"
	"os"
)

var (
    useDot bool
)

func init() {
	flag.BoolVar(&useDot, "useDot", false, "Set to true if your goconvey output was generated in Windows")
	flag.Parse()
}

func main() {
	// Read input
	report, err := Parse(os.Stdin, useDot)
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		os.Exit(1)
	}

	// Write xml
	err = JUnitReportXML(report, os.Stdout)
	if err != nil {
		fmt.Printf("Error writing XML: %s\n", err)
		os.Exit(1)
	}
}
