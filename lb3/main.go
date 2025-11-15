package main

import (
	"log"
	LinearEquation "shared/models/LinearEquation"
	NormNoise "shared/models/NormNoise"
	shared "shared/models/Point"
	Regression "shared/models/Regression"

	color "image/color"

	Plot "gonum.org/v1/plot"
	Plotter "gonum.org/v1/plot/plotter"
	vg "gonum.org/v1/plot/vg"
)

const (
	N    = 2
	k    = float32(2.45)
	b    = float32(3.11)
	seed = 1100
)

func CombinePoints(x, y []float32) []shared.Point {
	var result = []shared.Point{}

	for i := 0; i < len(x); i++ {
		result = append(result, shared.Point{X: x[i], Y: y[i]})
	}

	return result
}

func BuildPlot(coords []shared.Point, fns struct {
	Y_X func(x float32) float32
	X_Y func(y float32) float32
}, step float32, top int) {
	var actual Plotter.XYs
	var y_x_plot Plotter.XYs
	var x_y_plot Plotter.XYs
	var plot = Plot.New()

	for i := 0; i < len(coords); i++ {
		var coord = coords[i]
		actual = append(actual, Plotter.XY{X: float64(coord.X), Y: float64(coord.Y)})
	}

	for i_step := float32(0); i_step < float32(top); i_step += step {
		var i_y = fns.Y_X(i_step)
		var i_x = fns.X_Y(i_step)
		y_x_plot = append(y_x_plot, Plotter.XY{X: float64(i_step), Y: float64(i_y)})
		x_y_plot = append(x_y_plot, Plotter.XY{X: float64(i_x), Y: float64(i_step)})
	}

	y_x_line, _ := Plotter.NewLine(y_x_plot)
	// x_y_line, _ := Plotter.NewLine(x_y_plot)
	actual_line, _ := Plotter.NewLine(actual)
	// y_x_scatter, _ := Plotter.NewScatter(y_x_plot)
	// x_y_scatter, _ := Plotter.NewScatter(x_y_plot)

	y_x_line.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	y_x_line.LineStyle.Width = vg.Points(1)

	// x_y_line.Color = color.RGBA{R: 0, G: 255, B: 255, A: 255}
	// x_y_line.LineStyle.Width = vg.Points(3)

	// y_x_scatter.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	// y_x_scatter.Radius = 6

	// x_y_scatter.Color = color.RGBA{R: 0, G: 255, B: 255, A: 255}
	// x_y_scatter.Radius = 6

	plot.Add(y_x_line)
	plot.Add(actual_line)
	// plot.Add(x_y_line)
	// plot.Add(y_x_scatter)
	// plot.Add(x_y_scatter)

	if err := plot.Save(4*vg.Inch, 4*vg.Inch, "scatter.png"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// task specific code. has nothing related to real life
	var deviation = float32(N + 10) / 5
	var slice_length int = N + 10
	eq := LinearEquation.Make(k, b)

	// x,y - data that we have normally in real world (eg temp-date relation). always treated as noise
	noise := NormNoise.Make(deviation, int64(seed))
	var x = noise.GenerateLinearSequence(slice_length)
	var noised_values = noise.GenerateSlice(slice_length)
	var y = eq.ForSlice(x).ApplyNoise(noised_values)
	var coords = CombinePoints(x, y)

	// regress coordinates
	var linearRegression = Regression.Make(shared.PointArrayToIPointArray(coords))
	var fns = linearRegression.CalculateRegresionEquations()

	BuildPlot(coords, fns, float32(0.1), slice_length+1)
}
