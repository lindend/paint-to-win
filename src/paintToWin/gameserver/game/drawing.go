package game

type Point struct {
	X float32
	Y float32
}

type Color struct {
	R int
	G int
	B int
}

type Stroke struct {
	Start     Point
	End       Point
	Color     Color
	BrushSize int
	Pressure  float32
}

type Drawing struct {
	Strokes []Stroke
}

func NewDrawing() Drawing {
	return Drawing{[]Stroke{}}
}

func NewStroke(start Point, end Point, color Color, brushSize int, pressure float32) Stroke {
	return Stroke{start, end, color, brushSize, pressure}
}
