# FGEngine

A 2D fighting game engine written in Go, rendered with Ebitengine.

FGEngine simulates local 1v1 matches between two **data-driven characters**
at a fixed 60 FPS, on a **fully deterministic** simulation — designed to
support rollback netplay in the future (snapshots, state hashes, and input
logging are already reserved; no networking is implemented).

## Current state

The F0–F7 roadmap is complete. The engine now plays a full local match:

| Phase | Scope | Status |
|---|---|---|
| F0 | Determinism foundations + character data contract (TOML) | ✅ done |
| F1 | Hit resolution: damage, knockback, pushback, trades | ✅ done |
| F2 | Life & reactions: HP, hitstun, armor, invincibility | ✅ done |
| F3 | Guard: blockstun (ground/air), re-guard, landing conversion | ✅ done |
| F4 | Frame-data effects: knockdown, bounces, OTG, priority | ✅ done |
| F5 | Rounds, KO, match end, scene exit | ✅ done |
| F6 | Localized fight HUD (HP bars, round pips, timer) | ✅ done |
| F7 | End-to-end validation: scripted match scenario in CI, stable golden replays | ✅ done |

The first real character is in production: **Tiki**
(`assets/characters/tiki/` — design doc + concept art). `PlaceHolder`
remains the canonical test character used by the suite and CI.

## What's in the repo

| Path | What |
|---|---|
| `main.go` | Game runtime (entrypoint) |
| `cmd/editor/` | Character editor (imgui-based) |
| `cmd/test/` | Utility/smoke-test binary |
| `assets/characters/` | Character definitions (TOML) + art |
| `assets/text/` | Localized text (EN, BR) |
| `SPEC.md` | Normative source of truth for engine behavior |
| `AGENTS.md` | Agent onboarding & execution guide |

## Architecture in one minute

- **Determinism contract** (SPEC §3): the simulation is a pure function of
  (initial state, input stream). Fixed 60 steps/second; no wall clock, no
  map iteration, no unseeded RNG on the simulation path. Verified by golden
  replay tests (SPEC §3.6).
- **Data-driven characters** (SPEC §6): animations, frame data, boxes,
  cancels, and HP live in per-character TOML files with strict load-time
  validation. Go code contains no per-character values.
- **Frame pipeline** (SPEC §4): facing → input history/intents → physics →
  hit detection + pushbox → animation state decisions, in a normative order.

Key packages: `gameplay` (simulation loop), `input` (intents/sequences/SOCD,
ebiten-free), `character`/`animation` (data + playback), `scene`
(transitions), `graphics` (camera + rendering), `device` (keyboard/gamepad),
`editor` (content authoring).

## Getting started

Requirements: **Go 1.26.3** (see `go.mod`), desktop Linux.

```sh
go build .              # game runtime
go run .                # play

go test ./...           # headless test suite (no display needed)

go build ./cmd/editor   # character editor (Linux: requires cgo + GTK)
```

The editor's file picker on Linux uses cgo + GTK (`gtk+-3.0`, `-lX11`);
an active desktop session is required.

## CI

`.github/workflows/ci.yml` runs **only on version tags (`v*`)**: build, vet,
headless tests, and a gofmt gate. WIP pushes to `main` are deliberately not
gated — tag a release when the build has to be green.

## Non-goals for v1

- Online play (determinism reservations only)
- Audio
- Team/versus modes (1v1 only)
- Stage data format (solid-color stage)
- Mobile/touch input

## Documentation

- **`SPEC.md`** — normative specification: determinism, frame pipeline,
  input model, character data contract, combat rules, physics, roadmap.
- **`AGENTS.md`** — build commands, package layout, known constraints.
- **`lista de tarefas.md`** — Portuguese working task list (subordinate to
  SPEC §12).
