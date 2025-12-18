package scene

import (
	"fgengine/input"
	"os"
)

type Scene int

const (
	MainMenu = iota
	Match
	ControllerSelect
	MatchEnd
	Pause
	Options
	CharacterSelect
	Exit
)

type SceneManager struct {
	currentScene  Scene
	previousScene Scene
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
	sm.LoadScene(MainMenu)
	return sm
}

func (sm *SceneManager) setMenuItems() {

	match := MenuItem{
		Text: "match",
		action: func() {
			sm.LoadScene(Match)
		},
	}
	options := MenuItem{
		Text: "Options",
		action: func() {
			sm.LoadScene(Options)
		},
	}
	exit := MenuItem{
		Text: "Exit",
		action: func() {
			sm.LoadScene(Exit)
		},
	}

	sm.menuItems = []MenuItem{match, options, exit}
}

func (sm *SceneManager) LoadScene(scene Scene) {
	if sm.currentScene != -1 {
		sm.previousScene = sm.currentScene
	}

	sm.currentScene = scene

	switch scene {
	case MainMenu:
		// stuff here
	case ControllerSelect:
		// socorro
	case Exit:
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
