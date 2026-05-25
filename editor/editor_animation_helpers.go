package editor

import (
	"fmt"
	"sort"
	"strings"

	"fgengine/animation"
	"fgengine/types"
)

const defaultNewAnimationName = "new_animation"

func (ed *CharacterEditor) player() *animation.AnimationPlayer {
	if ed.char == nil || ed.char.StateMachine == nil {
		return nil
	}
	return ed.char.StateMachine.AnimPlayer
}

func (ed *CharacterEditor) activeAnimation() *animation.Animation {
	p := ed.player()
	if p == nil {
		return nil
	}
	return p.ActiveAnimation
}

func (ed *CharacterEditor) currentFrameData() *animation.FrameData {
	anim := ed.activeAnimation()
	if anim == nil || len(anim.FrameData) == 0 {
		return nil
	}
	if ed.selectedFrame < 0 {
		ed.selectedFrame = 0
	}
	if ed.selectedFrame >= len(anim.FrameData) {
		ed.selectedFrame = len(anim.FrameData) - 1
	}
	return &anim.FrameData[ed.selectedFrame]
}

func (ed *CharacterEditor) animationNames() []string {
	p := ed.player()
	if p == nil || p.Animations == nil {
		return nil
	}
	names := make([]string, 0, len(p.Animations))
	for name := range p.Animations {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (ed *CharacterEditor) activeAnimationIndex(names []string) int32 {
	if ed.activeAnimationName == "" {
		p := ed.player()
		if p != nil {
			ed.activeAnimationName = p.ActiveAnimationName()
		}
	}
	for i, name := range names {
		if name == ed.activeAnimationName {
			return int32(i)
		}
	}
	return 0
}

func (ed *CharacterEditor) setActiveAnimation(name string) {
	p := ed.player()
	if p == nil {
		return
	}
	p.SetAnimation(name)
	ed.activeAnimationName = name
	if ed.char != nil {
		ed.renameCharacterTo = ed.char.Name
	}
	ed.renameAnimationTo = name
	ed.newAnimationName = defaultNewAnimationName
	ed.selectedFrame = 0
	ed.jumpToFrame(0)
	fd := ed.currentFrameData()
	if fd != nil {
		ed.cancelTypes = strings.Join(fd.CancelTypes, ",")
	}
}

func (ed *CharacterEditor) renameCharacter(newName string) {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		ed.statusLine = "Character name cannot be empty"
		return
	}
	if ed.char == nil {
		return
	}
	if ed.char.Name == newName {
		ed.statusLine = "Character name is unchanged"
		return
	}

	ed.char.Name = newName
	ed.renameCharacterTo = newName
	ed.statusLine = "Character renamed"
	ed.markDirty()
}

func (ed *CharacterEditor) createEmptyAnimation(newName string) {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		ed.statusLine = "New animation name cannot be empty"
		return
	}

	p := ed.player()
	if p == nil || p.Animations == nil {
		ed.statusLine = "No character loaded"
		return
	}
	if _, exists := p.Animations[newName]; exists {
		ed.statusLine = "Animation name already exists"
		return
	}

	empty := &animation.Animation{
		Name:      newName,
		Sprites:   []*animation.Sprite{},
		FrameData: []animation.FrameData{{Duration: 6, SpriteIndex: 0}},
	}

	p.Animations[newName] = empty
	ed.setActiveAnimation(newName)
	ed.applyDefaultIdleAnchorToAnimationSprites(empty)
	ed.statusLine = "Empty animation created"
	ed.markDirty()
}

func (ed *CharacterEditor) renameActiveAnimation(newName string) {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		ed.statusLine = "Animation name cannot be empty"
		return
	}

	p := ed.player()
	if p == nil || p.Animations == nil {
		return
	}

	if _, exists := p.Animations[newName]; exists {
		ed.statusLine = "Animation name already exists"
		return
	}

	oldName := p.ActiveAnimationName()
	anim, ok := p.Animations[oldName]
	if !ok {
		ed.statusLine = "No active animation to rename"
		return
	}

	delete(p.Animations, oldName)
	p.Animations[newName] = anim
	anim.Name = newName
	p.SetAnimation(newName)
	ed.activeAnimationName = newName
	ed.statusLine = "Animation renamed"
	ed.markDirty()
}

func (ed *CharacterEditor) deleteActiveAnimation() {
	p := ed.player()
	if p == nil || p.Animations == nil {
		ed.statusLine = "No active animation to delete"
		return
	}

	oldName := p.ActiveAnimationName()
	if oldName == "" || oldName == "none" {
		ed.statusLine = "No active animation to delete"
		return
	}
	if _, ok := p.Animations[oldName]; !ok {
		ed.statusLine = "No active animation to delete"
		return
	}

	delete(p.Animations, oldName)
	names := ed.animationNames()

	if len(names) > 0 {
		ed.setActiveAnimation(names[0])
		ed.statusLine = fmt.Sprintf("Animation '%s' deleted", oldName)
		ed.markDirty()
		return
	}

	p.ActiveAnimation = nil
	p.FrameIndex = 0
	p.FrameTimeLeft = 0
	ed.activeAnimationName = ""
	ed.renameAnimationTo = ""
	ed.selectedFrame = 0
	ed.cancelTypes = ""
	ed.statusLine = fmt.Sprintf("Animation '%s' deleted. Character has no animations.", oldName)
	ed.markDirty()
}

func (ed *CharacterEditor) frameCount() int {
	anim := ed.activeAnimation()
	if anim == nil {
		return 0
	}
	return len(anim.FrameData)
}

func (ed *CharacterEditor) jumpToFrame(idx int) {
	p := ed.player()
	anim := ed.activeAnimation()
	if p == nil || anim == nil || len(anim.FrameData) == 0 {
		return
	}
	if idx < 0 {
		idx = 0
	}
	if idx >= len(anim.FrameData) {
		idx = len(anim.FrameData) - 1
	}
	p.FrameIndex = idx
	p.FrameTimeLeft = anim.FrameData[idx].Duration
	ed.selectedFrame = idx
	ed.cancelTypes = strings.Join(anim.FrameData[idx].CancelTypes, ",")
}

func (ed *CharacterEditor) totalAnimationDuration(anim *animation.Animation) int {
	total := 0
	for _, fd := range anim.FrameData {
		total += fd.Duration
	}
	return total
}

func (ed *CharacterEditor) elapsedFrames(anim *animation.Animation, frameIndex int, frameTimeLeft int) int {
	elapsed := 0
	for i := 0; i < frameIndex && i < len(anim.FrameData); i++ {
		elapsed += anim.FrameData[i].Duration
	}
	if frameIndex >= 0 && frameIndex < len(anim.FrameData) {
		elapsed += anim.FrameData[frameIndex].Duration - frameTimeLeft
	}
	if elapsed < 0 {
		return 0
	}
	return elapsed
}

func splitAndTrim(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trim := strings.TrimSpace(p)
		if trim != "" {
			out = append(out, trim)
		}
	}
	return out
}

func deleteRectAt(rects []types.Rect, idx int) []types.Rect {
	if idx < 0 || idx >= len(rects) {
		return rects
	}
	return append(rects[:idx], rects[idx+1:]...)
}

func cloneFrameData(src animation.FrameData) animation.FrameData {
	cloned := src

	if src.CancelTypes != nil {
		cloned.CancelTypes = append([]string(nil), src.CancelTypes...)
	}

	if src.Boxes != nil {
		cloned.Boxes = make(map[types.BoxType][]types.Rect, len(src.Boxes))
		for boxType, boxes := range src.Boxes {
			cloned.Boxes[boxType] = append([]types.Rect(nil), boxes...)
		}
	}

	return cloned
}

func (ed *CharacterEditor) copyCurrentFrameDataToFollowingFrames() {
	anim := ed.activeAnimation()
	fd := ed.currentFrameData()
	if anim == nil || fd == nil {
		return
	}
	if ed.selectedFrame < 0 || ed.selectedFrame >= len(anim.FrameData)-1 {
		ed.statusLine = "No following frames to copy to"
		return
	}

	for i := ed.selectedFrame + 1; i < len(anim.FrameData); i++ {
		anim.FrameData[i] = cloneFrameData(*fd)
	}

	ed.statusLine = "Current framedata copied to following frames"
	ed.markDirty()
}

func (ed *CharacterEditor) idleAnchor() (types.Vector2, bool) {
	p := ed.player()
	if p == nil || p.Animations == nil {
		return types.Vector2{}, false
	}

	idle, ok := p.Animations["idle"]
	if !ok || idle == nil || len(idle.Sprites) == 0 {
		return types.Vector2{}, false
	}

	for _, spr := range idle.Sprites {
		if spr != nil {
			return spr.Anchor, true
		}
	}

	return types.Vector2{}, false
}

func (ed *CharacterEditor) applyDefaultIdleAnchorToAnimationSprites(anim *animation.Animation) {
	if anim == nil {
		return
	}
	anchor, ok := ed.idleAnchor()
	if !ok {
		return
	}

	for _, spr := range anim.Sprites {
		if spr == nil {
			continue
		}
		spr.Anchor = anchor
	}
}
