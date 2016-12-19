package coverage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	c := New(1, 3)

	assert.Equal(t, len(c.values), 2)
	c.Add(2)
	assert.Equal(t, len(c.values), 3)
	c.Add(2)
	assert.Equal(t, len(c.values), 3)
}

func TestUnion(t *testing.T) {
	c1 := New(1, 2, 3, 4)
	c2 := New(2, 3, 4, 5, 6)

	c3 := c1.Union(c2)
	assert.Equal(t, 6, len(c3.values))
}

func TestDifferenceCount(t *testing.T) {
	c1 := New(5, 6)
	c2 := New(1, 2, 3, 4, 5, 6)

	assert.Equal(t, int64(0), DifferenceCount(c1, c2))
	assert.Equal(t, int64(4), DifferenceCount(c2, c1))
}
