//go:build !js && !wasm
// +build !js,!wasm

package editorimgui

import "slices"

func (g *Game) writeLog(text string) {
	if len(g.uiVariables.logBuf) > 0 {
		g.uiVariables.logBuf += "\n"
	}
	g.uiVariables.logBuf += text
}

func (g *Game) resetCharacterState() {
	g.character = nil
	g.writeLog("There was a character loaded, cleared current state")
}

func (g *Game) animationNames() []string {
	if g.animations() == nil {
		return nil
	}
	var names []string
	for name := range g.animations() {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
