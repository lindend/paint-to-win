package game

type Point struct {
	X float32
	Y float32
}

type Color struct {
	Red   int
	Green int
	Blue  int
}

type Stroke struct {
	Start     Point
	End       Point
	Color     Color
	BrushSize int
}

type Drawing struct {
	Strokes []Stroke
}

func NewDrawing() Drawing {
	return Drawing{[]Stroke{}}
}

func NewStroke(start Point, end Point, color Color, brushSize int) Stroke {
	return Stroke{start, end, color, brushSize}
}
