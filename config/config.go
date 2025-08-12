package config

var (
	WindowHeight, WindowWidth int
	ControllerDeadzone        float64
	WorldWidth, WorldHeight   float64 // this should not be here, making just a global const seems fine
	CameraWidth, CameraHeight int     // this too
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
