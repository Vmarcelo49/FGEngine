# FGEngine

A 2D fighting game engine written in Go, rendered with Ebitengine.

FGEngine simulates local 1v1 matches between two **data-driven characters**
at a fixed 60 FPS, on a **fully deterministic** simulation — designed to
support rollback netplay in the future (snapshots, state hashes, and input
logging are already reserved; no networking is implemented).

## Current state

The F0–F7 engine roadmap is complete. The runtime now includes a playable
roster, character selection, local two-player controls, and a CPU opponent:

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
| F8 | Four-fighter roster, character select, local versus and CPU mode | ✅ done |

The playable roster is **Brasa**, **Muralha**, **Íris**, and **Vulto**. Each
fighter has an independent sprite atlas and frame data, standing, crouching,
and aerial attacks, a unique special move, and dedicated guard and hit-reaction
poses. `PlaceHolder` remains the canonical test character used by the suite
and CI.

## What's in the repo

| Path | What |
|---|---|
| `main.go` | Game runtime (entrypoint) |
| `cmd/editor/` | Character editor (imgui-based) |
| `cmd/test/` | Utility/smoke-test binary |
| `assets/characters/` | Character definitions (TOML) + art |
| `assets/characters/*/sequences/` | Per-character transparent sprite sheets, assembled into that fighter's atlas |
| `cmd/generate_fighter_art/` | Rebuilds the four shipped transparent sprite atlases |
| `assets/text/` | Localized text (EN, BR) |
| `docs/SPEC.md` | Normative source of truth for engine behavior |
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

Requirements: **Go 1.26.3 or newer** (see `go.mod`), desktop Linux.
The Ebitengine Linux runtime uses cgo and X11/GL; on Debian/Ubuntu install
`libasound2-dev`, `libgl1-mesa-dev`, `libx11-dev`, `libxcursor-dev`, `libxi-dev`,
`libxinerama-dev`, `libxrandr-dev`, and `libxxf86vm-dev` before building.

```sh
go build .              # game runtime
go run .                # play

go test ./...           # headless test suite (no display needed)

go build ./cmd/editor   # character editor (Linux: requires cgo + GTK)
```

The editor's file picker on Linux uses cgo + GTK (`gtk+-3.0`, `-lX11`);
an active desktop session is required.

## Playing

The first screen lets you join a controller. With a keyboard, press Left or
Right to join a side, then press U (or N) to continue. The main menu offers
**Arcade / vs CPU**, **Local Versus**, controller setup, and quit.

P1 uses **W/A/S/D** to move and **U/I/O/J** for A/B/C/D attacks. P2 uses the
arrow keys and **N/M/comma/period**. A gamepad uses its D-pad or left stick
and four face buttons. In Local Versus, both players lock in a character; in
Arcade, P1 selects a fighter and then the CPU opponent.

All four fighters have separate standing, crouching, and aerial attack sets:
hold down with A/B/C/D for crouching normals, or attack with A/B/C/D while
airborne for that fighter's air normals. Each fighter's source art is kept in
its own `assets/characters/<fighter>/sequences/` directory. The four 4×4
transparent sheets contain standing attacks (rows A–D), crouching attacks
(rows A–D), aerial attacks (rows A–D), and special/defense art (unique special,
standing guard, crouching guard, then four hit reactions). No sequence sheet is
shared between fighters.

Specials use the motion shown on each character card: Brasa `236A`, Muralha
`623B`, Íris `236C`, and Vulto `214D`. The numbers are relative to the way
your fighter faces; `236` is down, down-forward, forward, then the listed
attack button. Press **F1** during a match to show hitboxes and engine guides.

### Headless Linux runtime (Ebitengine)

Ebitengine 2.9.9's [Linux CI workflow](https://github.com/hajimehoshi/ebiten/blob/v2.9.9/.github/workflows/test.yml)
starts Xvfb and sets `DISPLAY=:99.0`. Its runtime does not expose a separate
`EBITENGINE_HEADLESS` backend. To launch FGEngine without a physical monitor,
install `xvfb` and `xauth`, build normally, and run it inside a virtual X server:

```sh
go build -o fgengine .
xvfb-run -a -s "-screen 0 1024x768x24" ./fgengine
```

This uses the regular Ebitengine Linux window and graphics backend with a
virtual display. The pure simulation/data test suite remains display-free;
that is separate from running the game window headlessly. Building still
requires the Linux cgo/X11/GL development packages listed above.

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

- **`docs/SPEC.md`** — normative specification: determinism, frame pipeline,
  input model, character data contract, combat rules, physics, roadmap.
- **`AGENTS.md`** — build commands, package layout, known constraints
  (kept at the repo root, where agent tooling discovers it).
- **`docs/lista-de-tarefas.md`** — Portuguese working task list (subordinate
  to SPEC §12).
