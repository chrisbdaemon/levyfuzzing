package coverage

import "fmt"

// Coverage contains a list of all the blocks reached
// by some number of test cases. Using hash table for
// the lookup speed.
type Coverage struct {
	values map[int]bool
}

// New creates a new Coverage object and returns it
func New(insVals ...int) *Coverage {
	c := &Coverage{}
	c.values = make(map[int]bool, 0)

	for _, insVal := range insVals {
		c.Add(insVal)
	}

	return c
}

// DifferenceCount returns the number of unique values in c1
// that do not exist in c2
func DifferenceCount(c1, c2 *Coverage) int64 {
	newC := c1.Union(c2)
	return int64(len(newC.values) - len(c2.values))
}

// Add instrumentation value to list of coverage items
func (c *Coverage) Add(insVal int) {
	// No duplicates
	if _, exists := c.values[insVal]; exists {
		return
	}

	c.values[insVal] = true
}

// Len returns the number of values
func (c *Coverage) Len() int64 { return int64(len(c.values)) }

// Union combines two Coverage objects and returns a new one
func (c *Coverage) Union(other *Coverage) (newCoverage *Coverage) {
	newCoverage = &Coverage{}
	newCoverage.values = make(map[int]bool, len(c.values))

	// copy the first set of values
	for insVal := range c.values {
		newCoverage.values[insVal] = true
	}

	for insVal := range other.values {
		newCoverage.Add(insVal)
	}

	return
}

func (c *Coverage) String() string {
	keys := make([]int, len(c.values))

	i := 0
	for k := range c.values {
		keys[i] = k
		i++
	}

	return fmt.Sprintf("<Coverage values=%v>", keys)
}
