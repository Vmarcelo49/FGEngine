package editor

import (
	"fmt"
	"path/filepath"
	"strings"

	"fgengine/animation"
	"fgengine/types"

	imgui "github.com/gabstv/cimgui-go"
)

func (ed *CharacterEditor) drawTopMenuBar() {
	if !imgui.BeginMainMenuBar() {
		return
	}
	defer imgui.EndMainMenuBar()

	if imgui.BeginMenu("File") {
		if imgui.MenuItemBoolV("Create New Character", "", false, true) {
			ed.showCreateWindow = true
		}
		if imgui.MenuItemBoolV("Load Character", "", false, true) {
			ed.showLoadWindow = true
		}
		if imgui.MenuItemBoolV("Save Character", "", false, true) {
			ed.showSaveWindow = true
		}
		if imgui.MenuItemBoolV("Exit", "", false, true) {
			ed.requestExit()
		}
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Options") {
		if imgui.BeginMenu("Scaling") {
			for scale := 1; scale <= 6; scale++ {
				if imgui.MenuItemBoolV(fmt.Sprintf("%dx", scale), "", ed.previewScale == scale, true) {
					ed.previewScale = scale
				}
			}
			imgui.EndMenu()
		}
		imgui.EndMenu()
	}

	if ed.char != nil {
		editLabel := strings.TrimSpace(ed.char.Name)
		if editLabel == "" {
			editLabel = "Character"
		}
		if imgui.BeginMenu("Edit " + editLabel) {
			if imgui.MenuItemBoolV("Change Character Name", "", false, true) {
				ed.showChangeCharacterNameWindow = true
				if ed.renameCharacterTo == "" {
					ed.renameCharacterTo = ed.char.Name
				}
			}
			if imgui.MenuItemBoolV("Rename Active Animation", "", false, ed.activeAnimation() != nil) {
				ed.showRenameAnimationWindow = true
				if ed.renameAnimationTo == "" {
					ed.renameAnimationTo = ed.activeAnimationName
				}
			}
			if imgui.MenuItemBoolV("Delete This Animation", "", false, ed.activeAnimation() != nil) {
				ed.showDeleteAnimationWindow = true
			}
			imgui.EndMenu()
		}
	}
}

func (ed *CharacterEditor) drawCharacterWindow() {
	open := true
	if !imgui.BeginV("Character && Animation", &open, imgui.WindowFlags(0)) {
		imgui.End()
		return
	}
	defer imgui.End()

	if ed.char == nil {
		imgui.Text("No character loaded.")
		return
	}

	names := ed.animationNames()
	if len(names) == 0 {
		imgui.Text("No animations available.")
	} else {
		current := ed.activeAnimationIndex(names)
		if imgui.ComboStrarrV("Active Animation", &current, names, int32(len(names)), -1) {
			ed.setActiveAnimation(names[current])
		}
	}

	if ed.newAnimationName == "" {
		ed.newAnimationName = defaultNewAnimationName
	}
	imgui.InputTextWithHint("New Animation", "new animation name", &ed.newAnimationName, 0, nil)
	if imgui.Button("Create Empty Animation") {
		ed.createEmptyAnimation(ed.newAnimationName)
	}

	imgui.Separator()
	imgui.Text("Edit names and delete animation from the top menu")
	imgui.Text(fmt.Sprintf("Current Animation: %s", ed.activeAnimationName))
	imgui.Text(fmt.Sprintf("Frames: %d", ed.frameCount()))
	imgui.Text(ed.statusLine)
}

func (ed *CharacterEditor) drawAnimationPlayerWindow() {
	open := true
	if !imgui.BeginV("Animation Player", &open, imgui.WindowFlags(0)) {
		imgui.End()
		return
	}
	defer imgui.End()

	anim := ed.activeAnimation()
	if anim == nil {
		imgui.Text("No active animation.")
		return
	}

	if ed.paused {
		if imgui.Button("Resume") {
			ed.paused = false
		}
	} else {
		if imgui.Button("Pause") {
			ed.paused = true
		}
	}

	player := ed.player()
	if player != nil {
		imgui.SameLine()
		if player.ActiveAnimation.LoopFrames == nil { // TODO: should this be auto-created when creating an animation?
			player.ActiveAnimation.LoopFrames = &animation.LoopFrame{}
		}
		loopStartFrame := int32(player.ActiveAnimation.LoopFrames.Start)
		loopEndFrame := int32(player.ActiveAnimation.LoopFrames.End)
		if imgui.InputInt("Loop Start Frame", &loopStartFrame) {
			if loopStartFrame < 0 {
				loopStartFrame = 0
			}
			player.ActiveAnimation.LoopFrames.Start = int(loopStartFrame)
		}
		if imgui.InputInt("Loop End Frame", &loopEndFrame) {
			if loopEndFrame < 0 {
				loopEndFrame = 0
			}
			if loopEndFrame < loopStartFrame {
				loopEndFrame = loopStartFrame
			}
			player.ActiveAnimation.LoopFrames.End = int(loopEndFrame)
		}

		imgui.Text(fmt.Sprintf("Frame %d / %d", player.FrameIndex+1, len(anim.FrameData)))
		imgui.Text(fmt.Sprintf("Frame Time Left: %d", player.FrameTimeLeft))
	}

	total := ed.totalAnimationDuration(anim)
	if total > 0 && player != nil {
		elapsed := ed.elapsedFrames(anim, player.FrameIndex, player.FrameTimeLeft)
		imgui.ProgressBar(float32(elapsed) / float32(total))
	}

	imgui.SeparatorText("Timeline")
	for i := range anim.FrameData {
		selected := i == ed.selectedFrame
		label := fmt.Sprintf("Frame %d (dur %d)", i, anim.FrameData[i].Duration)
		if imgui.SelectableBoolV(label, selected, 0, imgui.Vec2{}) {
			ed.jumpToFrame(i)
		}
	}
}

func (ed *CharacterEditor) drawFrameDataWindow() {
	open := true
	if !imgui.BeginV("FrameData", &open, imgui.WindowFlags(0)) {
		imgui.End()
		return
	}
	defer imgui.End()

	fd := ed.currentFrameData()
	if fd == nil {
		imgui.Text("No frame selected.")
		return
	}

	dur := int32(fd.Duration)
	if imgui.InputInt("Duration", &dur) {
		if dur < 1 {
			dur = 1
		}
		fd.Duration = int(dur)
		ed.markDirty()
	}

	spr := int32(fd.SpriteIndex)
	if imgui.InputInt("Sprite Index", &spr) {
		if spr < 0 {
			spr = 0
		}
		fd.SpriteIndex = int(spr)
		ed.markDirty()
	}

	vx := float32(fd.IncVelocityX)
	if imgui.InputFloat("Velocity X", &vx) {
		fd.IncVelocityX = float64(vx)
		ed.markDirty()
	}

	vy := float32(fd.IncVelocityY)
	if imgui.InputFloat("Velocity Y", &vy) {
		fd.IncVelocityY = float64(vy)
		ed.markDirty()
	}

	if imgui.InputTextWithHint("Switch Animation", "optional", &fd.AnimationSwitch, 0, nil) {
		ed.markDirty()
	}

	if ed.cancelTypes == "" {
		ed.cancelTypes = strings.Join(fd.CancelTypes, ",")
	}
	if imgui.InputTextWithHint("Cancel Types", "jump,attack,dash", &ed.cancelTypes, 0, nil) {
		fd.CancelTypes = splitAndTrim(ed.cancelTypes)
		ed.markDirty()
	}

	priority := int32(fd.Priority)
	if imgui.InputInt("Priority", &priority) {
		fd.Priority = int(priority)
		ed.markDirty()
	}
	damage := int32(fd.Damage)
	if imgui.InputInt("Damage", &damage) {
		fd.Damage = int(damage)
		ed.markDirty()
	}
	hitstun := int32(fd.Hitstun)
	if imgui.InputInt("Hitstun", &hitstun) {
		fd.Hitstun = int(hitstun)
		ed.markDirty()
	}
	blockstun := int32(fd.Blockstun)
	if imgui.InputInt("Blockstun", &blockstun) {
		fd.Blockstun = int(blockstun)
		ed.markDirty()
	}
	pushback := int32(fd.Pushback)
	if imgui.InputInt("Pushback", &pushback) {
		fd.Pushback = int(pushback)
		ed.markDirty()
	}
	knockback := int32(fd.Knockback)
	if imgui.InputInt("Knockback", &knockback) {
		fd.Knockback = int(knockback)
		ed.markDirty()
	}
	knockup := int32(fd.Knockup)
	if imgui.InputInt("Knockup", &knockup) {
		fd.Knockup = int(knockup)
		ed.markDirty()
	}

	imgui.SeparatorText("Current Frame Sprite Anchor")
	anim := ed.activeAnimation()
	if anim == nil || fd.SpriteIndex < 0 || fd.SpriteIndex >= len(anim.Sprites) || anim.Sprites[fd.SpriteIndex] == nil {
		imgui.Text("No current frame sprite to edit anchor.")
	} else {
		spr := anim.Sprites[fd.SpriteIndex]
		ax := float32(spr.Anchor.X)
		if imgui.DragFloat("Anchor X", &ax) {
			spr.Anchor.X = float64(ax)
			ed.markDirty()
		}

		ay := float32(spr.Anchor.Y)
		if imgui.DragFloat("Anchor Y", &ay) {
			spr.Anchor.Y = float64(ay)
			ed.markDirty()
		}

		if imgui.Button("Apply Anchor to Next Sprites") {
			ed.applyAnchorNextSprites()
		}
	}

	if imgui.Button("Copy Current FrameData to Next Frames") {
		ed.copyCurrentFrameDataToFollowingFrames()
	}

	if imgui.Checkbox("Can Hard Knockdown", &fd.CanHardKnockdown) {
		ed.markDirty()
	}
	if imgui.Checkbox("Can Wall Bounce", &fd.CanWallBounce) {
		ed.markDirty()
	}
	if imgui.Checkbox("Can Ground Bounce", &fd.CanGroundBounce) {
		ed.markDirty()
	}
	if imgui.Checkbox("Can OTG", &fd.CanOTG) {
		ed.markDirty()
	}
	if imgui.Checkbox("Is Invincible", &fd.IsInvincible) {
		ed.markDirty()
	}
	if imgui.Checkbox("Has Armor", &fd.HasArmor) {
		ed.markDirty()
	}
}

func (ed *CharacterEditor) drawBoxEditorWindow() {
	open := true
	if !imgui.BeginV("Box Editor", &open, imgui.WindowFlags(0)) {
		imgui.End()
		return
	}
	defer imgui.End()

	fd := ed.currentFrameData()
	if fd == nil {
		imgui.Text("No frame selected.")
		return
	}

	if fd.Boxes == nil {
		fd.Boxes = make(map[types.BoxType][]types.Rect)
	}

	boxTypeNames := []string{"Collision", "Hit", "Hurt"}
	typeIndex := int32(ed.selectedBoxType)
	if imgui.ComboStrarrV("Box Type", &typeIndex, boxTypeNames, int32(len(boxTypeNames)), -1) {
		ed.selectedBoxType = types.BoxType(typeIndex)
		ed.selectedBoxIndex = 0
	}

	boxes := fd.Boxes[ed.selectedBoxType]
	imgui.SeparatorText("Boxes")
	for i, box := range boxes {
		selected := i == ed.selectedBoxIndex
		if imgui.SelectableBoolV(fmt.Sprintf("%d -> X%.1f Y%.1f W%.1f H%.1f", i, box.X, box.Y, box.W, box.H), selected, 0, imgui.Vec2{}) {
			ed.selectedBoxIndex = i
		}
	}

	if imgui.Button("Add New Box") {
		fd.Boxes[ed.selectedBoxType] = append(fd.Boxes[ed.selectedBoxType], types.Rect{W: 32, H: 32})
		ed.selectedBoxIndex = len(fd.Boxes[ed.selectedBoxType]) - 1
		boxes = fd.Boxes[ed.selectedBoxType]
		ed.markDirty()
	}

	if len(boxes) == 0 {
		imgui.Text("No boxes in selected type.")
		return
	}

	if ed.selectedBoxIndex < 0 {
		ed.selectedBoxIndex = 0
	}
	if ed.selectedBoxIndex >= len(boxes) {
		ed.selectedBoxIndex = len(boxes) - 1
	}

	b := &fd.Boxes[ed.selectedBoxType][ed.selectedBoxIndex]
	bx := float32(b.X)
	if imgui.InputFloat("X", &bx) {
		b.X = float64(bx)
		ed.markDirty()
	}
	by := float32(b.Y)
	if imgui.InputFloat("Y", &by) {
		b.Y = float64(by)
		ed.markDirty()
	}
	bw := float32(b.W)
	if imgui.InputFloat("W", &bw) {
		if bw < 0 {
			bw = 0
		}
		b.W = float64(bw)
		ed.markDirty()
	}
	bh := float32(b.H)
	if imgui.InputFloat("H", &bh) {
		if bh < 0 {
			bh = 0
		}
		b.H = float64(bh)
		ed.markDirty()
	}

	imgui.SeparatorText("Move To Type")
	targetTypeIndex := int32(ed.targetBoxType)
	if imgui.ComboStrarrV("Target Type", &targetTypeIndex, boxTypeNames, int32(len(boxTypeNames)), -1) {
		ed.targetBoxType = types.BoxType(targetTypeIndex)
	}

	if imgui.Button("Change Current Box Type") {
		box := fd.Boxes[ed.selectedBoxType][ed.selectedBoxIndex]
		fd.Boxes[ed.selectedBoxType] = deleteRectAt(fd.Boxes[ed.selectedBoxType], ed.selectedBoxIndex)
		fd.Boxes[ed.targetBoxType] = append(fd.Boxes[ed.targetBoxType], box)
		ed.selectedBoxType = ed.targetBoxType
		ed.selectedBoxIndex = len(fd.Boxes[ed.selectedBoxType]) - 1
		ed.markDirty()
	}

	if imgui.Button("Delete Selected Box") {
		fd.Boxes[ed.selectedBoxType] = deleteRectAt(fd.Boxes[ed.selectedBoxType], ed.selectedBoxIndex)
		if ed.selectedBoxIndex > 0 {
			ed.selectedBoxIndex--
		}
		ed.markDirty()
	}
}

func (ed *CharacterEditor) drawCreateCharacterWindow() {
	if !ed.showCreateWindow {
		return
	}

	if !imgui.BeginV("Create New Character", &ed.showCreateWindow, imgui.WindowFlags(0)) {
		imgui.End()
		return
	}
	defer imgui.End()

	imgui.InputTextWithHint("Name", "Character name", &ed.newCharacterName, 0, nil)
	if imgui.Button("Create") {
		ed.createNewCharacter(ed.newCharacterName)
		ed.markDirty()
		ed.showCreateWindow = false
	}
}

func (ed *CharacterEditor) drawLoadCharacterWindow() {
	if !ed.showLoadWindow {
		return
	}

	if !imgui.BeginV("Load Character", &ed.showLoadWindow, imgui.WindowFlags(0)) {
		imgui.End()
		return
	}
	defer imgui.End()

	imgui.InputTextWithHint("YAML Path", "./assets/characters/Name.yaml", &ed.loadPath, 0, nil)
	if imgui.Button("Browse...") {
		picked, err := ed.pickCharacterWithDialog()
		if err != nil {
			ed.statusLine = "Character picker failed: " + err.Error()
		} else if strings.TrimSpace(picked) != "" {
			ed.loadPath = picked
			if err := ed.loadCharacterFromPath(ed.loadPath); err != nil {
				ed.statusLine = "Load failed: " + err.Error()
			} else {
				ed.statusLine = "Character loaded"
				ed.showLoadWindow = false
			}
		}
	}
	imgui.SameLine()
	if imgui.Button("Load") {
		if err := ed.loadCharacterFromPath(ed.loadPath); err != nil {
			ed.statusLine = "Load failed: " + err.Error()
		} else {
			ed.statusLine = "Character loaded"
			ed.showLoadWindow = false
		}
	}
}

func (ed *CharacterEditor) drawSaveCharacterWindow() {
	if !ed.showSaveWindow {
		return
	}

	wasOpen := ed.showSaveWindow
	if !imgui.BeginV("Save Character", &ed.showSaveWindow, imgui.WindowFlags(0)) {
		imgui.End()
		if wasOpen && !ed.showSaveWindow {
			ed.exitAfterSave = false
		}
		return
	}
	defer imgui.End()

	if ed.savePath == "" && ed.char != nil {
		ed.savePath = filepath.Join("./assets/characters", ed.char.Name+".yaml")
	}

	imgui.InputTextWithHint("Save Path", "./assets/characters/Name.yaml", &ed.savePath, 0, nil)
	if imgui.Button("Save") {
		if err := ed.saveCharacterToPath(ed.savePath); err != nil {
			ed.statusLine = "Save failed: " + err.Error()
		} else {
			ed.statusLine = "Character saved"
			ed.showSaveWindow = false
			ed.clearDirty()
			ed.ignoreWindowClose = false
			if ed.exitAfterSave {
				ed.exitAfterSave = false
				ed.exitEditor = true
			}
		}
	}
}
