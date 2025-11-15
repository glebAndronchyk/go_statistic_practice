package shared

import (
	"math/rand"
	"math"
)

type NormNoise struct {
	deviation float32
	rnd       *rand.Rand
}

func Make(deviation float32, seed int64) NormNoise {
	var rnd = rand.New(rand.NewSource(seed))
	return NormNoise{deviation: deviation, rnd: rnd}
}

func (noise *NormNoise) GenerateSingle() float32 {
	return float32(math.Abs(noise.rnd.NormFloat64()))*noise.deviation
}

func (noise *NormNoise) GenerateSlice(end int) []float32 {
	var result = []float32{}

	for i := 0; i < end; i++ {
		result = append(result, noise.GenerateSingle())
	}

	return result
}

func (noise *NormNoise) GenerateLinearSequence(end int) []float32 {
	var result = []float32{}

	for i := 0; i < end; i++ {
		result = append(result, float32(i))
	}

	return result
}


