# AGENTS.md - FGEngine

This file is the agent-focused onboarding and execution guide for this repository.
It complements README.MD with practical context for coding agents.

## 1) Repository Purpose

FGEngine is a 2D fighting game engine written in Go, using Ebitengine for rendering/input.
The repository currently has three active surfaces:
- Runtime game executable (`main.go`)
- Character editor executable (`cmd/editor-imgui`)
- Utility/test executable (`cmd/test`)

Current project state from README.MD: major rewrites are in progress.

## 2) Tech Stack

### Language and runtime
- Go 1.25.0 (`go.mod`)

### Core libraries
- `github.com/hajimehoshi/ebiten/v2` (game loop, rendering, input)
- `github.com/ebitengine/debugui` (runtime debug overlay)
- `github.com/gabstv/ebiten-imgui/v3` and `github.com/gabstv/cimgui-go` (editor UI)
- `gopkg.in/yaml.v3` (character and language serialization)

### Platform and tooling notes
- Desktop native target via regular Go build/test.
- Editor file pickers on Linux use cgo + GTK (`editor/image_picker_gtk_linux.go`), with:
  - `#cgo pkg-config: gtk+-3.0`
  - `#cgo LDFLAGS: -lX11`

### Data and assets
- Character definitions: YAML under `assets/characters`
- Localized text: YAML under `assets/text`
- Shared/stage art: `assets/common`, `assets/stages`

## 3) Runtime and Editor Entrypoints

### Game runtime
- `main.go` -> `config.InitGameConfig()` -> `ebiten.RunGame(scene.NewSceneManager())`

### Character editor
- `cmd/editor-imgui/main.go` -> `ebiten.RunGame(editor.NewCharacterEditor())`
- Window defaults to 1920x1080 and resizing enabled.

### Utility binary
- `cmd/test/main.go` currently calls `language.ImportYAML("./ptbr.yaml")`.
- This is currently path-fragile and fails from repo root unless that file exists there.

## 4) High-Level Architecture

### Scene system
- `scene.Scene` interface:
  - `Update([2]input.GameInput) SceneStatus`
  - `Draw(*ebiten.Image)`
- `scene.SceneManager` owns current scene and routes transitions.
- Input carry-over prevention is implemented using a release-to-accept gate (`waitNeutral`) after scene swaps.

### Gameplay flow
- `scene/gameplay.go` creates two placeholder characters, camera, stage, and `gameplay.GameState`.
- `gameplay.GameState.Update` loop currently does:
  1. facing resolution
  2. input history + intent extraction
  3. frame velocity + physics
  4. hit detection + pushbox collision resolve
  5. post-physics animation state decisions + animation player update

### Character and animation model
- `character.Character` stores `Name` and `StateMachine`.
- `animation.StateMachine` stores runtime combat state (position/velocity/facing) plus `ActiveAnim`.
- `animation.AnimationPlayer` manages active animation, frame stepping, loop behavior, and frame data access.

### Rendering and debug model
- Camera abstraction in `graphics/camera.go` centralizes world/screen transforms.
- Gameplay scene draws stage, character sprites, debug boxes, and world guide lines.
- Runtime debug telemetry is wired through `debugui` update path in gameplay scene.

### Editor model
- `editor.CharacterEditor` is an `ebiten.Game` that edits:
  - character name
  - animation list and active animation
  - per-frame frame data
  - boxes (collision/hit/hurt)
  - sprite/image sets for animations
- Editor keeps a dirty-state workflow and unsaved-change prompts.
- Editor character I/O (`editor/editor_character_io.go`) normalizes sprite paths and can convert absolute paths to relative paths on save.

## 5) Project Layout Map (Where to Change What)

- `animation/`: animation playback, frame data, state machine
- `character/`: character loading, drawing, box rendering
- `cmd/editor-imgui/`: editor executable entrypoint
- `cmd/test/`: utility/smoke executable
- `config/`: runtime window/layout/lang settings
- `constants/`: world/camera/input constants
- `editor/`: imgui-based character editor
- `gameplay/`: game update loop, collision, hit detection
- `graphics/`: camera and image cache
- `input/`: keyboard/gamepad mapping, polling, input intent helpers
- `language/`: i18n YAML import model
- `scene/`: scene manager and scene implementations
- `stage/`: stage rendering/backdrop generation
- `types/`: shared vectors/rects/box types

Note: there is no top-level `collision/` package in the current workspace.

## 6) Data Contracts (Important)

### Character YAML contract
Loader/editor flows expect:
- `name`
- `stateMachine.activeAnim.animations` map

Minimal animation expectations:
- each animation has `sprites` and `framedata`
- framedata includes `duration`
- `spriteIndex` is used to select visual frame

Path behavior:
- Loader resolves relative sprite paths against character YAML file location.
- Editor save may rewrite absolute sprite paths to relative paths against save destination.

Reference file:
- `assets/characters/PlaceHolder.yaml`

### Language YAML contract
- `language.Language` fields:
  - `lang`
  - `game_text` map

Reference files:
- `assets/text/EN.yaml`
- `assets/text/BR.yaml`

## 7) Commands and Validation (Observed on Linux, 2026-04-23)

### Confirmed working
- Build runtime:
  - `go build .`
- Build editor binary:
  - `go build ./cmd/editor-imgui`
- Build utility binary:
  - `go build ./cmd/test`
- Validate all packages:
  - `go test ./...`

### Currently failing/stale commands
- Old targeted test command from previous AGENTS revisions is stale because it includes `./collision`, which no longer exists.
  - Fails with: `stat .../collision: directory not found`
- Utility run command is currently path-fragile:
  - `go run ./cmd/test` fails from repo root with `open ./ptbr.yaml: no such file or directory`

## 8) Current Risks and Constraints for Agents

1. Automated tests are mostly compile-level
- Most packages report `[no test files]`.

2. Refactor in progress
- README states broad rewrites are underway.
- Prefer small, localized edits unless explicitly asked for larger changes.

3. Editor has Linux-specific native picker dependency
- Linux picker path requires cgo + GTK and active desktop session.
- Consider this when debugging editor behavior in headless CI/containers.

4. Desktop-first assumptions remain
- Runtime/editor workflows prioritize desktop behavior and desktop input models.

5. Repository contains generated/runtime artifacts
- Root currently includes an `editor-imgui` ELF binary artifact; do not treat it as source.

## 9) Code Patterns and Conventions

1. Keep runtime-vs-serialized boundaries explicit
- Preserve `yaml:"-"` runtime-only tags in gameplay state fields.

2. Respect package ownership
- `input`: polling and normalization
- `scene`: scene transitions and gating
- `gameplay`: frame simulation coordination
- `character`/`animation`: per-character behavior and assets
- `editor`: content authoring workflow

3. Preserve deterministic frame behavior
- Maintain update ordering assumptions in `gameplay.GameState.Update`.
- Keep scene transition input-neutral gating behavior intact unless intentionally changed.

4. Prefer explicit nil checks and fallbacks
- Existing code frequently protects nullable runtime fields.

5. Keep camera-space math centralized
- Use camera helpers instead of ad-hoc transform logic.

## 10) Agent Workflow Checklist for This Repo

When implementing a change:
1. Identify the surface first: runtime, editor, shared systems, or utility binary.
2. Make the smallest viable change and keep package boundaries.
3. Run `go test` for changed packages; run `go test ./...` for cross-cutting changes.
4. If touching editor GTK integration, validate on Linux with cgo enabled.
5. Document any command failures with exact errors.

When changing YAML/data flow:
1. Update loaders and save/export paths together.
2. Preserve backward compatibility when practical.
3. Keep sprite path handling consistent (relative vs absolute behavior).

## 11) Suggested Next Stabilization Tasks

1. Fix `cmd/test` to use a stable asset path (for example under `assets/text`).
2. Add smoke tests for character YAML load/save and animation playback edge cases.
3. Replace stale command snippets in docs that reference removed packages.
4. Decide policy for checked-in binary artifacts like root `editor-imgui`.

## 12) Maintenance Rule

Treat this file as living documentation.
Whenever build commands, architecture, or known breakages change, update AGENTS.md in the same PR.
