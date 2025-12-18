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
			if g.character == nil {
				ctx.Text("No character loaded")
				return
			}

			ctx.SetGridLayout([]int{-1, -1}, nil)
			ctx.Text("Name:")
			ctx.TextField(&g.character.Name)

			ctx.Text("ID:")
			ctx.NumberField(&g.character.ID, 1)

			ctx.Text("Friction:")
			ctx.NumberFieldF(&g.character.Friction, 1, 2)

			ctx.Text("Jump Height:")
			ctx.NumberFieldF(&g.character.JumpHeight, 1, 2)
		})
		ctx.SetGridLayout([]int{-1}, nil)
		if g.character != nil {
			ctx.Button("Create New Animation").On(func() {
				g.menuBarNewAnim()
			})
		}

		animNames := g.animationNames()
		if len(animNames) == 0 {
			ctx.Text("No animations found")
			return
		}

		ctx.SetGridLayout([]int{-1, -1}, nil)
		ctx.Text("Select Animation:")
		ctx.Dropdown(&g.uiVariables.animationSelectionIndex, animNames).On(func() {
			g.character.SetAnimation(animNames[g.uiVariables.animationSelectionIndex])
		})
		ctx.Text("Animation Name:")
		oldName := animNames[g.uiVariables.animationSelectionIndex]
		ctx.TextField(&g.ActiveAnimation().Name).On(func() {
			if g.ActiveAnimation().Name == "" {
				g.ActiveAnimation().Name = "noName" // animations with no names can cause issues
			}

			newName := g.ActiveAnimation().Name

			if oldName != newName {
				anim := g.character.Animations[oldName]
				delete(g.character.Animations, oldName)
				g.character.Animations[newName] = anim
				g.writeLog(fmt.Sprintf("Animation renamed from '%s' to '%s'", oldName, newName))

				newAnimNames := g.animationNames()
				for i, name := range newAnimNames {
					if name == newName {
						g.uiVariables.animationSelectionIndex = i
						break
					}
				}
			}
		})
		g.boxEditor(ctx)

	})
}

func (g *Game) animationNames() []string {
	if g.character == nil {
		return nil
	}
	var anims []string
	for name := range g.character.Animations {
		anims = append(anims, name)
	}
	slices.Sort(anims) // consistent ordering for dropdown
	return anims
}

func (g *Game) menuBarNewAnim() {
	newAnim, err := g.newAnimationFileDialog()
	if err != nil {
		g.writeLog(fmt.Sprintf("Error creating new animation: %v", err))
		return
	}

	g.character.Animations[newAnim.Name] = newAnim
	g.character.AnimationPlayer.ActiveAnimation = newAnim
	g.writeLog("New animation created successfully")
}
