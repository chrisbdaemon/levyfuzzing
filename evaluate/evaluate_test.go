package evaluate

import (
	"testing"

	"github.com/chrisbdaemon/levyfuzzing/coverage"
	"github.com/stretchr/testify/assert"
)

type fakeTestCase struct {
	coverageData *coverage.Coverage
}

func (t *fakeTestCase) Coverage() *coverage.Coverage { return t.coverageData }

func TestScore(t *testing.T) {
	/*
		var testCases, oldTestCases []TestCase

		score, err := Score(testCases, oldTestCases, "id")
		assert.Nil(t, err)
		assert.True(t, score >= 0, "score should be positive")
	*/
}

func TestCompileCoverage(t *testing.T) {
	testCases := []TestCase{
		&fakeTestCase{
			coverage.New(1, 2, 3, 4),
		},
		&fakeTestCase{
			coverage.New(3, 4, 5, 6),
		},
	}

	c := compileCoverage(testCases)
	assert.EqualValues(t, 6, c.Len())
}
