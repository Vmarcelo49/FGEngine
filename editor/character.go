package editor

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/collision"
	"fgengine/state"
	"fgengine/types"
	"fmt"
	"path/filepath"
	"strings"
)

func (g *Game) createCharacter() {
	g.checkIfResetNeeded()
	g.activeCharacter = &character.Character{
		Animations:   make(map[string]*animation.Animation),
		Name:         "character",
		StateMachine: &state.StateMachine{}, //needed for the rect
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
	g.activeCharacter.StateMachine = &state.StateMachine{}

	// Set initial sprite if there's an animation available
	if len(g.activeCharacter.Animations) > 0 {
		for _, anim := range g.activeCharacter.Animations {
			if len(anim.FrameData) > 0 && len(anim.Sprites) > 0 {
				spriteIndex := anim.FrameData[0].SpriteIndex
				if spriteIndex >= 0 && spriteIndex < len(anim.Sprites) {
					g.activeCharacter.ActiveSprite = anim.Sprites[spriteIndex]
				} else if len(anim.Sprites) > 0 {
					g.activeCharacter.ActiveSprite = anim.Sprites[0] // Fallback
				}
				break
			}
		}
	}

	// Ensure animations map is initialized
	if character.Animations == nil {
		character.Animations = make(map[string]*animation.Animation)
	}

	idleAnim, ok := character.Animations["idle"]
	if !ok {
		// Create a placeholder idle animation using notFound.png
		g.writeLog("No 'idle' animation found, creating placeholder...")
		idleAnim = g.createPlaceholderIdleAnimation()
		character.Animations["idle"] = idleAnim
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
