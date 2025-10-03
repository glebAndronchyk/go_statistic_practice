package shared

import (
	"fmt"
	"strings"
)

type Desmos[V int | float32] struct{}

func MakeDesmos[V int | float32]() Desmos[V] {
	return Desmos[V]{}
}

func (d Desmos[V]) PlotPoints(a1, a2 []V) string {
	// Handle edge cases
	if len(a1) == 0 || len(a2) == 0 {
		return ""
	}

	// Use the minimum length if arrays differ in size
	minLen := len(a1)
	if len(a2) < minLen {
		minLen = len(a2)
	}

	var points []string
	for i := 0; i < minLen; i++ {
		point := fmt.Sprintf("(%g,%g)", a1[i], a2[i])
		points = append(points, point)
	}

	return strings.Join(points, ",")
}
