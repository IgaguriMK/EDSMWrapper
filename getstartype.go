package main

import (
	"fmt"
	"log"

	"github.com/IgaguriMK/planetStat/cache"
	"github.com/IgaguriMK/planetStat/cube"
	"github.com/IgaguriMK/planetStat/vec"
)

const EmptyResultLimit = 20

func main() {
	log.SetFlags(log.Lshortfile)

	cc, err := cache.NewController(".cache")
	if err != nil {
		log.Fatal(err)
	}

	cube := cube.FromCenter(vec.Vec3{25, -20, 25899}, vec.One.Scalar(1000))

	systems, err := cube.GetSystems(cc)
	if err != nil {
		log.Fatal("API load error:", err)
	}

	for _, s := range systems {
		t := s.PrimaryStar.Type
		if t != "" {
			fmt.Println(t)
		}
	}
}
