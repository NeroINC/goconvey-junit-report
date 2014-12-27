goconvey-junit-report
=====================

Converts a convey `go test` output to an xml report, suitable for applications that
expect junit xml reports (e.g. [Jenkins](http://jenkins-ci.org)).

Tests must be written in classic BDD format

	Given ...
	{And ...}*
	When ...
	Then ...
	{And ...}*

Also make sure you do not have any debug/error output on the stdout or it will affect the report generation.

Installation
------------

	go get github.com/NeroINC/goconvey-junit-report

Usage
-----

In Windows, you need to add the flag **-useDot=true** because goconvey uses different output on Windows OS:

	go test -v | goconvey-junit-report -useDot=true > report.xml

If you are inside a package that has multiple sub-packages, you can test all of them in a single run with

	go test -v ./... | goconvey-junit-report -useDot=true > report.xml


In linux, you need to set the environment TERM needs to be set to **sh** because goconvey adds colors to the output and that affects the report

    TERM=sh go test -v | goconvey-junit-report > report.xml

Or

    TERM=sh go test -v ./... | goconvey-junit-report > report.xml

//TODOs
------

1. When multiple packages are tested at the same time using `go test -v ./...` the JUnit Report is not informing to which subpackage each test belongs.
2. Add a test for linux reports.
3. Add examples of _test.go constructed in BDD format.

Credits
-------

This work is inspired on the go-junit-report tool created by jstemmer:
https://github.com/jstemmer/go-junit-report

What we added was more output information to the final report about each step of the BDD test and the assertions.
