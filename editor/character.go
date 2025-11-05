package editor

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/collision"
	"fgengine/state"
	"fgengine/types"
	"path/filepath"
	"strings"
)

func (g *Game) createCharacter() {
	g.checkIfResetNeeded()
	g.character = &character.Character{
		Animations:      make(map[string]*animation.Animation),
		Name:            "character",
		StateMachine:    &state.StateMachine{},
		AnimationPlayer: &animation.AnimationPlayer{},
	}
	g.writeLog("New character created")
}

func (g *Game) updateAnimationFrame() {
	if g.uiVariables.playingAnim && g.character != nil {
		g.character.AnimationPlayer.FrameCounter++
	}
	animPlayer := g.character.AnimationPlayer // Just to reduce line length

	if animPlayer.ShouldLoop {
		if animPlayer.FrameCounter >= animPlayer.ActiveAnimation.Duration() {
			animPlayer.FrameCounter = 0
			return
		}
	}
	if animPlayer.FrameCounter >= animPlayer.ActiveAnimation.Duration() {
		animPlayer.FrameCounter = animPlayer.ActiveAnimation.Duration() - 1
	}
}

func (g *Game) getActiveAnimation() *animation.Animation {
	if g.character == nil {
		return nil
	}
	if g.character.AnimationPlayer.ActiveAnimation != nil {
		return g.character.AnimationPlayer.ActiveAnimation
	}
	return nil
}

func (g *Game) loadCharacter() {
	g.checkIfResetNeeded()
	character, err := loadCharacterFromYAMLDialog()
	if err != nil {
		g.writeLog("Failed to load character: " + err.Error())
		return
	}
	g.character = character

	if len(g.character.Animations) > 0 {
		// Try to set idle animation, or first available animation
		if idleAnim, exists := g.character.Animations["idle"]; exists {
			g.character.AnimationPlayer.ActiveAnimation = idleAnim
		} else {
			for _, anim := range g.character.Animations {
				g.character.AnimationPlayer.ActiveAnimation = anim
				break
			}
		}
	}
	g.writeLog("Character loaded successfully")
}

// when creating or loading a new character, check if we need to reset the current state
func (g *Game) checkIfResetNeeded() { // TODO rename this to deleteCurrentCharacter
	if g.character != nil && g.getActiveAnimation() != nil {
		g.resetCharacterState()
	}
}

func (g *Game) saveCharacter() {
	if g.character == nil {
		g.writeLog("Failed to save character: No active character to save")
		return
	}

	err := exportCharacterToYAML(g.character)
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

// createPlaceholderIdleAnimation creates a default idle animation using notFound.png
func (g *Game) createPlaceholderIdleAnimation() *animation.Animation {
	placeholderSprite := &animation.Sprite{
		ImagePath: "..\\common\\notFound.png",
		Rect: types.Rect{
			W: 64, // Default width
			H: 64, // Default height
		},
		Boxes: make(map[collision.BoxType][]types.Rect),
	}

	placeholderFrame := animation.FrameData{
		Duration:    60, // 1 second at 60 FPS
		SpriteIndex: 0,  // Reference to the first (and only) sprite
	}

	return &animation.Animation{
		Name:      "idle",
		Sprites:   []*animation.Sprite{placeholderSprite},
		FrameData: []animation.FrameData{placeholderFrame},
	}
}
