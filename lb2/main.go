package main

import (
	"fmt"
	"shared/interfaces"
	desmos_constructor "shared/models/Desmos"
	sq "shared/models/Sequence"
	sd "shared/models/StatisticalDistribution"
)

const (
	Low  = 1
	High = 12
	N    = 12
)

func int64_to_float32(rx []int) []float32 {
	var tx = make([]float32, len(rx))

	for i, v := range rx {
		tx[i] = float32(v)
	}

	return tx
}

func main() {
	sequence := sq.Random(N, Low, High)
	var castedSequence interfaces.ISequence = sequence
	distr := sd.Complete(castedSequence, High)
	float_desmos := desmos_constructor.MakeDesmos[float32]()

	var mode = distr.GetVariantsMode()
	var median = sequence.GetVariationsMedian()
	var avg = sequence.GetAverage()

	fmt.Println("Ordered sequence:")
	fmt.Print(sequence.Variations, "\n\n")

	fmt.Println("Variants/Occurences")
	fmt.Print(float_desmos.PlotPoints(int64_to_float32(distr.Variants), int64_to_float32(distr.Occurences)), "\n\n")

	fmt.Println("High peak variants mode:")
	fmt.Print(float_desmos.PlotPoints(mode.GetVariants(mode.High), int64_to_float32(mode.GetOccurences(mode.High))), "\n\n")

	fmt.Println("Low peak variants mode:")
	fmt.Print(float_desmos.PlotPoints(mode.GetVariants(mode.Low), int64_to_float32(mode.GetOccurences(mode.Low))), "\n\n")

	fmt.Println("Median:")
	fmt.Print(median, "\n\n")

	fmt.Println("Avarage:")
	fmt.Print(avg, "\n\n")
}
