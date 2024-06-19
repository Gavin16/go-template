package service

// Shape
// interface usage demo
type Shape interface {
	Area() float64
}

type Center interface {
	Position() (float64, float64)
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
	XPos   float64
	YPos   float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func (c Circle) Position() (float64, float64) {
	return c.XPos, c.YPos
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
