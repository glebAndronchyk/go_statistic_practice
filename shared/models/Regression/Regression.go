package shared

import (
	"math"
	"shared/interfaces"
)

type Regression struct {
	points   []interfaces.IPoint
	c_x      float32
	c_y      float32
	c_xy     float32
	c_xx     float32
	c_yy     float32
	s_x      float32
	s_y      float32
	cor_coef float32
	// y->x equation coefs(ax+b)
	a float32
	b float32
	// x->y equation coefs(cx+d)
	c float32
	d float32
}

func Make(coords []interfaces.IPoint) Regression {
	return Regression{points: coords}
}

func (r *Regression) CalculateRegresionEquations() struct {
	Y_X func(x float32) float32
	X_Y func(y float32) float32
} {
	r.estimateCoefficients()
	r.calculateQuadraticDeviation()
	r.calculateCorrelationCoeficient()
	r.calculateLinearEquationCoeficients()

	return struct {
		Y_X func(x float32) float32
		X_Y func(y float32) float32
	}{
		Y_X: func(x float32) float32 {
			return float32(r.a*x + r.b)
		},
		X_Y: func(y float32) float32 {
			return float32(r.c*y + r.d)
		},
	}
}

func (r *Regression) estimateCoefficients() {
	var c_x = float32(0)
	var c_y = float32(0)
	var c_xy = float32(0)
	var c_xx = float32(0)
	var c_yy = float32(0)

	for i := 0; i < len(r.points); i++ {
		var point = r.points[i]

		c_x += point.GetX()
		c_y += point.GetY()
		c_xy += point.GetX() * point.GetY()
		c_xx += point.GetX() * point.GetX()
		c_yy += point.GetY() * point.GetY()
	}

	r.c_x = c_x / float32(len(r.points))
	r.c_y = c_y / float32(len(r.points))
	r.c_xy = c_xy / float32(len(r.points))
	r.c_xx = c_xx / float32(len(r.points))
	r.c_yy = c_yy / float32(len(r.points))
}

func (r *Regression) calculateQuadraticDeviation() {
	r.s_x = float32(math.Sqrt(float64(r.c_xx) - float64(r.c_x)*float64(r.c_x)))
	r.s_y = float32(math.Sqrt(float64(r.c_yy) - float64(r.c_y)*float64(r.c_y)))
}

func (r *Regression) calculateCorrelationCoeficient() {
	r.cor_coef = (r.c_xy - r.c_x*r.c_y) / (r.s_x * r.s_y)
}

func (r *Regression) calculateLinearEquationCoeficients() {
	// main regression coefs
	r.a = r.cor_coef * r.s_y / r.s_x
	r.b = r.c_y - r.a*r.c_x
	// "clarifying" regression coefs
	r.c = r.cor_coef * r.s_x / r.s_y
	r.d = r.c_x - r.c*r.c_y
}
