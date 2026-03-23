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

const (
	modernThemeVarCount   = 8
	modernThemeColorCount = 18
)

func (g *Game) updateImGui() error {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()

	g.pushModernTheme()
	defer g.popModernTheme()

	g.imguiMainLayout()

	return nil
}

func imguiInputText(label string, value *string) bool {
	return imgui.InputTextWithHint(label, "", value, 0, nil)
}

func (g *Game) pushModernTheme() {
	imgui.StyleColorsDark()

	imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 10)
	imgui.PushStyleVarFloat(imgui.StyleVarFrameRounding, 7)
	imgui.PushStyleVarFloat(imgui.StyleVarGrabRounding, 7)
	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 1)
	imgui.PushStyleVarFloat(imgui.StyleVarFrameBorderSize, 1)
	imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 14, Y: 14})
	imgui.PushStyleVarVec2(imgui.StyleVarFramePadding, imgui.Vec2{X: 10, Y: 8})
	imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{X: 10, Y: 8})

	imgui.PushStyleColorVec4(imgui.ColWindowBg, imgui.Vec4{X: 0.08, Y: 0.09, Z: 0.11, W: 0.97})
	imgui.PushStyleColorVec4(imgui.ColChildBg, imgui.Vec4{X: 0.10, Y: 0.11, Z: 0.14, W: 0.96})
	imgui.PushStyleColorVec4(imgui.ColMenuBarBg, imgui.Vec4{X: 0.06, Y: 0.07, Z: 0.09, W: 0.98})
	imgui.PushStyleColorVec4(imgui.ColBorder, imgui.Vec4{X: 0.22, Y: 0.25, Z: 0.31, W: 0.55})
	imgui.PushStyleColorVec4(imgui.ColText, imgui.Vec4{X: 0.92, Y: 0.94, Z: 0.98, W: 1.00})
	imgui.PushStyleColorVec4(imgui.ColTextDisabled, imgui.Vec4{X: 0.58, Y: 0.62, Z: 0.70, W: 1.00})
	imgui.PushStyleColorVec4(imgui.ColFrameBg, imgui.Vec4{X: 0.14, Y: 0.16, Z: 0.20, W: 0.95})
	imgui.PushStyleColorVec4(imgui.ColFrameBgHovered, imgui.Vec4{X: 0.18, Y: 0.21, Z: 0.27, W: 1.00})
	imgui.PushStyleColorVec4(imgui.ColFrameBgActive, imgui.Vec4{X: 0.22, Y: 0.27, Z: 0.36, W: 1.00})
	imgui.PushStyleColorVec4(imgui.ColTitleBg, imgui.Vec4{X: 0.09, Y: 0.10, Z: 0.13, W: 1.00})
	imgui.PushStyleColorVec4(imgui.ColTitleBgActive, imgui.Vec4{X: 0.10, Y: 0.13, Z: 0.19, W: 1.00})
	imgui.PushStyleColorVec4(imgui.ColButton, imgui.Vec4{X: 0.20, Y: 0.43, Z: 0.92, W: 0.80})
	imgui.PushStyleColorVec4(imgui.ColButtonHovered, imgui.Vec4{X: 0.28, Y: 0.52, Z: 0.99, W: 0.95})
	imgui.PushStyleColorVec4(imgui.ColButtonActive, imgui.Vec4{X: 0.15, Y: 0.34, Z: 0.76, W: 1.00})
	imgui.PushStyleColorVec4(imgui.ColHeader, imgui.Vec4{X: 0.16, Y: 0.22, Z: 0.34, W: 0.82})
	imgui.PushStyleColorVec4(imgui.ColHeaderHovered, imgui.Vec4{X: 0.23, Y: 0.31, Z: 0.45, W: 0.90})
	imgui.PushStyleColorVec4(imgui.ColHeaderActive, imgui.Vec4{X: 0.13, Y: 0.19, Z: 0.30, W: 1.00})
	imgui.PushStyleColorVec4(imgui.ColSeparator, imgui.Vec4{X: 0.26, Y: 0.31, Z: 0.39, W: 0.70})
}

func (g *Game) popModernTheme() {
	imgui.PopStyleColorV(modernThemeColorCount)
	imgui.PopStyleVarV(modernThemeVarCount)
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

		imgui.SameLine()
		imgui.Text("|")
		imgui.SameLine()
		imgui.Text(fmt.Sprintf("TPS %.2f  FPS %.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
		imgui.EndMainMenuBar()
	}

	leftFlags := imgui.WindowFlags(imgui.WindowFlagsNoMove | imgui.WindowFlagsNoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoSavedSettings)
	imgui.SetNextWindowBgAlpha(0.92)
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: menuY})
	imgui.SetNextWindowSize(imgui.Vec2{X: leftW, Y: mainPanelsH})
	if imgui.BeginV("Project##LeftPanel", nil, leftFlags) {
		g.imguiCharacterPanel()
		g.imguiAnimationPanel()
	}
	imgui.End()

	rightFlags := imgui.WindowFlags(imgui.WindowFlagsNoMove | imgui.WindowFlagsNoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoSavedSettings)
	imgui.SetNextWindowBgAlpha(0.92)
	imgui.SetNextWindowPos(imgui.Vec2{X: display.X - rightW, Y: menuY})
	imgui.SetNextWindowSize(imgui.Vec2{X: rightW, Y: mainPanelsH})
	if imgui.BeginV("Frame Tools##RightPanel", nil, rightFlags) {
		g.imguiFramePanel()
	}
	imgui.End()

	centerW := display.X - leftW - rightW
	if centerW > 220 {
		logFlags := imgui.WindowFlags(imgui.WindowFlagsNoMove | imgui.WindowFlagsNoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoSavedSettings)
		imgui.SetNextWindowBgAlpha(0.88)
		imgui.SetNextWindowPos(imgui.Vec2{X: leftW, Y: display.Y - bottomH})
		imgui.SetNextWindowSize(imgui.Vec2{X: centerW, Y: bottomH})
		if imgui.BeginV("Log##BottomCenter", nil, logFlags) {
			g.imguiLogPanel()
		}
		imgui.End()
	}
}

func (g *Game) imguiCharacterPanel() {
	imgui.SeparatorText("Character")
	if imgui.BeginChildStrV("CharacterCard", imgui.Vec2{X: 0, Y: 84}, true, imgui.WindowFlagsNone) {
		if g.character == nil {
			imgui.Text("No character loaded")
		} else {
			_ = imguiInputText("Character Name", &g.character.Name)
			imgui.Checkbox("Enable Box Mouse Edit", g.uiVariables.enableMouseInput)
		}
	}
	imgui.EndChild()
	imgui.Spacing()

	if g.character != nil {
		imgui.TextDisabled("Camera: Right mouse drag")
	}

	if g.character == nil {
		return
	}

	imgui.Spacing()
	imgui.SeparatorText("Quick Actions")
	if imgui.Button("Save Character") {
		g.saveCharacter()
	}
	imgui.SameLine()
	if imgui.Button("Load Character") {
		g.loadCharacter()
	}
}

func (g *Game) imguiAnimationPanel() {
	imgui.SeparatorText("Animations")
	if g.character == nil {
		imgui.TextDisabled("Create or load a character first")
		return
	}

	if g.uiVariables.newAnimationName == "" {
		g.uiVariables.newAnimationName = "newAnimation"
	}

	if imgui.BeginChildStrV("AnimCreateCard", imgui.Vec2{X: 0, Y: 92}, true, imgui.WindowFlagsNone) {
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
	}
	imgui.EndChild()

	imgui.Spacing()
	if imgui.BeginChildStrV("AnimListCard", imgui.Vec2{X: 0, Y: 160}, true, imgui.WindowFlagsNone) {
		imgui.TextDisabled("Animation list")
		for _, animName := range g.animationNames() {
			if imgui.Button("Select: " + animName) {
				g.setActiveAnimation(animName, true)
				g.uiVariables.renameAnimationName = animName
			}
		}
	}
	imgui.EndChild()

	if g.ActiveAnimation() == nil {
		return
	}

	if g.uiVariables.renameAnimationName == "" {
		g.uiVariables.renameAnimationName = g.ActiveAnimation().Name
	}

	imgui.Spacing()
	if imgui.BeginChildStrV("AnimRenameCard", imgui.Vec2{X: 0, Y: 86}, true, imgui.WindowFlagsNone) {
		_ = imguiInputText("Rename Active Animation", &g.uiVariables.renameAnimationName)
		if imgui.Button("Apply Animation Rename") {
			g.renameActiveAnimation(g.uiVariables.renameAnimationName)
		}
	}
	imgui.EndChild()
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
	imgui.SeparatorText("Timeline")
	if g.ActiveAnimation() == nil {
		imgui.TextDisabled("No active animation")
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

	if imgui.BeginChildStrV("FrameNavCard", imgui.Vec2{X: 0, Y: 136}, true, imgui.WindowFlagsNone) {
		imgui.Text(fmt.Sprintf("Frame %d/%d", player.FrameIndex+1, frameCount))

		if imgui.Button("Prev") {
			if player.FrameIndex > 0 {
				player.FrameIndex--
				player.FrameTimeLeft = g.ActiveAnimation().FrameData[player.FrameIndex].Duration
			}
		}
		imgui.SameLine()
		if imgui.Button("Next") {
			if player.FrameIndex < frameCount-1 {
				player.FrameIndex++
				player.FrameTimeLeft = g.ActiveAnimation().FrameData[player.FrameIndex].Duration
			}
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
	}
	imgui.EndChild()

	imgui.Spacing()
	if imgui.BeginChildStrV("FrameOpsCard", imgui.Vec2{X: 0, Y: 108}, true, imgui.WindowFlagsNone) {
		if imgui.Button("Add Frame by Image") {
			g.AddImageToAnimation()
		}
		imgui.SameLine()
		if imgui.Button("Duplicate") {
			g.duplicateLastFrameData()
		}
		imgui.SameLine()
		if imgui.Button("Remove") {
			g.uiVariables.frameDataIndex = player.FrameIndex
			g.removeFrame()
		}
	}
	imgui.EndChild()

	currentFrame := &g.ActiveAnimation().FrameData[player.FrameIndex]
	if g.uiVariables.frameDurationInput == "" {
		g.uiVariables.frameDurationInput = fmt.Sprintf("%d", currentFrame.Duration)
	}

	imgui.Spacing()
	if imgui.BeginChildStrV("FramePropsCard", imgui.Vec2{X: 0, Y: 136}, true, imgui.WindowFlagsNone) {
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
	}
	imgui.EndChild()

	imgui.Spacing()
	imgui.SeparatorText("Boxes")
	if imgui.Button("Add Collision") {
		g.uiVariables.boxDropdownTypeIndex = int(collision.Collision)
		g.addBox()
	}
	imgui.SameLine()
	if imgui.Button("Add Hurt") {
		g.uiVariables.boxDropdownTypeIndex = int(collision.Hurt)
		g.addBox()
	}
	imgui.SameLine()
	if imgui.Button("Add Hit") {
		g.uiVariables.boxDropdownTypeIndex = int(collision.Hit)
		g.addBox()
	}
}

func (g *Game) imguiLogPanel() {
	imgui.SeparatorText("Activity")
	imgui.TextDisabled("Recent operations")
	imgui.Separator()
	imgui.Text(g.uiVariables.logBuf)
}
