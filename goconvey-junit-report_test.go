package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"
)

type TestCase struct {
	name       string
	reportName string
	report     *Report
}

var testCases []TestCase = []TestCase{
	{
		name:       "01-pass.txt",
		reportName: "01-report.xml",
		report: &Report{
			Packages: []Package{
				{
					Name: "TestScenarioOne",
					Time: 10,
					Tests: []*Test{
						{
							Name:   "Given some pre-conditions",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "And other pre-conditions",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "When something happens",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "Then everything is fine",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "And we are cool bro",
							Result: PASS,
							Output: []string{},
						},
					},
				},
				{
					Name: "TestScenarioTwo",
					Time: 100,
					Tests: []*Test{
						{
							Name:   "Given some new pre-conditions",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "When something else happens",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "Then nothing broke here",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "And here was everything fine",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "And this one was skipped",
							Result: SKIP,
							Output: []string{},
						},
						{
							Name:   "And the last one OK too",
							Result: PASS,
							Output: []string{},
						},
					},
				},
			},
		},
	},
	{
		name:       "02-fail.txt",
		reportName: "02-report.xml",
		report: &Report{
			Packages: []Package{
				{
					Name: "TestScenarioOne",
					Time: 10,
					Tests: []*Test{
						{
							Name:   "Given some pre-conditions",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "And other pre-conditions",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "When something happens",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "Then everything is fine",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "And we are cool bro",
							Result: PASS,
							Output: []string{},
						},
					},
				},
				{
					Name: "TestScenarioTwo",
					Time: 100,
					Tests: []*Test{
						{
							Name:   "Given some new pre-conditions",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "When something else happens",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "Then something broke here",
							Result: FAIL,
							Output: []string{
								"* /home/my project/project_test.go",
								"Line 100:",
								"Expected: '-2'",
								"Actual:   '0'",
								"(Should be equal)",
								"* /home/my project/project_test.go",
								"Line 120:",
								"Expected: 'true'",
								"Actual:   'false'",
								"(Should be equal)",
							},
						},
						{
							Name:   "And here was everything fine",
							Result: PASS,
							Output: []string{},
						},
						{
							Name:   "And here not cool bro",
							Result: FAIL,
							Output: []string{
								"* /home/my project/project_test.go",
								"Line 300:",
								"Expected: 'false'",
								"Actual:   'true'",
								"(Should be equal)",
							},
						},
						{
							Name:   "And the last one OK",
							Result: PASS,
							Output: []string{},
						},
					},
				},
			},
		},
	},
}

func init() {
    //our test examples where created in windows
    useDot = true
}

func TestParser(t *testing.T) {
	for _, testCase := range testCases {
		file, err := os.Open("tests/" + testCase.name)
		if err != nil {
			t.Fatal(err)
		}

		report, err := Parse(file, true)
		if err != nil {
			t.Fatalf("error parsing: %s", err)
		}

		if report == nil {
			t.Fatalf("Report == nil")
		}

		expected := testCase.report
		if len(report.Packages) != len(expected.Packages) {
			t.Fatalf("Report packages == %d, want %d", len(report.Packages), len(expected.Packages))
		}

		for i, pkg := range report.Packages {
			expPkg := expected.Packages[i]

			if pkg.Name != expPkg.Name {
				t.Errorf("Package.Name == %s, want %s", pkg.Name, expPkg.Name)
			}

			if pkg.Time != expPkg.Time {
				t.Errorf("Package.Time == %d, want %d", pkg.Time, expPkg.Time)
			}

			if len(pkg.Tests) != len(expPkg.Tests) {
				t.Fatalf("Package Tests == %d, want %d", len(pkg.Tests), len(expPkg.Tests))
			}

			for j, test := range pkg.Tests {
				expTest := expPkg.Tests[j]

				if test.Name != expTest.Name {
					t.Errorf("Test.Name == %s, want %s", test.Name, expTest.Name)
				}

				if test.Result != expTest.Result {
					t.Errorf("Test.Result == %d, want %d on %s", test.Result, expTest.Result, test.Name)
				}

				testOutput := strings.Join(test.Output, "\n")
				expTestOutput := strings.Join(expTest.Output, "\n")
				if testOutput != expTestOutput {
					t.Errorf("Test.Output ==\n%s\n, want\n%s", testOutput, expTestOutput)
				}
			}
		}
	}
}

func TestJUnitFormatter(t *testing.T) {
	for _, testCase := range testCases {
		report, err := loadTestReport(testCase.reportName)
		if err != nil {
			t.Fatal(err)
		}

		var junitReport bytes.Buffer

		if err = JUnitReportXML(testCase.report, &junitReport); err != nil {
			t.Fatal(err)
		}

		if string(junitReport.Bytes()) != report {
			t.Fatalf("Report xml ==\n%s, want\n%s", string(junitReport.Bytes()), report)
		}
	}
}

func loadTestReport(name string) (string, error) {
	contents, err := ioutil.ReadFile("tests/" + name)
	if err != nil {
		return "", err
	}

	// replace value="1.0" With actual version
	report := strings.Replace(string(contents), `value="1.0"`, fmt.Sprintf(`value="%s"`, runtime.Version()), 1)

	return report, nil
}
