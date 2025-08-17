package editor

import (
	"fgengine/animation"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sqweek/dialog"
)

func (g *Game) checkIfResetNeeded() {
	if g.activeCharacter != nil && g.editorManager.activeAnimation != nil {
		g.resetCharacterState()
	}
}

func (g *Game) createCharacter() {
	g.checkIfResetNeeded()
	g.activeCharacter = &animation.Character{
		Animations: make(map[string]*animation.Animation),
		Name:       "character",
	}
}

func (g *Game) loadCharacter() {
	g.checkIfResetNeeded()
	character, err := LoadCharacterFromYAML()
	if err != nil {
		g.writeLog("Failed to load character: " + err.Error())
		return
	}
	g.activeCharacter = character
}

func (g *Game) saveCharacter() {
	if g.activeCharacter == nil {
		g.writeLog("Failed to save character: No active character to save")
		return
	}
	if g.editorManager.activeAnimation != nil && g.editorManager.activeAnimation.Name != "" {
		animCopy := deepCopyAnimation(g.editorManager.activeAnimation)
		g.activeCharacter.Animations[g.editorManager.activeAnimation.Name] = animCopy
		g.writeLog(fmt.Sprintf("Including current animation '%s' in character", g.editorManager.activeAnimation.Name))
	}

	var path string
	var err error

	// If we have a saved file path, use it directly
	if g.activeCharacter.FilePath != "" {
		path = g.activeCharacter.FilePath
	} else {
		// Ask for file path if we don't have one
		path, err = dialog.File().Filter(".yaml", "yaml").Save()
		if err != nil {
			g.writeLog("Failed to save character: " + err.Error())
			return
		}
		// Ensure the path has the correct .yaml extension
		path = ensureExtension(path, "yaml")
		// Remember the file path for future saves
		g.activeCharacter.FilePath = path
	}

	err = ExportCharacterToYaml(g.activeCharacter, path)
	if err != nil {
		g.writeLog("Failed to export character: " + err.Error())
	} else {
		g.writeLog("Character saved successfully!")
	}
}

func ensureExtension(path, extension string) string {
	extension = strings.TrimPrefix(extension, ".")

	currentExt := strings.ToLower(filepath.Ext(path))
	expectedExt := "." + strings.ToLower(extension)

	if currentExt == expectedExt {
		return path
	}

	return path + "." + extension
}
