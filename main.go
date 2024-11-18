package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"v4lweb/camera"
	"v4lweb/web"
)

type Indicator struct{}

func (ind *Indicator) StreamOff() { log.Println("StreamOff") }
func (ind *Indicator) StreamOn()  { log.Println("StreamOn") }

func main() {
	config := &camera.VideoConfig{
		CameraType: camera.V4L_CAMERA,
		Path:       "/dev/video0",
		Codec:      "MJPG",
		Width:      1920,
		Height:     1080,
		FPS:        30,
	}

	cfg, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	mux := &http.ServeMux{}
	httpServer := &http.Server{
		Handler:      mux,
		Addr:         cfg.HttpUrl,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	var indicator = &Indicator{}
	srv, err := web.NewCameraServer(0, config, indicator)
	if err != nil {
		log.Fatal(err)
	}
	web.ServeCamera(mux, srv)

	httpErr := make(chan error, 1)
	go func() {
		httpErr <- httpServer.ListenAndServe()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-httpErr:
		log.Printf("failed to serve http: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

}

type Config struct {
	HttpUrl string
	Camera  *camera.VideoConfig
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
