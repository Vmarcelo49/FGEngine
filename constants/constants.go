package constants

import (
	"fgengine/types"
	"image/color"
)

const (
	WorldWidth  float64 = 768 // 640 = camera width  * 1.2
	WorldHeight float64 = 432 // 360 = camera height * 1.2

	CameraWidth  float64 = 640
	CameraHeight float64 = 360

	Gravity         float64 = 1
	MaxInputHistory int     = 30

	GroundLevelY float64 = WorldHeight - 50
)

const (
	LayerBG = iota
	LayerHUD
	LayerPlayer
	LayerEffects
)

const LayerCount = 4

var StageColor = color.RGBA{R: 100, G: 149, B: 237, A: 255} // Cornflower Blue

var World = types.Rect{X: 0, Y: 0, W: WorldWidth, H: WorldHeight} // camera * 1.2
var Camera = types.Rect{X: 0, Y: 0, W: CameraWidth, H: CameraHeight}

type Scene int

const (
	SceneMainMenu = iota
	SceneMatch
	SceneControllerSelect
	SceneMatchEnd
	ScenePause
	SceneOptions
	SceneOptions_ControllerSetup
	SceneCharacterSelect
	SceneExit
)
