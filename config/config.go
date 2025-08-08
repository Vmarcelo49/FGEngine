package config

var (
	WindowHeight, WindowWidth int
	Zoom                      float64
	// worldSize
	// ingameRes 640x360 seems good
)

func GetZoom() float64 {
	if Zoom <= 0 {
		return 1.0 // Default zoom level
	}
	return Zoom
}

func InitDefaultConfig() {
	WindowHeight = 1600
	WindowWidth = 900
	Zoom = 1.0
}
