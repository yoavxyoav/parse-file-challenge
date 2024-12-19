package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	var numOfLines int
	fmt.Print("how many points to generate? (100000000 ~=  1.2G): ")
	fmt.Scanln(&numOfLines)

	if numOfLines == 0 {
		numOfLines = 100000000
	}

	pf, err := os.Create("points.txt")
	if err != nil {
		panic(err)
	}
	defer pf.Close()

	pfv, err := os.Create("points-verify.txt")
	if err != nil {
		panic(err)
	}
	defer pfv.Close()

	min := -99.99
	max := 99.99

	sum1 := 0.0
	sum2 := 0.0

	pointsBuf := bufio.NewWriter(pf)

	fmt.Printf("Generating %d lines\n", numOfLines)
	for i := 0; i < numOfLines; i++ {
		r1 := min + rand.Float64()*(max-min)
		// r2 := min + rand.Float64()*(max-min)
		r2 := r1 + rand.Float64()*(max-r1)
		sr1 := fmt.Sprintf("%.2f", r1)
		r1, _ = strconv.ParseFloat(sr1, 8)
		sr2 := fmt.Sprintf("%.2f", r2)
		r2, _ = strconv.ParseFloat(sr2, 8)
		sum1 = sum1 + r1
		sum2 = sum2 + r2
		pointsBuf.WriteString(fmt.Sprintf("%.2f,%.2f\n", r1, r2))
	}

	if err := pointsBuf.Flush(); err != nil {
		panic(err)
	}

	// fmt.Println(fmt.Sprintf("%.10f", sum1), fmt.Sprintf("%.10f", sum2), numOfLines)
	s1 := math.Round(sum1*100) / 100
	s2 := math.Round(sum2*100) / 100
	fmt.Println(fmt.Sprintf("%.6f", s1), fmt.Sprintf("%.6f", s2), numOfLines)

	// pfv.WriteString(fmt.Sprintf("%s,%s,%d\n", s1[:len(s1)-8], s2[:len(s2)-8], numOfLines))
	pfv.WriteString(fmt.Sprintf("%.2f,%.2f,%d\n", s1, s2, numOfLines))
}
