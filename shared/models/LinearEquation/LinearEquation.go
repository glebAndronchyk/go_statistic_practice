package shared

type LinearEquation struct {
	k float32
	b float32
}

func Make(k, b float32) LinearEquation {
	return LinearEquation{k: k, b: b}
}

func (eq LinearEquation) GetSingle(x float32) float32 {
	return eq.k*x + eq.b
}

func (eg LinearEquation) GetSlice(x []float32) []float32 {
	var result = []float32{}

	for i := 0; i < len(x); i++ {
		result = append(result, eg.GetSingle(x[i]))
	}

	return result
}
