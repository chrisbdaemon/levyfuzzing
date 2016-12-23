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
	c1 := coverage.New(1, 2, 3, 4)
	c2 := []*coverage.Coverage{
		coverage.New(1, 2, 3, 4, 5, 6),
		coverage.New(9, 2, 6, 4),
		coverage.New(200),
	}

	score, err := Score(c2, c1)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, score)
}
