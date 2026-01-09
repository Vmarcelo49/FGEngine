# Fighting Game Implementation Notes

**Date:** January 9, 2026  
**Goal:** Transform the FGEngine skeleton into a playable 2D fighting game

---

## What Was Already There

The engine had a solid foundation with several systems partially implemented:

### Core Systems (Pre-existing)
- **State Machine**: Bitflag-based state system with composite states (idle, walk, dash, jump, attack, hitstun, etc.)
- **Animation System**: Frame-based animation player with sprite management and frame data
- **Graphics**: Camera system with world-to-screen transforms, render queue with layering, image caching
- **Input**: GameInput bitflags, input sequences for special moves, SOCD cleaning, gamepad support
- **Character**: YAML-based character loading, basic sprite rendering
- **Stage**: Grid/solid/image stage types with proper rendering
- **Types**: Rect and Vector2 with collision helpers
- **Constants**: World bounds, camera size, gravity, ground level

### What Was Missing/Incomplete
- No actual game loop connecting input → logic → update → collision
- Character movement was hardcoded and not request-based
- No hit detection or damage system
- No player-vs-player setup
- Missing character asset files (helmet.yaml didn't exist)
- Input handling mutated states directly instead of queuing requests
- No hitstun recovery logic
- Camera wasn't following players
- Scene manager was stubbed but never wired up properly

---

## Implementation Changes

### 1. State Machine Refactor
**Files Modified:** `state/stateMachine.go`, `state/stateGroups.go`

#### Added Request-Based Input Queue
```go
MoveInput       int  // -1 left, 0 idle, 1 right
JumpRequested   bool
AttackRequested bool
DashRequested   bool
HitstunFrames   int  // countdown timer
```

**Why:** Direct state mutation in `HandleInput` caused characters to get stuck in states. Now inputs are *requests* that the character update loop processes based on current state.

#### Added ClearState Helper
```go
func (sm *StateMachine) ClearState(flags State) {
    sm.PreviousState = sm.ActiveState
    sm.ActiveState &^= flags
}
```

**Why:** `RemoveState` had guardrails that auto-toggled grounded/airborne, which prevented clean transitions during hitstun or attacks.

#### Input Processing Flow
1. Reset all request flags each frame
2. Check for input sequences (dash)
3. Set directional MoveInput (-1, 0, 1)
4. Queue jump/attack/dash flags
5. Guard with `IsInactable()` check **after** resets (prevents stale movement)

---

### 2. Character Update Logic
**File Modified:** `character/drawAndUpdate.go`

Completely rewrote the `Update()` method to process queued inputs:

#### Update Order
1. **Hitstun countdown** → clear hitstun state when timer expires
2. **Attack processing** → start attack animation, reset hit flag
3. **Dash/Walk movement** → apply speed based on MoveInput and facing
4. **Jump processing** → set airborne, apply jump velocity
5. **Gravity** → apply constant downward force when airborne
6. **Animation advance** → update frame timers
7. **Frame velocity overrides** → apply ChangeXSpeed/ChangeYSpeed from frame data
8. **Position update** → apply velocity to position
9. **Falling detection** → set StateFalling when velocity.Y > 0
10. **Ground collision** → clamp to ground level, clear airborne states
11. **Wall collision** → clamp to world bounds
12. **Animation selection** → pick animation based on final state

#### Key Features
- **Movement speed:** 3.2 units/frame walk, 6.0 units/frame dash
- **Friction:** Applied when MoveInput is 0, decays velocity gradually
- **Facing-relative movement:** MoveInput direction respects CharacterOrientation
- **Animation guards:** `ensureAnimation()` prevents restarting same animation

---

### 3. Character Assets & Defaults
**Files Modified:** `character/character.go`  
**File Created:** `assets/characters/helmet.yaml`

#### Added Defaults in `initialize()`
```go
if c.Friction == 0 {
    c.Friction = 0.65
}
if c.JumpHeight == 0 {
    c.JumpHeight = 14
}
c.StateMachine.HP = 100
```

#### Created Minimal Character Definition
The helmet.yaml defines 6 animations (idle, walk, dash, jump, fall, attack) all using the same sprite placeholder. Attack has:
- **Startup:** 6 frames (Phase = 1)
- **Active:** 10 frames (Phase = 2, Damage = 12)

#### Sprite Size Probing
Added `ensureSpriteRect()` to auto-detect image dimensions if rect is unspecified, using `image.DecodeConfig()` for fast header parsing.

#### BoundingBox Helper
```go
func (c *Character) BoundingBox() types.Rect
```
Returns sprite-sized box at character position for collision checks.

---

### 4. Input System Improvements
**File Modified:** `input/mapping.go`

#### Added Custom Keyboard Bindings Constructor
```go
func NewInputManagerWithKeyboardBindings(bindings map[GameInput][]ebiten.Key) *InputManager
```
Allows per-player keyboard controls while keeping default gamepad mappings.

#### Added Poll Method
```go
func (im *InputManager) Poll() GameInput
```
Aggregates keyboard + assigned gamepad inputs for a single player, replaces the global `LocalInputsFromIDS`.

---

### 5. Main Game Loop
**File Modified:** `main.go`

#### Core Game Loop (Update)
```go
1. input.Update() // scene manager hook
2. updateFacings() // orient characters toward each other
3. Poll inputs from both players
4. logic.UpdateByInputs() // run state machines + character updates
5. resolveHits() // check for attack collisions
6. updateCamera() // center on both players
7. debugui.Update() // show stats
```

#### Two-Player Setup
- **P1 Controls:** WASD movement, FGH attacks
- **P2 Controls:** Arrow keys movement, NM, attacks
- Each player gets custom InputManager with distinct keyboard bindings
- Characters spawn 120 units apart from world center at ground level

#### Camera System
- Centers on midpoint between both players
- Locks to world bounds
- Debug zoom with Q/E, reset with R

#### Hit Detection (`applyHit`)
1. Check attacker has active frame (Phase = Active)
2. Check attacker.IsAttacking() and hasn't hit yet this swing
3. AABB overlap check: `attackerChar.BoundingBox().IsOverlapping(defender.BoundingBox())`
4. Apply damage, set hitstun duration (default 30 frames)
5. Apply knockback velocity based on relative positions
6. Clear defender's attack states, add StateOnHitsun
7. Mark `AttackHasHit = true` to prevent multi-hitting

---

### 6. Render Queue Setup
**File Modified:** `main.go`

Added all drawables to render queue:
- **LayerBG:** Stage (grid or image fallback)
- **LayerPlayer:** Both characters
- **LayerEffects:** BoxDrawable for both characters (debug hitboxes)

---

## Current Game Features

### ✅ Working
- **Two-player local versus** with independent keyboard controls
- **Movement:** Walk left/right with facing-relative speed
- **Dash:** Double-tap forward (66 input sequence)
- **Jump:** Press up while grounded, applies jump velocity
- **Attack:** Press A/F/N, plays attack animation
- **Hit detection:** Bounding box overlap during active frames
- **Damage & HP:** Deals damage, displays in debug UI
- **Hitstun:** Locks victim for 30 frames, prevents further actions
- **Knockback:** Launches victim backward and upward
- **Gravity & ground collision:** Characters fall and land properly
- **Camera tracking:** Follows both players, centers on midpoint
- **Animation system:** Smooth frame-based playback with looping
- **Facing:** Characters always face each other

### 🚧 Incomplete / TODO

#### High Priority
1. **Proper hitbox/hurtbox collision**
   - Current implementation uses sprite bounding box
   - Should use frame-specific boxes from `frameData.Boxes[Hit]` and `Boxes[Hurt]`
   - `collision/detection.go` is still a placeholder

2. **Input buffer**
   - No input buffering during hitstun/recovery
   - Should allow queuing next action before current one ends

3. **Attack animations**
   - All animations use single placeholder sprite
   - Need proper attack art with startup/active/recovery frames
   - Need multiple attacks (light/medium/heavy, standing/crouching/aerial)

4. **Special moves**
   - Input sequences exist (236A, 214A) but aren't wired to attacks
   - Need to detect sequences and trigger special states

5. **Blocking**
   - StateBlock exists but no input handling
   - Need to detect backward+attack for block

6. **Combo system**
   - No cancel windows or juggle states
   - Attack completion just returns to idle

7. **Dash mechanics**
   - Dash detection works but no dash animation duration
   - Should have startup/active/recovery phases
   - No backdash (44 sequence) yet

#### Medium Priority
8. **Sound system**
   - FrameData has `CommonAudioID` and `UniqueAudioID` fields
   - No audio playback implemented

9. **Stage hazards/bounds**
   - Characters stop at world bounds but no corner push
   - No wall bounce/ground bounce mechanics

10. **Health bars & UI**
    - HP shown in debug window only
    - Need proper HUD with health bars, timer, round counter

11. **Round system**
    - No win condition
    - No round start/end animations
    - StateWinAnimation/StateRoundStartAnimation exist but unused

12. **AI opponent**
    - Only local multiplayer implemented
    - Need basic CPU AI for single-player

#### Low Priority
13. **Scene manager integration**
    - Scene manager exists but game boots directly to match
    - No main menu, character select, or options

14. **Projectiles**
    - No projectile system
    - FrameData suggests support but not implemented

15. **Grab system**
    - StateGrabbing/StateGrabbed/StateGrabTech exist
    - No grab input or throw mechanics

16. **Replay/netplay**
    - logic.UpdateByInputs is deterministic (good for netplay)
    - No recording or network layer

---

## Known Issues

### Critical
- **No mirrored dash:** 66 works for both players but doesn't respect facing (always dashes "right")
- **Hitstun spam:** Attacker can hit defender again immediately after first hitstun expires if still overlapping
- **Sprite flipping:** Characters don't flip their sprites based on facing direction

### Minor
- **Debug hitboxes always show:** BoxDrawable should toggle with a debug key
- **Camera zoom persists:** Q/E zoom isn't very useful, should be removed or clamped
- **No pause:** Can't pause the game

---

## Architecture Notes

### Strengths
- **Bitflag states** allow clean composite state checking (e.g., `StateAirborne | StateAttack`)
- **Frame data system** is powerful and extensible (supports boxes, cancels, damage, etc.)
- **Request-based input** prevents state corruption and enables buffering
- **Deterministic update** makes replays/netplay feasible
- **Render queue layers** keep draw order clean

### Design Patterns Used
- **State pattern:** StateMachine with composite states
- **Component pattern:** Character has StateMachine, AnimationPlayer, Input
- **Observer pattern:** RenderQueue collects Drawables
- **Data-driven:** Characters defined in YAML with frame data

### Potential Improvements
- **ECS architecture:** Current OOP design works but ECS would scale better for many entities
- **Animation state machine:** Animation selection is a giant switch statement, could be data-driven
- **Collision layers:** All characters collide with all; need teams/projectiles
- **Input recording:** Store inputs in a ring buffer for replay analysis

---

## Testing Checklist

### Manual Tests Performed ✅
- [x] Both players can move left/right
- [x] Jump works and returns to ground
- [x] Dash triggers on double-forward
- [x] Attack animation plays on button press
- [x] Hit detection triggers on overlap during active frames
- [x] Damage is applied and HP decreases
- [x] Hitstun locks victim and launches them
- [x] Camera follows both players
- [x] Characters face each other
- [x] World bounds prevent falling off stage
- [x] Debug UI shows position, state, HP

### Not Tested ❌
- [ ] Gamepad input (only keyboard tested)
- [ ] Special move input sequences (236A, 214A)
- [ ] Multiple attacks in rapid succession
- [ ] Corner pressure behavior
- [ ] Simultaneous attacks (trades)
- [ ] Animation cancels
- [ ] Long hitstun strings

---

## File Change Summary

### Modified Files
- `main.go` - Complete game loop rewrite, two-player setup, hit detection
- `state/stateMachine.go` - Added request fields, hitstun timer, ClearState
- `state/stateGroups.go` - Request-based input handling
- `character/character.go` - Defaults, sprite probing, bounding box
- `character/drawAndUpdate.go` - Complete update rewrite with proper physics
- `input/mapping.go` - Poll method, custom keyboard bindings

### Created Files
- `assets/characters/helmet.yaml` - Minimal playable character definition
- `docs/fighting-game-implementation.md` - This document

### Unchanged Systems (Still Functional)
- `animation/` - Frame data, sprite management
- `graphics/` - Camera, render queue, image cache
- `collision/` - Type definitions (detection still stubbed)
- `types/` - Rect, Vector2 math
- `constants/` - World/camera dimensions
- `config/` - Window size, deadzone
- `stage/` - Background rendering

---

## Next Steps (Priority Order)

1. **Fix sprite flipping** - Add horizontal flip based on CharacterOrientation
2. **Implement proper hitbox collision** - Use FrameData.Boxes instead of sprite bounds
3. **Add blocking** - Detect backward input, add StateBlock handling
4. **Create more attacks** - Standing light/medium/heavy, crouching attacks
5. **Add health bars** - Proper UI overlay with HP bars
6. **Implement win condition** - Detect HP = 0, trigger round end
7. **Add attack variety** - Different damage, hitstun, knockback per attack
8. **Implement special moves** - Wire 236A/214A sequences to special attacks
9. **Add combo system** - Cancel windows, hit confirms
10. **Create proper character sprites** - Replace placeholder with animated art

---

## Conclusion

The engine foundation was excellent - the state system, animation framework, and rendering pipeline all worked as designed. The main gaps were:

1. **No deterministic input → logic → collision flow**
2. **Character movement logic was incomplete**
3. **No hit detection or damage system**
4. **Missing example character assets**

All of these have been addressed. The game is now **playable** with basic fighting mechanics (move, jump, dash, attack, take damage). The next phase is polish and depth - proper animations, more attacks, combos, blocking, and UI.

The codebase is in good shape for expansion. The request-based input system prevents state bugs, the frame data structure supports complex moves, and the deterministic update loop is netplay-ready.
