package main

import (
	"fmt"
	"log"

	"github.com/IgaguriMK/planetStat/cache"
	"github.com/IgaguriMK/planetStat/cube"
	"github.com/IgaguriMK/planetStat/vec"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cc, err := cache.NewController(".cache")
	if err != nil {
		log.Fatal(err)
	}

	cube := cube.FromCenter(vec.Vec3{0, 0, 0}, vec.One.Scalar(20))

	systems, err := cube.GetSystems(cc)
	if err != nil {
		log.Fatal("API load error:", err)
	}

	for _, s := range systems {
		info, err := s.GetSystemInfo(cc)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("(%d) %s, stars:%d, bodys:%d\n", s.ID, s.Name, info.StarCount(), info.PlanetCount())
	}
}
