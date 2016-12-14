package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/chrisbdaemon/levyfuzzing/testcase"
)

func main() {
	var b1, b2 float64
	var seedFilename = flag.String("seed", "", "input seed")
	var outputDir = flag.String("output", "", "directory to store generated test cases")
	var roundSize = flag.Uint("size", 500, "size of each iteration")
	var cmd = flag.String("cmd", "", "command under testing")
	var segmentCount = flag.Uint("segment-count", 0, "number of segments")
	flag.Parse()

	required := []string{"seedFilename", "outputDir", "roundSize", "cmd", "segmentCount"}
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required -%s argument/flag\n", req)
			flag.Usage()
			os.Exit(2)
		}
	}

	seed, err := testcase.New(*seedFilename)
	if err != nil {
		log.Fatalln("Unable to build seedFilename:", err)
	}

	var testCases []*testcase.TestCase
	var newTestCases []*testcase.TestCase
	a1, a2 := seedParams()
	segmentOffset := rand.Int63n(int64(*segmentCount))

	testCases = append(testCases, seed)
	for {
		seed = testCases[len(testCases)-1]
		newTestCases, err = testcase.GenerateNew(seed, a1, a2, segmentOffset, int64(*roundSize))
		if err != nil {
			log.Fatal("Unable to create test cases:", err)
		}

		score := evaluateTestCases(newTestCases, testCases)
		a1, a2 = updateParameters(score, a1, a2, b1, b2)

		testCases = append(testCases, newTestCases...)
	}
}

func updateParameters(score int64, a1, a2, b1, b2 float64) (a1New, a2New float64) {
	return
}

func evaluateTestCases(new, old []*testcase.TestCase) (score int64) {
	return
}

func seedParams() (a1, a2 float64) {
	rand.Seed(time.Now().Unix())
	a1 = rand.Float64() + float64(rand.Int63n(2))
	a2 = rand.Float64() + float64(rand.Int63n(2))
	return
}
