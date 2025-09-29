package constants

import "image/color"

const (
	WorldWidth  float64 = 960 // 640 = camera width  * 1.5
	WorldHeight float64 = 540 // 360 = camera height * 1.5

	CameraWidth  float64 = 640
	CameraHeight float64 = 360

	Gravity         float64 = 1
	MaxInputHistory int     = 30
)

var StageColor = color.RGBA{R: 100, G: 149, B: 237, A: 255} // Cornflower Blue
