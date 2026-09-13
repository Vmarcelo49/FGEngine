# FGEngine — Engine Specification

This document is the **source of truth for behavior**. Code, docs (including AGENTS.md),
and future development (human or AI-agent) must conform to it. When code and this spec
disagree, the spec wins — fix the code, or change the spec deliberately first.

Status legend used throughout:
`[done]` implemented and working · `[partial]` implemented, contract not yet fully met · `[todo]` specified, not implemented.

---

## 1. Overview

### 1.1 What FGEngine is

FGEngine is a 2D fighting-game engine written in idiomatic Go, rendered with
Ebitengine. It simulates 1v1 matches between two data-driven characters at a
fixed 60 FPS, with a deterministic simulation designed to support rollback
netplay in the future.

### 1.2 Design principles (normative)

1. **P1 — Clean, idiomatic Go.** No premature abstractions; small packages with
   clear ownership; explicit nil checks and fallbacks.
2. **P2 — Data-driven characters.** Everything character-specific — animations,
   frame data, boxes, cancels, transitions, HP — lives in the character data
   file. Go code must not contain per-character values. Go code may only assume:
   - the data contract defined in §6, and
   - the canonical state names defined in §6.7.
3. **P3 — Completely deterministic.** The simulation is a pure function of
   (initial state, input stream). See §3 for the normative rules.
4. **P4 — Fixed timestep.** 60 simulation steps per second. No real-time delta
   anywhere in the simulation. All durations are integer frames.
5. **P5 — Rollback-ready.** The engine reserves (but does not implement) the
   pieces needed for rollback netplay: state snapshots, per-frame state hashes,
   input logging. See §10.

### 1.3 In scope for v1

- 1v1 local matches (two controllers/keyboard), fixed 60 FPS.
- Full combat pipeline: hit detection → damage/hitstun → guard → frame-data
  effects → rounds/KO → HUD (phases F0–F7, §12).
- Character authoring data files (TOML) and the canonical data contract (§6).
- Determinism contract and verification tooling (§3).
- Rollback *reservations* only (snapshots, hashes, input log) — §10.

### 1.4 Non-goals (explicit)

- **Online play.** No networking is implemented in v1; only the
  determinism/snapshot reservations of §10.
- **Full character editor.** The editor is a non-goal until the groundwork in
  F0–F5 is solid. The editor (existing imgui tool) must eventually round-trip
  the §6 contract, but its feature set is out of scope here.
- **Audio.** Frame data carries sound IDs (§6.5) but v1 has no audio system.
- **Team/versus modes.** v1 is strictly 1v1.
- **Stage data format.** v1 ships the solid-color stage; a stage data format
  is future work (§9).
- **Mobile/touch input.** Desktop (keyboard + gamepad) only.

---

## 2. Architecture

### 2.1 Package ownership

| Package | Responsibility |
|---|---|
| `animation` | Animation playback, frame data, per-character runtime state machine |
| `character` | Character loading and validation (drawing lives in `graphics`) |
| `gameplay` | Frame simulation: pipeline order, physics, hit detection, body collision |
| `input` | Input intents, sequences, SOCD — pure logic, no ebiten (testable headless) |
| `device` | Gamepad/keyboard polling, mappings, device ownership |
| `scene` | Scene manager, scene implementations, input-neutral gating between scenes |
| `graphics` | Camera/world-screen transforms, image cache, character sprite/box rendering |
| `stage` | Stage rendering |
| `language` | i18n text (EN/BR) |
| `config` | User config (window, deadzone, language) |
| `constants` | World/camera/input constants |
| `types` | Shared vector/rect/box types |
| `editor`, `cmd/*` | Editor executable and utility binaries (out of v1 scope, §1.4) |

### 2.2 Runtime flow

`main` → `config.InitGameConfig()` → `ebiten.RunGame(scene.NewSceneManager())`.

Scene flow: `controller select → main menu → gameplay → (match end)`,
with pause overlaying gameplay `[todo]`. Pause (when implemented)
suspends simulation steps entirely (no partial updates); resume continues
with the exact next step — determinism-safe by construction.
Scene transitions are gated by a release-to-accept rule: after a swap, inputs
are discarded until both players are neutral (prevents button carry-over).
This is input plumbing, not simulation, and is exempt from §3 ordering rules.

### 2.3 Rendering

- `Draw` is presentation only. It must never mutate simulation state.
- Camera math is centralized in `graphics/camera.go` (world ↔ screen).
  Ad-hoc transforms elsewhere are forbidden.
- World-space box math has a single shared helper used by both sim and
  debug draw (no parallel implementations).
- Debug overlays (debugui, guide lines, box rendering) are presentation.

---

## 3. Determinism contract

This section is normative. Violating any rule below is a defect.

### 3.1 Timestep

- The game runs at ebiten's default **60 TPS** (ticks per second): one
  `Update` == one simulation step. The default is used as-is, no pinning
  required; changing the tick rate is a determinism-scope event (re-run
  golden hash tests, §3.6).
- One `Update` == one simulation step. All frame counts (durations, stun,
  buffers, timers) are integers in units of these steps.
- No real-time delta, no `time.Now`, no frame-rate compensation in simulation.

### 3.2 Floating point

- Simulation math uses IEEE-754 `float64`.
- **Guarantee scope:** identical results for identical (Go toolchain version,
  platform, build flags). This is the determinism level required for rollback
  between equivalent targets. Cross-platform bitwise equality is *not*
  guaranteed and is explicitly out of contract.
- Arithmetic order is part of the contract: the update pipeline order in
  §4.1 is fixed and must not be reordered, parallelized, or made
  input-dependent.

### 3.3 Forbidden sources of nondeterminism

In any code path reachable from the simulation (`GameState.Update` and
everything it calls), it is forbidden to:

1. **Range over Go maps** where the result can affect simulation behavior.
   Map iteration order is random. Load-time and presentation-only uses
   (asset indexes, debug dumps) are exempt, and polling/mapping iteration
   is covered separately by §5.6 (input layer outside the contract);
   anything reachable from `GameState.Update` that branches, selects, or
   accumulates in iteration order is forbidden.
   - *Known violation (must be fixed in F0):* `input.CheckInputSequences`
     iterates the `InputSequences` map and takes the last detected sequence.
     It must evaluate sequences in the fixed priority order of §5.5.
2. Use wall-clock time or real elapsed time.
3. Use unseeded randomness (`math/rand` global, `crypto/rand`, OS entropy).
4. Use goroutines, channels, or async state in the simulation path.
   Simulation is single-threaded and synchronous.

### 3.4 Random number generation

- All simulation randomness (if/when introduced) must come from a single
  deterministic PRNG owned by `GameState` (e.g. SplitMix64 or PCG32), seeded
  with the match seed.
- The PRNG state is part of the state snapshot (§3.5).

### 3.5 State snapshot and hashing (rollback reservation)

- `GameState` must be **fully snapshot-able**: a snapshot captures every
  byte of runtime state — both `StateMachine`s (position, velocity, facing,
  HP, `StunFrames`, `IgnoreGravityFrames`, `AnimationQueue` contents,
  animation name, frame index, frame time left), both input histories, the
  connect ledger, the PRNG state, and all match counters (round, wins,
  timer, phase, freeze countdown). `TotalDuration` is a derived cache of
  file data and is excluded. This list is authoritative; §7.1 and §12
  defer to it.
- After each simulation step the engine computes a **state hash**:
  canonical byte encoding of the snapshot → FNV-1a 64-bit. (FNV-1a is a
  fast, non-cryptographic hash function, not an encoder; canonical means
  fixed field order with floats hashed by their IEEE-754 bits.)
- The per-step hash is computed every step over the canonical encoding;
  full snapshot capture is on demand (tests, rollback rewind points).
  Hash cost must stay negligible; snapshots must not affect simulation
  cost measurably.

### 3.6 Verification (mandatory test)

- A replay test must exist: a recorded list of `[2]input.GameInput` frames is
  applied from a fixed initial state for N frames; the final state hash is
  asserted against a golden value, and re-running the replay must produce the
  same hash. `[todo]` (F0)
- The test must fail if any §3 rule is violated in a way observable through
  state (map-order behavior, timing, etc.).

---

## 4. Frame model and update order

### 4.1 The simulation pipeline (normative order)

`GameState.Update(inputs)` performs, in exactly this order:

1. **Facing resolution.** Grounded characters only: P1 faces P2 and vice
   versa based on X positions. Airborne characters keep their last facing.
   Exact X-equality resolves P1-right / P2-left.
2. **Input intake** (P1, then P2, in that fixed order):
   - decrement `StunFrames` (floor 0) — stun countdowns tick before
     anything else can set them this frame (§7.6);
   - append `inputs[i]` to the player's input history
     (max length `MaxInputHistory = 30`, oldest dropped);
   - compute the player's **intent** (animation name) from the
     *facing-corrected* history (§5).
3. **Per-player step** (P1, then P2, in that fixed order):
   a. `ApplyVelocity` — add frame-data velocity increments
      (`changeXSpeed` flipped when facing left; `changeYSpeed`).
   b. **Cancel check** — if the intent is allowed by the current frame's
      `cancelTypes`, switch to the intent animation. (Before physics because
      some cancels modify velocity.)
   c. `ApplyPhysics` — friction (grounded), gravity, integration
      (velocity → position), world-bound clamping (§8).
4. **Hit detection** — `checkhit(P1→P2)` then `checkhit(P2→P1)`
   (§7.1). Both directions on the same frame is a **trade** (§7.3).
   The connect ledger is pruned at the start of every `Update`, before
   hit detection.
5. **Body collision** — pushbox overlap resolution, once, after both players
   have integrated (§8.4).
6. **Post-physics animation decisions** (P1, then P2): landing / fall /
   idle / intent-driven transitions (canonical runtime rules, §6.7), then
   advance the animation player (frame step, loop handling).
7. **Match flow** (F5): in `PhaseFight` this runs last — round-end
   detection; in other phases it runs alone — freeze countdown
   (`PhaseRoundEnd`) or no-op (`PhaseMatchEnd`). See §7.7.

The order is fixed. Steps 3 and 6 alternate per-player in lockstep; they must
not be merged, reordered, or made conditional on input values (aside from the
explicit intent logic above).

### 4.2 Symmetry rule

There is no "active" player. P1 and P2 are processed in the same fixed
order every frame. Any rule that distinguishes players must be symmetric
except where explicitly stated (e.g. initial positions, §8.5).

---

## 5. Input model

### 5.1 Devices

- Keyboard and/or gamepad per player; mapping is data (config), not code.
- `GamepadID(-1)` denotes the keyboard. Multiple devices may be OR-merged
  per player side.

### 5.2 Logical inputs

- Directions: `Up`, `Down`, `Left`, `Right`. SOCD: opposing pairs
  (`Left`+`Right`, `Up`+`Down`) clear to neutral — neither direction wins.
  The filter applies after the per-player device merge, so conflicts from
  merged devices resolve the same way.
- Buttons: `A`, `B`, `C`, `D` (four attack inputs).
- **Numpad naming convention** (facing-relative, §5.3):
  `7 8 9` top row, `4 . 6` middle row, `1 2 3` bottom row.

### 5.3 Input history and facing correction

- Each player keeps the last `MaxInputHistory = 30` raw `GameInput` frames.
- Before intent extraction, the history is **facing-corrected**: when the
  character faces left, `Left`/`Right` are swapped entry by entry.
  Consequently, *after correction*, `Right` always means **forward** and
  `Left` always means **back**, regardless of which side of the screen the
  character is on. All intent names and all authored movement data use this
  corrected, facing-relative frame.
- Scene-transition input gating (§2.2) discards inputs, not history entries
  of the previous scene (histories are reset with the match).

### 5.4 Single-input intents

From the last corrected frame, highest-priority intent wins: buttons
first, then direction combos, then single directions:

`D > C > B > A > (diagonals 9/7/3/1) > (singles 8/2/4/6) > (none)`

Buttons beat any direction (e.g. Down+A resolves to `A`). Diagonals and
singles are each mutually exclusive after SOCD filtering and combo
checks, so no ordering applies within those groups. Numpad intents:
`7` jump back, `8` jump, `9` jump forward, `1` crouch back, `2` crouch,
`3` crouch forward, `4` walk back, `6` walk forward.

### 5.5 Motion sequences (specials) and buffering

- A sequence is a list of expected corrected inputs matched **backwards**
  from the newest history entry, with a **buffer** of 10 *non-neutral*
  intervening inputs allowed between steps (`NoInput` slots in a pattern match
  any non-directional button press). Neutral (`NoInput`) frames do not
  consume buffer, so a sequence can span an arbitrarily long neutral gap
  (tap forward, wait, tap forward still dashes) — current behavior, kept
  for v1; an absolute-frame window is possible later tuning. Matching is
  superset-lenient: extra held buttons don't break a step (`Down+Right+A`
  satisfies an expected `Down`).
- Sequences are defined in data with an **explicit priority order**:
  supers/specials first, then normals, then movement. Evaluation order is the
  defined order — never map iteration order (§3.3). `[todo]` — F0 migrates
  `InputSequences` from a map to an ordered list.
- When a sequence matches, its name is the intent. The engine-standard
  sequence set is the full input vocabulary of §6.7: dashes `44`, `66`
  and the seven special motions (`236x`, `214x`, `246x`, `642x`, `623x`,
  `423x`, `22x`), each with x ∈ {A,B,C,D}. The input layer recognizes
  all of them; each character implements only the motions it uses
  (unimplemented ones fizzle per the §6.7 missing-intent rule). Each
  motion's exact input pattern is defined in the ordered sequence table
  (data); standard
  readings: `236x` = Down, Down-Forward, Forward + x; `214x` = Down,
  Down-Back, Back + x; `623x` = Forward, Down, Down-Forward + x;
  `423x` = Back, Down, Down-Back + x; `22x` = Down, Down + x;
  `246x` = Down, Back, Forward + x and `642x` = Forward, Back, Down + x
  (simplified half-circles — no corner/diagonal checks). The split is
  deliberate: directional QCF/QCB/DP-family motions check corners for
  execution leniency, while half-circles stay simplified.
- **Discrete intents** fire **once**. The discrete set is exactly the
  attack/motion intents: `A`, `B`, `C`, `D`, dashes `44`/`66`, and all
  special motions (`236x`, `214x`, `246x`, `642x`, `623x`, `423x`, `22x`
  × A–D). If the previous frame's history also resolves to the same
  discrete intent, the intent is suppressed (no re-trigger while holding).
  Everything else (movement, reactions, `fall`, `landing`, `idle`, ...) is
  continuous. Holding a button therefore plays its animation through once
  instead of looping it (§6.4): with the intent suppressed, the animation
  is no longer *held* and advances past its loop range.

### 5.6 Input as simulation input

`GameState` receives `[2]input.GameInput` per frame. Everything after the
history push (§4.1 step 2) derives from those inputs alone. The input layer
(polling) is not part of the determinism contract.

---

## 6. Character data contract

### 6.1 Storage format policy

- **Format: TOML, exclusively.** YAML was fully removed (characters, config,
  and language files all migrated). (TOML over JSON because data files are
  hand-edited and benefit from comments; a binary format is deliberately
  not adopted — authorability matters more than size here.)
- The **logical contract below is authoritative**; it is format-independent.
  F0 delivers the TOML loader/writer, converts all assets, and removes the
  YAML loader.
- Files are UTF-8. Unknown fields are a **load error** (strict mode), so
  typos surface instead of silently no-op'ing.

### 6.2 Character file layout (logical)

```
Character
├── name: string                     # display name, required
├── properties: CharacterProperties  # §6.3, required
└── animations: map[name → Animation] # §6.4
```

TOML mapping: `[properties]` table; each animation is `[animations."<name>"]`
with `[[animations."<name>".sprites]]` and
`[[animations."<name>".framedata]]` arrays (document order preserved).

### 6.3 Character properties

| Field | Type | Meaning |
|---|---|---|
| `maxHP` | int | Health pool. **Required.** (Replaces the hardcoded 10000 in `character.initialize`.) |

Extension rule: new properties are added here, never in code. (Future
candidates, not yet contract: weight/push priority, walk speed — see §13.)

### 6.4 Animations

```
Animation
├── name: string            # the map key
├── sprites: [Sprite]       # required, ≥1
├── framedata: [FrameData]  # required, ≥1
└── loopFrames: {start, end} # optional — frame index range looped while held
```

- `sprites` entries: `{ imgPath, rect (optional sub-rect), anchor }`.
  `imgPath` is **relative to the character file** (resolved at load).
- `framedata` entries are consumed sequentially; `duration` is the number of
  simulation steps the entry stays active.
- `loopFrames {start, end}`: while the animation is *held* (current intent
  equals the animation name, or the name is `idle`), the frame index wraps
  from `end` back to `start`. On release, the index jumps past `end` to the
  first non-looped frame. The same range doubles as the stun-hold pose for
  reaction states while `StunFrames > 0` (§7.6); `start == end` holds a
  single frame (runtime must accept it).

### 6.5 FrameData fields

Status: `[active]` drives the runtime today · `[F#]` activated by the named
phase (§12) · `[reserved]` parsed and stored, no runtime effect yet.

| Field | Type | Meaning | Status |
|---|---|---|---|
| `duration` | int | Steps this entry stays active. Required, ≥1 | `[active]` |
| `spriteIndex` | int | Which sprite to draw (default 0) | `[active]` |
| `changeXSpeed` | float | Velocity increment X, **facing-relative** (+ = forward) | `[active]` |
| `changeYSpeed` | float | Velocity increment Y (− = up) | `[active]` |
| `cancelTypes` | [string] | Animations cancellable *into* during this frame. Empty = no cancels. `"any"` = everything | `[active]` |
| `boxes` | map[BoxType → [Rect]] | `collision` / `hit` / `hurt` boxes, anchor-relative (§6.6). BoxType is a string-kind enum; unknown keys are a load error (§6.8 rule 10) | `[active]` |
| `damage` | int | HP removed on hit | `[F2]` |
| `hitstun` | int | Stun duration: value loaded into the defender's `StunFrames` in the selected hitstun state (§7.2, §7.6) | `[F2]` |
| `pushback` | int | Horizontal position push applied on hit | `[F2]` |
| `knockback` | int | Horizontal velocity impulse on hit | `[F2]` |
| `knockup` | int | Vertical velocity impulse on hit (launch if >0) | `[F2]` |
| `blockstun` | int | Stun duration: value loaded into the defender's `StunFrames` in the selected blockstun state (§7.4, §7.6) | `[F3]` |
| `isInvincible` | bool | Hits pass through this frame | `[F2]` |
| `hasArmor` | bool | Absorption, checked on the **defender's** frame: a hit on this frame deals damage but causes no hitstun/knockback (§7.2) | `[F2]` |
| `isRecovery` | bool | No cancels out of this frame, even if `cancelTypes` lists some | `[F4]` |
| `isHoldable` | bool | Frame continues while its trigger input is held | `[F4]` |
| `animationSwitch` | string | Forced transition to this animation when this frame ends (data-driven transition) | `[F4]` |
| `canHardKnockdown` | bool | Hit causes knockdown (fall down + get-up) | `[F4]` |
| `canWallBounce` | bool | Knocked character bounces off stage walls | `[F4]` |
| `canGroundBounce` | bool | Launched character bounces off the ground | `[F4]` |
| `canOTG` | bool | OTG = Off The Ground: attack may connect with a downed (`knockdown`) opponent. Hits without `canOTG` pass through downed opponents (§7.5) | `[F4]` |
| `priority` | int | Trade resolution priority (higher wins; equal = both hit) | `[F4]` |
| `soundID` | int | Common sound effect id (0 = none) | `[reserved]` |
| `uniqueSoundID` | int | Character-unique sound effect id | `[reserved]` |

### 6.6 Boxes and coordinate convention

- Box types: `collision` (pushbox — body), `hit` (deals damage), `hurt`
  (receives damage).
- All box rects are authored in the **facing-right reference frame**,
  relative to the **sprite anchor** of the current sprite.
- World transform (normative):
  - facing right: `worldX = pos.X + box.X − anchor.X`
  - facing left:  `worldX = pos.X − box.X − box.W + anchor.X`
  - both: `worldY = pos.Y + box.Y − anchor.Y`
- **All** `collision` boxes of a frame participate in pushbox resolution
  (the current implementation only uses the first — this is a known
  limitation; `[todo]` F0). Multiple `hit`/`hurt` boxes already all
  participate.

### 6.7 Canonical state names

Animation names are strings in the file, but the runtime depends on a
**canonical set** — and the canonical set is **complete and required**:
every character file must define every *required* state in the table
below. A missing required state is a load error (§6.8), not a style
choice. The single exception is the Specials family: the table lists the
full engine-standard input vocabulary (all possible special input
combinations), but each character implements only the motions it actually
uses — absent special states are not an error (see Missing-intent rule).
Anything beyond the table is a character-specific extension and must
never be hardcoded in Go.

Variant scheme — two axes:

1. **Facing axis** (front / back): actions done while moving or leaning
   toward vs. away from the opponent. Where a family has a facing axis,
   all members are required (crouch, walk, dash, jump).
2. **Posture axis** (standing / crouch / air): reaction families where the
   reaction exists for that posture. Air guarding is in v1, so the block
   family has an air variant (`airBlock`); knockdown has no variants (one
   lying state).

| Family | Names | Meaning |
|---|---|---|
| Neutral | `idle` | Neutral standing. No variants. |
| Crouch | `1` `2` `3` | Crouch. |
| Walk | `4` `6` | Walk back / walk front (facing axis). |
| Dash | `44` `66` | Dash back / dash front (facing axis). |
| Jump | `7` `8` `9` | Jump back / up / front (facing axis). |
| Normals | `A` `B` `C` `D` | Four normal attack slots, single posture (per-posture normals are an open question, §13). |
| Specials (optional subset) | `236x` `214x` `246x` `642x` `623x` `423x` `22x`, x ∈ {A,B,C,D} | Engine-standard special-motion vocabulary (28 possible states). Each character implements only the motions it uses; absent ones are not an error. Motion names are engine data; what each move does is per character. |
| Hitstun | `hurt` | Standing hitstun. |
| Hitstun | `crouchHurt` | Hitstun from any crouch state. |
| Hitstun | `airHurt` | Hitstun while airborne. |
| Block | `blockHit` | Standing blockstun (grounded guard, not crouching). |
| Block | `crouchBlock` | Blockstun from a crouch state. |
| Block | `airBlock` | Blockstun while airborne (air guard, §7.4). |
| Fall | `fall` `landing` | Falling / grounded landing recovery. |
| Knockdown | `knockdown` `getup` | Lying after a knockdown / rising out of it (invincible while active). |
| Round | `ko` `win` | Defeated / victory pose. |

**Selection rule (normative).** The full set is required, so reaction
selection is unconditional — the runtime does not need to check that the
target state exists (defensive nil-safety may remain, as in the `landing`
rule (item 1 above), but the contract guarantees the target):

- Hit reaction: defender airborne → `airHurt`; in a crouch state
  (`1`/`2`/`3`) → `crouchHurt`;
  otherwise → `hurt`.
- Block reaction: defender airborne → `airBlock`; in a crouch state
  (`1`/`2`/`3`) → `crouchBlock`; otherwise → `blockHit`.

**Air guarding is in v1.** Holding back while airborne guards air hits
(§7.4) into `airBlock`. Landing converts remaining `airBlock` frames to
`blockHit`.

**Missing-intent rule.** If an input intent names an animation absent from
the character file, the intent is ignored and the current animation
continues (no state change, no error). This is the normal case for
unimplemented special motions — not a bug. (`cancelTypes` and
`animationSwitch` targets are still validated at load, §6.8 rule 3, so
only live input intents can hit this path.)

Extension rule: any additional names are allowed in files (e.g.
`"41236A"`, `"360A"`). `cancelTypes` entries and `animationSwitch` values may
reference any existing animation name or the special token `"any"`.

**Canonical runtime transitions** (Go-side, uniform for all characters,
documented here so they are not invisible hardcoding):

1. Landing (airborne → grounded): if in `airBlock` with stun frames
   remaining → `blockHit` (counter preserved, §7.4); otherwise, while the
   current animation is not `landing` → `landing`.
2. Airborne and current animation finished → `fall`.
3. `landing`/`fall` finished: an active intent (other than the current name)
   → that intent; otherwise → `idle`.
4. Any finished non-looping animation with no intent → `idle`, except:
5. Explicit exits (these states never fall through to rule 4):
   - `knockdown` finished → `getup`.
   - `getup` finished → `idle`.
   - `ko` / `win` finished → stay (terminal states, no exit).

### 6.8 Validation rules (at load, strict)

A character file is rejected if:

1. `name` or `properties` is missing; `maxHP` missing or ≤ 0.
2. Any *required* state in the §6.7 table is missing (everything except
   the Specials family, which is implemented per character). The required
   set is required in full from F0 on — when a state's mechanic arrives
   later (e.g. `knockdown` in F4), a minimal placeholder animation is
   acceptable until real content replaces it.
3. Any `cancelTypes` entry or `animationSwitch` value references a
   non-existent animation name (other than `"any"`).
4. `framedata` empty, or any `duration` ≤ 0.
5. Any `spriteIndex` out of range; any sprite `imgPath` missing on disk.
6. Any `loopFrames.start/end` out of frame range or `start > end`.
   (`start == end` is legal and means hold a single frame.)
7. A reaction state (`hurt`, `crouchHurt`, `airHurt`, `blockHit`,
   `crouchBlock`, `airBlock`) without `loopFrames` (the stun-hold
   pose/range, §7.6).
8. A reaction-state frame carrying `cancelTypes` (reaction states are
   never cancellable, §7.6).
9. Unknown fields/keys present (§6.1 strict mode).
10. A `boxes` entry keyed by anything other than `collision` / `hit` /
    `hurt`.

### 6.9 Path resolution

- `imgPath` is always interpreted **relative to the character file's
  directory** (absolute paths are allowed but discouraged).
- Saving/exporting must produce relative paths against the save destination
  (the editor's existing normalization rule, preserved).

---

## 7. Combat rules

### 7.1 Hit detection `[partial]`

- For each ordered pair (attacker, defender): every `hit` box of the
  attacker's active frame is tested against every `hurt` box of the
  defender's active frame, in world coordinates (§6.6).
- Overlap = hit **candidate**.
- Hit candidates are resolved per §7.2.
- A given active frame's hitboxes may connect **at most once** (no multi-hit
  per frame entry; multi-hit moves are authored as several frame entries).
  Enforced by a connect ledger: `GameState` records (attacker, animation,
  frame index, defender) tuples that already connected; entries clear when
  the attacker's frame advances. The ledger is snapshot state (§3.5).
- `[todo]` F1: apply all effects below; today only overlap is detected.

### 7.2 Hit resolution `[todo]` (F1/F2)

On a hit candidate, in this order. Velocity always applies by increment
— but hit resolution zeroes the defender's velocity first, so a hit
overwrites momentum:

0. **Downed opponents (OTG):** if the defender is in `knockdown`, only
   hits whose active frame has `canOTG: true` connect — everything else
   passes through with no effect. Downed defenders cannot guard. An OTG
   connect zeroes the defender's velocity, applies damage (step 3), pops
   the defender up (`vx += ±knockback` away from the attacker,
   `vy += −knockup`) into `airHurt` with `StunFrames = hitstun` (the
   attack's value itself, §7.6), then proceeds to steps 6–7. `getup` is
   fully invincible.
1. **Invincibility:** defender's active frame `isInvincible` → no hit.
2. **Guard** (§7.4): if guarded → blockstun path, stop.
3. **Damage:** `defender.HP −= damage` (clamped at 0).
4. **Armor:** if the defender's active frame has `hasArmor` → skip 5–6
   (no hitstun, no knockback/pushback) but damage still applied.
5. **Reaction:** zero the defender's velocity, then apply
   `vx += ±knockback` (away from the attacker in world X) and
   `vy += −knockup`; defender enters the selected hitstun state per the
   §6.7 selection rule (`airHurt` / `crouchHurt` / `hurt`) with
   `StunFrames = hitstun`. Grounded hits with `knockup = 0` leave the
   defender grounded.
6. **Pushback:** one-time position displacement apart along X — defender
   by `pushback` away from the attacker, attacker by `pushback/2`
   (integer division) away from the defender.
7. **KO check** (§7.7).

### 7.3 Trades `[todo]` (F1)

- Both directions are checked every frame (§4.1 step 4), so **both hits can
  land on the same frame** (trade): both players apply §7.2 independently.
- `priority` (when activated, F4): if priorities differ, only the higher
  priority hits; equal (or both 0) = trade.
- No hitlag, no hitstop in v1.

### 7.4 Guard `[todo]` (F3)

A hit is **guarded** when the defender is holding **back** (corrected
history shows `Left`-family input toward the attacker on the hit frame)
and either guard posture holds:

- **Ground guard:** defender is grounded and in a guardable state —
  `idle`, `4`, crouch states (`1`, `2`, `3`), or continuing blockstun
  (`blockHit`, `crouchBlock`). Attack/recovery and hitstun states are not
  guardable.
- **Air guard:** defender is airborne in a guardable airborne state —
  `7`/`8`/`9`, `fall`, or continuing blockstun (`airBlock`). Hitstun
  states (`hurt`, `crouchHurt`, `airHurt`) and `ko` are not guardable.

On ground guard: no damage (v1 chip = 0), defender enters `blockHit` or
`crouchBlock` per the §6.7 selection rule with `StunFrames = blockstun`
(§7.6), pushback applied (§7.2 step 6).

On air guard: no damage, defender enters `airBlock` with
`StunFrames = blockstun` (§7.6), pushback applied, and keeps falling.
Landing while `airBlock` is active converts the remaining blockstun to
`blockHit` (grounded). When blockstun expires mid-air, the defender
returns to `fall` via the post-physics rules (§6.7).

### 7.5 Knockdown, bounce, OTG `[todo]` (F4)

- `canHardKnockdown`: hit sends defender into a knockdown: launched per
  §7.2, falls, on ground contact enters `knockdown` (lying), then
  `getup` (fully invincible while active).
- **OTG (Off The Ground):** while lying in `knockdown`, the defender is
  hittable **only** by attacks whose active frame has `canOTG: true`
  (§7.2 step 0); all other hits pass through. Hitting a launched opponent
  before ground contact (juggle) is normal airborne combat and is not
  governed by this rule.
- `canWallBounce`: when a knocked character hits a stage wall with
  significant velocity, velocity is reflected horizontally with damping
  (v1: `vx = −vx × 0.6`) and a bounce animation frame is used.
- `canGroundBounce`: same on landing from a launch (`vy = −vy × 0.4`,
  once per launch).
  (`canOTG` is covered by the OTG rule above.)

### 7.6 State interactions (normative summary)

- Reaction states (`hurt`, `crouchHurt`, `airHurt`, `blockHit`,
  `crouchBlock`, `airBlock`) are **not cancellable** (their frames have no
  `cancelTypes`). They are still hittable: a new hit (or blocked hit)
  during stun simply re-enters the applicable state with the new frame
  counts (hitstun chaining beyond this is out of scope for v1). Re-guard:
  blockstun states re-evaluate guard (§7.4); hitstun states never do.
- `knockdown` takes only OTG hits (§7.5); `getup` and `ko` take none.
- **Stun hold:** while `StunFrames > 0`, the player holds within the
  active reaction animation's `loopFrames` range (looping if
  `start < end`, holding a single frame if `start == end`); the animation
  never advances past the range and no transition (including rule 4)
  fires. When the counter reaches 0, normal advance and transitions
  resume. The counter ticks once per step (§4.1), is part of the state
  snapshot (§3.5), and is always loaded from the attack's `hitstun` /
  `blockstun` value itself — never a separate tunable.
- `isInvincible` frames ignore all hits; `hasArmor` frames absorb the
  reaction, not the damage.
- A hit while the defender is in `ko` is impossible (round already ended).

### 7.7 Rounds and match flow `[todo]` (F5)

- **Round timer:** `GameState.TimerFrames int`, initialized to
  `RoundTimerFrames = 99 * 60` (constants, single source) and decremented
  once per simulation step during the fight phase. Displayed seconds are
  ceiling: `(TimerFrames + 59) / 60`. Only adds/subs ever apply — no
  other timer math exists in v1.
- **Round ends** when: KO (HP ≤ 0 → loser enters `ko`, winner enters
  `win`), or timer hits 0 (higher HP wins: loser `ko`, winner `win`;
  tie = no round awarded, full reset, round number unchanged).
  Double KO (both HP ≤ 0 the same frame, only possible via trade) ties:
  no round awarded, full reset, round number unchanged.
- **Ownership:** the match phase (fight → round-end freeze → reset →
  next fight), the freeze countdown, and round-win counting are
  `GameState`-owned simulation state and transitions — scenes only
  observe (§10, item 5).
- **Match:** best of 3 (first to 2 rounds). Round wins are `GameState`
  fields (snapshot state); the HUD only reads them (§7.8).
- After round end: full freeze (v1: 60 frames — `Update` ticks only the
  freeze countdown; simulation, inputs, and animations hold), then next
  round (positions reset to §8.5 starts, HP refilled, stun and launch
  flags cleared, histories and ledger cleared, timer reset, both to
  `idle`) or match end. Richer freeze behavior is future work.
- Match end → scene transition out of gameplay (back to menu).

### 7.8 HUD `[todo]` (F6)

- HP bars: P1 left, P2 right, mirrored drain; show current/max HP from
  snapshot state.
- Round-win pips per player; center round timer.
- KO / round-end text, localized via `language` (EN/BR).
- HUD reads simulation state (presentation), never mutates it.

---

## 8. Physics

### 8.1 Constants (single source: `constants`)

| Constant | Value | Meaning |
|---|---|---|
| `WorldWidth` / `WorldHeight` | 768 × 432 | Simulated world extent |
| `CameraWidth` / `CameraHeight` | 640 × 360 | Viewport |
| `GroundLevelY` | `WorldHeight − 50` | Ground plane (Y grows downward) |
| `Gravity` | 1 (per step²) | Applied while airborne or falling |
| `maxVerticalSpeedY` | 10 | Terminal fall speed |
| `maxHorizontalSpeedX` | 999 | Horizontal velocity clamp (grounded) |
| `horizontalFriction` | 0.8 | Grounded velocity multiplier per step |
| `minHorizontalSpeed` | 0.05 | Below this, grounded vx is zeroed |
| `MaxInputHistory` | 30 | Input history depth (§5.3) |

### 8.2 Integration order (normative, per step, per player)

1. Frame-data velocity increments (§4.1 step 3a).
2. Friction (grounded only): `vx *= 0.8`; zero if `|vx| < 0.05`.
3. Gravity: applied when airborne **or** `vy < 0`; capped at terminal
   speed.
4. Integrate: `pos += vel` (X then Y).
5. Bounds clamping (§8.5).

### 8.3 Airborne rules

- `IsAirborne() == pos.Y < GroundLevelY`.
- Jump-start animations (`7`, `8`, `9`) are disallowed while already
  airborne (canonical rule, §6.7 runtime behavior).
- Facing is frozen in the air (§4.1 step 1).

### 8.4 Body collision (pushbox) `[partial]`

- After both players integrate, the first-overlapping `collision` box pair
  (first in slice order — deterministic) is resolved along the axis of
  least penetration. Axis ties (`overlapX == overlapY`) resolve along Y;
  a zero delta component takes the else branch.
- Separation is split by relative speed (velocity-weighted); the faster
  character moves more. Speeds equal → 50/50.
- Velocity transfer: the slower character adopts the faster character's
  velocity along the resolution axis, and the faster character's velocity
  is damped by ×0.9.
- `[todo]` F0: consider **all** collision boxes per frame, not just the
  first (§6.6).

### 8.5 World bounds

- X clamped to `[0, WorldWidth]`, touching a wall zeroes `vx`.
- Y clamped to `[0, GroundLevelY]` (Y grows downward): ground contact
  (`Y > GroundLevelY`) clamps and zeroes `vy`; ceiling contact (`Y ≤ 0`)
  zeroes `vy` only while moving up (`vy < 0`).
- Initial positions (normative): P1 at `WorldWidth/4` facing right, P2 at
  `3·WorldWidth/4` facing left, both at `GroundLevelY` (grounded — no
  opening fall).

---

## 9. Stage

- v1 stage: solid color, static, rendered in the background layer.
- The stage contributes the world bounds of §8.5 and the ground plane.
- A stage data format (boundaries, visuals, hazards, camera limits) is
  future work and explicitly not part of this contract yet.

---

## 10. Online readiness (rollback) — reservations only

No networking in v1 (§1.4). The engine must keep these invariants so rollback
can be added later without re-architecture:

1. Simulation is a pure function of (initial state, input frames) — §3.
2. Full snapshot of `GameState` obtainable on demand — §3.5.
3. Per-step state hash — §3.5.
4. Input log: the exact `[2]input.GameInput` per step is the replayable
   input stream, owned by `GameState` as an append-only `InputLog`
   (snapshot-exempt — it replays from the initial state — but
   test-visible). The 30-frame ring (§5.3) is for intent detection only,
   not replay.
5. All match-affecting counters (round, timer, HP, snapshot-visible UI
   state) live inside `GameState`, never in scene/UI code.

Rollback plan (future): rewind to a snapshot, re-simulate the buffered input
frames, resync UI from state. Nothing here requires implementation in v1.

---

## 11. Engineering rules

1. **Idiomatic Go, minimal ceremony.** No reflection, no interface soup.
   Interfaces only where a second implementation exists or is proven
   needed.
2. **Data-driven discipline (P2).** If a value differs per character, it
   belongs in the character file. If it differs per match, it belongs in
   `GameState`. Go constants are for *engine* physics (§8.1) and *engine*
   behavior only.
3. **No simulation state in scenes/UI.** Scenes and HUD may read
   `GameState`; only the §4.1 pipeline may write it.
4. **Strict data loading** (§6.8): load-time validation with actionable
   errors (file + field named).
5. **Tests (mandatory, `[todo]` F0 onward):**
   - Replay determinism test with golden state hash (§3.6).
   - Character file round-trip test (load → validate → save → load,
     hash-equal).
   - Input intent/sequence tests (buffering, suppression, priority order).
   - Compile-level tests are not sufficient for `gameplay`/`input` changes.
6. **Docs stay in sync.** AGENTS.md tracks build/architecture facts; SPEC.md
   tracks behavior. A PR that changes behavior updates SPEC.md in the same
   PR.
7. **Go toolchain:** pinned to the `go.mod` version (Go + ebiten); upgrading
   either is a determinism-scope event (re-run golden hash tests).

---

## 12. Roadmap — phases and status

Phases are ordered; each phase's acceptance criteria must pass before the
next phase starts (a phase may be *started* early only on explicit decision).
(`lista de tarefas.md` is the Portuguese working task list and is subordinate
to this section — on conflict, §12 wins.)

### F0 — Foundations: determinism & data contract `[todo]` — **current**

1. Remove map-iteration nondeterminism: ordered sequence table with fixed
   priority (§5.5), and an audit that no other map iteration reaches the
   simulation (§3.3).
2. Deterministic PRNG owned by `GameState` (SplitMix64, seeded per match,
   §3.4); state snapshot + per-step hash (§3.5) and golden replay test
   (§3.6). Snapshot includes `StunFrames`, the connect ledger, input
   histories, PRNG state, and match counters (§7.1, §7.6, §7.7).
3. Adopt the data contract (§6): `properties.maxHP` in file; strict load
   validation (§6.8, incl. the `loopFrames`-on-reactions rule); canonical
   name checks.
4. TOML loader/writer (pick an actively maintained TOML 1.0 Go library
   and validate quoted numeric keys like `"44"` against it); convert
   `PlaceHolder` (and all assets) to TOML; remove the YAML loader; editor
   save path updated to write TOML (editor itself remains non-goal,
   §1.4). Extend the placeholder character to the **full** §6.7 set: it
   currently lacks jumps (`7`/`8`/`9`), crouch (`1`/`2`/`3`), `fall`,
   `landing`, `44`, and all reaction/round states (`hurt`, `crouchHurt`,
   `airHurt`, `blockHit`, `crouchBlock`, `airBlock`, `ko`, `win`,
   `knockdown`, `getup`). Special motions only as implemented (none
   required initially). Minimal 1–2 frame placeholder animations are
   acceptable until real content replaces them. Held states (`idle`,
   walks) need valid `loopFrames`.
5. Pushbox: all collision boxes per frame (§6.6).
6. Input: ordered sequence table covering dashes (`44`, `66`) and the
   full standard motion set with exact patterns (§5.5, §6.7), replacing
   the map outright (including the non-standard `426A` entry); SOCD
   filtering applies after the per-player device merge (§5.2).
- **Acceptance:** replay test green; `PlaceHolder.toml` loads, validates,
  plays; `config.toml` replaces `config.yaml`; `EN.toml`/`BR.toml`
  replace the language YAML. No YAML remains.
- **Note:** the game relies on ebiten's default 60 TPS; no tick-rate
  pinning is required (§3.1).

### F1 — Finish hit application `[todo]`

- Requires a checked-in, valid, self-contained test character with
  hittable frames (placeholder fixtures with missing sprites/boxes don't
  qualify).
- Apply hit effects per §7.2 (damage, knockback/knockup, pushback,
  one-hit-per-frame rule).
- Trade behavior per §7.3.
- **Acceptance:** a scripted A-attack reduces defender HP and applies
  knockback/pushback displacement; both hits in a trade land.
  (Reaction-state entry belongs to F2.)

### F2 — Life & reactions `[todo]`

- HP from `properties.maxHP` (replaces hardcoded 10000), HP clamping,
  KO detection at 0.
- `hurt` states with proper frame data for the test character;
  invincibility and armor honored.
- Stun-hold via `loopFrames` on reaction states (runtime must accept
  `start == end` as single-frame hold, §6.4, §7.6); `StunFrames`
  countdown per §4.1/§7.2.
- **Acceptance:** hit → HP drop → `hurt` → back to `idle` cycle is stable
  and repeatable in replay.

### F3 — Guard `[todo]`

- Guard rule §7.4; `blockHit` state; blockstun.
- **Acceptance:** holding back turns hits into blockstun with no damage,
  grounded and airborne (into `blockHit`/`crouchBlock`/`airBlock`);
  walking into an attack without holding back gets hit.

### F4 — Frame-data effects `[todo]`

- Wire remaining fields: `isRecovery`, `isHoldable`, `animationSwitch`,
  `canHardKnockdown`, `canWallBounce`, `canGroundBounce`, `canOTG`,
  `priority`.
- Bounce mechanics need per-character launch flags (e.g.
  `hasBouncedThisLaunch`, snapshot state, reset on neutral ground
  contact) and a defined bounce frame source — specify before
  implementing.
- **Acceptance:** each field demonstrably changes runtime behavior in a
  unit test.

### F5 — Rounds, KO, match end `[todo]`

- Timer, round/match flow §7.7, `win`/`ko` states, scene transition out.
- **Acceptance:** a full 2-round match ends and returns to the menu.

### F6 — Fight HUD `[todo]`

- §7.8 HUD, localized.
- **Acceptance:** HP, rounds, timer readable without debug overlays.

### F7 — End-to-end validation `[todo]`

- A dedicated test character (full required set plus at least one
  implemented special motion) and a scripted scenario exercising: walk,
  jump, normal, special, hit, block, trade, KO, round end, match end.
- **Acceptance:** scenario passes unattended in CI; replay hashes stable.

---

## 13. Open questions (decide before the owning phase starts)

1. **Chip damage** — v1 assumes 0 chip on block (§7.4). Confirm.
2. **Throw** — no throw state/input in the model yet. Out of v1; define
   trigger + states later.
3. **Sprite sheets** — v1 uses per-frame image files; sheet support is an
   editor/format extension, not yet contract.
4. **Weight** — pushbox weighting is velocity-based (§8.4). A per-character
   weight property is a candidate for §6.3 later.
5. **Hitstop/hitlag** — explicitly out of v1; likely a later feel phase.
6. **Crouch/air normals** — the input model maps Down+A and air+A to the
   same `A` intents as standing, so normals are single-posture in v1.
   Per-posture normals (`2A`-style crouch normals, air normals) are an
   input-model + state-table extension; decide before F4 whether they
   enter the contract.

---

*End of specification.*
