package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"v4lweb/camera"
)

type Config struct {
	HttpUrl string
	Camera  *camera.VideoConfig
	Drivers map[string][]*camera.ControlKey
}

func LoadConfig(filename string) (cfg *Config, err error) {
	cfg = &Config{}
	var f *os.File
	f, err = os.Open(filename)
	if err != nil {
		log.Println("config.Load Open", err)
		return
	}
	defer f.Close()
	var buf []byte
	buf, err = io.ReadAll(f)
	if err != nil {
		log.Println("config.Load ReadAll", err)
		return
	}
	err = json.Unmarshal(buf, cfg)
	if err != nil {
		log.Println("config.Load Unmarshal", err)
		return
	}

	log.Println(cfg.HttpUrl, cfg.Camera.Codec)
	return
}
