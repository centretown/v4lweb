package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"v4lweb/config"
	"v4lweb/web"
)

type Indicator struct{}

func (ind *Indicator) StreamOff() { log.Println("StreamOff") }
func (ind *Indicator) StreamOn()  { log.Println("StreamOn") }

func main() {
	// videoConfig := &camera.VideoConfig{
	// 	CameraType: camera.V4L_CAMERA,
	// 	Path:       "/dev/video0",
	// 	Codec:      "MJPG",
	// 	Width:      1920,
	// 	Height:     1080,
	// 	FPS:        30,
	// }

	cfg, err := config.LoadConfig("config.json")
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
	srv, err := web.NewCameraServer(0, cfg.Camera, indicator)
	if err != nil {
		log.Fatal(err)
	}

	web.ServeCamera(cfg, mux, srv)

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
