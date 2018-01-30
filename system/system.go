package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/IgaguriMK/edsmWrapper/cache"
	"github.com/IgaguriMK/edsmWrapper/vec"
)

const SystemInfoCacheVer = 1
const RetryCount = 20

var (
	apiCallWaitDefault = time.Millisecond * 4000
	apiCallWait        = apiCallWaitDefault
)

var (
	ErrNotFound  = errors.New("API returns nil response")
	ErrAPILocked = errors.New("API locked")
)

var apiLock sync.Mutex

func init() {
	cache.AddCacheType("systemInfo")
}

func CheckAPILocked() (bool, error) {
	url := "https://www.edsm.net/api-system-v1/bodies?systemId=27"
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
	apiLock.Lock()
	defer apiLock.Unlock()

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

	var bytes []byte

	tryCount := 0
	for {
		var err error
		var ok bool

		time.Sleep(apiCallWait)
		bytes, ok, err = getAPI(url)

		if err != nil {
			return nil, err
		}
		if ok {
			apiCallWait = apiCallWait * 4 / 5
			if apiCallWait < apiCallWaitDefault {
				apiCallWait = apiCallWaitDefault
			} else {
				log.Println("API call wait is", float64(apiCallWait)/1e9)
			}
			break
		}
		if tryCount > RetryCount {
			return nil, ErrAPILocked
		}

		locked, err := CheckAPILocked()
		if err != nil {
			return nil, err
		}
		if !locked {
			break
		}
		apiCallWait = apiCallWait * 2
		log.Println("API call wait is", float64(apiCallWait)/1e9)
		tryCount++
		log.Println("Retry", tryCount)
		time.Sleep(apiCallWait)
	}

	var v []System
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func GetSystemByName(name string) (*System, error) {
	apiLock.Lock()
	defer apiLock.Unlock()

	params := url.Values{}
	params.Add("systemName", name)
	params.Add("showId", "1")
	params.Add("showCoordinates", "1")
	params.Add("showPermit", "1")
	params.Add("showPrimaryStar", "1")

	url := "https://www.edsm.net/api-v1/system?" + params.Encode()
	log.Println(url)

	var bytes []byte

	tryCount := 0
	for {
		var err error
		var ok bool

		time.Sleep(apiCallWait)
		bytes, ok, err = getAPI(url)

		if err != nil {
			return nil, err
		}
		if ok {
			apiCallWait = apiCallWait * 4 / 5
			if apiCallWait < apiCallWaitDefault {
				apiCallWait = apiCallWaitDefault
			} else {
				log.Println("API call wait is", float64(apiCallWait)/1e9)
			}
			break
		}
		if tryCount > RetryCount {
			return nil, ErrAPILocked
		}

		locked, err := CheckAPILocked()
		if err != nil {
			return nil, err
		}
		if !locked {
			break
		}
		apiCallWait = apiCallWait * 2
		log.Println("API call wait is", float64(apiCallWait)/1e9)
		tryCount++
		log.Println("Retry", tryCount)
		time.Sleep(apiCallWait)
	}

	var v System
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (sys System) GetSystemInfo(cc *cache.CacheController) (*SystemInfo, error) {
	apiLock.Lock()
	defer apiLock.Unlock()

	if sys.systemInfoCache != nil {
		return sys.systemInfoCache, nil
	}

	var systemInfo SystemInfo
	if cc.Find(SystemInfoCacheVer, fmt.Sprintf("systemInfo/%d", sys.ID), &systemInfo) {
		return &systemInfo, nil
	}

	url := fmt.Sprintf("https://www.edsm.net/api-system-v1/bodies?systemId=%d", sys.ID)
	log.Println(url)

	var bytes []byte

	tryCount := 0
	for {
		var err error
		var ok bool

		time.Sleep(apiCallWait)
		bytes, ok, err = getAPI(url)

		if err != nil {
			return nil, err
		}
		if ok {
			apiCallWait = apiCallWait * 4 / 5
			if apiCallWait < apiCallWaitDefault {
				apiCallWait = apiCallWaitDefault
			} else {
				log.Println("API call wait is", float64(apiCallWait)/1e9)
			}
			break
		}
		if tryCount > RetryCount {
			return nil, ErrAPILocked
		}

		locked, err := CheckAPILocked()
		if err != nil {
			return nil, err
		}
		if !locked {
			return nil, ErrNotFound
		}
		apiCallWait = apiCallWait * 2
		log.Println("API call wait is", float64(apiCallWait)/1e9)
		tryCount++
		log.Println("Retry", tryCount)
		time.Sleep(apiCallWait)
	}

	err := json.Unmarshal(bytes, &systemInfo)
	if err != nil {
		return nil, err
	}

	cc.Store(SystemInfoCacheVer, systemInfo)

	return &systemInfo, nil
}

func getAPI(url string) ([]byte, bool, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, false, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, false, err
	}

	if string(bytes) == "[]" {
		return bytes, false, nil
	}

	return bytes, true, nil
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
		InnerRadius float64 `json:"innerRadius"`
		Mass        string  `json:"mass"`
		Name        string  `json:"name"`
		OuterRadius float64 `json:"outerRadius"`
		Type        string  `json:"type"`
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

func (b Body) Terraformable() string {
	if b.TerraformingState == "Candidate for terraforming" {
		return "T"
	} else {
		return "F"
	}
}

var shortTypes = map[string]string{
	"A (Blue-White) Star":     "V_A",
	"B (Blue-White) Star":     "V_B",
	"Black Hole":              "X_BH",
	"F (White) Star":          "V_F",
	"G (White-Yellow) Star":   "V_G",
	"Herbig Ae/Be Star":       "P_AeBe",
	"K (Yellow-Orange) Star":  "V_K",
	"L (Brown dwarf) Star":    "BD_L",
	"M (Red dwarf) Star":      "V_M",
	"M (Red giant) Star":      "III_M",
	"MS-type Star":            "M_MS",
	"Neutron Star":            "X_N",
	"O (Blue-White) Star":     "V_O",
	"S-type Star":             "M_S",
	"Supermassive Black Hole": "X_SMBH",
	"T (Brown dwarf) Star":    "BD_R",
	"T Tauri Star":            "P_TTS",
	"Y (Brown dwarf) Star":    "BD_Y",

	"Ammonia world":                     "AW",
	"Class I gas giant":                 "GG1",
	"Class II gas giant":                "GG2",
	"Class III gas giant":               "GG3",
	"Class IV gas giant":                "GG4",
	"Class V gas giant":                 "GG5",
	"Earth-like world":                  "ELW",
	"Gas giant with ammonia-based life": "GGABL",
	"Gas giant with water-based life":   "GGWBL",
	"High metal content world":          "HMC",
	"Icy body":                          "Icy",
	"Metal-rich body":                   "Metal",
	"Rocky Ice world":                   "RockyIce",
	"Rocky body":                        "Rocky",
	"Water giant":                       "WG",
	"Water world":                       "WW",

	"": "None",
}

var ShortTypeUnknown []string

func ShortType(longType string) string {
	s, ok := shortTypes[longType]
	if ok {
		return s
	}

	ShortTypeUnknown = append(ShortTypeUnknown, longType)

	s = longType
	s = strings.Replace(s, " ", "_", -1)
	s = strings.Replace(s, "(", "_", -1)
	s = strings.Replace(s, ")", "_", -1)

	return s
}

func (b Body) ShortSubType() string {
	longType := b.SubType
	s, ok := shortTypes[longType]
	if ok {
		return s
	}

	ShortTypeUnknown = append(ShortTypeUnknown, longType)

	s = longType
	s = strings.Replace(s, " ", "_", -1)
	s = strings.Replace(s, "(", "_", -1)
	s = strings.Replace(s, ")", "_", -1)

	return s
}
