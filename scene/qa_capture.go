//go:build fgengineqa

package scene

import (
	"fgengine/config"
	"fgengine/constants"
	"fgengine/input"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type qaScenario struct {
	fighter string
	action  string
}

type qaCapture struct {
	outDir            string
	scenarios         []qaScenario
	index             int
	tick              int
	match             *GameplayScene
	expected          string
	motion            []input.GameInput
	attackSent        bool
	sawAnimation      bool
	framePending      bool
	impactReached     bool
	impactShotPending bool
	finishAfterDraw   bool
	afterImpact       int
	frameNumber       int
	failures          []string
	drawErr           error
}

// CaptureRosterFrames drives the real gameplay scene with scripted inputs.
// It writes each rendered frame and marks verified hits and blocks separately.
// The fgengineqa build tag keeps this driver out of the normal game build.
func CaptureRosterFrames(outDir string) error {
	for _, subdir := range []string{"frames", "impacts"} {
		if err := os.MkdirAll(filepath.Join(outDir, subdir), 0o755); err != nil {
			return err
		}
	}
	config.ActiveConfig.Language = "EN"
	game := &qaCapture{outDir: outDir}
	for _, fighter := range []string{"Brasa", "Muralha", "Iris", "Vulto"} {
		for _, action := range []string{"standing", "crouching", "airborne", "special", "guard"} {
			game.scenarios = append(game.scenarios, qaScenario{fighter: fighter, action: action})
		}
	}
	ebiten.SetWindowSize(int(constants.CameraWidth), int(constants.CameraHeight))
	ebiten.SetWindowTitle("FGEngine visual QA")
	return ebiten.RunGame(game)
}

func (q *qaCapture) startScenario() {
	scenario := q.scenarios[q.index]
	other := "Muralha"
	if scenario.fighter == other {
		other = "Brasa"
	}
	attacker, defender := scenario.fighter, other
	if scenario.action == "guard" {
		attacker, defender = other, scenario.fighter
	}
	q.match = MakeGameplayScene(attacker, defender, false).(*GameplayScene)
	q.match.gamestate.Characters[0].StateMachine.Position.X = 300
	q.match.gamestate.Characters[1].StateMachine.Position.X = 360
	q.tick = 0
	q.attackSent = false
	q.sawAnimation = false
	q.impactReached = false
	q.impactShotPending = false
	q.finishAfterDraw = false
	q.afterImpact = 0
	q.motion = nil
	switch scenario.action {
	case "standing", "guard":
		q.expected = "A"
	case "crouching":
		q.expected = "2A"
	case "airborne":
		q.expected = "airA"
	case "special":
		q.expected = q.match.gamestate.Characters[0].Properties.SpecialInput
		q.motion = cpuSpecialPattern(q.expected)
	}
}

func (q *qaCapture) Update() error {
	if q.drawErr != nil {
		return q.drawErr
	}
	if q.index >= len(q.scenarios) {
		if len(q.failures) > 0 {
			return fmt.Errorf("%d gameplay scenarios failed: %s", len(q.failures), strings.Join(q.failures, "; "))
		}
		return ebiten.Termination
	}
	if q.framePending {
		return nil
	}
	if q.match == nil {
		q.startScenario()
	}
	scenario := q.scenarios[q.index]
	attacker := q.match.gamestate.Characters[0].StateMachine
	defender := q.match.gamestate.Characters[1].StateMachine
	var inputs [2]input.GameInput
	switch scenario.action {
	case "standing":
		if q.tick == 0 {
			inputs[0] = input.A
		}
	case "crouching":
		if q.tick == 0 {
			inputs[0] = input.Down | input.A
		}
	case "airborne":
		if q.tick == 0 {
			inputs[0] = input.Up
		} else if attacker.IsAirborne() && !q.attackSent {
			inputs[0] = input.A
			q.attackSent = true
		}
	case "special":
		if q.tick < len(q.motion) {
			inputs[0] = q.motion[q.tick]
		}
	case "guard":
		inputs[1] = input.Right // The right-side defender faces left.
		if q.tick == 5 {
			inputs[0] = input.A
		}
	}
	q.match.Update(inputs)
	if attacker.AnimPlayer.ActiveAnimationName() == q.expected {
		q.sawAnimation = true
	}
	verified := false
	if scenario.action == "guard" {
		verified = defender.AnimPlayer.ActiveAnimationName() == "blockHit" && defender.StunFrames > 0 && defender.HP == defender.MaxHP
	} else {
		verified = q.sawAnimation && attacker.AnimPlayer.ActiveAnimationName() == q.expected && defender.HP < defender.MaxHP
	}
	if verified && !q.impactReached {
		q.impactReached = true
		q.impactShotPending = true
	}
	if q.impactReached {
		q.afterImpact++
	}
	q.tick++
	if q.afterImpact >= 18 || q.tick >= 60 {
		if !q.impactReached {
			q.failures = append(q.failures, fmt.Sprintf("%s/%s: no verified impact in %d ticks (attacker=%s, defender=%s, HP=%d/%d)",
				scenario.fighter, scenario.action, q.tick, attacker.AnimPlayer.ActiveAnimationName(),
				defender.AnimPlayer.ActiveAnimationName(), defender.HP, defender.MaxHP))
		}
		q.finishAfterDraw = true
	}
	q.framePending = true
	return nil
}

func (q *qaCapture) Draw(screen *ebiten.Image) {
	if q.match == nil {
		return
	}
	q.match.Draw(screen)
	if !q.framePending {
		return
	}
	scenario := q.scenarios[q.index]
	img := image.NewRGBA(screen.Bounds())
	screen.ReadPixels(img.Pix)
	path := filepath.Join(q.outDir, "frames", fmt.Sprintf("%06d.png", q.frameNumber))
	if err := writeQAFrame(path, img); err != nil {
		q.drawErr = err
		return
	}
	q.frameNumber++
	if q.impactShotPending {
		impactPath := filepath.Join(q.outDir, "impacts", scenario.fighter+"_"+scenario.action+".png")
		if err := writeQAFrame(impactPath, img); err != nil {
			q.drawErr = err
			return
		}
		fmt.Printf("verified and captured %s/%s\n", scenario.fighter, scenario.action)
		q.impactShotPending = false
	}
	q.framePending = false
	if q.finishAfterDraw {
		q.index++
		q.match = nil
		q.finishAfterDraw = false
	}
}

func writeQAFrame(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("capture %s: %w", path, err)
	}
	err = png.Encode(f, img)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return fmt.Errorf("capture %s: %w", path, err)
	}
	return nil
}

func (q *qaCapture) Layout(_, _ int) (int, int) {
	return int(constants.CameraWidth), int(constants.CameraHeight)
}
