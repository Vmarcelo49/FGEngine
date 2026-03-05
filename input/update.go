package input

import "fgengine/constants"

func Update(activeScene *constants.Scene) {
	UpdateGamepads()

	switch *activeScene {
	case constants.SceneOptions_ControllerSetup:

	}

}
