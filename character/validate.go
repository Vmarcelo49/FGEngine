package character

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

// RequiredStates is the full §6.7 required set: every state except the
// Specials family (implemented per character, absent = fizzled intent).
var RequiredStates = []string{
	"idle", "fall", "landing",
	"hurt", "crouchHurt", "airHurt",
	"blockHit", "crouchBlock", "airBlock",
	"knockdown", "getup", "ko", "win",
	"1", "2", "3", "4", "6", "44", "66", "7", "8", "9",
	"A", "B", "C", "D",
}

// ReactionStates must define loopFrames (the stun-hold pose) and carry no
// cancelTypes (§6.8 rules 7–8, §7.6).
var ReactionStates = []string{
	"hurt", "crouchHurt", "airHurt",
	"blockHit", "crouchBlock", "airBlock",
}

// IsSpecialMotion reports whether name is an engine-standard special motion
// (7 motions × A–D, §6.7).
func IsSpecialMotion(name string) bool {
	if len(name) < 3 {
		return false
	}
	motion, button := name[:len(name)-1], name[len(name)-1:]
	switch motion {
	case "236", "214", "623", "423", "22", "246", "642":
	default:
		return false
	}
	switch button {
	case "A", "B", "C", "D":
		return true
	}
	return false
}

// Validate enforces the §6.8 load rules with file-agnostic, field-named
// errors. Sprite paths must already be resolved (relative → absolute)
// before calling, so rule 5 checks them on disk directly.
// The returned errors are sorted for stable output.
func (c *Character) Validate() []error {
	var errs []error
	fail := func(format string, args ...any) {
		errs = append(errs, fmt.Errorf(format, args...))
	}

	if c.Name == "" {
		fail("name: missing")
	}
	if c.Properties.MaxHP <= 0 {
		fail("properties.maxHP: must be > 0, got %d", c.Properties.MaxHP)
	}
	if c.StateMachine == nil || c.StateMachine.AnimPlayer == nil {
		fail("stateMachine.activeAnim: missing")
		return errs
	}
	anims := c.StateMachine.AnimPlayer.Animations

	for _, name := range RequiredStates {
		if anims[name] == nil {
			fail("animations: required state %q missing", name)
		}
	}

	names := make([]string, 0, len(anims))
	for name := range anims {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		anim := anims[name]
		if anim == nil {
			fail("animations.%s: null entry", name)
			continue
		}
		if len(anim.FrameData) == 0 {
			fail("animations.%s: framedata empty", name)
			continue
		}
		isReaction := slices.Contains(ReactionStates, name)
		if isReaction && anim.LoopFrames == nil {
			fail("animations.%s: reaction state must define loopFrames", name)
		}
		for i := range anim.FrameData {
			fd := &anim.FrameData[i]
			if fd.Duration <= 0 {
				fail("animations.%s.framedata[%d]: duration must be >= 1, got %d", name, i, fd.Duration)
			}
			if fd.SpriteIndex < 0 || fd.SpriteIndex >= len(anim.Sprites) {
				fail("animations.%s.framedata[%d]: spriteIndex %d out of range (0..%d)", name, i, fd.SpriteIndex, len(anim.Sprites)-1)
			}
			for _, ct := range fd.CancelTypes {
				if ct != "any" && anims[ct] == nil {
					fail("animations.%s.framedata[%d]: cancelTypes references unknown animation %q", name, i, ct)
				}
			}
			if fd.AnimationSwitch != "" && anims[fd.AnimationSwitch] == nil {
				fail("animations.%s.framedata[%d]: animationSwitch references unknown animation %q", name, i, fd.AnimationSwitch)
			}
			if isReaction && len(fd.CancelTypes) > 0 {
				fail("animations.%s.framedata[%d]: reaction states must not carry cancelTypes", name, i)
			}
		}
		for i, sprite := range anim.Sprites {
			if sprite == nil {
				fail("animations.%s.sprites[%d]: null entry", name, i)
				continue
			}
			if sprite.ImagePath == "" {
				fail("animations.%s.sprites[%d]: imgPath missing", name, i)
			} else if _, err := os.Stat(sprite.ImagePath); err != nil {
				fail("animations.%s.sprites[%d]: imgPath not on disk: %s", name, i, sprite.ImagePath)
			}
		}
		if lf := anim.LoopFrames; lf != nil {
			if lf.Start < 0 || lf.End >= len(anim.FrameData) || lf.Start > lf.End {
				fail("animations.%s: loopFrames {%d,%d} out of frame range 0..%d", name, lf.Start, lf.End, len(anim.FrameData)-1)
			}
		}
	}

	slices.SortFunc(errs, func(a, b error) int {
		return strings.Compare(a.Error(), b.Error())
	})
	return errs
}
