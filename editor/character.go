package editor

import (
	"fgengine/animation"
	"fgengine/character"
	"fmt"
	"path/filepath"
	"strings"
)

func (g *Game) createCharacter() {
	g.checkIfResetNeeded()
	g.activeCharacter = &character.Character{
		Animations: make(map[string]*animation.Animation),
		Name:       "character",
	}
	g.writeLog("New character created")
}

func (g *Game) loadCharacter() {
	g.checkIfResetNeeded()
	character, err := loadCharacterFromYAMLDialog()
	if err != nil {
		g.writeLog("Failed to load character: " + err.Error())
		return
	}
	g.activeCharacter = character

	// Set initial sprite if there's an animation available
	if len(g.activeCharacter.Animations) > 0 {
		for _, anim := range g.activeCharacter.Animations {
			if len(anim.Sprites) > 0 {
				g.activeCharacter.ActiveSprite = anim.Sprites[0]
				break
			}
		}
	}
	idleAnim, ok := character.Animations["idle"]
	if !ok {
		panic("Character must have an 'idle' animation")
	}
	g.editorManager.activeAnimation = idleAnim
	g.writeLog("Character loaded successfully")
}

// when creating or loading a new character, check if we need to reset the current state
func (g *Game) checkIfResetNeeded() {
	if g.activeCharacter != nil && g.editorManager.activeAnimation != nil {
		g.resetCharacterState()
	}
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

	err := exportCharacterToYAML(g.activeCharacter)
	if err != nil {
		g.writeLog("Failed to export character: " + err.Error())
	} else {
		g.writeLog("Character saved successfully to assets/characters/!")
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
