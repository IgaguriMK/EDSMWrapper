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

	cube := cube.FromCenter(vec.Zero, vec.One.Scalar(10))

	systems, err := cube.GetSystems(cc)
	if err != nil {
		log.Fatal("API load error:", err)
	}

	for _, s := range systems {
		fmt.Printf("%s\n", s.Name)
	}
}
