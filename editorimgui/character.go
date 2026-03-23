//go:build !js && !wasm
// +build !js,!wasm

package editorimgui

import (
	"fgengine/animation"
	"fgengine/character"
	"fgengine/constants"
	"fgengine/types"
)

func (g *Game) activeAnimPlayer() *animation.AnimationPlayer {
	if g.character == nil || g.character.StateMachine == nil {
		return nil
	}
	return g.character.StateMachine.ActiveAnim
}

func (g *Game) animations() map[string]*animation.Animation {
	player := g.activeAnimPlayer()
	if player == nil {
		return nil
	}
	return player.Animations
}

func (g *Game) setActiveAnimation(name string, loop bool) {
	player := g.activeAnimPlayer()
	if player == nil {
		return
	}
	player.SetAnimation(name, loop)
}

func (g *Game) setupEditorPreviewCharacter() {
	if g.character == nil {
		return
	}
	if g.character.StateMachine == nil {
		g.character.StateMachine = &animation.StateMachine{}
	}
	if g.character.StateMachine.ActiveAnim == nil {
		g.character.StateMachine.ActiveAnim = &animation.AnimationPlayer{}
	}
	if g.character.StateMachine.ActiveAnim.Animations == nil {
		g.character.StateMachine.ActiveAnim.Animations = make(map[string]*animation.Animation)
	}

	g.character.StateMachine.Position = types.Vector2{X: constants.WorldWidth / 2, Y: constants.GroundLevelY}
	g.character.StateMachine.Velocity = types.Vector2{}
	g.character.StateMachine.IgnoreGravityFrames = 0

	if g.camera != nil {
		g.camera.SetPosition(types.Vector2{X: -constants.Camera.W - 150, Y: -200})
	}
}

func (g *Game) createCharacter() {
	g.deleteCurrentCharacter()
	g.character = &character.Character{
		Name: "character",
		StateMachine: &animation.StateMachine{
			ActiveAnim: &animation.AnimationPlayer{
				Animations: make(map[string]*animation.Animation),
			},
		},
	}
	g.setupEditorPreviewCharacter()
	g.renderQueue.Clear()
	g.renderQueue.Add(g.character, constants.LayerPlayer)
	g.renderQueue.Add(&character.BoxDrawable{Character: g.character}, constants.LayerHUD)
	g.writeLog("New character created")
}

func (g *Game) updateAnimationFrame() {
	if !g.uiVariables.playingAnim {
		return
	}
	player := g.activeAnimPlayer()
	if player == nil {
		return
	}
	player.Update()
}

func (g *Game) ActiveAnimation() *animation.Animation {
	player := g.activeAnimPlayer()
	if player != nil && player.ActiveAnimation != nil {
		return player.ActiveAnimation
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
	g.setupEditorPreviewCharacter()

	if len(g.animations()) > 0 {
		// Try to set idle animation, or first available animation
		if _, exists := g.animations()["idle"]; exists {
			g.setActiveAnimation("idle", true)
		} else {
			for name := range g.animations() {
				g.setActiveAnimation(name, true)
				break
			}
		}
	}
	g.writeLog("Character loaded successfully")
}

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
