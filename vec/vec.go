package Vec

import "math"

type Vec3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

const (
	Zero = Vec3{0, 0, 0}
	X    = Vec3{1, 0, 0}
	Y    = Vec3{0, 1, 0}
	Z    = Vec3{0, 0, 0}
)

func (v Vec3) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{
		X: v.X + u.X,
		Y: v.Y + u.Y,
		Z: v.Z + u.Z,
	}
}

func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{
		X: v.X - u.X,
		Y: v.Y - u.Y,
		Z: v.Z - u.Z,
	}
}

func (v Vec3) Scalar(k float64) Vec3 {
	return Vec3{
		X: k * v.X,
		Y: k * v.Y,
		Z: k * v.Z,
	}
}

func (v Vec3) Unit() Vec3 {
	if v.X == 0 && v.Y == 0 && v.Z == 0 {
		return Zero
	}

	a := 1 / v.Abs()

	return Vec3{
		X: a * v.X,
		Y: a * v.Y,
		Z: a * v.Z,
	}
}

func (v Vec3) Dot(u Vec3) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}
