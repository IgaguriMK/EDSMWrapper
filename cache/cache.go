package cache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var EnableCacheLog bool = true

type CacheController struct {
	dirName string
	maxAge  int64
}

func NewController(dirName string) (*CacheController, error) {
	info, err := os.Stat(dirName)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dirName, os.ModePerm)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else if !info.IsDir() {
		return nil, errors.New("not directory")
	}

	return &CacheController{
		dirName: dirName,
		maxAge:  3600 * 24 * 30,
	}, nil
}

func (cc *CacheController) RegisterType(typeName string) error {
	path := cc.dirName + "/" + typeName
	info, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !info.IsDir() {
		return errors.New("found not directory")
	}

	return nil
}

func (cc *CacheController) Store(version int64, v Cacheable) {
	key := v.Key()

	con, err := json.Marshal(v)
	if err != nil {
		if EnableCacheLog {
			log.Println("[Cache] faled marshal content:", err)
		}
		return
	}

	c := cache{
		Timestamp: time.Now().Unix(),
		Version:   version,
		Content:   string(con),
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		if EnableCacheLog {
			log.Println("[Cache] faled marshal cache:", err)
		}
		return
	}

	f, err := os.Create(cc.dirName + "/" + key + ".json")
	if err != nil {
		if EnableCacheLog {
			log.Println("[Cache] faled create file:", err)
		}
		return
	}
	defer f.Close()

	_, err = f.Write(bytes)
	if err != nil {
		if EnableCacheLog {
			log.Println("[Cache] faled write to file:", err)
		}
		return
	}
}

func (cc *CacheController) Find(version int64, key string, v Cacheable) bool {
	path := cc.dirName + "/" + key + ".json"
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	var c cache
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		if EnableCacheLog {
			log.Println("[Cache] faled parse cache:", err, "cache:", string(bytes))
		}
		return false
	}

	if c.Version != version {
		os.Remove(path)
		return false
	}

	err = json.Unmarshal([]byte(c.Content), &v)
	if err != nil {
		if EnableCacheLog {
			log.Println("[Cache] faled parse content:", err)
		}
		return false
	}

	return true
}

type Cacheable interface {
	Key() string
}

type cache struct {
	Timestamp int64  `json:"timestamp"`
	Version   int64  `json:"version"`
	Content   string `json:"content"`
}
