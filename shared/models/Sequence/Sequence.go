package shared

import (
	// "fmt"
	"math/rand"
	// "shared/interfaces"
	// mode_pkg "shared/models/Mode"
	"sort"
)

type Sequence struct {
	Source            []int
	Variations        []int
	Source_length     int
	Variations_length int
}

// / reducer
type reducerFn[T any, U any] func(U, T, int) U

func Reduce[T any, U any](slice []T, reducer reducerFn[T, U], initialValue U) U {
	var acc U = initialValue

	for index, item := range slice {
		acc = reducer(acc, item, index)
	}

	return acc
}

////

func Random(n int, low int, high int) Sequence {
	var sequence []int = []int{}
	var variations []int = []int{}

	for i := 0; i < n; i++ {
		var rndVal = rand.Intn(high) + low
		sequence = append(sequence, rndVal)
	}

	variations = make([]int, len(sequence))
	copy(variations, sequence)

	sort.Slice(variations, func(i, j int) bool {
		return variations[i] < variations[j]
	})

	return Sequence{Source: sequence, Variations: variations, Source_length: len(sequence), Variations_length: len(variations)}
}

func (s Sequence) GetVariations() []int {
	return s.Variations
}

func (s Sequence) GetVariationsMedian() []int {
	var midIdx = s.Variations_length / 2
	var median = []int{}

	median = append(median, s.Variations[midIdx])
	if s.Variations_length%2 == 0 {
		median = append(median, s.Variations[midIdx-1])
	}

	return median
}

func (s Sequence) GetAverage() float32 {
	var total = Reduce(s.Source, func(acc int, curr int, _ int) int {
		return acc + curr
	}, 0)

	return float32(total) / float32(s.Source_length)
}
