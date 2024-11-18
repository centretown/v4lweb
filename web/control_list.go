package web

import (
	"log"
	"net/http"
	"v4lweb/camera"
)

type ControlList struct {
	webcam   *camera.Webcam
	Id       int
	Handlers []*WebcamHandler
}

func NewControlList(mux *http.ServeMux, webcam *camera.Webcam, id int, handlers []*WebcamHandler) *ControlList {
	ctll := &ControlList{
		webcam:   webcam,
		Id:       id,
		Handlers: make([]*WebcamHandler, 0, len(handlers)),
	}
	for _, ctl := range handlers {
		ctll.AddHandler(mux, ctl)
	}
	return ctll
}

func (ctll *ControlList) AddHandler(mux *http.ServeMux, ctlh *WebcamHandler) {
	var err error
	if ctlh == nil {
		log.Fatalln("AddControl control is nil")
	}
	ctlh.webcam = ctll.webcam
	ctlh.Info, err = ctlh.webcam.GetControlInfo(ctlh.Key)
	ctlh.Value = ctlh.webcam.GetControlValue(ctlh.Key)

	ctll.Handlers = append(ctll.Handlers, ctlh)
	if err != nil {
		log.Println("AddControl", err)
	}

	for _, ctl := range ctlh.Controls {
		mux.Handle(ctl.Url, ctlh)
	}

}

func (ctlh *ControlList) ResetControls() {
	for _, ctl := range ctlh.Handlers {
		ctl.webcam.SetValue(ctl.Key, ctl.Info.Default)
		log.Println("ResetControls", ctl.Key, ctl.Info.Default)
		ctl.Value = ctl.Info.Default
	}
}
