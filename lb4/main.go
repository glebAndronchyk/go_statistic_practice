package main

import (
	"fmt"
	math "math"
	rand "math/rand"
	"sort"
	"strings"
)

type Interval struct {
	Lower     float64
	Upper     float64
	Density   float64
	Frequency float64
}

const (
	N = 20
	C = 20
	L = 2
	M = 0
	S = 1
)

func expDistr(lambda float64, z float64) float64 {
	return -1 / lambda * math.Log(z)
}

func normDistr(m float64, s float64) float64 {
	var sum = 0.0

	for i := 0; i < 12; i++ {
		sum += rand.Float64()
	}

	return m + s*(sum-6.0)
}

func getDensityIntervals(data []float64) []Interval {
	var sorted = make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)

	var min = sorted[0]
	var max = sorted[len(sorted)-1]
	var n = len(sorted)

	intervalsAmount := int(math.Ceil(1 + 3.322*math.Log10(float64(n))))
	intervalWidth := (max - min) / float64(intervalsAmount)
	var intervals = make([]Interval, intervalsAmount)

	for i := 0; i < intervalsAmount; i++ {
		intervals[i].Lower = min + float64(i)*intervalWidth
		intervals[i].Upper = min + float64(i+1)*intervalWidth
		intervals[i].Density = 0.0
		intervals[i].Frequency = 0.0
	}

	for _, value := range sorted {
		var intervalIndex = int((value - min) / intervalWidth)
		if intervalIndex >= intervalsAmount {
			intervalIndex = intervalsAmount - 1
		}

		intervals[intervalIndex].Frequency++
	}

	for i := range intervals {
		intervals[i].Density = float64(intervals[i].Frequency / intervalWidth)
	}

	return intervals
}

func drawASCIIHistogram(intervals []Interval) {
	maxFreq := 0.0
	for _, iv := range intervals {
		if iv.Frequency > maxFreq {
			maxFreq = iv.Frequency
		}
	}

	fmt.Println(strings.Repeat("=", 70))

	for _, iv := range intervals {
		barWidth := int(50 * float64(iv.Frequency) / float64(maxFreq))
		bar := strings.Repeat("█", barWidth)

		fmt.Printf("[%7.3f, %7.3f) |%s %.2f (%.2f)\n",
			iv.Lower, iv.Upper, bar, iv.Frequency, iv.Density)
	}

	fmt.Println(strings.Repeat("=", 70))
}

func main() {
	var count = N * C
	var expSlice = []float64{}
	var normSlice = []float64{}

	for i := 0; i < count; i++ {
		expSlice = append(expSlice, expDistr(float64(L), rand.Float64()))
		normSlice = append(normSlice, normDistr(float64(M), float64(S)))
	}

	var expIntervals = getDensityIntervals(expSlice)
	var normIntervals = getDensityIntervals(normSlice)

	fmt.Println("Exponential distribution:")
	drawASCIIHistogram(expIntervals);
	fmt.Println("Normal distribution:")
	drawASCIIHistogram(normIntervals);
}
