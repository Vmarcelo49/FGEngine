package editor

import (
	"fmt"

	"fgengine/animation"

	imgui "github.com/gabstv/cimgui-go"
)

func (ed *CharacterEditor) markDirty() {
	ed.dirty = true
}

func (ed *CharacterEditor) clearDirty() {
	ed.dirty = false
}

func (ed *CharacterEditor) requestExit() {
	if ed.exitEditor {
		return
	}
	ed.ignoreWindowClose = false
	if ed.dirty {
		ed.showExitWindow = true
		return
	}
	ed.exitEditor = true
}

func (ed *CharacterEditor) drawImportImagesAsAnimationWindow() {
	if !ed.showImportWindow {
		return
	}

	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(ed.width) / 2, Y: float32(ed.height) / 2}, imgui.CondAppearing, imgui.Vec2{X: 0.5, Y: 0.5})
	flags := imgui.WindowFlags(imgui.WindowFlagsAlwaysAutoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoMove | imgui.WindowFlagsNoSavedSettings)
	if !imgui.BeginV("Import Images", &ed.showImportWindow, flags) {
		imgui.End()
		return
	}
	defer imgui.End()

	imgui.Text("Deseja importar essas imagens como animacao?")
	imgui.Separator()

	if imgui.Button("Sim") {
		ed.importImagesAsAnimation(ed.pendingImportPaths)
		ed.pendingImportPaths = nil
		ed.showImportWindow = false
	}
	imgui.SameLine()
	if imgui.Button("Nao") {
		ed.addImagesToActiveAnimation(ed.pendingImportPaths)
		ed.statusLine = fmt.Sprintf("Added %d image(s)", len(ed.pendingImportPaths))
		ed.pendingImportPaths = nil
		ed.showImportWindow = false
	}
}

func (ed *CharacterEditor) drawNotificationWindow(message string) {
	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(ed.width) / 2, Y: float32(ed.height) / 2}, imgui.CondAppearing, imgui.Vec2{X: 0.5, Y: 0.5})
	flags := imgui.WindowFlags(imgui.WindowFlagsAlwaysAutoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoMove | imgui.WindowFlagsNoSavedSettings)
	if !imgui.BeginV("Notification", nil, flags) {
		imgui.End()
		return
	}
	defer imgui.End()

	imgui.Text(message)
	if imgui.Button("OK") {
		// Just close the window on OK
	}
}

func (ed *CharacterEditor) drawUnsavedChangesWindow() {
	if !ed.showExitWindow {
		return
	}

	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(ed.width) / 2, Y: float32(ed.height) / 2}, imgui.CondAppearing, imgui.Vec2{X: 0.5, Y: 0.5})
	flags := imgui.WindowFlags(imgui.WindowFlagsAlwaysAutoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoMove | imgui.WindowFlagsNoSavedSettings)
	if !imgui.BeginV("Unsaved Changes", &ed.showExitWindow, flags) {
		imgui.End()
		return
	}
	defer imgui.End()

	imgui.Text("O projeto nao foi salvo.")
	imgui.Text("Deseja salvar antes de fechar?")
	imgui.Separator()

	if imgui.Button("Salvar") {
		ed.showSaveWindow = true
		ed.exitAfterSave = true
		ed.showExitWindow = false
		ed.ignoreWindowClose = true
	}
	imgui.SameLine()
	if imgui.Button("Descartar") {
		ed.exitAfterSave = false
		ed.exitEditor = true
		ed.showExitWindow = false
		ed.ignoreWindowClose = false
	}
	imgui.SameLine()
	if imgui.Button("Cancelar") {
		ed.exitAfterSave = false
		ed.showExitWindow = false
		ed.ignoreWindowClose = false
	}
}

func (ed *CharacterEditor) importImagesAsAnimation(paths []string) {
	anim := ed.activeAnimation()
	if anim == nil {
		return
	}

	cleanPaths := dedupeNonEmpty(paths)
	if len(cleanPaths) == 0 {
		return
	}

	sprites := make([]*animation.Sprite, 0, len(cleanPaths))
	frames := make([]animation.FrameData, 0, len(cleanPaths))
	for idx, path := range cleanPaths {
		sprites = append(sprites, &animation.Sprite{ImagePath: path})
		frames = append(frames, animation.FrameData{
			Duration:    6,
			SpriteIndex: idx,
		})
	}

	anim.Sprites = sprites
	anim.FrameData = frames
	ed.applyDefaultIdleAnchorToAnimationSprites(anim)
	anim.TotalDuration = 0
	ed.setActiveAnimation(anim.Name)
	ed.jumpToFrame(0)
	ed.selectedFrame = 0
	ed.statusLine = fmt.Sprintf("Imported %d image(s) as animation", len(cleanPaths))
	ed.markDirty()
}

func (ed *CharacterEditor) applyAnchorNextSprites() {
	anim := ed.activeAnimation()
	fd := ed.currentFrameData()
	if anim == nil || fd == nil {
		return
	}
	if fd.SpriteIndex < 0 || fd.SpriteIndex >= len(anim.Sprites) {
		return
	}

	source := anim.Sprites[fd.SpriteIndex]
	if source == nil {
		return
	}

	for idx := fd.SpriteIndex + 1; idx < len(anim.Sprites); idx++ {
		if anim.Sprites[idx] == nil {
			continue
		}
		anim.Sprites[idx].Anchor = source.Anchor
	}

	ed.statusLine = "Anchor propagated to following sprites"
	ed.markDirty()
}

func (ed *CharacterEditor) drawChangeCharacterNameWindow() {
	if !ed.showChangeCharacterNameWindow {
		return
	}

	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(ed.width) / 2, Y: float32(ed.height) / 2}, imgui.CondAppearing, imgui.Vec2{X: 0.5, Y: 0.5})
	flags := imgui.WindowFlags(imgui.WindowFlagsAlwaysAutoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoMove | imgui.WindowFlagsNoSavedSettings)
	if !imgui.BeginV("Change Character Name", &ed.showChangeCharacterNameWindow, flags) {
		imgui.End()
		return
	}
	defer imgui.End()

	imgui.InputTextWithHint("Name", "new character name", &ed.renameCharacterTo, 0, nil)
	if imgui.Button("Apply") {
		ed.renameCharacter(ed.renameCharacterTo)
		ed.showChangeCharacterNameWindow = false
	}
	imgui.SameLine()
	if imgui.Button("Cancel") {
		ed.showChangeCharacterNameWindow = false
	}
}

func (ed *CharacterEditor) drawRenameActiveAnimationWindow() {
	if !ed.showRenameAnimationWindow {
		return
	}

	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(ed.width) / 2, Y: float32(ed.height) / 2}, imgui.CondAppearing, imgui.Vec2{X: 0.5, Y: 0.5})
	flags := imgui.WindowFlags(imgui.WindowFlagsAlwaysAutoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoMove | imgui.WindowFlagsNoSavedSettings)
	if !imgui.BeginV("Rename Active Animation", &ed.showRenameAnimationWindow, flags) {
		imgui.End()
		return
	}
	defer imgui.End()

	if ed.activeAnimation() == nil {
		imgui.Text("No active animation.")
		if imgui.Button("Close") {
			ed.showRenameAnimationWindow = false
		}
		return
	}

	imgui.InputTextWithHint("Name", "new animation name", &ed.renameAnimationTo, 0, nil)
	if imgui.Button("Apply") {
		ed.renameActiveAnimation(ed.renameAnimationTo)
		ed.showRenameAnimationWindow = false
	}
	imgui.SameLine()
	if imgui.Button("Cancel") {
		ed.showRenameAnimationWindow = false
	}
}

func (ed *CharacterEditor) drawDeleteAnimationWindow() {
	if !ed.showDeleteAnimationWindow {
		return
	}

	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(ed.width) / 2, Y: float32(ed.height) / 2}, imgui.CondAppearing, imgui.Vec2{X: 0.5, Y: 0.5})
	flags := imgui.WindowFlags(imgui.WindowFlagsAlwaysAutoResize | imgui.WindowFlagsNoCollapse | imgui.WindowFlagsNoMove | imgui.WindowFlagsNoSavedSettings)
	if !imgui.BeginV("Delete This Animation", &ed.showDeleteAnimationWindow, flags) {
		imgui.End()
		return
	}
	defer imgui.End()

	animName := ed.activeAnimationName
	if animName == "" && ed.activeAnimation() != nil {
		animName = ed.activeAnimation().Name
	}
	if animName == "" {
		imgui.Text("No active animation.")
	} else {
		imgui.Text(fmt.Sprintf("Delete animation '%s'?", animName))
	}
	imgui.Separator()

	if imgui.Button("Delete") {
		ed.deleteActiveAnimation()
		ed.showDeleteAnimationWindow = false
	}
	imgui.SameLine()
	if imgui.Button("Cancel") {
		ed.showDeleteAnimationWindow = false
	}
}
