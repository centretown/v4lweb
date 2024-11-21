package web

import (
	"html/template"
	"log"
	"net/http"
	"v4lweb/camera"
	"v4lweb/config"
)

func NewCameraServer(id int, vcfg *camera.VideoConfig,
	indicator camera.StreamIndicator) (cameraServer *camera.Server, err error) {

	var source camera.VideoSource
	switch vcfg.CameraType {
	case camera.V4L_CAMERA:
		source = camera.NewWebcam(vcfg.Path)
	case camera.IP_CAMERA:
		source = camera.NewIpcam(vcfg.Path)
	default:
		return
	}
	cameraServer = camera.NewVideoServer(id, source, vcfg, indicator)
	err = cameraServer.Open()
	return
}

func ServeCamera(cfg *config.Config, mux *http.ServeMux, camServer *camera.Server) {

	log.Println("serveCamera", camServer.Url())

	tmpl, err := template.New("layout.response").
		Parse(`<div id="response-div" class="fade-it">{{.}}</div>`)
	if err != nil {
		log.Fatal(err)
	}
	webcamHandlers := CreateNexigoHandlers(cfg, tmpl)

	mux.Handle(camServer.Url(), camServer.Stream())
	source := camServer.Source
	webcam, isWebcam := source.(*camera.Webcam)
	if isWebcam {
		ctll := NewControlList(mux, webcam, 0, webcamHandlers)
		mux.HandleFunc("/resetcontrols",
			func(w http.ResponseWriter, r *http.Request) {
				ctll.ResetControls()
			})
	}

	go camServer.Serve()
	log.Printf("Serving %s\n", camServer.Url())
}
