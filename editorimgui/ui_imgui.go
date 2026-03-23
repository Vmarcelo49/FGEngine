//go:build !js && !wasm
// +build !js,!wasm

package editorimgui

import (
	"fgengine/collision"
	"fmt"
	"strconv"

	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) updateImGui() error {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()

	g.imguiMainLayout()

	return nil
}

func imguiInputText(label string, value *string) bool {
	return imgui.InputTextWithHint(label, "", value, 0, nil)
}

func (g *Game) imguiMainLayout() {
	io := imgui.CurrentIO()
	display := io.DisplaySize()
	menuY := float32(28)
	leftW := float32(360)
	rightW := float32(420)
	bottomH := float32(210)

	if display.Y < menuY+bottomH+200 {
		bottomH = (display.Y - menuY) * 0.35
	}
	if bottomH < 120 {
		bottomH = 120
	}

	mainPanelsH := display.Y - menuY
	if mainPanelsH < 120 {
		mainPanelsH = 120
	}

	if imgui.BeginMainMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItemBool("New Character") {
				g.createCharacter()
			}
			if imgui.MenuItemBool("Load Character") {
				g.loadCharacter()
			}
			if imgui.MenuItemBool("Save Character") {
				g.saveCharacter()
			}
			imgui.EndMenu()
		}

		if imgui.BeginMenu("View") {
			if imgui.MenuItemBool("Toggle Box Mouse Edit") {
				current := *g.uiVariables.enableMouseInput
				*g.uiVariables.enableMouseInput = !current
			}
			imgui.EndMenu()
		}

		imgui.Text(fmt.Sprintf("TPS: %.2f | FPS: %.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
		imgui.EndMainMenuBar()
	}

	leftFlags := imgui.WindowFlags(imgui.WindowFlagsNoMove | imgui.WindowFlagsNoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoSavedSettings)
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: menuY})
	imgui.SetNextWindowSize(imgui.Vec2{X: leftW, Y: mainPanelsH})
	if imgui.BeginV("Project##LeftPanel", nil, leftFlags) {
		g.imguiCharacterPanel()
		g.imguiAnimationPanel()
	}
	imgui.End()

	rightFlags := imgui.WindowFlags(imgui.WindowFlagsNoMove | imgui.WindowFlagsNoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoSavedSettings)
	imgui.SetNextWindowPos(imgui.Vec2{X: display.X - rightW, Y: menuY})
	imgui.SetNextWindowSize(imgui.Vec2{X: rightW, Y: mainPanelsH})
	if imgui.BeginV("Frame Tools##RightPanel", nil, rightFlags) {
		g.imguiFramePanel()
	}
	imgui.End()

	centerW := display.X - leftW - rightW
	if centerW > 220 {
		logFlags := imgui.WindowFlags(imgui.WindowFlagsNoMove | imgui.WindowFlagsNoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoSavedSettings)
		imgui.SetNextWindowPos(imgui.Vec2{X: leftW, Y: display.Y - bottomH})
		imgui.SetNextWindowSize(imgui.Vec2{X: centerW, Y: bottomH})
		if imgui.BeginV("Log##BottomCenter", nil, logFlags) {
			g.imguiLogPanel()
		}
		imgui.End()
	}
}

func (g *Game) imguiCharacterPanel() {
	if g.character == nil {
		imgui.Text("No character loaded")
		return
	}

	imgui.Separator()
	imgui.Text("Character")
	_ = imguiInputText("Character Name", &g.character.Name)
	imgui.Checkbox("Enable Box Mouse Edit", g.uiVariables.enableMouseInput)
}

func (g *Game) imguiAnimationPanel() {
	if g.character == nil {
		return
	}

	if g.uiVariables.newAnimationName == "" {
		g.uiVariables.newAnimationName = "newAnimation"
	}

	imgui.Separator()
	imgui.Text("Animations")
	_ = imguiInputText("New Animation Name", &g.uiVariables.newAnimationName)
	if imgui.Button("Create Animation From PNG") {
		newAnim, err := g.newAnimationFileDialog()
		if err != nil {
			g.writeLog(fmt.Sprintf("Error creating animation: %v", err))
			return
		}
		if g.uiVariables.newAnimationName != "" {
			newAnim.Name = g.uiVariables.newAnimationName
		}
		if g.animations() == nil {
			return
		}
		g.animations()[newAnim.Name] = newAnim
		g.setActiveAnimation(newAnim.Name, true)
		g.uiVariables.renameAnimationName = newAnim.Name
		g.writeLog(fmt.Sprintf("Animation '%s' created", newAnim.Name))
	}

	for _, animName := range g.animationNames() {
		if imgui.Button("Select: " + animName) {
			g.setActiveAnimation(animName, true)
			g.uiVariables.renameAnimationName = animName
		}
	}

	if g.ActiveAnimation() == nil {
		return
	}

	if g.uiVariables.renameAnimationName == "" {
		g.uiVariables.renameAnimationName = g.ActiveAnimation().Name
	}

	_ = imguiInputText("Rename Active Animation", &g.uiVariables.renameAnimationName)
	if imgui.Button("Apply Animation Rename") {
		g.renameActiveAnimation(g.uiVariables.renameAnimationName)
	}
}

func (g *Game) renameActiveAnimation(newName string) {
	if g.ActiveAnimation() == nil || newName == "" {
		return
	}
	oldName := g.ActiveAnimation().Name
	if oldName == newName {
		return
	}
	anim := g.animations()[oldName]
	if anim == nil {
		return
	}
	delete(g.animations(), oldName)
	anim.Name = newName
	g.animations()[newName] = anim
	g.setActiveAnimation(newName, true)
	g.writeLog(fmt.Sprintf("Animation renamed from '%s' to '%s'", oldName, newName))
}

func (g *Game) imguiFramePanel() {
	if g.ActiveAnimation() == nil {
		return
	}
	player := g.activeAnimPlayer()
	if player == nil {
		return
	}
	frameCount := len(g.ActiveAnimation().FrameData)
	if frameCount == 0 {
		imgui.Text("No frame data")
		return
	}

	if player.FrameIndex < 0 || player.FrameIndex >= frameCount {
		player.FrameIndex = 0
	}

	imgui.Separator()
	imgui.Text(fmt.Sprintf("Frame %d/%d", player.FrameIndex+1, frameCount))

	if imgui.Button("Prev Frame") {
		if player.FrameIndex > 0 {
			player.FrameIndex--
			player.FrameTimeLeft = g.ActiveAnimation().FrameData[player.FrameIndex].Duration
		}
	}
	imgui.SameLine()
	if imgui.Button("Next Frame") {
		if player.FrameIndex < frameCount-1 {
			player.FrameIndex++
			player.FrameTimeLeft = g.ActiveAnimation().FrameData[player.FrameIndex].Duration
		}
	}

	if imgui.Button("Add Frame by Image") {
		g.AddImageToAnimation()
	}
	imgui.SameLine()
	if imgui.Button("Duplicate Frame") {
		g.duplicateLastFrameData()
	}
	imgui.SameLine()
	if imgui.Button("Remove Frame") {
		g.uiVariables.frameDataIndex = player.FrameIndex
		g.removeFrame()
	}

	if !g.uiVariables.playingAnim {
		if imgui.Button("Play") {
			g.uiVariables.playingAnim = true
		}
	} else {
		if imgui.Button("Stop") {
			g.uiVariables.playingAnim = false
		}
	}

	currentFrame := &g.ActiveAnimation().FrameData[player.FrameIndex]
	if g.uiVariables.frameDurationInput == "" {
		g.uiVariables.frameDurationInput = fmt.Sprintf("%d", currentFrame.Duration)
	}
	if imguiInputText("Frame Duration", &g.uiVariables.frameDurationInput) {
		if dur, err := strconv.Atoi(g.uiVariables.frameDurationInput); err == nil && dur > 0 {
			currentFrame.Duration = dur
		}
	}

	if imgui.Button("Sprite -") {
		if currentFrame.SpriteIndex > 0 {
			currentFrame.SpriteIndex--
		}
	}
	imgui.SameLine()
	if imgui.Button("Sprite +") {
		maxSprite := len(g.ActiveAnimation().Sprites) - 1
		if currentFrame.SpriteIndex < maxSprite {
			currentFrame.SpriteIndex++
		}
	}
	imgui.Text(fmt.Sprintf("Sprite Index: %d", currentFrame.SpriteIndex))

	if imgui.Button("Add Collision Box") {
		g.uiVariables.boxDropdownTypeIndex = int(collision.Collision)
		g.addBox()
	}
	imgui.SameLine()
	if imgui.Button("Add Hurt Box") {
		g.uiVariables.boxDropdownTypeIndex = int(collision.Hurt)
		g.addBox()
	}
	imgui.SameLine()
	if imgui.Button("Add Hit Box") {
		g.uiVariables.boxDropdownTypeIndex = int(collision.Hit)
		g.addBox()
	}
}

func (g *Game) imguiLogPanel() {
	imgui.Separator()
	imgui.Text("Log")
	imgui.Text(g.uiVariables.logBuf)
}
