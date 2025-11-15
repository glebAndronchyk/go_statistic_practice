package shared

type LinearEquation struct {
	k float32
	b float32

	forSlice []float32
}

func Make(k, b float32) LinearEquation {
	return LinearEquation{k: k, b: b}
}

func (eq LinearEquation) GetSingle(x float32) float32 {
	return eq.k*x + eq.b
}

func (eg LinearEquation) ForSlice(x []float32) LinearEquation {
	var result = []float32{}

	for i := 0; i < len(x); i++ {
		result = append(result, eg.GetSingle(x[i]))
	}

	eg.forSlice = result

	return eg
}

func (eg LinearEquation) ApplyNoise(noise []float32) []float32 {

	for i := 0; i < len(eg.forSlice); i++ {
		eg.forSlice[i] = eg.forSlice[i] + noise[i]
	}

	return eg.forSlice
}
