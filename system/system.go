package system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/IgaguriMK/planetStat/vec"
)

var apiCallMux sync.Mutex

type System struct {
	Coords       vec.Vec3 `json:"coords"`
	CoordsLocked bool     `json:"coordsLocked"`
	Name         string   `json:"name"`
	ID           int      `json:"id"`
	ID64         int64    `json:"id64"`
	PermitName   string   `json:"permitName"`
	PrimaryStar  struct {
		IsScoopable bool   `json:"isScoopable"`
		Name        string `json:"name"`
		Type        string `json:"type"`
	} `json:"primaryStar"`
	RequirePermit bool `json:"requirePermit"`
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

	apiCallMux.Lock()
	res, err := http.Get(url)
	apiCallMux.Unlock()

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
