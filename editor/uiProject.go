package editor

import (
	"fgengine/config"
	"fmt"
	"image"
	"slices"

	"github.com/ebitengine/debugui"
)

func (g *Game) uiProjectPanel(ctx *debugui.Context) {
	ctx.Window("Project", image.Rect(leftPanelX, toolbarHeight, panelWidth, config.WindowHeight), func(layout debugui.ContainerLayout) {
		ctx.Header("Character", true, func() {
			if g.activeCharacter == nil {
				ctx.Text("No character loaded")
				return
			}

			ctx.SetGridLayout([]int{-1, -1}, nil)
			ctx.Text("Name:")
			ctx.TextField(&g.activeCharacter.Name).On(func() {
				if g.activeCharacter.Name == "" {
					g.activeCharacter.Name = "character"
				}
			})

			ctx.Text("ID:")
			ctx.NumberField(&g.activeCharacter.ID, 1)

			ctx.Text("Friction:")
			ctx.NumberField(&g.activeCharacter.Friction, 1)

			ctx.Text("Jump Height:")
			ctx.NumberField(&g.activeCharacter.JumpHeight, 1)
		})
		ctx.SetGridLayout([]int{-1}, nil)
		if g.activeCharacter != nil {
			ctx.Button("Create New Animation").On(func() {
				g.menuBarNewAnim()
			})
		}
		animNames := g.getAnimationNames()
		if len(animNames) == 0 {
			ctx.Text("No animations found")

			return
		}
		ctx.SetGridLayout([]int{-1, -1}, nil)
		ctx.Text("Select Animation:")
		ctx.Dropdown(&g.editorManager.animationSelectionIndex, animNames).On(func() {
			g.editorManager.activeAnimation = g.activeCharacter.Animations[animNames[g.editorManager.animationSelectionIndex]]
			g.editorManager.previousAnimationName = g.editorManager.activeAnimation.Name
		})
		ctx.Text("Animation Name:")
		if g.editorManager.previousAnimationName == "" && g.editorManager.activeAnimation != nil {
			g.editorManager.previousAnimationName = g.editorManager.activeAnimation.Name
		}
		ctx.TextField(&g.editorManager.activeAnimation.Name).On(func() {
			if g.editorManager.previousAnimationName != "" {
				delete(g.activeCharacter.Animations, g.editorManager.previousAnimationName)
			}

			if g.editorManager.activeAnimation.Name == "" {
				g.editorManager.activeAnimation.Name = "default"
			}
			g.activeCharacter.Animations[g.editorManager.activeAnimation.Name] = g.editorManager.activeAnimation
			g.writeLog(fmt.Sprintf("Animation name changed from '%s' to '%s'", g.editorManager.previousAnimationName, g.editorManager.activeAnimation.Name))
			g.writeLog(fmt.Sprintf("Current animations: %s", g.getAnimationNames()))

			g.editorManager.previousAnimationName = g.editorManager.activeAnimation.Name
		})
		g.boxEditor(ctx)
	})
}

func (g *Game) getAnimationNames() []string {
	if g.activeCharacter == nil {
		return nil
	}
	var anims []string
	for name := range g.activeCharacter.Animations {
		anims = append(anims, name)
	}
	slices.Sort(anims) // consistent ordering for dropdown
	return anims
}

func (g *Game) menuBarNewAnim() {
	newAnim, err := g.editorManager.newAnimationFileDialog()
	if err != nil {
		g.writeLog(fmt.Sprintf("Error creating new animation: %v", err))
		return
	}
	if g.activeCharacter != nil {
		g.activeCharacter.Animations[newAnim.Name] = newAnim
	}
	g.editorManager.activeAnimation = g.activeCharacter.Animations[newAnim.Name]
	g.editorManager.previousAnimationName = newAnim.Name
	g.writeLog("New animation created successfully")
	g.activeCharacter.ActiveSprite = newAnim.Sprites[0]
}
