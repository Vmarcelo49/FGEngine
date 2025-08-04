package config

var (
	WindowHeight, WindowWidth int
	Zoom                      float64
)

func GetZoom() float64 {
	if Zoom <= 0 {
		return 1.0 // Default zoom level
	}
	return Zoom
}
