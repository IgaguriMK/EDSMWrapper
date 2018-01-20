package main

import (
	"fmt"
	"log"

	"github.com/IgaguriMK/planetStat/cache"
	"github.com/IgaguriMK/planetStat/cube"
	"github.com/IgaguriMK/planetStat/system"
	"github.com/IgaguriMK/planetStat/vec"
)

const EmptyResultLimit = 20

func main() {
	log.SetFlags(log.Lshortfile)

	cc, err := cache.NewController(".cache")
	if err != nil {
		log.Fatal(err)
	}

	cube := cube.FromCenter(vec.Vec3{25, -20, 25899}, vec.One.Scalar(200))

	systems, err := cube.GetSystems(cc)
	if err != nil {
		log.Fatal("API load error:", err)
	}

	emptyCount := 0
	fmt.Println("StarType\tStarTemp\tDistance")
	for _, s := range systems {
		info, err := s.GetSystemInfo(cc)
		if err != nil {
			if err == system.ErrNotFound {
				emptyCount++
				log.Println("Empty result", emptyCount)

				if emptyCount >= EmptyResultLimit {
					locked, err := system.CheckAPILocked()
					if err != nil {
						log.Fatal("API call error", err)
					}
					if locked {
						log.Fatal("API call")
					} else {
						emptyCount = 0
					}
				}
				continue
			}
			log.Fatal("API call error", err)
		}
		emptyCount = 0

		if info.StarCount() > 1 {
			continue
		}

		star := info.Stars()[0]
		starType := system.ShortType(star.SubType)
		for _, b := range info.Planets() {
			terraformable := "F"
			if b.TerraformingState == "Candidate for terraforming" {
				terraformable = "T"
			}
			fmt.Printf("%s\t%f\t%f\t%s\n", starType, star.SurfaceTemperature, b.DistanceToArrival, terraformable)
		}
	}
}
