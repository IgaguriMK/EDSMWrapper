package cube

import "math"

const ChunkSize = 10
const MaxSize = 200

type Cube struct {
	Pos  vec.Vec3
	Size vec.Vec3
}

func FromCenter(center, size vec.Vec3) Cube {
	hs = size.Scalar(0.5)
	return Cube{center.Sub(hs), center.Add(hs)}
}

func (c Cube) WrapBoundary() Cube {
	p = c.Pos
	q = c.Pos.Add(c.Size)

	p.X, q.X = allign(p.X, q.X)
	p.Y, q.Y = allign(p.Y, q.Y)
	p.Z, q.Z = allign(p.Z, q.Z)

	return Cube{p, q}
}

func allign(l, r float64) (float64, float64) {
	return ChunkSize * math.Floor(l/ChunkSize), ChunkSize * math.Ceil(r/ChunkSize)
}
