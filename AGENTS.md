# AGENTS.md - FGEngine

This file is the agent-focused onboarding and execution guide for this repository.
It complements README.MD with practical context for coding agents.

## 1) Repository Purpose

FGEngine is a 2D fighting game engine written in Go, using Ebitengine for rendering and input.
The project currently has two primary surfaces:
- Game runtime (main executable)
- Utility/test binary (cmd/test)

Current status from README.MD: major rewrites are in progress.

## 2) Tech Stack

### Language and runtime
- Go 1.25.0 (go.mod)

### Core libraries
- github.com/hajimehoshi/ebiten/v2 (game loop, rendering, input, gamepads)
- github.com/ebitengine/debugui (runtime debug UI tooling)
- gopkg.in/yaml.v3 (serialization for characters, language data)

### Build targets and platform strategy
- Desktop native target via regular Go build/test

### Data and assets
- Character definitions: YAML under assets/characters
- Localized text: YAML under assets/text
- Sprites and stage art: assets/common, assets/stages, etc.

## 3) High-Level Architecture

### Runtime entrypoints
- main.go -> config.InitGameConfig() -> ebiten.RunGame(scene.NewSceneManager())
- cmd/test/main.go -> language YAML loading smoke program

### Scene system
- scene.Scene interface:
  - Update([2]input.GameInput) SceneStatus
  - Draw(*ebiten.Image)
- scene.SceneManager owns current scene and routes transitions with SceneStatus.
- Scene transitions:
  - Controller selection
  - Main menu
  - Gameplay

### Gameplay flow
- scene/gameplay.go creates:
  - 2 characters via character.LoadCharacter("PlaceHolder", side)
  - camera
  - stage
  - gameplay.GameState
- gameplay.GameState.Update delegates to each character state machine.

### Character and animation model
- character.Character has:
  - Name
  - StateMachine *animation.StateMachine
- animation.StateMachine contains runtime-only combat state and an ActiveAnim player.
- animation.AnimationPlayer handles:
  - Active animation
  - Frame timing
  - Loop/non-loop behavior
  - Sprite selection via frame data

### Input model
- input.GameInput is a bitmask byte.
- Directional and button states are merged and polled every frame.
- SOCD cleaning exists (Left+Right and Up+Down neutralization).
- Global input ownership supports P1/P2 grouping.

### Rendering model
- Camera abstraction in graphics/camera.go:
  - WorldToScreen and CameraTransform
  - optional world bounds locking
- Stage rendering supports:
  - solid color
  - grid
  - image-based stages
- Image cache in graphics/imageCache.go uses mutex-protected map and fallback sprite.

## 4) Project Layout Map (Where to Change What)

- animation/: animation playback, frame data, state machine
- character/: character loading (YAML), drawing, box rendering
- collision/: hit/hurt/collision box types and detection
- config/: runtime window/layout/lang settings
- constants/: world size, camera size, scene constants
- gameplay/: game state update
- graphics/: camera, cache, render queue abstractions
- input/: keyboard/gamepad mapping, poll/update helpers, special input
- language/: i18n YAML import model
- rollback/: rollback/netcode experiments (currently empty placeholder)
- scene/: scene manager and scene implementations
- stage/: stage visuals and background generation
- types/: shared types (vectors, rects)

## 5) Data Contracts (Important)

### Character YAML contract
Loader expects:
- name
- stateMachine.activeAnim.animations map

Minimal animation expectations:
- each animation has sprites and framedata
- frame data supports duration, spriteIndex, optional velocity deltas and metadata

Example source of truth:
- assets/characters/PlaceHolder.yaml

Path behavior:
- character loader resolves sprite relative paths against the character YAML path.

### Language YAML contract
- language.Language:
  - lang
  - game_text map

Files:
- assets/text/EN.yaml
- assets/text/BR.yaml

## 6) Commands and Validation (Observed on Linux, 2026-04-10)

### Works now
- Validate core packages:
  - go test ./animation ./character ./collision ./config ./constants ./gameplay ./graphics ./input ./language ./scene ./stage ./types . ./cmd/test
- Validate full repository:
  - go test ./...


## 7) Current Risks and Constraints for Agents

1. No stable automated test suite yet
- Most packages return "[no test files]"; current validation is mostly compile-level.

2. Refactor in progress
- README indicates active rewrites.
- Prefer minimal, localized changes; avoid broad architectural rewrites unless requested.

3. Scope is desktop-first for now
- Prioritize Linux/Windows desktop behavior in runtime and tooling.
- Defer non-desktop targets unless explicitly requested.

4. Rollback package is not wired in runtime flow yet
- rollback/ exists but is currently empty in the workspace snapshot.
- Treat rollback integration as future work unless explicitly requested.

## 8) Code Patterns and Conventions to Follow

1. Keep data/runtime split through struct tags
- Runtime-only state uses yaml:"-" where needed.
- Preserve serialization boundaries when changing structs.

2. Favor package-local ownership of responsibilities
- input handles polling and normalization.
- scene handles state transitions.
- gameplay coordinates entities.
- character/animation handle per-entity behavior.

3. Prefer explicit nil checks and graceful fallback
- Existing code frequently checks nil before access.
- image cache fallback to default image is an established pattern.

4. Keep camera-space transforms centralized
- Use graphics.CameraTransform and WorldToScreen instead of custom per-call math.

5. Use bitmask-safe input logic
- Use GameInput helpers (IsPressed, JustPressed, JustReleased).
- Preserve SOCD cleanup behavior when extending input systems.

## 9) Agent Workflow Checklist for This Repo

When implementing a change:
1. Identify surface: runtime game, shared systems, or utility binaries.
2. Make smallest viable change and keep package boundaries.
3. Run targeted go test commands for changed packages.
4. If change is runtime-wide, run the validated core package command from section 6.
5. Document any command failure and exact error in your final report.

When adding features:
1. Add/update YAML schema with backward compatibility in mind.
2. Update loaders and default/fallback behavior together.
3. Keep scene transitions and input semantics deterministic frame-to-frame.

## 10) Suggested Next Stabilization Tasks (High Value)

1. Define rollback package integration points with scene/gameplay flow.
2. Introduce at least smoke tests for YAML load paths and animation playback.
3. Add a repo-native cross-platform build script for Linux/macOS/Windows desktop workflows.
4. Document one canonical run command for desktop game and one for utility/smoke binaries.

## 11) Why This AGENTS.md Is Structured This Way

This structure follows practical guidance from widely adopted agent-doc patterns:
- Keep a single predictable root file with execution-critical information.
- Prioritize build/test/run commands that are actually validated.
- Include architecture map so agents search less and edit the right package first.
- Explicitly list known breakages and constraints to reduce failed attempts.
- Keep guidance concise, actionable, and maintenance-friendly.

External references consulted while preparing this file:
- https://agents.md/
- https://docs.github.com/en/copilot/how-tos/configure-custom-instructions/add-repository-instructions
- https://llmstxt.org/
- https://developers.openai.com/api/docs/guides/prompt-engineering

## 12) Maintenance Rule

Treat this file as living documentation.
Whenever build commands, architecture, or known breakages change, update AGENTS.md in the same PR.
