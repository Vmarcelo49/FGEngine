# AGENTS.md - FGEngine

This file is the agent-focused onboarding and execution guide for this repository.
It complements README.MD with practical context for coding agents.

> **Behavior source of truth:** `docs/SPEC.md` is normative for all
> engine behavior — determinism contract, frame pipeline, input model,
> character data contract, combat rules, physics, and the F0–F7 roadmap.
> When this file and `docs/SPEC.md` disagree, `docs/SPEC.md` wins. Read the relevant
> `docs/SPEC.md` sections before changing behavior, and update `docs/SPEC.md` in the
> same PR when behavior changes (SPEC §11.6).

## 1) Repository Purpose

FGEngine is a 2D fighting game engine written in Go, using Ebitengine for rendering/input.
The repository currently has three active surfaces:
- Runtime game executable (`main.go`)
- Character editor executable (`cmd/editor`)
- Utility/test executable (`cmd/test`)

Current project state from README.MD: F0–F8 are complete; the runtime has
four selectable fighters, local versus, a CPU opponent, and full match flow.

## 2) Tech Stack

### Language and runtime
- Go 1.26.3 minimum (`go.mod`).

### Core libraries
- `github.com/hajimehoshi/ebiten/v2` (game loop, rendering, input)
- `github.com/ebitengine/debugui` (runtime debug overlay)
- `github.com/gabstv/ebiten-imgui/v3` and `github.com/gabstv/cimgui-go` (editor UI)
- `github.com/pelletier/go-toml/v2` (config, character, and language serialization — SPEC §6.1)

### Platform and tooling notes
- Desktop native target via regular Go build/test.
- Editor file pickers on Linux use cgo + GTK (`editor/image_picker_gtk_linux.go`), with:
  - `#cgo pkg-config: gtk+-3.0`
  - `#cgo LDFLAGS: -lX11`

### Data and assets
- Character definitions: TOML under `assets/characters`; the playable roster
  uses transparent atlases under `assets/characters/{brasa,iris,muralha,vulto}/`.
- Per-character source sprite sheets live in `assets/characters/<fighter>/sequences/`
  and are packed into atlas cells by `cmd/generate_fighter_art/`.
- Every playable fighter has its own four source sheets: 4×4 standing attacks,
  crouching attacks, aerial attacks, and a special/defense sheet. The first
  three use rows A–D; the final sheet uses rows special, high guard, crouch
  guard, and standing/crouch/air/block hit reactions. Do not point one fighter
  at another fighter's art or reuse a sequence sheet across the roster.
- `4` is the back/guard animation and still needs `changeXSpeed` so holding
  back continues to move the fighter; `blockHit` is the standing blockstun
  reaction. Keep their art and movement data distinct.
- Localized text: TOML under `assets/text`
- Shared/stage art: `assets/common`, `assets/stages`

## 3) Runtime and Editor Entrypoints

### Game runtime
- `main.go` -> `config.InitGameConfig()` -> `ebiten.RunGame(scene.NewSceneManager())`

### Character editor
- `cmd/editor/main.go` -> `ebiten.RunGame(editor.NewCharacterEditor())`
- Window defaults to 1920x1080 and resizing enabled.

### Utility binary
- `cmd/test/main.go` calls `language.LoadLang(language.English)` (stable path under `assets/text/`, works from repo root).

## 4) High-Level Architecture

### Scene system
- `scene.Scene` interface:
  - `Update([2]input.GameInput) SceneStatus`
  - `Draw(*ebiten.Image)`
- `scene.SceneManager` owns current scene and routes transitions.
- Input carry-over prevention is implemented using a release-to-accept gate (`waitNeutral`) after scene swaps.

### Gameplay flow
- Scene flow is controller setup → main menu → mode select → character select → gameplay → main menu.
- `scene/characterSelect.go` loads top-level character TOMLs (excluding `PlaceHolder`) and returns selected file names to the scene manager.
- `scene/gameplay.go` loads the selected fighters, camera, stage, and `gameplay.GameState`.
- Versus-CPU mode generates deterministic P2 input from the current state and round timer, then feeds both input frames into `GameState.Update`.
- `gameplay.GameState.Update` loop currently does:
  1. facing resolution
  2. input history + intent extraction
  3. frame velocity + physics
  4. hit detection + pushbox collision resolve
  5. post-physics animation state decisions + animation player update

### Character and animation model
- `character.Character` stores `Name` and `StateMachine`.
- `animation.StateMachine` stores runtime combat state (position/velocity/facing) plus `AnimPlayer`.
- `animation.AnimationPlayer` manages active animation, frame stepping, loop behavior, and frame data access.

### Rendering and debug model
- Camera abstraction in `graphics/camera.go` centralizes world/screen transforms.
- Gameplay scene draws atlas-cropped character sprites, stage backdrop, and HUD.
- Hitboxes and world guides are hidden by default; press F1 to toggle the debug view.
- Runtime debug telemetry is wired through `debugui` in the gameplay scene.

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
- `character/`: character loading and validation
- `cmd/editor/`: editor executable entrypoint
- `cmd/test/`: utility/smoke executable
- `config/`: runtime window/layout/lang settings
- `constants/`: world/camera/input constants
- `device/`: gamepad/keyboard polling, mappings, device ownership
- `editor/`: imgui-based character editor
- `gameplay/`: game update loop, collision, hit detection
- `graphics/`: camera, image cache, character sprite/box rendering
- `input/`: input intents, sequences, SOCD (ebiten-free)
- `language/`: i18n TOML import model
- `scene/`: scene manager and scene implementations
- `stage/`: stage rendering/backdrop generation
- `types/`: shared vectors/rects/box types
- `docs/`: `SPEC.md` (normative behavior spec) and `lista-de-tarefas.md`
  (Portuguese working task list)

Note: there is no top-level `collision/` package in the current workspace.

## 6) Data Contracts (Important)

The authoritative contracts live in `docs/SPEC.md` §6 (character files,
frame data, boxes, canonical states, strict validation). What follows is
current-state orientation only — on conflict, `docs/SPEC.md` wins.

### Character files (TOML per SPEC §6.1)
Loader/editor flows expect the flat SPEC §6.2 layout:
- `name` + `[properties]`
- `[animations."<name>"]` tables with `sprites`/`framedata` arrays

Minimal animation expectations:
- each animation has `sprites` and `framedata`
- framedata includes `duration`
- `spriteIndex` is used to select visual frame

Path behavior:
- Loader resolves relative sprite paths against the character file location.
- Editor save may rewrite absolute sprite paths to relative paths against save destination.

Reference file:
- `assets/characters/PlaceHolder.toml` (fixture art in `assets/characters/placeholder/`)

Playable roster files are `Brasa.toml`, `Muralha.toml`, `Iris.toml`, and
`Vulto.toml`. Their selection metadata (`role`, `tagline`, `specialName`,
`specialInput`, `accentColor`) is presentation-only; combat balance remains
in animation frame data. Sprite `rect` values crop tiles from each fighter's
PNG atlas.

### Local keyboard controls
- P1: WASD + U/I/O/J.
- P2: arrow keys + N/M/comma/period.
- After the shared keyboard joins either side on the controller screen, both
  layouts are routed independently for local-versus play.

### Language TOML contract
- `language.Language` fields:
  - `lang`
  - `game_text` map

Reference files:
- `assets/text/EN.toml`
- `assets/text/BR.toml`

## 7) Commands and Validation (Observed on Linux, 2026-04-23)

### Confirmed working
- Build runtime:
  - `go build .`
- Build editor binary:
  - `go build ./cmd/editor`
- Build utility binary:
  - `go build ./cmd/test`
- Validate all packages:
  - `go test ./...`
- Headless suite: sim/data packages carry no display dependency, so the
  full suite runs without Xvfb or a display server. CI runs it on version
  tags only (`.github/workflows/ci.yml`, triggered by `v*`: build, vet,
  test, gofmt gate) — WIP pushes to main are deliberately not gated, so
  tag when the build has to be green. CI installs the cgo build deps
  explicitly — Ebitengine's Linux X11/GL link set plus GTK for the editor
  picker (§8 item 3) — because `./...` still builds the cgo editor.

### Currently failing/stale commands
- Old targeted test command from previous AGENTS revisions is stale because it includes `./collision`, which no longer exists.
  - Fails with: `stat .../collision: directory not found`

## 8) Current Risks and Constraints for Agents

1. Display-linked packages carry no tests
- `scene`, `graphics`, `device`, `stage`, `editor` need a display and are
  covered by build/vet plus review; sim/data packages (`gameplay`,
  `animation`, `character`, `input`, `config`, `types`, `language`) have
  headless unit/integration tests per SPEC §11.5.

2. Post-rewrite stabilization
- The F0–F7 roadmap is complete (see README.MD); the focus is now feel
  tuning, content (first real character: Tiki), and docs sync.
- Prefer small, localized edits unless explicitly asked for larger changes.

3. Editor has Linux-specific native picker dependency
- Linux picker path requires cgo + GTK and active desktop session.
- Consider this when debugging editor behavior in headless CI/containers.

4. Desktop-first assumptions remain
- Runtime/editor workflows prioritize desktop behavior and desktop input models.

5. Runtime artifacts are gitignored
- `config.toml` (created on first run) and `imgui.ini` (editor runtime
  state) are local artifacts and must not be committed.

## 9) Code Patterns and Conventions

1. Keep runtime-vs-serialized boundaries explicit
- Preserve `toml:"-"` runtime-only tags in gameplay state fields.

2. Respect package ownership
- `input`: intent extraction and normalization (ebiten-free)
- `device`: hardware polling and mapping
- `scene`: scene transitions and gating
- `gameplay`: frame simulation coordination
- `character`/`animation`: per-character data and behavior
- `editor`: content authoring workflow

3. Preserve deterministic frame behavior
- `docs/SPEC.md` §3 (determinism contract) and §4 (frame pipeline order) are normative — do not reorder, parallelize, or introduce map iteration / wall-clock / unseeded RNG into the simulation path.
- Maintain update ordering assumptions in `gameplay.GameState.Update`.
- Keep scene transition input-neutral gating behavior intact unless intentionally changed.

4. Prefer explicit nil checks and fallbacks
- Existing code frequently protects nullable runtime fields.

5. Keep camera-space math centralized
- Use camera helpers instead of ad-hoc transform logic.

## 10) Agent Workflow Checklist for This Repo

When implementing a change:
1. Identify the surface first: runtime, editor, shared systems, or utility binary.
2. Read the governing `docs/SPEC.md` sections before changing behavior (determinism §3, pipeline §4, data §6, combat §7 as applicable).
3. Make the smallest viable change and keep package boundaries.
4. Run `go test` for changed packages; run `go test ./...` for cross-cutting changes.
5. If touching editor GTK integration, validate on Linux with cgo enabled.
6. Document any command failures with exact errors.
7. Update `docs/SPEC.md` in the same PR when behavior changes (SPEC §11.6).

When changing character/data flow (TOML per SPEC §6.1):
1. Update loaders and save/export paths together.
2. Preserve backward compatibility when practical.
3. Keep sprite path handling consistent (relative vs absolute behavior).

## 11) Current Follow-up Notes

1. Pushbox resolution still uses the first overlapping collision box pair
   (§8.4); multiple hit and hurt boxes are supported.
2. Pause, audio, and richer stage authoring remain future work in `docs/SPEC.md`.

## 12) Maintenance Rule

Treat this file as living documentation.
Whenever build commands, architecture, or known breakages change, update AGENTS.md in the same PR.
Behavior changes update `docs/SPEC.md` in the same PR instead (SPEC §11.6) — this file tracks build/architecture/breakage facts only.
