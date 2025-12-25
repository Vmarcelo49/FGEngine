package input

import "fgengine/constants"

func Update(activeScene *constants.Scene) {
	checkForGamepadsConnections()

	switch *activeScene {
	case constants.SceneOptions_ControllerSetup:

	}

}
