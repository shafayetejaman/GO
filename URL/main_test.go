package main

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	type testCase struct {
		email    string
		expected string
	}

	runCases := []testCase{
		{"wayne.lagner@dev.boot", "mailto:wayne.lagner@dev.boot"},
		{"heckmann@what.de", "mailto:heckmann@what.de"},
		{"a.liar@pants.fire", "mailto:a.liar@pants.fire"},
	}

	submitCases := append(runCases, []testCase{
		{"", "mailto:"},
	}...)

	testCases := runCases
	if withSubmit {
		testCases = submitCases
	}

	skipped := len(submitCases) - len(testCases)

	passCount := 0
	failCount := 0

	for _, test := range testCases {
		output := getMailtoLinkForEmail(test.email)
		if output != test.expected {
			failCount++
			t.Errorf(`---------------------------------
Email:		%v
Expecting:  %v
Actual:     %v
Fail
`, test.email, test.expected, output)
		} else {
			passCount++
			fmt.Printf(`---------------------------------
Email:		%v
Expecting:  %v
Actual:     %v
Pass
`, test.email, test.expected, output)
		}
	}

	fmt.Println("---------------------------------")
	if skipped > 0 {
		fmt.Printf("%d passed, %d failed, %d skipped\n", passCount, failCount, skipped)
	} else {
		fmt.Printf("%d passed, %d failed\n", passCount, failCount)
	}
}

// withSubmit is set at compile time depending
// on which button is used to run the tests
var withSubmit = true
