package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fgengine/animation"
	"fgengine/character"
	"fgengine/types"

	"gopkg.in/yaml.v3"
)

func (ed *CharacterEditor) createNewCharacter(name string) {
	if strings.TrimSpace(name) == "" {
		name = "NewCharacter"
	}

	idleFrame := animation.FrameData{
		Duration:    6,
		SpriteIndex: 0,
		Boxes: map[types.BoxType][]types.Rect{
			types.Collision: {{X: 12, Y: 8, W: 40, H: 56}},
			types.Hurt:      {{X: 12, Y: 8, W: 40, H: 56}},
		},
		CancelTypes: []string{"any"},
	}

	idleAnim := &animation.Animation{
		Name:      "idle",
		Sprites:   []*animation.Sprite{},
		FrameData: []animation.FrameData{idleFrame},
	}
	ed.normalizeAnimationSprites(idleAnim)

	player := &animation.AnimationPlayer{Animations: map[string]*animation.Animation{"idle": idleAnim}}
	player.SetAnimation("idle")

	ed.char = &character.Character{
		Name: name,
		StateMachine: &animation.StateMachine{
			AnimPlayer: player,
		},
	}

	ed.setActiveAnimation("idle")
	ed.renameCharacterTo = name
	ed.newAnimationName = defaultNewAnimationName
	ed.renameAnimationTo = "idle"
	ed.selectedFrame = 0
	ed.selectedBoxType = types.Collision
	ed.selectedBoxIndex = 0
	ed.cancelTypes = strings.Join(idleFrame.CancelTypes, ",")
	ed.statusLine = "Created new character"
}

func (ed *CharacterEditor) loadCharacterFromPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path cannot be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	loaded := &character.Character{}
	if err := yaml.Unmarshal(data, loaded); err != nil {
		return fmt.Errorf("yaml parse failed: %w", err)
	}

	if loaded.StateMachine == nil || loaded.StateMachine.AnimPlayer == nil || loaded.StateMachine.AnimPlayer.Animations == nil {
		return fmt.Errorf("missing stateMachine.activeAnim.animations")
	}

	for name, anim := range loaded.StateMachine.AnimPlayer.Animations {
		if anim == nil {
			continue
		}
		anim.Name = name
		for _, spr := range anim.Sprites {
			if spr == nil || spr.ImagePath == "" {
				continue
			}
			if filepath.IsAbs(spr.ImagePath) {
				continue
			}
			spr.ImagePath = filepath.Clean(filepath.Join(filepath.Dir(path), spr.ImagePath))
		}
		ed.normalizeAnimationSprites(anim)
	}

	ed.char = loaded
	if _, ok := loaded.StateMachine.AnimPlayer.Animations["idle"]; ok {
		ed.setActiveAnimation("idle")
	} else {
		names := ed.animationNames()
		if len(names) > 0 {
			ed.setActiveAnimation(names[0])
		}
	}

	ed.savePath = path
	if ed.char != nil {
		ed.renameCharacterTo = ed.char.Name
	}
	ed.newAnimationName = defaultNewAnimationName
	ed.clearDirty()
	ed.exitAfterSave = false
	ed.showExitWindow = false
	ed.ignoreWindowClose = false
	return nil
}

func (ed *CharacterEditor) saveCharacterToPath(path string) error {
	if ed.char == nil {
		return fmt.Errorf("there is no character to save")
	}
	if ed.char.StateMachine == nil || ed.char.StateMachine.AnimPlayer == nil || ed.char.StateMachine.AnimPlayer.Animations == nil {
		return fmt.Errorf("character is missing stateMachine.activeAnim.animations")
	}
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed creating directories: %w", err)
	}
	// make loopframes nil if invalid or empty
	for i, anim := range ed.char.StateMachine.AnimPlayer.Animations {
		if anim == nil {
			continue
		}
		loopFrames := anim.LoopFrames
		if loopFrames == nil {
			continue
		}

		totalDuration := ed.totalAnimationDuration(anim)
		invalidLoop := loopFrames.Start == loopFrames.End || (loopFrames.Start == 0 && loopFrames.End == 0)
		if !invalidLoop && (loopFrames.Start > totalDuration || loopFrames.End > totalDuration) {
			invalidLoop = true
			ed.statusLine = fmt.Sprintf("Warning: animation '%s' has invalid loop frames and they were removed", anim.Name)
		}

		if invalidLoop {
			ed.char.StateMachine.AnimPlayer.Animations[i].LoopFrames = nil
		}
	}

	originalPaths := make(map[*animation.Sprite]string)
	for _, anim := range ed.char.StateMachine.AnimPlayer.Animations {
		if anim == nil {
			continue
		}
		for _, spr := range anim.Sprites {
			if spr == nil || spr.ImagePath == "" {
				continue
			}
			if _, seen := originalPaths[spr]; seen {
				continue
			}
			originalPaths[spr] = spr.ImagePath
			if filepath.IsAbs(spr.ImagePath) {
				rel, err := filepath.Rel(filepath.Dir(path), spr.ImagePath)
				if err == nil {
					spr.ImagePath = rel
				}
			}
		}
	}

	defer func() {
		for spr, p := range originalPaths {
			spr.ImagePath = p
		}
	}()

	out, err := yaml.Marshal(ed.char)
	if err != nil {
		return fmt.Errorf("yaml marshal failed: %w", err)
	}

	if err := os.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	return nil
}

func (ed *CharacterEditor) normalizeAnimationSprites(anim *animation.Animation) {
	if anim == nil {
		return
	}

	if len(anim.Sprites) == 0 {
		for i := range anim.FrameData {
			anim.FrameData[i].SpriteIndex = 0
		}
		return
	}

	oldToNew := make(map[int]int, len(anim.Sprites))
	cleaned := make([]*animation.Sprite, 0, len(anim.Sprites))
	for oldIndex, spr := range anim.Sprites {
		if spr == nil || strings.TrimSpace(spr.ImagePath) == "" {
			continue
		}
		oldToNew[oldIndex] = len(cleaned)
		cleaned = append(cleaned, spr)
	}

	anim.Sprites = cleaned

	if len(anim.Sprites) == 0 {
		for i := range anim.FrameData {
			anim.FrameData[i].SpriteIndex = 0
		}
		return
	}

	for i := range anim.FrameData {
		oldIndex := anim.FrameData[i].SpriteIndex
		if mappedIndex, ok := oldToNew[oldIndex]; ok {
			anim.FrameData[i].SpriteIndex = mappedIndex
			continue
		}
		anim.FrameData[i].SpriteIndex = 0
	}
}
