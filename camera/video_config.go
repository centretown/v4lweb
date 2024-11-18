package camera

type CameraType int

const (
	V4L_CAMERA CameraType = iota
	IP_CAMERA
	CAMERATYPE_COUNT
)

var cameraType = []string{
	"V4L_CAMERA",
	"IP_CAMERA",
}

func (ct CameraType) String() string {
	if ct >= CAMERATYPE_COUNT {
		return "Unknown"
	}
	return cameraType[ct]
}

// {"Format":1448695129,"Width":1280,"Height":720,"FPS":{"N":10,"D":1}}
type VideoConfig struct {
	CameraType CameraType
	Path       string
	Codec      string
	Width      int
	Height     int
	FPS        uint32
}
