package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func vanillaReadAndSum(filePath string) (float64, float64, int64) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	var sumX, sumY float64
	var lines int64

	buffer := make([]byte, 1024)
	line := ""

	for {
		n, err := file.Read(buffer)
		if n == 0 {
			break
		}

		if err != nil {
			fmt.Println("Error reading file:", err)
			return 0, 0, 0
		}

		data := string(buffer[:n])
		for _, char := range data {
			if char == '\n' {
				commaIdx := strings.Index(line, ",")
				x, _ := strconv.ParseFloat(line[:commaIdx], 64)
				y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)

				sumX += x
				sumY += y
				lines++
				line = ""
			} else {
				line += string(char)
			}
		}
	}

	if len(line) > 0 {
		commaIdx := strings.Index(line, ",")
		x, _ := strconv.ParseFloat(line[:commaIdx], 64)
		y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)

		sumX += x
		sumY += y
		lines++
	}

	return sumX, sumY, lines
}

// Concurrent implementation with workers
func concurrentReadAndSum(filePath string) (float64, float64, int64) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	lines := make(chan string, 100)
	results := make(chan [3]float64, 100)
	done := make(chan bool)

	worker := func() {
		for line := range lines {
			commaIdx := strings.Index(line, ",")
			x, _ := strconv.ParseFloat(line[:commaIdx], 64)
			y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
			results <- [3]float64{x, y, 1}
		}
		done <- true
	}

	// Start workers
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Read file and send lines to workers
	go func() {
		buffer := make([]byte, 1024)
		line := ""
		for {
			n, err := file.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				fmt.Println("Error reading file:", err)
				close(lines)
				return
			}
			data := string(buffer[:n])
			for _, char := range data {
				if char == '\n' {
					lines <- line
					line = ""
				} else {
					line += string(char)
				}
			}
		}
		if len(line) > 0 {
			lines <- line
		}
		close(lines)
	}()

	// Close results channel when workers are done
	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(results)
	}()

	// Aggregate results
	var sumX, sumY float64
	var linesCount int64
	for res := range results {
		sumX += res[0]
		sumY += res[1]
		linesCount += int64(res[2])
	}

	return sumX, sumY, linesCount
}

// Optimized implementation with byte-level processing
func optimizedConcurrentReadAndSum(filePath string) (float64, float64, int64) {
	const bufferSize = 32768 // Large buffer for efficient I/O
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	lines := make(chan string, 100)  // Channel for lines
	results := make(chan [3]float64) // Channel for results
	done := make(chan struct{})      // Channel to signal completion

	// Worker function to process lines
	worker := func() {
		for line := range lines {
			commaIdx := strings.Index(line, ",")
			if commaIdx == -1 {
				continue // Skip malformed lines
			}
			x, _ := strconv.ParseFloat(line[:commaIdx], 64)
			y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
			results <- [3]float64{x, y, 1}
		}
		done <- struct{}{} // Signal that the worker is done
	}

	// Start workers
	numWorkers := 16
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Goroutine to read the file and send lines to workers
	go func() {
		buffer := make([]byte, bufferSize)
		line := ""

		for {
			n, err := file.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				fmt.Println("Error reading file:", err)
				close(lines)
				return
			}

			data := string(buffer[:n])
			for _, char := range data {
				if char == '\n' {
					lines <- line
					line = ""
				} else {
					line += string(char)
				}
			}
		}

		// Send the last line if it doesn't end with a newline
		if len(line) > 0 {
			lines <- line
		}
		close(lines) // Close the lines channel to signal no more input
	}()

	// Goroutine to close the results channel after all workers finish
	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done // Wait for each worker to signal completion
		}
		close(results) // Close the results channel
	}()

	// Aggregate results
	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}

// Better optimized implementation with channel aggregation
func betterOptimizedConcurrentReadAndSum(filePath string) (float64, float64, int64) {
	const bufferSize = 4096
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	lines := make(chan string, 100)
	results := make(chan [3]float64, 100)
	done := make(chan bool)

	worker := func() {
		for line := range lines {
			commaIdx := strings.Index(line, ",")
			x, _ := strconv.ParseFloat(line[:commaIdx], 64)
			y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
			results <- [3]float64{x, y, 1}
		}
		done <- true
	}

	// Start workers
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Read file and send lines to workers
	go func() {
		buffer := make([]byte, bufferSize)
		line := ""
		for {
			n, err := file.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				fmt.Println("Error reading file:", err)
				close(lines)
				return
			}
			data := string(buffer[:n])
			for _, char := range data {
				if char == '\n' {
					lines <- line
					line = ""
				} else {
					line += string(char)
				}
			}
		}
		if len(line) > 0 {
			lines <- line
		}
		close(lines)
	}()

	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(results)
	}()

	// Aggregate results
	var sumX, sumY float64
	var lineCount int64
	for res := range results {
		sumX += res[0]
		sumY += res[1]
		lineCount += int64(res[2])
	}

	return sumX, sumY, lineCount
}

func streamingReadAndSum(filePath string) (float64, float64, int64) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	lines := make(chan string)       // Channel to stream lines to workers
	results := make(chan [3]float64) // Channel for results
	done := make(chan struct{})      // Done channel for worker coordination

	// Worker function to process lines
	worker := func() {
		for line := range lines {
			commaIdx := strings.Index(line, ",")
			if commaIdx == -1 {
				continue // Skip malformed lines
			}
			x, _ := strconv.ParseFloat(line[:commaIdx], 64)
			y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
			results <- [3]float64{x, y, 1}
		}
		done <- struct{}{} // Signal that the worker is done
	}

	// Start workers
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Goroutine to read the file line by line and send lines to workers
	go func() {
		buffer := make([]byte, 1024) // Small buffer for efficient reads
		line := ""                   // Line accumulator
		for {
			n, err := file.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				fmt.Println("Error reading file:", err)
				close(lines)
				return
			}

			// Process the buffer content
			data := string(buffer[:n])
			for _, char := range data {
				if char == '\n' {
					lines <- line
					line = ""
				} else {
					line += string(char)
				}
			}
		}

		// Handle the last line if it doesn't end with a newline
		if len(line) > 0 {
			lines <- line
		}
		close(lines) // Signal that no more lines will be sent
	}()

	// Goroutine to close the results channel after all workers finish
	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(results)
	}()

	// Aggregate results
	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}

func optimizedStreamingReadAndSum(filePath string) (float64, float64, int64) {
	const bufferSize = 65536 // Large buffer for efficient file reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	lines := make(chan []string, 10) // Channel for batches of lines
	results := make(chan [3]float64) // Channel for aggregated results
	done := make(chan struct{})      // Done channel for workers

	// Worker function to process a batch of lines
	worker := func() {
		var localSumX, localSumY float64
		var localLines int64

		for batch := range lines {
			for _, line := range batch {
				commaIdx := strings.Index(line, ",")
				if commaIdx == -1 {
					continue // Skip malformed lines
				}
				x, _ := strconv.ParseFloat(line[:commaIdx], 64)
				y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
				localSumX += x
				localSumY += y
				localLines++
			}
		}

		// Send aggregated results for this worker
		results <- [3]float64{localSumX, localSumY, float64(localLines)}
		done <- struct{}{} // Signal that the worker is done
	}

	// Start workers
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Read file in chunks and send batches of lines to workers
	go func() {
		buffer := make([]byte, bufferSize)
		line := ""
		batch := []string{}

		for {
			n, err := file.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				fmt.Println("Error reading file:", err)
				close(lines)
				return
			}

			data := string(buffer[:n])
			for _, char := range data {
				if char == '\n' {
					batch = append(batch, line)
					line = ""

					// Send a batch of 100 lines to workers
					if len(batch) >= 100 {
						lines <- batch
						batch = []string{}
					}
				} else {
					line += string(char)
				}
			}
		}

		// Handle remaining lines and the last line
		if len(line) > 0 {
			batch = append(batch, line)
		}
		if len(batch) > 0 {
			lines <- batch
		}
		close(lines) // Signal that no more batches will be sent
	}()

	// Close results channel after all workers finish
	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(results)
	}()

	// Aggregate results
	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}

func optimizedReadAndSum(filePath string) (float64, float64, int64) {
	const bufferSize = 65536 // Large buffer for efficient file reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	// Shared accumulators for results
	var totalSumX, totalSumY float64
	var totalLines int64

	// Buffer to read chunks of the file
	buffer := make([]byte, bufferSize)
	line := ""

	for {
		n, err := file.Read(buffer)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			return 0, 0, 0
		}

		// Process the buffer content
		data := string(buffer[:n])
		for _, char := range data {
			if char == '\n' {
				// Process the completed line
				commaIdx := strings.Index(line, ",")
				if commaIdx != -1 {
					x, _ := strconv.ParseFloat(line[:commaIdx], 64)
					y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
					totalSumX += x
					totalSumY += y
					totalLines++
				}
				line = ""
			} else {
				line += string(char)
			}
		}
	}

	// Process the last line if it doesn't end with a newline
	if len(line) > 0 {
		commaIdx := strings.Index(line, ",")
		if commaIdx != -1 {
			x, _ := strconv.ParseFloat(line[:commaIdx], 64)
			y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
			totalSumX += x
			totalSumY += y
			totalLines++
		}
	}

	return totalSumX, totalSumY, totalLines
}

func fastReadAndSumWithChannels(filePath string) (float64, float64, int64) {
	const bufferSize = 65536 // Large buffer for efficient reading
	const batchSize = 100    // Number of lines per batch
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	lines := make(chan []string, 10) // Channel for batches of lines
	results := make(chan [3]float64) // Channel for aggregated results
	done := make(chan struct{})      // Done channel for workers

	// Worker function to process batches of lines
	worker := func() {
		var localSumX, localSumY float64
		var localLineCount int64

		for batch := range lines {
			for _, line := range batch {
				commaIdx := strings.Index(line, ",")
				if commaIdx == -1 {
					continue // Skip malformed lines
				}
				x, _ := strconv.ParseFloat(line[:commaIdx], 64)
				y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
				localSumX += x
				localSumY += y
				localLineCount++
			}
		}

		// Send aggregated results
		results <- [3]float64{localSumX, localSumY, float64(localLineCount)}
		done <- struct{}{} // Signal completion
	}

	// Start workers
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Goroutine to read the file and send batches of lines to workers
	go func() {
		buffer := make([]byte, bufferSize)
		line := ""
		batch := []string{}

		for {
			n, err := file.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				fmt.Println("Error reading file:", err)
				close(lines)
				return
			}

			data := string(buffer[:n])
			for _, char := range data {
				if char == '\n' {
					batch = append(batch, line)
					line = ""
					if len(batch) >= batchSize {
						lines <- batch // Send a full batch
						batch = []string{}
					}
				} else {
					line += string(char)
				}
			}
		}

		// Send remaining lines
		if len(line) > 0 {
			batch = append(batch, line)
		}
		if len(batch) > 0 {
			lines <- batch
		}

		close(lines) // Signal no more lines
	}()

	// Goroutine to close the results channel after all workers finish
	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(results)
	}()

	// Aggregate final results
	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}

func optimizedParsingAndSum(filePath string) (float64, float64, int64) {
	const bufferSize = 65536 // Large buffer for efficient reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	var totalSumX, totalSumY float64
	var totalLines int64

	buffer := make([]byte, bufferSize)
	lineStart := 0                          // Track the start of the current line
	leftover := make([]byte, 0, bufferSize) // Buffer for leftover partial lines

	for {
		n, err := file.Read(buffer)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			return 0, 0, 0
		}

		// Combine leftover from previous buffer with new data
		data := append(leftover, buffer[:n]...)
		lineStart = 0 // Reset lineStart for the combined data

		for i := 0; i < len(data); i++ {
			if data[i] == '\n' {
				// Parse the line
				line := data[lineStart:i]
				commaIdx := findComma(line)
				if commaIdx != -1 {
					x := parseFloat(line[:commaIdx])
					y := parseFloat(line[commaIdx+1:])
					totalSumX += x
					totalSumY += y
					totalLines++
				}
				lineStart = i + 1 // Move to the start of the next line
			}
		}

		// Handle leftover partial line
		if lineStart < len(data) {
			leftover = append(leftover[:0], data[lineStart:]...) // Save leftover data
		} else {
			leftover = leftover[:0] // Clear leftover buffer if no partial line
		}
	}

	// Process any remaining line if the file doesn't end with a newline
	if len(leftover) > 0 {
		line := leftover
		commaIdx := findComma(line)
		if commaIdx != -1 {
			x := parseFloat(line[:commaIdx])
			y := parseFloat(line[commaIdx+1:])
			totalSumX += x
			totalSumY += y
			totalLines++
		}
	}

	return totalSumX, totalSumY, totalLines
}

//func findComma(line []byte) int {
//	for i, b := range line {
//		if b == ',' {
//			return i
//		}
//	}
//	return -1
//}
//
//func parseFloat(b []byte) float64 {
//	var result float64
//	var decimalPlace float64 = 1
//	var isFraction bool
//
//	for _, c := range b {
//		if c == '.' {
//			isFraction = true
//			continue
//		}
//		if c >= '0' && c <= '9' {
//			digit := float64(c - '0')
//			if isFraction {
//				decimalPlace /= 10
//				result += digit * decimalPlace
//			} else {
//				result = result*10 + digit
//			}
//		}
//	}
//	return result
//}

func optimizedParsingWithPointers(filePath string) (float64, float64, int64) {
	const bufferSize = 65536 // Large buffer for efficient reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	var totalSumX, totalSumY float64
	var totalLines int64

	buffer := make([]byte, bufferSize)
	lineStart := new(int)                   // Pointer to track the start of the line
	leftover := make([]byte, 0, bufferSize) // Buffer for leftover partial lines

	for {
		n, err := file.Read(buffer)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			return 0, 0, 0
		}

		// Combine leftover from previous buffer with new data
		data := append(leftover, buffer[:n]...)
		*lineStart = 0 // Reset lineStart for the combined data

		for i := 0; i < len(data); i++ {
			if data[i] == '\n' {
				// Parse the line
				line := data[*lineStart:i]
				commaIdx := findComma(line)
				if commaIdx != -1 {
					x := parseFloat(line[:commaIdx])
					y := parseFloat(line[commaIdx+1:])
					totalSumX += x
					totalSumY += y
					totalLines++
				}
				*lineStart = i + 1 // Move to the start of the next line
			}
		}

		// Handle leftover partial line
		if *lineStart < len(data) {
			leftover = append(leftover[:0], data[*lineStart:]...) // Save leftover data
		} else {
			leftover = leftover[:0] // Clear leftover buffer if no partial line
		}
	}

	// Process any remaining line if the file doesn't end with a newline
	if len(leftover) > 0 {
		line := leftover
		commaIdx := findComma(line)
		if commaIdx != -1 {
			x := parseFloat(line[:commaIdx])
			y := parseFloat(line[commaIdx+1:])
			totalSumX += x
			totalSumY += y
			totalLines++
		}
	}

	return totalSumX, totalSumY, totalLines
}
func combinedOptimizedParsing(filePath string) (float64, float64, int64) {
	const bufferSize = 65536 // Large buffer for efficient reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	var totalSumX, totalSumY float64
	var totalLines int64

	buffer := make([]byte, bufferSize)
	var leftover []byte // Slice to store leftover partial lines
	lineStart := 0      // Track the start of the current line

	for {
		n, err := file.Read(buffer)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			return 0, 0, 0
		}

		// Combine leftover with the current buffer
		data := append(leftover, buffer[:n]...)
		lineStart = 0 // Reset lineStart for the combined data

		for i := 0; i < len(data); i++ {
			if data[i] == '\n' {
				// Parse the line
				if lineStart < i { // Ensure valid range for slicing
					line := data[lineStart:i]
					commaIdx := findComma(line)
					if commaIdx != -1 {
						x := parseFloat(line[:commaIdx])
						y := parseFloat(line[commaIdx+1:])
						totalSumX += x
						totalSumY += y
						totalLines++
					}
				}
				lineStart = i + 1 // Move to the start of the next line
			}
		}

		// Handle leftover partial line
		if lineStart < len(data) {
			leftover = append([]byte{}, data[lineStart:]...) // Save leftover data
		} else {
			leftover = nil // Clear leftover if no partial line
		}
	}

	// Process any remaining line if the file doesn't end with a newline
	if len(leftover) > 0 {
		commaIdx := findComma(leftover)
		if commaIdx != -1 {
			x := parseFloat(leftover[:commaIdx])
			y := parseFloat(leftover[commaIdx+1:])
			totalSumX += x
			totalSumY += y
			totalLines++
		}
	}

	return totalSumX, totalSumY, totalLines
}

func findComma(line []byte) int {
	for i, b := range line {
		if b == ',' {
			return i
		}
	}
	return -1
}

func parseFloat(b []byte) float64 {
	var result float64
	var decimalPlace float64 = 1
	var isFraction bool

	for _, c := range b {
		if c == '.' {
			isFraction = true
			continue
		}
		if c >= '0' && c <= '9' {
			digit := float64(c - '0')
			if isFraction {
				decimalPlace /= 10
				result += digit * decimalPlace
			} else {
				result = result*10 + digit
			}
		}
	}
	return result
}

func optimizedParsingWithChannels(filePath string) (float64, float64, int64) {
	const bufferSize = 65536 // Large buffer for efficient reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	chunks := make(chan []byte, 10)  // Channel for file chunks
	results := make(chan [3]float64) // Channel for worker results
	done := make(chan struct{})      // Done channel for workers

	// Worker function to process chunks
	worker := func() {
		var localSumX, localSumY float64
		var localLines int64

		for chunk := range chunks {
			lineStart := 0
			for i := 0; i < len(chunk); i++ {
				if chunk[i] == '\n' {
					// Parse the line
					line := chunk[lineStart:i]
					commaIdx := findComma(line)
					if commaIdx != -1 {
						x := parseFloat(line[:commaIdx])
						y := parseFloat(line[commaIdx+1:])
						localSumX += x
						localSumY += y
						localLines++
					}
					lineStart = i + 1
				}
			}

			// Handle leftover partial line
			if lineStart < len(chunk) {
				chunk = append([]byte{}, chunk[lineStart:]...) // Carry forward partial line
			} else {
				chunk = nil
			}
		}

		// Send local results
		results <- [3]float64{localSumX, localSumY, float64(localLines)}
		done <- struct{}{}
	}

	// Start workers
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Goroutine to read file and send chunks to workers
	go func() {
		buffer := make([]byte, bufferSize)
		leftover := make([]byte, 0, bufferSize)

		for {
			n, err := file.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				fmt.Println("Error reading file:", err)
				close(chunks)
				return
			}

			// Combine leftover with the current buffer
			chunk := append(leftover, buffer[:n]...)
			leftover = nil // Reset leftover

			// Find last newline in the chunk
			lastNewline := -1
			for i := len(chunk) - 1; i >= 0; i-- {
				if chunk[i] == '\n' {
					lastNewline = i
					break
				}
			}

			if lastNewline != -1 {
				// Send complete lines to workers
				chunks <- chunk[:lastNewline+1]
				// Save leftover partial line
				leftover = append([]byte{}, chunk[lastNewline+1:]...)
			} else {
				// If no newline, entire chunk is leftover
				leftover = append([]byte{}, chunk...)
			}
		}

		// Handle leftover as the final line
		if len(leftover) > 0 {
			chunks <- leftover
		}

		close(chunks) // Signal no more chunks
	}()

	// Close results channel after all workers finish
	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(results)
	}()

	// Aggregate final results
	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}

func optimizedParsingWithChannels_2(filePath string) (float64, float64, int64) {
	const bufferSize = 65536 // Large buffer for efficient reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	chunks := make(chan []byte, 10)  // Channel for file chunks
	results := make(chan [3]float64) // Channel for aggregated results
	done := make(chan struct{})      // Done channel for workers

	// Worker function
	worker := func() {
		var localSumX, localSumY float64
		var localLines int64

		for chunk := range chunks {
			lineStart := 0
			for i := 0; i < len(chunk); i++ {
				if chunk[i] == '\n' {
					// Parse the line
					line := chunk[lineStart:i]
					commaIdx := findComma(line)
					if commaIdx != -1 {
						x := parseFloat(line[:commaIdx])
						y := parseFloat(line[commaIdx+1:])
						localSumX += x
						localSumY += y
						localLines++
					}
					lineStart = i + 1
				}
			}

			// Handle leftover partial line
			if lineStart < len(chunk) {
				leftover := append([]byte{}, chunk[lineStart:]...)
				chunks <- leftover // Carry over leftover to next chunk
			}
		}

		// Send aggregated results for this worker
		results <- [3]float64{localSumX, localSumY, float64(localLines)}
		done <- struct{}{}
	}

	// Start workers
	numWorkers := 4 // Adjust based on system capability
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Goroutine to read file and send chunks to workers
	go func() {
		buffer := make([]byte, bufferSize)
		var leftover []byte

		for {
			n, err := file.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				fmt.Println("Error reading file:", err)
				close(chunks)
				return
			}

			// Combine leftover from previous buffer with new data
			chunk := append(leftover, buffer[:n]...)
			leftover = nil // Reset leftover

			// Find the last newline in the chunk
			lastNewline := -1
			for i := len(chunk) - 1; i >= 0; i-- {
				if chunk[i] == '\n' {
					lastNewline = i
					break
				}
			}

			if lastNewline != -1 {
				// Send complete lines to workers
				chunks <- chunk[:lastNewline+1]
				// Save leftover partial line
				leftover = append([]byte{}, chunk[lastNewline+1:]...)
			} else {
				// If no newline, entire chunk is leftover
				leftover = append([]byte{}, chunk...)
			}
		}

		// Send any remaining leftover as the last line
		if len(leftover) > 0 {
			chunks <- leftover
		}
		close(chunks) // Signal no more chunks
	}()

	// Goroutine to close results channel after all workers finish
	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(results)
	}()

	// Aggregate final results
	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}

func syncReadAndSum(filePath string) (float64, float64, int64) {
	const bufferSize = 65536
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex

	var totalSumX, totalSumY float64
	var totalLines int64

	lines := make(chan string, 100)

	worker := func() {
		defer wg.Done()
		var localSumX, localSumY float64
		var localLines int64

		for line := range lines {
			commaIdx := strings.Index(line, ",")
			if commaIdx != -1 {
				x, _ := strconv.ParseFloat(line[:commaIdx], 64)
				y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
				localSumX += x
				localSumY += y
				localLines++
			}
		}

		mu.Lock()
		totalSumX += localSumX
		totalSumY += localSumY
		totalLines += localLines
		mu.Unlock()
	}

	numWorkers := 4
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	buffer := make([]byte, bufferSize)
	line := ""

	for {
		n, err := file.Read(buffer)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			close(lines)
			return 0, 0, 0
		}

		data := string(buffer[:n])
		for _, char := range data {
			if char == '\n' {
				lines <- line
				line = ""
			} else {
				line += string(char)
			}
		}
	}

	if len(line) > 0 {
		lines <- line
	}
	close(lines)
	wg.Wait()

	return totalSumX, totalSumY, totalLines
}

func bufioReadAndSum(filePath string) (float64, float64, int64) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var totalSumX, totalSumY float64
	var totalLines int64

	for scanner.Scan() {
		line := scanner.Text()
		commaIdx := strings.Index(line, ",")
		if commaIdx != -1 {
			x, _ := strconv.ParseFloat(line[:commaIdx], 64)
			y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
			totalSumX += x
			totalSumY += y
			totalLines++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return 0, 0, 0
	}

	return totalSumX, totalSumY, totalLines
}

func bufioWithSyncReadAndSum(filePath string) (float64, float64, int64) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex

	var totalSumX, totalSumY float64
	var totalLines int64

	lines := make(chan string, 100)

	worker := func() {
		defer wg.Done()
		var localSumX, localSumY float64
		var localLines int64

		for line := range lines {
			commaIdx := strings.Index(line, ",")
			if commaIdx != -1 {
				x, _ := strconv.ParseFloat(line[:commaIdx], 64)
				y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
				localSumX += x
				localSumY += y
				localLines++
			}
		}

		mu.Lock()
		totalSumX += localSumX
		totalSumY += localSumY
		totalLines += localLines
		mu.Unlock()
	}

	numWorkers := 4
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines <- scanner.Text()
	}

	close(lines)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return 0, 0, 0
	}

	return totalSumX, totalSumY, totalLines
}

func bufioWithChannelsReadAndSum(filePath string) (float64, float64, int64) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	lines := make(chan string, 100)  // Channel to pass lines to workers
	results := make(chan [3]float64) // Channel for aggregated results
	done := make(chan struct{})      // Channel to signal worker completion

	// Worker function to process lines
	worker := func() {
		var localSumX, localSumY float64
		var localLines int64

		for line := range lines {
			commaIdx := strings.Index(line, ",")
			if commaIdx != -1 {
				x, _ := strconv.ParseFloat(line[:commaIdx], 64)
				y, _ := strconv.ParseFloat(line[commaIdx+1:], 64)
				localSumX += x
				localSumY += y
				localLines++
			}
		}

		// Send the local results to the results channel
		results <- [3]float64{localSumX, localSumY, float64(localLines)}
		done <- struct{}{} // Signal this worker is done
	}

	// Start workers
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Goroutine to read file line by line and send lines to workers
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines) // Close the lines channel when done
	}()

	// Goroutine to close the results channel after all workers finish
	go func() {
		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(results)
	}()

	// Aggregate final results
	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}

func optimizedParsingWithReadAt(filePath string) (float64, float64, int64) {
	const bufferSize = 65536 // Large buffer for efficient reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	stat, _ := file.Stat()
	fileSize := stat.Size()
	numWorkers := 4
	chunkSize := fileSize / int64(numWorkers)

	results := make(chan [3]float64, numWorkers)
	var wg sync.WaitGroup

	worker := func(offset, size int64) {
		defer wg.Done()
		buffer := make([]byte, size)
		_, err := file.ReadAt(buffer, offset)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file chunk:", err)
			return
		}

		var localSumX, localSumY float64
		var localLines int64

		lineStart := 0
		for i := 0; i < len(buffer); i++ {
			if buffer[i] == '\n' {
				line := buffer[lineStart:i]
				commaIdx := findComma(line)
				if commaIdx != -1 {
					x := parseFloat(line[:commaIdx])
					y := parseFloat(line[commaIdx+1:])
					localSumX += x
					localSumY += y
					localLines++
				}
				lineStart = i + 1
			}
		}

		// Handle leftover line in the chunk
		if lineStart < len(buffer) {
			line := buffer[lineStart:]
			commaIdx := findComma(line)
			if commaIdx != -1 {
				x := parseFloat(line[:commaIdx])
				y := parseFloat(line[commaIdx+1:])
				localSumX += x
				localSumY += y
				localLines++
			}
		}

		results <- [3]float64{localSumX, localSumY, float64(localLines)}
	}

	// Spawn workers to process file chunks
	for i := 0; i < numWorkers; i++ {
		offset := int64(i) * chunkSize
		size := chunkSize
		if i == numWorkers-1 { // Last worker takes remaining bytes
			size += fileSize % int64(numWorkers)
		}
		wg.Add(1)
		go worker(offset, size)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate final results
	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}

func optimizedParsingWithReadAtEnhanced(filePath string) (float64, float64, int64) {
	//const bufferSize = 65536
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	stat, _ := file.Stat()
	fileSize := stat.Size()
	numWorkers := runtime.NumCPU()
	chunkSize := fileSize / int64(numWorkers)

	results := make(chan [3]float64, numWorkers)
	var wg sync.WaitGroup

	worker := func(offset, size int64) {
		defer wg.Done()
		buffer := make([]byte, size)
		_, err := file.ReadAt(buffer, offset)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file chunk:", err)
			return
		}

		var localSumX, localSumY float64
		var localLines int64

		lineStart := 0
		for i := 0; i < len(buffer); i++ {
			if buffer[i] == '\n' {
				line := buffer[lineStart:i]
				commaIdx := findComma(line)
				if commaIdx != -1 {
					x := parseFloat(line[:commaIdx])
					y := parseFloat(line[commaIdx+1:])
					localSumX += x
					localSumY += y
					localLines++
				}
				lineStart = i + 1
			}
		}

		if lineStart < len(buffer) {
			line := buffer[lineStart:]
			commaIdx := findComma(line)
			if commaIdx != -1 {
				x := parseFloat(line[:commaIdx])
				y := parseFloat(line[commaIdx+1:])
				localSumX += x
				localSumY += y
				localLines++
			}
		}

		results <- [3]float64{localSumX, localSumY, float64(localLines)}
	}

	for i := 0; i < numWorkers; i++ {
		offset := int64(i) * chunkSize
		size := chunkSize
		if i == numWorkers-1 {
			size += fileSize % int64(numWorkers)
		}
		wg.Add(1)
		go worker(offset, size)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}
func optimizedParsingWithReadAtAndBuffer(filePath string) (float64, float64, int64) {
	const bufferSize = 65536
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0
	}
	defer file.Close()

	stat, _ := file.Stat()
	fileSize := stat.Size()
	numWorkers := runtime.NumCPU()
	chunkSize := fileSize / int64(numWorkers)

	results := make(chan [3]float64, numWorkers)
	var wg sync.WaitGroup

	worker := func(offset, size int64) {
		defer wg.Done()
		buffer := make([]byte, bufferSize)
		var localSumX, localSumY float64
		var localLines int64

		bytesRead := int64(0)
		leftover := make([]byte, 0)

		for bytesRead < size {
			toRead := bufferSize
			if size-bytesRead < int64(bufferSize) {
				toRead = int(size - bytesRead)
			}

			n, err := file.ReadAt(buffer[:toRead], offset+bytesRead)
			if err != nil && err != io.EOF {
				fmt.Println("Error reading file chunk:", err)
				return
			}

			bytesRead += int64(n)
			data := append(leftover, buffer[:n]...)
			leftover = nil

			lineStart := 0
			for i := 0; i < len(data); i++ {
				if data[i] == '\n' {
					line := data[lineStart:i]
					commaIdx := findComma(line)
					if commaIdx != -1 {
						x := parseFloat(line[:commaIdx])
						y := parseFloat(line[commaIdx+1:])
						localSumX += x
						localSumY += y
						localLines++
					}
					lineStart = i + 1
				}
			}

			if lineStart < len(data) {
				leftover = data[lineStart:]
			}
		}

		if len(leftover) > 0 {
			commaIdx := findComma(leftover)
			if commaIdx != -1 {
				x := parseFloat(leftover[:commaIdx])
				y := parseFloat(leftover[commaIdx+1:])
				localSumX += x
				localSumY += y
				localLines++
			}
		}

		results <- [3]float64{localSumX, localSumY, float64(localLines)}
	}

	for i := 0; i < numWorkers; i++ {
		offset := int64(i) * chunkSize
		size := chunkSize
		if i == numWorkers-1 {
			size += fileSize % int64(numWorkers)
		}
		wg.Add(1)
		go worker(offset, size)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var totalSumX, totalSumY float64
	var totalLines int64

	for res := range results {
		totalSumX += res[0]
		totalSumY += res[1]
		totalLines += int64(res[2])
	}

	return totalSumX, totalSumY, totalLines
}
func parse() (float64, float64, int64) {
	//return optimizedParsingWithReadAtAndBuffer("points.txt")
	return optimizedParsingWithReadAtEnhanced("points.txt")
	//return optimizedParsingWithReadAt("points.txt")
	//return bufioWithChannelsReadAndSum("points.txt")
	//return bufioWithSyncReadAndSum("points.txt")
	//return bufioReadAndSum("points.txt")
	//return syncReadAndSum("points.txt")
	//return optimizedParsingWithChannels_2("points.txt")
	//return optimizedParsingWithChannels("points.txt")
	//return combinedOptimizedParsing("points.txt")
	//return optimizedParsingWithPointers("points.txt")
	//return optimizedParsingAndSum("points.txt")
	//return fastReadAndSumWithChannels("points.txt")
	//return optimizedReadAndSum("points.txt")
	//return optimizedConcurrentReadAndSum("points.txt")
	//return streamingReadAndSum("points.txt")
	//return betterOptimizedConcurrentReadAndSum("points.txt")
	//return optimizedConcurrentReadAndSum("points.txt")
	//return concurrentReadAndSum("points.txt")
	//return vanillaReadAndSum("points.txt")

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

	//runAllFunctionsAndMeasureAvg(5)

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
