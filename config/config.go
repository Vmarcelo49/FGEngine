package config

var (
	WindowHeight, WindowWidth int
	WorldWidth, WorldHeight   float64
	CameraWidth, CameraHeight int
	ControllerDeadzone        float64
	// worldSize
	// ingameRes 640x360 seems good
)

func InitDefaultConfig() {
	WindowWidth = 1600
	WindowHeight = 900

	WorldWidth = 640
	WorldHeight = 360

	CameraWidth = 640
	CameraHeight = 360

	ControllerDeadzone = 0.3
}
