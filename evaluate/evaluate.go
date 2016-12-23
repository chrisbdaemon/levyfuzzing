package evaluate

import (
	"fmt"

	"github.com/chrisbdaemon/levyfuzzing/coverage"
)

// TestCase implements the pieces of testcase.TestCase that
// we really care about
type TestCase interface {
	Coverage() *coverage.Coverage
}

// Score uses the new set of coverage data and compares each entry with
// the coverage previously obtained and returns the average number of new
// instrumentation points the new coverage shows.
func Score(newCoverages []*coverage.Coverage, oldCoverage *coverage.Coverage) (score int, err error) {

	if len(newCoverages) == 0 {
		err = fmt.Errorf("unable to determine score, no coverage data provided")
		return
	}

	sumDifferences := int64(0)

	for _, newCoverage := range newCoverages {
		sumDifferences += coverage.DifferenceCount(newCoverage, oldCoverage)
	}

	score = int(sumDifferences / int64(len(newCoverages)))

	return
}
