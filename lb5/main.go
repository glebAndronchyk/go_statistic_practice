package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// GoldsteinCoefficients stores the coefficients for Goldstein approximation
type GoldsteinCoefficients struct {
	a, b, c float64
}

var goldsteinCoeffs = []GoldsteinCoefficients{
	{1.0000886, -0.2237368, -0.01513904},
	{0.4713941, 0.02607083, -0.008986007},
	{0.0001348028, 0.01128186, 0.02277679},
	{-0.008553069, -0.01153761, -0.01323293},
	{0.00312558, 0.00516965, -0.006950356},
	{-0.0008426812, 0.00253001, 0.001060438},
	{0.0000978049, -0.00145011, 0.001565326},
}

func main() {
	sample, alpha := readInput()

	fmt.Printf("\nSample size: %d\n", len(sample))
	fmt.Printf("Significance level α: %.3f\n\n", alpha)

	histogram := createHistogram(sample, 0, 20)

	fmt.Println("Histogram:")
	displayHistogram(histogram)
	fmt.Println()

	fmt.Println("=== Testing for Normal Distribution ===")
	testNormalDistribution(histogram, alpha, len(sample))
	fmt.Println()

	fmt.Println("=== Testing for Uniform Distribution ===")
	testUniformDistribution(histogram, alpha, len(sample))
}

func readInput() ([]int, float64) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter integer numbers in range [0, 20] separated by spaces:")
	scanner.Scan()
	input := scanner.Text()

	parts := strings.Fields(input)
	sample := make([]int, 0)

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		if num < 0 || num > 20 {
			continue
		}
		sample = append(sample, num)
	}

	if len(sample) == 0 {
		os.Exit(1)
	}

	var alpha float64
	for {
		fmt.Print("Enter significance level α (e.g., 0.05): ")
		scanner.Scan()
		alphaStr := scanner.Text()

		var err error
		alpha, err = strconv.ParseFloat(alphaStr, 64)
		if err != nil || alpha <= 0 || alpha >= 1 {
			fmt.Println("Error: α must be a number between 0 and 1")
			continue
		}
		break
	}

	return sample, alpha
}

type Histogram struct {
	Bins       []Bin
	TotalCount int
}

type Bin struct {
	Lower int
	Upper int
	Count int
}

func createHistogram(sample []int, min, max int) *Histogram {
	freq := make(map[int]int)
	for _, val := range sample {
		freq[val]++
	}

	bins := make([]Bin, 0)
	for i := min; i <= max; i++ {
		if freq[i] > 0 {
			bins = append(bins, Bin{
				Lower: i,
				Upper: i,
				Count: freq[i],
			})
		}
	}

	bins = mergeBins(bins)

	return &Histogram{
		Bins:       bins,
		TotalCount: len(sample),
	}
}

func mergeBins(bins []Bin) []Bin {
	if len(bins) == 0 {
		return bins
	}

	merged := make([]Bin, 0)
	current := bins[0]

	for i := 1; i < len(bins); i++ {
		if current.Count < 5 {
			current.Upper = bins[i].Upper
			current.Count += bins[i].Count
		} else {
			merged = append(merged, current)
			current = bins[i]
		}
	}

	if len(merged) > 0 && current.Count < 5 {
		merged[len(merged)-1].Upper = current.Upper
		merged[len(merged)-1].Count += current.Count
	} else {
		merged = append(merged, current)
	}

	return merged
}

func displayHistogram(hist *Histogram) {
	fmt.Println("Bin Range\t\tCount")
	fmt.Println("─────────────────────────────")
	for _, bin := range hist.Bins {
		if bin.Lower == bin.Upper {
			fmt.Printf("[%d]\t\t\t%d\n", bin.Lower, bin.Count)
		} else {
			fmt.Printf("[%d - %d]\t\t%d\n", bin.Lower, bin.Upper, bin.Count)
		}
	}
}

func calculateChiSquaredCritical(alpha float64, r int) float64 {
	n := float64(r)

	var d float64
	if alpha >= 0.5 && alpha <= 0.999 {
		d = 2.0637 * math.Pow(math.Log(1.0/(1.0-alpha))-0.16, 0.4274) - 1.5774
	} else if alpha >= 0.001 && alpha <= 0.5 {
		d = -2.0637 * math.Pow(math.Log(1.0/alpha)-0.16, 0.4274) + 1.5774
	}

	sum := 0.0
	for i := 0; i <= 6; i++ {
		coeff := goldsteinCoeffs[i]
		power := math.Pow(n, -float64(i)/2.0)
		dPower := math.Pow(d, float64(i))
		term := coeff.a + coeff.b/n + coeff.c/(n*n)
		sum += power * dPower * term
	}

	chiSquared := n * math.Pow(sum, 3)
	return chiSquared
}

func testNormalDistribution(hist *Histogram, alpha float64, n int) {
	mean := 0.0
	variance := 0.0

	values := make([]float64, 0, n)
	for _, bin := range hist.Bins {
		midpoint := float64(bin.Lower+bin.Upper) / 2.0
		for j := 0; j < bin.Count; j++ {
			values = append(values, midpoint)
		}
	}

	for _, v := range values {
		mean += v
	}
	mean /= float64(n)

	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(n - 1)
	stdDev := math.Sqrt(variance)

	fmt.Printf("Estimated mean: %.4f\n", mean)
	fmt.Printf("Estimated standard deviation: %.4f\n", stdDev)
	fmt.Println()

	theoreticalFreq := make([]float64, len(hist.Bins))
	for i, bin := range hist.Bins {
		lower := float64(bin.Lower) - 0.5
		upper := float64(bin.Upper) + 0.5

		prob := normalCDF((upper-mean)/stdDev) - normalCDF((lower-mean)/stdDev)
		theoreticalFreq[i] = prob * float64(n)
	}

	chiSquared := 0.0
	fmt.Println("Bin\t\tObserved\tExpected\tContribution")
	fmt.Println("─────────────────────────────────────────────────────")
	for i, bin := range hist.Bins {
		observed := float64(bin.Count)
		expected := theoreticalFreq[i]
		contribution := 0.0
		if expected > 0 {
			contribution = math.Pow(observed-expected, 2) / expected
		}
		chiSquared += contribution

		binStr := fmt.Sprintf("[%d-%d]", bin.Lower, bin.Upper)
		if bin.Lower == bin.Upper {
			binStr = fmt.Sprintf("[%d]", bin.Lower)
		}
		fmt.Printf("%-12s\t%d\t\t%.2f\t\t%.4f\n", binStr, bin.Count, expected, contribution)
	}

	degreesOfFreedom := len(hist.Bins) - 2 - 1

	if degreesOfFreedom <= 0 {
		fmt.Println("\nError: Not enough bins for valid test (degrees of freedom <= 0)")
		return
	}

	chiSquaredCritical := calculateChiSquaredCritical(1-alpha, degreesOfFreedom)

	fmt.Println()
	fmt.Printf("χ²_empirical = %.4f\n", chiSquared)
	fmt.Printf("χ²_critical(α=%.3f, r=%d) = %.4f\n", alpha, degreesOfFreedom, chiSquaredCritical)
	fmt.Println()

	if chiSquared < chiSquaredCritical {
		fmt.Printf("χ²_emp < χ²_crit → Hypothesis ACCEPTED\n")
		fmt.Printf("Data follows Normal distribution at significance level α=%.3f\n", alpha)
	} else {
		fmt.Printf("χ²_emp ≥ χ²_crit → Hypothesis REJECTED\n")
		fmt.Printf("Data does NOT follow Normal distribution at significance level α=%.3f\n", alpha)
	}
}

func testUniformDistribution(hist *Histogram, alpha float64, n int) {
	minVal := hist.Bins[0].Lower
	maxVal := hist.Bins[len(hist.Bins)-1].Upper
	rangeSize := float64(maxVal - minVal + 1)

	theoreticalFreq := make([]float64, len(hist.Bins))
	for i, bin := range hist.Bins {
		binSize := float64(bin.Upper - bin.Lower + 1)
		prob := binSize / rangeSize
		theoreticalFreq[i] = prob * float64(n)
	}

	chiSquared := 0.0
	fmt.Println("Bin\t\tObserved\tExpected\tContribution")
	fmt.Println("─────────────────────────────────────────────────────")
	for i, bin := range hist.Bins {
		observed := float64(bin.Count)
		expected := theoreticalFreq[i]
		contribution := 0.0
		if expected > 0 {
			contribution = math.Pow(observed-expected, 2) / expected
		}
		chiSquared += contribution

		binStr := fmt.Sprintf("[%d-%d]", bin.Lower, bin.Upper)
		if bin.Lower == bin.Upper {
			binStr = fmt.Sprintf("[%d]", bin.Lower)
		}
		fmt.Printf("%-12s\t%d\t\t%.2f\t\t%.4f\n", binStr, bin.Count, expected, contribution)
	}

	degreesOfFreedom := len(hist.Bins) - 0 - 1

	if degreesOfFreedom <= 0 {
		fmt.Println("\nError: Not enough bins for valid test (degrees of freedom <= 0)")
		return
	}

	chiSquaredCritical := calculateChiSquaredCritical(1-alpha, degreesOfFreedom)

	fmt.Println()
	fmt.Printf("χ²_empirical = %.4f\n", chiSquared)
	fmt.Printf("χ²_critical(α=%.3f, r=%d) = %.4f\n", alpha, degreesOfFreedom, chiSquaredCritical)
	fmt.Println()

	if chiSquared < chiSquaredCritical {
		fmt.Printf("χ²_emp < χ²_crit → Hypothesis ACCEPTED\n")
		fmt.Printf("Data follows Uniform distribution at significance level α=%.3f\n", alpha)
	} else {
		fmt.Printf("χ²_emp ≥ χ²_crit → Hypothesis REJECTED\n")
		fmt.Printf("Data does NOT follow Uniform distribution at significance level α=%.3f\n", alpha)
	}
}

func normalCDF(x float64) float64 {
	return 0.5 * (1.0 + math.Erf(x/math.Sqrt2))
}
