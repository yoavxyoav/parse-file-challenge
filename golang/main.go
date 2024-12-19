package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func parse() (float64, float64, int64) {
	return 0.0, 0.0, 0
}

func compFloat(f1 float64, f2 float64) bool {
	i1 := int64(f1*100 + 0.5)
	i2 := int64(f1*100 + 0.5)
	return i1 == i2
}

func run() time.Duration {
	start := time.Now()
	s1, s2, lines := parse()
	elapsed := time.Since(start)

	data, err := os.ReadFile("points-verify.txt")
	if err != nil {
		panic(err)
	}

	parts := strings.Split(string(data[:len(data)-1]), ",")
	vl, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		panic(err)
	}
	f1, err := strconv.ParseFloat(parts[0], 8)
	if err != nil {
		panic(err)
	}
	f2, err := strconv.ParseFloat(parts[1], 8)
	if err != nil {
		panic(err)
	}

	if lines != vl {
		panic(fmt.Sprintf("Expected number of lines to be: %d got %d\n", vl, lines))
	}

	if !compFloat(s1, f1) {
		panic(fmt.Sprintf("Expected first number to be: %.2f got %.2f\n", f1, s1))
	}

	// if fmt.Sprintf("%.2f", s2) != fmt.Sprintf("%.2f", f2) {
	if !compFloat(s2, f2) {
		panic(fmt.Sprintf("Expected second number to be: %.2f got %.2f\n", f2, s2))
	}

	return elapsed
}

func main() {
	bestTime, err := time.ParseDuration("1h")
	if err != nil {
		panic(err)
	}

	for true {
		execTime := run()
		if execTime.Milliseconds() < bestTime.Milliseconds() {
			bestTime = execTime
			fmt.Printf("Execution time: %s\n", bestTime)
		}
	}
}
