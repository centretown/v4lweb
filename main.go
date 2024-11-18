package main

import (
	"log"
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
	var indicator = &Indicator{}
	srv, err := web.NewCameraServer(0, config, indicator)
	if err != nil {
		log.Fatal(err)
	}
	go srv.Serve()
}
