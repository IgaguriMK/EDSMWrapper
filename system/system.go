package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/IgaguriMK/planetStat/cache"
	"github.com/IgaguriMK/planetStat/vec"
)

const SystemInfoCacheVer = 1

var (
	apiCallWait = time.Second
)

var (
	ErrNotFound = errors.New("API returns nil response")
)

func init() {
	cache.AddCacheType("systemInfo")
}

type System struct {
	Coords       vec.Vec3 `json:"coords"`
	CoordsLocked bool     `json:"coordsLocked"`
	Name         string   `json:"name"`
	ID           int32    `json:"id"`
	ID64         int64    `json:"id64"`
	PermitName   string   `json:"permitName"`
	PrimaryStar  struct {
		IsScoopable bool   `json:"isScoopable"`
		Name        string `json:"name"`
		Type        string `json:"type"`
	} `json:"primaryStar"`
	RequirePermit   bool `json:"requirePermit"`
	systemInfoCache *SystemInfo
}

func Get(x, y, z, size float64) ([]System, error) {
	params := url.Values{}
	params.Add("x", fmt.Sprint(x))
	params.Add("y", fmt.Sprint(y))
	params.Add("z", fmt.Sprint(z))
	params.Add("size", fmt.Sprint(size))
	params.Add("showId", "1")
	params.Add("showCoordinates", "1")
	params.Add("showPermit", "1")
	params.Add("showPrimaryStar", "1")

	url := "https://www.edsm.net/api-v1/cube-systems?" + params.Encode()
	log.Println(url)

	time.Sleep(apiCallWait)
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var v []System
	err = json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (sys System) GetSystemInfo(cc *cache.CacheController) (*SystemInfo, error) {
	if sys.systemInfoCache != nil {
		return sys.systemInfoCache, nil
	}

	var systemInfo SystemInfo
	if cc.Find(SystemInfoCacheVer, fmt.Sprintf("systemInfo/%d", sys.ID), &systemInfo) {
		return &systemInfo, nil
	}

	url := fmt.Sprintf("https://www.edsm.net/api-system-v1/bodies?systemId=%d", sys.ID)
	log.Println(url)

	time.Sleep(apiCallWait)
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if string(bytes) == "[]" {
		return nil, ErrNotFound
	}

	err = json.Unmarshal(bytes, &systemInfo)
	if err != nil {
		return nil, err
	}

	cc.Store(SystemInfoCacheVer, systemInfo)

	return &systemInfo, nil
}

type SystemInfo struct {
	Bodies []Body `json:"bodies"`
	ID     int32  `json:"id"`
	ID64   int64  `json:"id64"`
	Name   string `json:"name"`
}

func (info SystemInfo) Key() string {
	return fmt.Sprintf("systemInfo/%d", info.ID)
}

func (info *SystemInfo) StarCount() int {
	count := 0
	for _, b := range info.Bodies {
		if b.Type == "Star" {
			count++
		}
	}
	return count
}

func (info *SystemInfo) PlanetCount() int {
	count := 0
	for _, b := range info.Bodies {
		if b.Type == "Planet" {
			count++
		}
	}
	return count
}

func (info *SystemInfo) Stars() []Body {
	bodies := make([]Body, 0)
	for _, b := range info.Bodies {
		if b.Type == "Star" {
			bodies = append(bodies, b)
		}
	}
	return bodies
}

func (info *SystemInfo) Planets() []Body {
	bodies := make([]Body, 0)
	for _, b := range info.Bodies {
		if b.Type == "Planet" {
			bodies = append(bodies, b)
		}
	}
	return bodies
}

type Body struct {
	AbsoluteMagnitude     float64           `json:"absoluteMagnitude"`
	Age                   int64             `json:"age"`
	ArgOfPeriapsis        float64           `json:"argOfPeriapsis"`
	AtmosphereComposition map[string]string `json:"atmosphereComposition"`
	AtmosphereType        string            `json:"atmosphereType"`
	AxialTilt             float64           `json:"axialTilt"`
	Belts                 []struct {
		InnerRadius int64  `json:"innerRadius"`
		Mass        string `json:"mass"`
		Name        string `json:"name"`
		OuterRadius int64  `json:"outerRadius"`
		Type        string `json:"type"`
	} `json:"belts"`
	DistanceToArrival   float64           `json:"distanceToArrival"`
	EarthMasses         float64           `json:"earthMasses"`
	Gravity             float64           `json:"gravity"`
	ID                  int32             `json:"id"`
	ID64                int64             `json:"id64"`
	IsLandable          bool              `json:"isLandable"`
	IsMainStar          bool              `json:"isMainStar"`
	IsScoopable         bool              `json:"isScoopable"`
	Luminosity          string            `json:"luminosity"`
	Materials           map[string]string `json:"materials"`
	Name                string            `json:"name"`
	Offset              int64             `json:"offset"`
	OrbitalEccentricity float64           `json:"orbitalEccentricity"`
	OrbitalInclination  float64           `json:"orbitalInclination"`
	OrbitalPeriod       float64           `json:"orbitalPeriod"`
	Radius              float64           `json:"radius"`
	Rings               []struct {
		InnerRadius float64 `json:"innerRadius"`
		Mass        string  `json:"mass"`
		Name        string  `json:"name"`
		OuterRadius float64 `json:"outerRadius"`
		Type        string  `json:"type"`
	} `json:"rings"`
	RotationalPeriod              float64 `json:"rotationalPeriod"`
	RotationalPeriodTidallyLocked bool    `json:"rotationalPeriodTidallyLocked"`
	SemiMajorAxis                 float64 `json:"semiMajorAxis"`
	SolarMasses                   float64 `json:"solarMasses"`
	SolarRadius                   float64 `json:"solarRadius"`
	SubType                       string  `json:"subType"`
	SurfacePressure               float64 `json:"surfacePressure"`
	SurfaceTemperature            float64 `json:"surfaceTemperature"`
	TerraformingState             string  `json:"terraformingState"`
	Type                          string  `json:"type"`
	UpdateTime                    string  `json:"updateTime"`
	VolcanismType                 string  `json:"volcanismType"`
}

func CheckAPILocked() (bool, error) {
	url := "https://www.edsm.net/api-system-v1/bodies?systemId=27"
	time.Sleep(apiCallWait)
	res, err := http.Get(url)

	if err != nil {
		return false, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	return string(bytes) == "[]", nil
}
