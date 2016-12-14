package testcase

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// TestCase is the object representing a single
// test case. It holds coverage results and other
// data specific to each test case.
type TestCase struct {
	filename     string
	coverage     []int64
	segmentCount int64
	segmentSize  int64
}

// New creates and returns a new TestCase object after
// verifying the given filename exists and is accessible
func New(filename string, segmentCount int64) (t *TestCase, err error) {
	// ensure file exists and is readable
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	stats, err := f.Stat()
	if err != nil {
		return
	}

	if segmentCount <= 0 {
		err = fmt.Errorf("Segment count must be greater than 0")
		return
	}

	t = &TestCase{}
	t.filename = filename
	t.segmentCount = segmentCount
	t.segmentSize = stats.Size() / segmentCount

	return
}

// GenerateNew performs the levy flight-based mutations on the initial seed
// test case. Parameters should be self-explanatory.
func GenerateNew(seed *TestCase, outputDir string,
	a1, a2 float64, segmentOffset, count int64) (testCases []*TestCase, err error) {

	var newTestCase *TestCase
	var newFilename string

	for i := int64(0); i < count; i++ {
		segmentOffset += flight(seed.segmentCount, a1)
		segmentOffset = wrapValue(segmentOffset, seed.segmentCount)

		newFilename, err = copyFile(seed.filename, outputDir)
		if err != nil {
			err = fmt.Errorf("Unable to copy seed: %s", err)
			return
		}

		newTestCase, err = New(newFilename, seed.segmentCount)
		if err != nil {
			err = fmt.Errorf("Unable to create new testcase: %s", err)
			return
		}

		err = newTestCase.performMutation(a2, outputDir, segmentOffset)
		if err != nil {
			err = fmt.Errorf("Unable to mutate new test case: %s", err)
			return
		}

		testCases = append(testCases, newTestCase)
	}

	return
}

func copyFile(filename, outputDirectory string) (string, error) {
	origFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer origFile.Close()

	newFilename := generateFilename(outputDirectory)
	newFile, err := os.Create(newFilename)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	bytes := make([]byte, 255)
	for {
		nBytes, err := origFile.Read(bytes)
		if nBytes == 0 {
			break
		}

		_, err = newFile.Write(bytes[:nBytes])
		if err != nil {
			return "", err
		}
	}

	return newFilename, nil
}

func handleOutOfBoundsValues(nBytes int, segment []byte, newValue int16) {
	newValueIndex := nBytes - 1

	for newValue > 255 || newValue < 0 {
		segment[newValueIndex] = byte(wrapValue(int64(newValue), 255))

		newValueIndex--
		if newValueIndex < 0 {
			break
		}

		if newValue > 255 {
			newValue = int16(segment[newValueIndex]) + 1
		} else {
			newValue = int16(segment[newValueIndex]) - 1
		}
	}

	if newValueIndex >= 0 {
		segment[newValueIndex] = byte(newValue)
	}
}

func (t *TestCase) performMutation(diffusivity float64, outputDir string,
	segmentOffset int64) (err error) {

	file, err := os.OpenFile(t.filename, os.O_RDWR, 0)
	if err != nil {
		err = fmt.Errorf("unable to open seed file: %s", err)
		return
	}
	defer file.Close()

	// Pull the correct segment from the file
	segment := make([]byte, t.segmentSize)
	file.Seek(segmentOffset*t.segmentSize, 0)
	nBytes, err := file.Read(segment)
	if err != nil && err != io.EOF {
		return
	}

	newValue := int16(segment[nBytes-1])
	newValue += int16(flight(255, diffusivity))

	// Handle underflows and overflows for multi-byte segments
	if newValue < 255 || newValue < 0 {
		handleOutOfBoundsValues(nBytes, segment, newValue)
	}

	// Write the segment back out
	file.Seek(segmentOffset*t.segmentSize, 0)
	nBytes, err = file.Write(segment)
	if err != nil {
		return
	}

	return
}

func wrapValue(val, max int64) (newVal int64) {
	newVal = val
	finished := false
	for finished == false {
		if newVal > max {
			newVal -= (max + 1)
		} else if newVal < 0 {
			newVal += (max + 1)
		} else {
			finished = true
		}
	}

	return
}

func flight(maxVal int64, diffusivity float64) (val int64) {
	diffusivity++
	valMin := math.Pow(1.0, -diffusivity)
	valMax := math.Pow(float64(maxVal), -diffusivity)

	randVal := rand.Float64()*(valMax-valMin) + valMin
	val = int64(math.Pow(randVal, -1.0/diffusivity))
	if rand.Int()%2 == 0 {
		val = -val
	}
	return
}

// ripped from:
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func generateFilename(directory string) string {
	length := 10
	b := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax

		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--

		}
		cache >>= letterIdxBits
		remain--

	}

	return filepath.Join(directory, string(b))
}
