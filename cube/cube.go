package cube

import (
	"fmt"
	"math"

	"github.com/IgaguriMK/planetStat/cache"
	"github.com/IgaguriMK/planetStat/system"
	vec "github.com/IgaguriMK/planetStat/vec"
)

const ChunkSize = 100
const MaxSize = 200
const CubeCacheVer = 2

func init() {
	cache.AddCacheType("chunk")
}

type Cube struct {
	Pos  vec.Vec3
	Size vec.Vec3
}

func FromCenter(center, size vec.Vec3) Cube {
	hs := size.Scalar(0.5)
	return Cube{center.Sub(hs), size}
}

func (cube Cube) GetSystems(cacheController *cache.CacheController) ([]system.System, error) {
	cps := cube.Chunks()

	systems := make([]system.System, 0)
	containedSystems := make(map[int32]bool)

	for _, cp := range cps {
		var ch Chunk
		ok := cacheController.Find(CubeCacheVer, cp.PosStr(), &ch)

		if !ok {
			c := cp.Center()
			ss, err := system.Get(c.X, c.Y, c.Z, ChunkSize)
			if err != nil {
				return nil, err
			}

			ssc := make([]system.System, 0, len(ss))
			for _, s := range ss {
				if cp.Contains(s.Coords) {
					ssc = append(ssc, s)
				}
			}

			ch = Chunk{
				Pos:     cp,
				Systems: ssc,
			}

			cacheController.Store(CubeCacheVer, ch)
		}

		for _, s := range ch.Systems {
			if cube.Contains(s.Coords) && !containedSystems[s.ID] {
				systems = append(systems, s)
				containedSystems[s.ID] = true
			}
		}
	}

	return systems, nil
}

func (c Cube) Chunks() []ChunkPos {
	cw := c.WrapBoundary()
	bot := PosChunk(cw.Pos)
	top := PosChunk(cw.Pos.Add(cw.Size))

	chunks := make([]ChunkPos, 0)
	for x := bot.X; x <= top.X; x++ {
		for y := bot.Y; y <= top.Y; y++ {
			for z := bot.Z; z <= top.Z; z++ {
				chunks = append(chunks, ChunkPos{x, y, z})
			}
		}
	}

	return chunks
}

func (c Cube) WrapBoundary() Cube {
	p := c.Pos
	q := c.Pos.Add(c.Size)

	p.X, q.X = allign(p.X, q.X)
	p.Y, q.Y = allign(p.Y, q.Y)
	p.Z, q.Z = allign(p.Z, q.Z)

	return Cube{p, q}
}

func allign(l, r float64) (float64, float64) {
	return ChunkSize * math.Floor(l/ChunkSize), ChunkSize * math.Ceil(r/ChunkSize)
}

func (c Cube) Contains(v vec.Vec3) bool {
	p := c.Pos
	q := c.Pos.Add(c.Size)

	if v.X < p.X || q.X < v.X {
		return false
	}
	if v.Y < p.Y || q.Y < v.Y {
		return false
	}
	if v.Z < p.Z || q.Z < v.Z {
		return false
	}

	return true
}

type Chunk struct {
	Pos     ChunkPos        `json:"pos"`
	Systems []system.System `json:"systems"`
}

func (c Chunk) Key() string {
	return c.Pos.PosStr()
}

type ChunkPos struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

func PosChunk(v vec.Vec3) ChunkPos {
	return ChunkPos{
		X: int(math.Floor(v.X / ChunkSize)),
		Y: int(math.Floor(v.Y / ChunkSize)),
		Z: int(math.Floor(v.Z / ChunkSize)),
	}
}

func (pc ChunkPos) Contains(v vec.Vec3) bool {
	vpc := PosChunk(v)
	return pc == vpc
}

func (cp ChunkPos) Center() vec.Vec3 {
	return vec.Vec3{
		X: float64(cp.X)*ChunkSize + 0.5*ChunkSize,
		Y: float64(cp.Y)*ChunkSize + 0.5*ChunkSize,
		Z: float64(cp.Z)*ChunkSize + 0.5*ChunkSize,
	}
}

func (cp ChunkPos) PosStr() string {
	xAbs, xSig := cp.X, "p"
	yAbs, ySig := cp.Y, "p"
	zAbs, zSig := cp.Z, "p"

	if xAbs < 0 {
		xAbs = -xAbs
		xSig = "n"
	}
	if yAbs < 0 {
		yAbs = -yAbs
		ySig = "n"
	}
	if zAbs < 0 {
		zAbs = -zAbs
		zSig = "n"
	}

	return fmt.Sprintf("chunk/%s%d%s%d%s%d", xSig, xAbs, ySig, yAbs, zSig, zAbs)
}
