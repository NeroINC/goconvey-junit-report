package main

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Result int

const (
	PASS Result = iota
	FAIL
	SKIP
)

type Report struct {
	Packages []Package
}

type Package struct {
	Name  string
	Time  int
	Tests []*Test
}

type Test struct {
	Name   string
	Result Result
	Output []string
}

type FailedTest struct {
	test       *Test
	totalFails int
}

var (
	xSuccess                                                                  = `✔`
	xFailure                                                                  = `✘`
	xError                                                                    = `🔥`
	xSkip                                                                     = `⚠`
	xSpecial                                                                  = `✔✘🔥⚠`
	dotSuccess                                                                = `\.`
	dotFailure                                                                = `x`
	dotError                                                                  = `E`
	dotSkip                                                                   = `S`
	dotSpecial                                                                = `.xES`
	regexTest, regexTestSuccess, regexTestSkip, regexPackage, regexAssertions *regexp.Regexp
	success, failure, error_, skip, special                                   string
)

func initialize(useDot bool) {
	if useDot {
		success, failure, error_, skip, special = dotSuccess, dotFailure, dotError, dotSkip, dotSpecial
	} else {
		success, failure, error_, skip, special = xSuccess, xFailure, xError, xSkip, xSpecial
	}

	regexTest = regexp.MustCompile(`(?i)^\s+(Given|When|Then|And)`)
	regexTestSuccess = regexp.MustCompile(`\s(` + success + `)*$`)
	regexTestSkip = regexp.MustCompile(`\s` + skip + `$`)
	regexPackage = regexp.MustCompile(`^--- (PASS|FAIL): (.+) \((\d+\.\d+) seconds\)$`)
	regexAssertions = regexp.MustCompile(`^\d+ assertion(s)? thus far$`)
}

func Parse(r io.Reader, useDot bool) (*Report, error) {
	initialize(useDot)
	reader := bufio.NewReader(r)

	report := &Report{make([]Package, 0)}

	// keep track of tests we find
	tests := make([]*Test, 0)
	failedTests := make([]FailedTest, 0)
	currFailedTestIndex := -1

	// current test
	var test *Test

	var readingFailures bool

	// parse lines
	for {
		l, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		line := string(l)

		if regexTest.FindString(line) != "" {
			// start of a new test
			testName := strings.TrimSpace(strings.TrimRight(line, special))

			testResult := FAIL

			if regexTestSuccess.FindString(line) != "" {
				testResult = PASS
			}

			if regexTestSkip.FindString(line) != "" {
				testResult = SKIP
			}

			test = &Test{
				Name:   testName,
				Result: testResult,
				Output: make([]string, 0),
			}

			if testResult == FAIL {
				result := line[strings.LastIndex(line, " "):]
				failedTests = append(failedTests, FailedTest{test, strings.Count(result, failure)})
			}

			tests = append(tests, test)
		} else if matches := regexPackage.FindStringSubmatch(line); len(matches) == 4 {
			report.Packages = append(report.Packages, Package{
				Name:  matches[2],
				Time:  parseTime(matches[3]),
				Tests: tests,
			})

			tests = make([]*Test, 0)
			failedTests = make([]FailedTest, 0)
			currFailedTestIndex = -1
			readingFailures = false
		} else if test != nil {
			if regexAssertions.FindString(line) != "" {
				//no more errors
				readingFailures = false
			} else if strings.HasPrefix(line, "Failures:") {
				readingFailures = true
				currFailedTestIndex = 0
			} else if readingFailures {
				if strings.HasPrefix(strings.TrimSpace(line), "*") {
					failedTests[currFailedTestIndex].totalFails--

					if failedTests[currFailedTestIndex].totalFails < 0 {
						currFailedTestIndex++
						failedTests[currFailedTestIndex].totalFails--
					}
				}

				if len(line) > 0 {
					failedTests[currFailedTestIndex].test.Output =
						append(failedTests[currFailedTestIndex].test.Output, strings.TrimSpace(line))
				}
			}
		}
	}

	return report, nil
}

func parseTime(time string) int {
	t, err := strconv.Atoi(strings.Replace(time, ".", "", -1))
	if err != nil {
		return 0
	}
	return t
}
