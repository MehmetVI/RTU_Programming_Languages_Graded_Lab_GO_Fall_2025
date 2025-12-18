package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Name: Mehmet Kaan ULU
// Student ID: 231ADB102

// -----------------------------
// Entry point
// -----------------------------

func main() {
	n := flag.Int("r", -1, "generate N random integers (N >= 10)")
	flag.Parse()

	if *n == -1 {
		log.Fatal("Usage: gosort -r N")
	}

	if err := runRandom(*n); err != nil {
		log.Fatal(err)
	}
}

// -----------------------------
// -r mode logic
// -----------------------------

func runRandom(n int) error {
	if n < 10 {
		return errors.New("N must be >= 10")
	}

	numbers := generateRandomNumbers(n)

	fmt.Println("Original numbers:")
	fmt.Println(numbers)

	chunks := splitIntoChunks(numbers)

	fmt.Println("\nChunks before sorting:")
	printChunks(chunks)

	sortedChunks := sortChunksConcurrently(chunks)

	fmt.Println("\nChunks after sorting:")
	printChunks(sortedChunks)

	result := mergeSortedChunks(sortedChunks)

	fmt.Println("\nFinal sorted result:")
	fmt.Println(result)

	return nil
}

// -----------------------------
// Chunking logic
// -----------------------------

func splitIntoChunks(numbers []int) [][]int {
	n := len(numbers)

	numChunks := int(math.Ceil(math.Sqrt(float64(n))))
	if numChunks < 4 {
		numChunks = 4
	}
	if numChunks > n {
		numChunks = n
	}

	
	chunks := make([][]int, 0, numChunks)

	base := n / numChunks
	rem := n % numChunks

	start := 0
	for i := 0; i < numChunks; i++ {
		size := base
		if i < rem {
			size++
		}
		end := start + size
		if end > n {
			end = n
		}
		chunk := append([]int(nil), numbers[start:end]...)
		chunks = append(chunks, chunk)
		start = end
	}

	return chunks
}

// -----------------------------
// Concurrent sorting
// -----------------------------

func sortChunksConcurrently(chunks [][]int) [][]int {
	type result struct {
		idx   int
		chunk []int
	}

	out := make(chan result, len(chunks))
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for i := range chunks {
		i := i
		go func() {
			defer wg.Done()
			sort.Ints(chunks[i])
			out <- result{idx: i, chunk: chunks[i]}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	sorted := make([][]int, len(chunks))
	for r := range out {
		sorted[r.idx] = r.chunk
	}
	return sorted
}

// -----------------------------
// Merge logic
// -----------------------------

func mergeSortedChunks(chunks [][]int) []int {
	if len(chunks) == 0 {
		return nil
	}

	merged := chunks[0]
	for i := 1; i < len(chunks); i++ {
		merged = mergeTwoSorted(merged, chunks[i])
	}
	return merged
}

func mergeTwoSorted(a, b []int) []int {
	res := make([]int, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}
	res = append(res, a[i:]...)
	res = append(res, b[j:]...)
	return res
}

// -----------------------------
// Helpers
// -----------------------------

func generateRandomNumbers(n int) []int {
	rand.Seed(time.Now().UnixNano())

	nums := make([]int, n)
	for i := 0; i < n; i++ {
	
		nums[i] = rand.Intn(1000)
	}
	return nums
}

func printChunks(chunks [][]int) {
	for i, c := range chunks {
		fmt.Printf("Chunk %d: %v\n", i, c)
	}
}
