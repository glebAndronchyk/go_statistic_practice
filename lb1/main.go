package main

import (
	// "encoding/json"
	"fmt"
	"os"
	sq "shared/models/Sequence"
	sd "shared/models/StatisticalDistribution"
	"text/tabwriter"
)

const (
	Low  = 1
	High = 5
	N    = 12
)

func main() {
	sequence := sq.Random(N, Low, High)
	var sd = buildStatisticalDistribution(sequence.Variations)

	// sdJsonBytes, _ := json.MarshalIndent(statisticalDistribution, "", "  ") //json-like indentation

	fmt.Println(" \nStarting sequence:\n", sequence)
	fmt.Println()
	// fmt.Println(string(sdJsonBytes))

	const paddingAmount = 3
	const paddingChar = ' '
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, paddingAmount, paddingChar, tabwriter.Debug)

	fmt.Fprintln(writer, "  xi\t ni\t Частотність\t Інт. частота")
	for i := 0; i < len(sd.Variants); i++ {
		fmt.Fprintln(writer, " ", sd.Variants[i], "\t", sd.Occurences[i], "\t", sd.RelativeFrequencies[i], "\t", sd.IntegralFrequencies[i])
	}
	writer.Flush()
	fmt.Println()
}

// побудувати статистичний розподіл
func buildStatisticalDistribution(variationsSet []int) sd.StatisticalDistribution {
	s := sd.Incomplete(variationsSet, High)

	s.CalculateVariantsAndOccurences()
	s.CalculateRelativeFrequencies()
	s.CalculateRelativeIntegralFrequencies()

	return s
}
