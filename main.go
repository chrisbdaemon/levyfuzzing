package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/chrisbdaemon/levyfuzzing/evaluate"
	"github.com/chrisbdaemon/levyfuzzing/testcase"
)

func main() {
	var b1, b2 float64
	var seedFilename = flag.String("seed", "", "input seed")
	var outputDir = flag.String("output", "", "directory to store generated test cases")
	var roundSize = flag.Uint("size", 500, "size of each iteration")
	var cmd = flag.String("cmd", "", "command under testing")
	var segmentCount = flag.Uint("segment-count", 0, "number of segments")
	var showMapPath = flag.String("afl-showmap-path", "", "path to afl-showmap")
	flag.Parse()

	required := []string{"seed", "output", "size", "cmd", "segment-count"}
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required -%s argument/flag\n", req)
			flag.Usage()
			os.Exit(2)
		}
	}

	if *showMapPath != "" {
		testcase.ShowMapPath = *showMapPath
	}

	seed, err := testcase.New(*seedFilename, int64(*segmentCount))
	if err != nil {
		log.Fatalln("Unable to build seedFilename:", err)
	}

	err = seed.Execute(*cmd)
	if err != nil {
		log.Fatalln("Unable to execute binary:", err)
	}

	fmt.Println(seed.Coverage())

	var testCases []*testcase.TestCase
	var newTestCases []*testcase.TestCase
	a1, a2 := seedParams()
	segmentOffset := rand.Int63n(int64(*segmentCount))

	testCases = append(testCases, seed)
	for {
		seed = testCases[len(testCases)-1]
		newTestCases, err = testcase.GenerateNew(seed, *outputDir, a1, a2, segmentOffset, int64(*roundSize))
		if err != nil {
			log.Fatal("Unable to create test cases:", err)
		}

		err = executeTestCases(newTestCases, *cmd)
		if err != nil {
			log.Fatal("Unable to execute test cases:", err)
		}

		score, err := evaluate.Score(newTestCases, testCases)
		if err != nil {
			log.Fatalln("Unable to evaluate test cases:", err)
		}
		a1, a2 = updateParameters(int64(score), a1, a2, b1, b2)

		testCases = append(testCases, newTestCases...)
		break
	}
}

func executeTestCases(testCases []*testcase.TestCase, cmd string) (err error) {
	for _, testCase := range testCases {
		err = testCase.Execute(cmd)
		if err != nil {
			return
		}
	}
	return
}

func updateParameters(score int64, a1, a2, b1, b2 float64) (a1New, a2New float64) {
	a1New = a1
	a2New = a2

	// not yet implemented

	return
}

func seedParams() (a1, a2 float64) {
	rand.Seed(time.Now().Unix())

	// generate two floats.. [0,2)
	a1 = rand.Float64() + float64(rand.Int63n(2))
	a2 = rand.Float64() + float64(rand.Int63n(2))
	return
}
