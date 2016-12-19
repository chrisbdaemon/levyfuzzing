package evaluate

import (
	"log"

	"github.com/chrisbdaemon/levyfuzzing/coverage"
)

// TestCase implements the pieces of testcase.TestCase that
// we really care about
type TestCase interface {
	Coverage() *coverage.Coverage
}

// Score expects two []TestCase parameters, it compares the code coverage of each
// set of test cases, and returns the average number of coverage points that each element of testCases
// hit but were missed by oldTestCases
func Score(testCases, oldtestCases interface{}) (score int, err error) {
	sumDifferences := int64(0)
	oldCoverage := compileCoverage(oldtestCases.([]TestCase))

	for _, testCase := range testCases.([]TestCase) {
		sumDifferences += coverage.DifferenceCount(testCase.Coverage(), oldCoverage)
	}

	return
}

func compileCoverage(testCases []TestCase) (fullCoverage *coverage.Coverage) {
	if len(testCases) == 0 {
		log.Fatalln("no test cases to evaluate")
	}

	fullCoverage = testCases[0].Coverage()

	for _, testCase := range testCases {
		if testCase.Coverage() == nil {
			log.Fatalln("missing coverage data for:", testCase)
		}

		fullCoverage = fullCoverage.Union(testCase.Coverage())
	}

	return
}
