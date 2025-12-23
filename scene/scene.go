package scene

import (
	"fgengine/constants"
	"fgengine/input"
	"os"
)

type SceneManager struct {
	ActiveScene   constants.Scene
	PreviousScene constants.Scene
	menuItems     []MenuItem
	selectedIndex int
}

type MenuItem struct {
	Text   string
	action func()
}

func NewSceneManager() SceneManager {
	sm := SceneManager{}
	sm.setMenuItems()
	sm.LoadScene(constants.SceneMainMenu)
	return sm
}

func (sm *SceneManager) setMenuItems() {

	match := MenuItem{
		Text: "match",
		action: func() {
			sm.LoadScene(constants.SceneMatch)
		},
	}
	options := MenuItem{
		Text: "Options",
		action: func() {
			sm.LoadScene(constants.SceneOptions)
		},
	}
	exit := MenuItem{
		Text: "Exit",
		action: func() {
			sm.LoadScene(constants.SceneExit)
		},
	}

	sm.menuItems = []MenuItem{match, options, exit}
}

func (sm *SceneManager) LoadScene(scene constants.Scene) {
	if sm.ActiveScene != -1 {
		sm.PreviousScene = sm.ActiveScene
	}

	sm.ActiveScene = scene

	switch scene {
	case constants.SceneMainMenu:
		// stuff here
	case constants.SceneControllerSelect:
		// socorro
	case constants.SceneExit:
		os.Exit(0)
	}

}

func (sm *SceneManager) Update(allInputs []input.GameInput) {
	for _, in := range allInputs {
		if in.IsPressed(input.Down) {
			sm.selectedIndex++
			if sm.selectedIndex > len(sm.menuItems)-1 {
				sm.selectedIndex = 0
			}
		}
		if in.IsPressed(input.Up) {
			sm.selectedIndex--
			if sm.selectedIndex < 0 {
				sm.selectedIndex = len(sm.menuItems) - 1
			}
		}
		if in.IsPressed(input.A) || in.IsPressed(input.D) {
			sm.menuItems[sm.selectedIndex].action()
		}
	}
}

type MenuDrawable struct{}

func (md *MenuDrawable) Draw() {

}
