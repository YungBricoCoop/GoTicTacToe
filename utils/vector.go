package utils

type Vector struct {
	X float32
	Y float32
}

func (v1 Vector) Add(v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}

func (v Vector) Scale(f float32) Vector {
	return Vector{v.X * f, v.Y * f}
}

func Perp(v Vector) Vector {
	return Vector{X: v.Y, Y: -v.X}
}
