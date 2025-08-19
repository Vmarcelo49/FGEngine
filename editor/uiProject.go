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
		//g.boxEditor()
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
	slices.Sort(anims) // Sorting alphabetically for consistent ordering
	return anims
}

func (g *Game) menuBarNewAnim() {
	newAnim, err := g.editorManager.newAnimationFileDialog()
	if err != nil {
		g.writeLog(fmt.Sprintf("Error creating new animation: %v", err.Error()))
	} else {
		if g.activeCharacter != nil {
			g.activeCharacter.Animations[newAnim.Name] = newAnim
		}
		g.editorManager.activeAnimation = g.activeCharacter.Animations[newAnim.Name]
		g.editorManager.previousAnimationName = newAnim.Name
		g.writeLog("New animation created successfully")
	}
}

/*
func (g *Game) boxEditor() {
	ctx.Header("Box Editor", true, func() {
		frame := g.getCurrentFrame()
		if frame == nil {
			ctx.Text("No frame selected")
			return
		}

		if g.anim.boxRenderer == nil {
			g.loadBoxRenderer(frame)
		}

		ctx.SetGridLayout([]int{-1}, nil)
		ctx.Text(fmt.Sprintf("Boxes - Collision: %d, Hurt: %d, Hit: %d", len(frame.CollisionBoxes), len(frame.HurtBoxes), len(frame.HitBoxes)))
		ctx.Checkbox(&g.choiceShowAllBoxes, "Show All Boxes")

		if g.anim.boxRenderer.selectedBox != nil {
			boxTypeStr := boxTypes[g.anim.boxRenderer.selectedBoxType]
			ctx.Text(fmt.Sprintf("Selected: %s Box", boxTypeStr))

			ctx.SetGridLayout([]int{-1, -1}, nil)
			ctx.Button("Clear Selection").On(func() {
				g.clearBoxSelection()
			})
			ctx.Button("Delete Box").On(func() {
				g.deleteSelectedBox()
			})
		} else {
			ctx.Text("No box selected")
			ctx.SetGridLayout([]int{-1}, nil)
			ctx.Button("Clear Selection").On(func() {
				g.clearBoxSelection()
			})
		}

		ctx.SetGridLayout([]int{-1}, nil)
		ctx.Text("Create New Box:")
		ctx.SetGridLayout([]int{-2, -1}, nil)
		ctx.Dropdown(&g.boxActionIndex, boxTypes).On(func() {
			g.clearBoxSelection()
			g.anim.boxRenderer.currentBoxType = BoxType(g.boxActionIndex)
		})
		ctx.Button("Add Box").On(func() {
			g.addBox()
		})

		if g.anim.boxRenderer != nil && g.anim.boxRenderer.selectedBox != nil {
			ctx.SetGridLayout([]int{-1}, nil)
			ctx.Text("Box Properties:")
			ctx.SetGridLayout([]int{-1, -1}, nil)
			ctx.Text("X:")
			ctx.NumberFieldF(&g.anim.boxRenderer.selectedBox.rect.X, 1.0, 1)
			ctx.Text("Y:")
			ctx.NumberFieldF(&g.anim.boxRenderer.selectedBox.rect.Y, 1.0, 1)
			ctx.Text("Width:")
			ctx.NumberFieldF(&g.anim.boxRenderer.selectedBox.rect.W, 1.0, 1)
			ctx.Text("Height:")
			ctx.NumberFieldF(&g.anim.boxRenderer.selectedBox.rect.H, 1.0, 1)
		}
	})
}
*/
