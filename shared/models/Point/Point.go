package shared

import "shared/interfaces"

type Point struct {
	X float32
	Y float32
}

func PointArrayToIPointArray(points []Point) []interfaces.IPoint {
	ipoints := make([]interfaces.IPoint, len(points))
	for i, p := range points {
		ipoints[i] = p
	}
	return ipoints
}

func (p Point) GetX() float32 {
	return p.X
}

func (p Point) GetY() float32 {
	return p.Y
}

func (p Point) ToIPoint() interfaces.IPoint {
	return interfaces.IPoint(p)
}
