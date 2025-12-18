package editor

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/collision"
	"fgengine/constants"
	"fgengine/state"
	"fgengine/types"
)

func (g *Game) createCharacter() {
	g.deleteCurrentCharacter()
	g.character = &character.Character{
		Animations:      make(map[string]*animation.Animation),
		Name:            "character",
		StateMachine:    &state.StateMachine{},
		AnimationPlayer: &animation.AnimationPlayer{},
	}
	g.renderQueue.Clear()
	g.renderQueue.Add(g.character, constants.LayerPlayer)
	g.renderQueue.Add(&character.BoxDrawable{Character: g.character}, constants.LayerHUD)
	g.writeLog("New character created")
}

func (g *Game) updateAnimationFrame() {
	if !g.uiVariables.playingAnim || g.character == nil {
		return
	}
	g.character.AnimationPlayer.Update()
}

func (g *Game) ActiveAnimation() *animation.Animation {
	if g.character == nil {
		return nil
	}
	if g.character.AnimationPlayer.ActiveAnimation != nil {
		return g.character.AnimationPlayer.ActiveAnimation
	}
	return nil
}

func (g *Game) loadCharacter() {
	g.deleteCurrentCharacter()
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
func (g *Game) deleteCurrentCharacter() {
	if g.character != nil && g.ActiveAnimation() != nil {
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

// createPlaceholderIdleAnimation creates a default idle animation using notFound.png
func (g *Game) createPlaceholderIdleAnimation() *animation.Animation {
	placeholderSprite := &animation.Sprite{
		ImagePath: "..\\common\\notFound.png",
		Rect: types.Rect{
			W: 64,
			H: 64,
		},
	}

	placeholderFrame := animation.FrameData{
		Duration: 1,
		Boxes:    make(map[collision.BoxType][]types.Rect),
	}

	return &animation.Animation{
		Name:      "idle",
		Sprites:   []*animation.Sprite{placeholderSprite},
		FrameData: []animation.FrameData{placeholderFrame},
	}
}
