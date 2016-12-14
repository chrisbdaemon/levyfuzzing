package testcase

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New("fake", 2)
	assert.Error(t, err, "should fail from bad filename")

	testcase, err := New("testcase.go", 4)
	assert.Nil(t, err)
	assert.NotNil(t, testcase, "should return ptr to testcase")
}

func TestCopyFile(t *testing.T) {
	filename := "testcase.go"
	newFilename, err := copyFile(filename, "/tmp")
	assert.Nil(t, err)

	_, err = os.Stat(newFilename)
	assert.Nil(t, err)

	os.Remove(newFilename)
}

func TestGenerateNew(t *testing.T) {
	seed := &TestCase{}
	seed.segmentCount = 6
	seed.segmentSize = 4
	seed.filename = "testcase.go"

	// random, bogus parameters
	newTestCases, err := GenerateNew(seed, "/tmp", 1.41, 0.256, 3, 2)
	assert.IsType(t, []*TestCase{}, newTestCases, "GenerateNew should return valid slice")
	assert.Equal(t, 2, len(newTestCases), "should create new test cases")
	assert.Nil(t, err)

	// delete newly generated test cases
	for _, tc := range newTestCases {
		_, err := os.Stat(tc.filename)
		assert.Nil(t, err)

		//os.Remove(tc.filename)
	}
}

func TestHandleOutOfBoundsValues(t *testing.T) {
	var segment []byte

	segment = []byte{0, 0, 0, 255}
	handleOutOfBoundsValues(4, segment, 256)
	assert.True(t, byteSlicesAreEqual(segment, []byte{0, 0, 1, 0}))

	segment = []byte{0, 0, 1, 255}
	handleOutOfBoundsValues(4, segment, -1)
	assert.True(t, byteSlicesAreEqual(segment, []byte{0, 0, 0, 255}))

	segment = []byte{0, 1, 0, 255}
	handleOutOfBoundsValues(4, segment, -1)
	assert.True(t, byteSlicesAreEqual(segment, []byte{0, 0, 255, 255}))
}

func TestFlight(t *testing.T) {
	for i := 0; i < 500; i++ {
		a := flight(3, 0.142)
		assert.True(t, a <= 3 && a >= -3, fmt.Sprintf("should be between -3 and 3, got %d", a))
	}
}

func TestWrapValue(t *testing.T) {
	tests := make(map[int64]int64, 7)

	tests[-3] = 253
	tests[-2] = 254
	tests[-1] = 255
	tests[0] = 0
	tests[255] = 255
	tests[256] = 0
	tests[257] = 1

	for input, expected := range tests {
		retval := wrapValue(input, 255)
		assert.Equal(t, expected, retval, fmt.Sprintf("should failed with input: %d", input))
	}
}

func byteSlicesAreEqual(s1, s2 []byte) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}
