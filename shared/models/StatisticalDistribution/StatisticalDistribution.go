package shared

import (
	"shared/interfaces"
	mode_pkg "shared/models/Mode"
)

type StatisticalDistribution struct {
	// xj
	variations []int
	// j
	high int

	// xi
	Variants []int
	// mi
	Occurences []int
	// mimax
	IntegralFrequencies []int

	// fi
	RelativeFrequencies []float32
	// Fj
	RelativeIntegralFrequencies []float32
}

func Complete(sequence interfaces.ISequence, high int) StatisticalDistribution {
	sd := StatisticalDistribution{variations: sequence.GetVariations(), high: high}

	sd.CalculateVariantsAndOccurences()
	sd.CalculateRelativeFrequencies()
	sd.CalculateRelativeIntegralFrequencies()

	return sd
}

// конструктор структури
func Incomplete(variations []int, high int) StatisticalDistribution {
	return StatisticalDistribution{variations: variations, high: high}
}

func (s StatisticalDistribution) GetVariants() []int {
	return s.Variants
}

func (s StatisticalDistribution) GetOccurences() []int {
	return s.Occurences
}

// знаходження варіацій, та абсолютних диференціальних та інтегральних частот
func (s *StatisticalDistribution) CalculateVariantsAndOccurences() {
	var variants = make([]int, s.high)
	var occurences = make([]int, s.high)
	var frequencies = make([]int, s.high)

	for i := 0; i < len(s.variations); i++ {
		var x int = s.variations[i]
		variants[x-1] = x
		occurences[x-1]++
	}

	for i := 0; i < len(occurences); i++ {
		for j := 0; j <= i; j++ {
			frequencies[i] += occurences[j]
		}
	}

	s.Variants = variants
	s.Occurences = occurences
	s.IntegralFrequencies = frequencies
}

// знаходження відносних диференціальних частот
func (s *StatisticalDistribution) CalculateRelativeFrequencies() {
	var relDiffOccr = make([]float32, len(s.Variants))
	var n = len(s.variations)

	for i := 0; i < len(s.Occurences); i++ {
		var fi = float32(s.Occurences[i]) / float32(n)
		relDiffOccr[i] = fi
	}

	s.RelativeFrequencies = relDiffOccr
}

// знаходження відносних інтегральних частот
func (s *StatisticalDistribution) CalculateRelativeIntegralFrequencies() {
	var relIntegFrq = make([]float32, len(s.Variants))

	for i := 0; i < len(s.RelativeFrequencies); i++ {
		for j := 0; j <= i; j++ {
			relIntegFrq[i] += s.RelativeFrequencies[j]
		}
	}

	s.RelativeIntegralFrequencies = relIntegFrq
}

func (s *StatisticalDistribution) GetVariantsMode() struct {
	High          []mode_pkg.Mode
	Low           []mode_pkg.Mode
	GetVariants   func(mode []mode_pkg.Mode) []float32
	GetOccurences func(mode []mode_pkg.Mode) []int
} {
	var variants = s.Variants
	var occurences = s.Occurences

	var high = []mode_pkg.Mode{}
	var low = []mode_pkg.Mode{}

	// checks only first item in the high peaks slice. Its enough to check only first item because we will have same frequencies for each peak here
	compareWithHighPeaks := func(occurences int) bool {
		if len(high) == 0 || high[0].Occurences <= occurences {
			return true
		}
		return false
	}

	// when we receive new highest frequency, we should move current frequencies to local(eg low) peaks and set pnly current variant as highest
	invalidateHighPeaksWithMode := func(newPeak mode_pkg.Mode) {
		low = append(low, high...)
		high = []mode_pkg.Mode{newPeak}
	}

	// adds new high peak and invalidates them when needed
	addHighPeak := func(val float32, occurences int) mode_pkg.Mode {
		mode := mode_pkg.NewMode(val, occurences)

		if len(high) == 0 || high[0].Occurences == occurences {
			high = append(high, mode)
		} else {
			invalidateHighPeaksWithMode(mode)
		}

		return mode
	}

	// adds low peak
	addLowPeak := func(val float32, occurences int) mode_pkg.Mode {
		mode := mode_pkg.NewMode(val, occurences)
		low = append(low, mode)

		return mode
	}

	for i := 1; i < len(variants); i++ {
		// left variant in the beggining wont exist so treat it as -1
		var left = -1
		if i-1 >= 0 {
			left = variants[i-1]
		}

		var mid = variants[i]

		// right variant in the end wont exist so treat it as -1
		var right = -1
		if i+1 < len(variants) {
			right = variants[i+1]
		}

		var occurences_left = 0
		if left < 0 {
			occurences_left = -1
		} else if left != 0 {
			occurences_left = occurences[left-1]
		}

		var occurences_mid = 0
		if mid != 0 {
			occurences_mid = occurences[mid-1]
		}

		var occurences_right = 0
		if right < 0 {
			occurences_right = -1
			// 0 means -- entry doesnt exist in the sequence
		} else if right != 0 {
			occurences_right = occurences[right-1]
		}

		// only this case should be covered. others might be ignored because they usually mean plateau from both sides or recession
		var isLocalPeak = occurences_mid > occurences_left && occurences_mid > occurences_right
		var isPlateauFromTheLeft = occurences_mid == occurences_left && occurences_mid > occurences_right
		var isPlateauFromTheRight = occurences_mid == occurences_right && occurences_mid > occurences_left

		if isLocalPeak {
			var isHighPeak = compareWithHighPeaks(occurences_mid)
			if isHighPeak {
				addHighPeak(float32(mid), occurences_mid)
			} else {
				addLowPeak(float32(mid), occurences_mid)
			}
		}

		if isPlateauFromTheLeft {
			var isHighPeak = compareWithHighPeaks(occurences_mid)
			// find avg between mid and left value
			var avg = float32(mid+left) / 2
			if isHighPeak {
				addHighPeak(avg, occurences_mid)
			} else {
				addLowPeak(avg, occurences_mid)
			}
		}

		if isPlateauFromTheRight {
			var isHighPeak = compareWithHighPeaks(occurences_mid)
			// find avg between mid and right value
			var avg = float32(mid+right) / 2
			if isHighPeak {
				addHighPeak(avg, occurences_mid)
			} else {
				addLowPeak(avg, occurences_mid)
			}
		}
	}

	return struct {
		High          []mode_pkg.Mode
		Low           []mode_pkg.Mode
		GetVariants   func(mode []mode_pkg.Mode) []float32
		GetOccurences func(mode []mode_pkg.Mode) []int
	}{
		High: high,
		Low:  low,
		GetVariants: func(mode []mode_pkg.Mode) []float32 {
			var variants = make([]float32, len(mode))

			for i, val := range mode {
				variants[i] = val.Value
			}

			return variants
		},
		GetOccurences: func(mode []mode_pkg.Mode) []int {
			var occurences = make([]int, len(mode))

			for i, val := range mode {
				occurences[i] = val.Occurences
			}

			return occurences
		},
	}
}
