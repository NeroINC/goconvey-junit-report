package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"runtime"
	"strings"
)

type JUnitTestSuites struct {
	XMLName   xml.Name `xml:"testsuites"`
	Tests     int      `xml:"tests,attr"`
	Failures  int      `xml:"failures,attr,omitempty"`
	Errors    int      `xml:"errors,attr,omitempty"`
	Disabled  int      `xml:"disabled,attr,omitempty"`
	Time      string   `xml:"time,attr"`
	totaltime int
	Name      string `xml:"name,attr,omitempty"`
	TestSuits []JUnitTestSuite
}

type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr,omitempty"`
	Errors     int             `xml:"errors,attr,omitempty"`
	Disabled   int             `xml:"disabled,attr,omitempty"`
	Skipped    int             `xml:"skipped,attr,omitempty"`
	Time       string          `xml:"time,attr,omitempty"`
	Name       string          `xml:"name,attr,omitempty"`
	Package    string          `xml:"package,attr,omitempty"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase
}

type JUnitTestCase struct {
	XMLName    xml.Name      `xml:"testcase"`
	Classname  string        `xml:"classname,attr"`
	Name       string        `xml:"name,attr"`
	Assertions int           `xml:"assertions,attr,omitempty"`
	Skipped    *JUnitSkipped `xml:"skipped,omitempty"`
	Failure    *JUnitFailure `xml:"failure,omitempty"`
	Error      *JUnitError   `xml:"error,omitempty"`
}

type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type JUnitSkipped struct{}

type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

type JUnitError struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

func NewJUnitProperty(name, value string) JUnitProperty {
	return JUnitProperty{
		Name:  name,
		Value: value,
	}
}

// JUnitReportXML writes a junit xml representation of the given report to w
// in the format described at http://windyroad.org/dl/Open%20Source/JUnit.xsd
func JUnitReportXML(report *Report, w io.Writer) error {
	suites := &JUnitTestSuites{}
	suites.TestSuits = []JUnitTestSuite{}

	// convert Report to JUnit test suites
	for _, pkg := range report.Packages {
		ts := JUnitTestSuite{
			Tests:      len(pkg.Tests),
			Failures:   0,
			Time:       formatTime(pkg.Time),
			Name:       pkg.Name,
			Properties: []JUnitProperty{},
			TestCases:  []JUnitTestCase{},
		}

		suites.Tests += ts.Tests
		suites.totaltime += pkg.Time

		classname := pkg.Name
		if idx := strings.LastIndex(classname, "/"); idx > -1 && idx < len(pkg.Name) {
			classname = pkg.Name[idx+1:]
		}

		// properties
		ts.Properties = append(ts.Properties, NewJUnitProperty("go.version", runtime.Version()))

		// individual test cases
		for _, test := range pkg.Tests {
			testCase := JUnitTestCase{
				Classname: classname,
				Name:      test.Name,
				Failure:   nil,
			}

			switch test.Result {
			case FAIL:
				ts.Failures += 1
				suites.Failures += 1

				testCase.Failure = &JUnitFailure{
					Message:  "Failed",
					Type:     "",
					Contents: strings.Join(test.Output, "\n"),
				}
			case SKIP:
				ts.Skipped += 1
				testCase.Skipped = &JUnitSkipped{}
			}

			ts.TestCases = append(ts.TestCases, testCase)
		}

		suites.TestSuits = append(suites.TestSuits, ts)
	}

	suites.Time = formatTime(suites.totaltime)

	// to xml
	bytes, err := xml.MarshalIndent(suites, "", "\t")
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(w)

	// remove newline from xml.Header, because xml.MarshalIndent starts with a newline
	writer.WriteString(xml.Header[:len(xml.Header)-1])
	writer.Write(bytes)
	writer.WriteByte('\n')
	writer.Flush()

	return nil
}

func countFailures(tests []Test) (result int) {
	for _, test := range tests {
		if test.Result == FAIL {
			result += 1
		}
	}
	return
}

func formatTime(time int) string {
	return fmt.Sprintf("%.3f", float64(time)/1000.0)
}
