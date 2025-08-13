package config

// Only user-configurable settings should be here
var (
	WindowHeight, WindowWidth int
	ControllerDeadzone        float64
)

func InitDefaultConfig() {
	WindowWidth = 1600
	WindowHeight = 900
	ControllerDeadzone = 0.3
}
