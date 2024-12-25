package main

import (
	"fmt"
	"time"
)

type functionInfo struct {
	name     string
	function func() (float64, float64, int64)
}

// The measuring function
func runAllFunctionsAndMeasureAvg(n int) {
	// List all parsing functions with their names
	functions := []functionInfo{
		{"vanillaReadAndSum", func() (float64, float64, int64) { return vanillaReadAndSum("points.txt") }},
		{"concurrentReadAndSum", func() (float64, float64, int64) { return concurrentReadAndSum("points.txt") }},
		{"optimizedConcurrentReadAndSum", func() (float64, float64, int64) { return optimizedConcurrentReadAndSum("points.txt") }},
		{"betterOptimizedConcurrentReadAndSum", func() (float64, float64, int64) { return betterOptimizedConcurrentReadAndSum("points.txt") }},
		{"streamingReadAndSum", func() (float64, float64, int64) { return streamingReadAndSum("points.txt") }},
		{"optimizedStreamingReadAndSum", func() (float64, float64, int64) { return optimizedStreamingReadAndSum("points.txt") }},
		{"optimizedReadAndSum", func() (float64, float64, int64) { return optimizedReadAndSum("points.txt") }},
		{"fastReadAndSumWithChannels", func() (float64, float64, int64) { return fastReadAndSumWithChannels("points.txt") }},
		{"bufioReadAndSum", func() (float64, float64, int64) { return bufioReadAndSum("points.txt") }},
		{"bufioWithSyncReadAndSum", func() (float64, float64, int64) { return bufioWithSyncReadAndSum("points.txt") }},
		{"bufioWithChannelsReadAndSum", func() (float64, float64, int64) { return bufioWithChannelsReadAndSum("points.txt") }},
		{"syncReadAndSum", func() (float64, float64, int64) { return syncReadAndSum("points.txt") }},
		{"optimizedParsingWithReadAt", func() (float64, float64, int64) { return optimizedParsingWithReadAt("points.txt") }},
		{"optimizedParsingWithReadAtEnhanced", func() (float64, float64, int64) { return optimizedParsingWithReadAtEnhanced("points.txt") }},
		{"optimizedParsingWithReadAtAndBuffer", func() (float64, float64, int64) { return optimizedParsingWithReadAtAndBuffer("points.txt") }},
		{"optimizedParsingWithChannels", func() (float64, float64, int64) { return optimizedParsingWithChannels("points.txt") }},
		{"optimizedParsingWithChannels_2", func() (float64, float64, int64) { return optimizedParsingWithChannels_2("points.txt") }},
		{"optimizedParsingAndSum", func() (float64, float64, int64) { return optimizedParsingAndSum("points.txt") }},
		{"combinedOptimizedParsing", func() (float64, float64, int64) { return combinedOptimizedParsing("points.txt") }},
	}

	results := make(map[string]time.Duration)

	// Measure runtime for each function
	for _, fi := range functions {
		var totalDuration time.Duration

		for i := 0; i < n; i++ {
			start := time.Now()
			_, _, _ = fi.function() // Run the function
			elapsed := time.Since(start)
			totalDuration += elapsed
		}

		// Calculate and store average runtime
		results[fi.name] = totalDuration / time.Duration(n)
	}

	// Find the fastest and slowest functions
	var fastest, slowest string
	var fastestTime, slowestTime time.Duration

	for name, avgTime := range results {
		if fastest == "" || avgTime < fastestTime {
			fastest = name
			fastestTime = avgTime
		}
		if slowest == "" || avgTime > slowestTime {
			slowest = name
			slowestTime = avgTime
		}
	}

	// Print the results
	fmt.Println("Average Runtimes:")
	for name, avgTime := range results {
		fmt.Printf("%s: %v\n", name, avgTime)
	}

	fmt.Printf("\nFastest Function: %s with avg time %v\n", fastest, fastestTime)
	fmt.Printf("Slowest Function: %s with avg time %v\n", slowest, slowestTime)
}
