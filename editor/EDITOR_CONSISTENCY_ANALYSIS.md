# Editor Package Consistency Analysis

**Date:** November 3, 2025  
**Analyzed by:** AI Code Review  
**Package:** `fgengine/editor`

---

## Executive Summary

The editor package is generally well-structured but has several consistency issues, logic problems, and areas where it deviates from its own patterns. This document catalogs all findings with severity ratings and recommendations.

---

## 1. CRITICAL ISSUES

### 1.1 AnimationPlayer Not Updated in Editor Loop
**Location:** `editorMain.go`, `updateAnimationPlayback()`  
**Severity:** CRITICAL

**Problem:**  
The editor manually manages `frameIndex` and `frameCounter` but the `AnimationPlayer` (which is initialized on the character) is never updated. The comments say "AnimationPlayer automatically handles sprite selection" but this is misleading - the AnimationPlayer is never actually being used in the editor.

```go
// In editorMain.go:updateAnimationPlayback()
g.editorManager.frameCounter++
g.editorManager.frameIndex++
// Comment says: "AnimationPlayer automatically handles sprite selection"
// But AnimationPlayer.FrameCounter is never updated!
```

**Impact:**  
- The character's AnimationPlayer state doesn't match the editor's state
- When drawing, the character might not display correctly if graphics.Draw uses AnimationPlayer
- Inconsistency between editor preview and actual game behavior

**Recommendation:**  
Either:
1. Sync `g.activeCharacter.AnimationPlayer.FrameCounter` with `g.editorManager.frameCounter`
2. OR use AnimationPlayer directly instead of separate frameIndex/frameCounter tracking
3. Remove misleading comments

---

### 1.2 Sprite Index Out of Bounds Risk
**Location:** `manager.go`, `getCurrentSprite()`  
**Severity:** CRITICAL

**Problem:**  
The function checks bounds but doesn't handle the case where FrameData references a SpriteIndex that doesn't exist:

```go
func (e *EditorManager) getCurrentSprite() *animation.Sprite {
    // Checks frameIndex bounds...
    spriteIndex := e.activeAnimation.FrameData[e.frameIndex].SpriteIndex
    if spriteIndex < 0 || spriteIndex >= len(e.activeAnimation.Sprites) {
        return nil
    }
    return e.activeAnimation.Sprites[spriteIndex]
}
```

**Missing Validation:**  
When frames are removed in `uiTimeline.go:removeFrame()`, only FrameData is removed, but the Sprites array is never cleaned up. This can lead to:
- Orphaned sprites consuming memory
- SpriteIndex values becoming invalid after sprite removal
- No way to remove sprites from the Sprites array

**Recommendation:**  
1. Add sprite cleanup logic when removing frames
2. Add validation when setting SpriteIndex in UI
3. Consider adding a "Remove Unused Sprites" maintenance function

---

### 1.3 Box Editor State Inconsistency
**Location:** `box.go`, Multiple functions  
**Severity:** HIGH

**Problem:**  
BoxEditor maintains its own `boxes` map that points to sprite boxes, but this reference can become stale:

```go
type BoxEditor struct {
    boxes map[collision.BoxType][]types.Rect  // This is a reference
    // ...
}

// In loadBoxEditor:
g.editorManager.boxEditor = &BoxEditor{
    boxes: sprite.Boxes,  // Direct reference assignment
    // ...
}
```

When frames change, `refreshBoxEditor()` updates the reference, but there's no guarantee this happens consistently.

**Issues:**
- If `refreshBoxEditor()` isn't called, boxEditor points to old frame's boxes
- Editing boxes affects wrong frame
- `refreshBoxEditor()` is only called in timeline frame change, not elsewhere

**Found Inconsistencies:**
- `uiTimeline.go`: Calls `refreshBoxEditor()` after frame change ✓
- `box.go:deleteSelectedBox()`: Doesn't refresh - could affect wrong frame ✗
- `box.go:addBox()`: Doesn't refresh after adding ✗
- `uiProject.go`: When switching animations, sets boxEditor to nil but doesn't call refresh ✗

**Recommendation:**  
1. Call `refreshBoxEditor()` consistently after any frame change
2. OR: Always get boxes directly from `getCurrentSprite()` instead of caching
3. Add `validateBoxEditor()` function called before any box operation

---

## 2. LOGIC ERRORS

### 2.1 Animation Playback Loop Logic
**Location:** `editorMain.go:updateAnimationPlayback()`  
**Severity:** MEDIUM

**Problem:**  
The animation playback doesn't respect the animation's total duration or handle looping properly:

```go
if g.editorManager.frameCounter >= currentFrameDuration {
    g.editorManager.frameCounter = 0
    g.editorManager.frameIndex++
    
    if g.editorManager.frameIndex >= len(g.editorManager.activeAnimation.FrameData) {
        g.editorManager.frameIndex = 0  // Simple loop
    }
}
```

**Issues:**
- Doesn't respect `AnimationPlayer.ShouldLoop` setting
- Game logic uses `AnimationPlayer.GetActiveFrameData()` which has sophisticated looping
- Editor behavior differs from actual game behavior

**Recommendation:**  
Use the actual AnimationPlayer logic:
```go
func (g *Game) updateAnimationPlayback() {
    if !g.editorManager.playingAnim || g.activeCharacter == nil {
        return
    }
    g.activeCharacter.AnimationPlayer.FrameCounter++
    // Update editor display from AnimationPlayer state
}
```

---

### 2.2 Frame Duration Minimum Not Enforced Everywhere
**Location:** `uiTimeline.go`, `sprite.go`  
**Severity:** LOW

**Problem:**  
In `uiTimeline.go`, duration is clamped to minimum 1:
```go
if duration < 1 {
    duration = 1
}
```

But in `sprite.go:addSpriteByFile()`:
```go
newFrameData := animation.FrameData{
    Duration: 1,  // Default, but no validation
    SpriteIndex: len(e.activeAnimation.Sprites) - 1,
}
```

And `newSpriteFromImage()` creates default duration 60, while `newAnimationFileDialog()` also uses 60.

**Inconsistency:**
- `addSpriteByFile()`: Creates frames with duration 1
- `newAnimationFileDialog()`: Creates frames with duration 60
- No clear default policy

**Recommendation:**  
1. Define a constant `DefaultFrameDuration = 60` 
2. Use it consistently across all frame creation
3. Add validation function for FrameData

---

### 2.3 Character Loading Animation State Mismatch
**Location:** `character.go:loadCharacter()`  
**Severity:** MEDIUM

**Problem:**  
Complex initialization logic with redundant operations:

```go
func (g *Game) loadCharacter() {
    // ...
    character.AnimationPlayer = &animation.AnimationPlayer{}
    
    // Set the initial animation if available
    if len(character.Animations) > 0 {
        if idleAnim, exists := character.Animations["idle"]; exists {
            character.AnimationPlayer.ActiveAnimation = idleAnim
        } else {
            for _, anim := range character.Animations {
                character.AnimationPlayer.ActiveAnimation = anim
                break
            }
        }
    }
    
    // Set initial sprite if there's an animation available
    if len(character.Animations) > 0 {
        for _, anim := range character.Animations {
            if len(anim.FrameData) > 0 && len(anim.Sprites) > 0 {
                // AnimationPlayer will handle sprite selection automatically
                break
            }
        }
    }
    
    // Later...
    idleAnim, ok := character.Animations["idle"]
    if !ok {
        idleAnim = g.createPlaceholderIdleAnimation()
        character.Animations["idle"] = idleAnim
    }
    g.editorManager.activeAnimation = idleAnim
}
```

**Issues:**
1. Sets `AnimationPlayer.ActiveAnimation` early
2. Second loop does nothing (just breaks)
3. Then overwrites by setting `editorManager.activeAnimation = idleAnim`
4. Creates placeholder idle if missing, but already tried to set idle above

**Recommendation:**  
Simplify to:
```go
func (g *Game) loadCharacter() {
    g.checkIfResetNeeded()
    character, err := loadCharacterFromYAMLDialog()
    if err != nil {
        g.writeLog("Failed to load character: " + err.Error())
        return
    }
    
    g.activeCharacter = character
    g.activeCharacter.StateMachine = &state.StateMachine{}
    g.activeCharacter.AnimationPlayer = &animation.AnimationPlayer{}
    
    if character.Animations == nil {
        character.Animations = make(map[string]*animation.Animation)
    }
    
    // Get or create idle animation
    idleAnim, ok := character.Animations["idle"]
    if !ok {
        g.writeLog("No 'idle' animation found, creating placeholder...")
        idleAnim = g.createPlaceholderIdleAnimation()
        character.Animations["idle"] = idleAnim
    }
    
    // Set animation in both places
    g.editorManager.setActiveAnimation(idleAnim)
    g.activeCharacter.AnimationPlayer.ActiveAnimation = idleAnim
    
    g.writeLog("Character loaded successfully")
}
```

---

## 3. CONSISTENCY VIOLATIONS

### 3.1 Inconsistent Reset Pattern
**Location:** Multiple files  
**Severity:** LOW

**Problem:**  
When switching animations or frames, different reset patterns are used:

```go
// Pattern 1: In manager.go:setActiveAnimation()
e.frameIndex = 0
e.frameCounter = 0
e.playingAnim = false
e.boxEditor = nil

// Pattern 2: In uiProject.go (animation dropdown)
g.editorManager.frameIndex = 0
g.editorManager.frameCounter = 0
g.editorManager.boxEditor = nil
g.refreshBoxEditor()  // Extra step

// Pattern 3: In uiTimeline.go (frame slider)
g.editorManager.frameIndex = frameIndex
g.editorManager.frameCounter = 0
g.refreshBoxEditor()  // Extra step, no boxEditor = nil
```

**Issues:**
- `setActiveAnimation()` doesn't call `refreshBoxEditor()`
- Some places set `boxEditor = nil`, others call `refreshBoxEditor()`
- Inconsistent about stopping playback (`playingAnim = false`)

**Recommendation:**  
Create unified state management:
```go
func (e *EditorManager) setActiveAnimation(anim *animation.Animation) {
    e.activeAnimation = anim
    e.resetPlayback()
    e.boxEditor = nil
}

func (e *EditorManager) setFrame(index int) {
    e.frameIndex = index
    e.frameCounter = 0
    // Don't reset boxEditor, just refresh
}

func (e *EditorManager) resetPlayback() {
    e.frameIndex = 0
    e.frameCounter = 0
    e.playingAnim = false
}
```

---

### 3.2 Logging Inconsistency
**Location:** Multiple files  
**Severity:** LOW

**Problem:**  
Log messages have inconsistent format and detail level:

```go
// Good examples:
g.writeLog("Character loaded successfully")
g.writeLog(fmt.Sprintf("Added %s box", boxType.String()))

// Minimal examples:
g.writeLog("New character created")  // No details

// Over-detailed:
g.writeLog(fmt.Sprintf("Animation name changed from '%s' to '%s'", old, new))
g.writeLog(fmt.Sprintf("Current animations: %s", g.getAnimationNames()))

// Missing logs:
// - When copying frames (copyLastFrame)
// - When removing frames (removeFrame)
// - When frame duration changes
// - When box properties change
```

**Recommendation:**  
1. Log all state-changing operations
2. Use consistent format: "Action completed: details"
3. Consider log levels (INFO, DEBUG, ERROR)

---

### 3.3 Error Handling Inconsistency
**Location:** Multiple files  
**Severity:** MEDIUM

**Problem:**  
Some functions panic, some return errors, some silently fail:

```go
// Pattern 1: Return error (good)
func (e *EditorManager) addSpriteByFile(path string) error {
    if e == nil || e.activeAnimation == nil {
        return fmt.Errorf("no active animation available")
    }
    // ...
}

// Pattern 2: Silent fail (bad)
func (g *Game) deleteSelectedBox() {
    activeBox := g.getActiveBox()
    if activeBox == nil {
        return  // No error, no log
    }
    // ...
}

// Pattern 3: Check and log (good)
func (g *Game) addBox() {
    if g.editorManager.activeAnimation == nil {
        g.writeLog("Cannot add box: No active animation")
        return
    }
    // ...
}

// Pattern 4: Panic in game code (character.go)
if !ok {
    panic(fmt.Sprintf("Animation '%s' not found...", name))
}
```

**Recommendation:**  
Establish error handling pattern for editor:
1. Always log errors to the editor log
2. Never panic in editor code (use error returns)
3. UI actions should always provide feedback

---

### 3.4 Nil Check Patterns
**Location:** Multiple files  
**Severity:** LOW

**Problem:**  
Inconsistent nil checking patterns:

```go
// Pattern 1: Check Game fields
if g.editorManager.boxEditor == nil || g.editorManager.activeAnimation == nil {
    return
}

// Pattern 2: Check via helper
sprite := g.editorManager.getCurrentSprite()
if sprite == nil {
    ctx.Text("No sprite selected")
    return
}

// Pattern 3: Nested checks
if g.activeCharacter != nil && g.editorManager.activeAnimation != nil {
    // ...
}

// Pattern 4: No check (risky)
// In some UI code, assumes activeAnimation is not nil
g.editorManager.activeAnimation.Name = newName  // Could panic
```

**Recommendation:**  
1. Always check before dereferencing
2. Use early returns for cleaner code
3. Consider helper method: `isEditorReady() bool`

---

## 4. CODE STRUCTURE ISSUES

### 4.1 Mixed Responsibilities in Game Struct
**Location:** `editorMain.go`  
**Severity:** MEDIUM

**Problem:**  
The `Game` struct mixes several concerns:

```go
type Game struct {
    debugui         debugui.DebugUI
    activeCharacter *character.Character
    editorManager   *EditorManager
    inputManager    *input.InputManager
    lastMouseX      int          // Camera drag state
    lastMouseY      int          // Camera drag state
    isDragging      bool         // Camera drag state
    camera          *graphics.Camera
}
```

Camera drag state is in `Game`, but box drag state is in `BoxEditor`.

**Recommendation:**  
Move camera-related state to a `CameraController` struct:
```go
type CameraController struct {
    camera      *graphics.Camera
    lastMouseX  int
    lastMouseY  int
    isDragging  bool
}

func (cc *CameraController) HandleInput() {
    // Move handleCameraInput logic here
}
```

---

### 4.2 EditorManager Should Own More State
**Location:** `editorMain.go`, `manager.go`  
**Severity:** MEDIUM

**Problem:**  
Game has direct access to many editor states:
- `g.editorManager.frameIndex` accessed everywhere
- `g.editorManager.playingAnim` modified by UI and update loop
- Box editor methods are on Game, not EditorManager

**Recommendation:**  
Move editor operations to EditorManager:
```go
// Instead of:
g.editorManager.frameIndex++

// Use:
g.editorManager.NextFrame()

// Instead of:
if g.editorManager.playingAnim { ... }

// Use:
g.editorManager.Update()  // Handles playback internally
```

---

### 4.3 UI Files Have Business Logic
**Location:** `uiTimeline.go`, `uiProject.go`  
**Severity:** MEDIUM

**Problem:**  
UI callback functions contain business logic:

```go
// In uiTimeline.go
ctx.Button("Remove Frame").On(func() {
    g.removeFrame()  // This is in uiTimeline.go
})

func (g *Game) removeFrame() {
    // Business logic for removing frames
    if g.editorManager.activeAnimation == nil || len(g.editorManager.activeAnimation.FrameData) == 0 {
        return
    }
    // 15 lines of logic...
}
```

**Recommendation:**  
Move business logic to separate files:
- `editorOperations.go` - File operations (save, load, export)
- `animationOperations.go` - Animation/frame manipulation
- `boxOperations.go` - Box manipulation

Keep UI files only for UI layout and simple callbacks.

---

## 5. MISSING FEATURES / VALIDATION

### 5.1 No Undo/Redo System
**Location:** N/A  
**Severity:** LOW

**Problem:**  
Editor has no undo/redo capability. Accidental deletions can't be recovered.

**Recommendation:**  
Consider implementing a command pattern for reversible operations.

---

### 5.2 No Dirty State Tracking
**Location:** N/A  
**Severity:** LOW

**Problem:**  
No indication if character/animation has unsaved changes.

**Recommendation:**  
Add dirty flag and warn on exit/load without saving.

---

### 5.3 No Animation/Frame Validation
**Location:** Various  
**Severity:** MEDIUM

**Problem:**  
No validation that animations are well-formed:
- Can have frames with invalid SpriteIndex
- Can have empty animations
- Can have orphaned sprites
- No validation of frame properties (negative durations, etc.)

**Recommendation:**  
Add validation function called before save:
```go
func (a *Animation) Validate() []string {
    var errors []string
    if len(a.FrameData) == 0 {
        errors = append(errors, "Animation has no frames")
    }
    for i, fd := range a.FrameData {
        if fd.SpriteIndex < 0 || fd.SpriteIndex >= len(a.Sprites) {
            errors = append(errors, fmt.Sprintf("Frame %d has invalid sprite index", i))
        }
        if fd.Duration < 1 {
            errors = append(errors, fmt.Sprintf("Frame %d has invalid duration", i))
        }
    }
    return errors
}
```

---

### 5.4 No Sprite Reuse System
**Location:** `sprite.go`, `uiTimeline.go`  
**Severity:** LOW

**Problem:**  
Every frame creates a new sprite, even if the image is the same. This is wasteful for animations that hold a pose for multiple frames.

The `copyLastFrame()` function does copy the sprite, but there's no way to have multiple frames reference the same sprite intentionally.

**Current Design:**
```
Animation {
    Sprites: [S0, S1, S2, S3]  // Each sprite is a copy
    FrameData: [
        {SpriteIndex: 0, Duration: 10},
        {SpriteIndex: 1, Duration: 10},
        {SpriteIndex: 2, Duration: 10},
        {SpriteIndex: 3, Duration: 10},
    ]
}
```

**Better Design:**
```
Animation {
    Sprites: [S0, S1]  // Unique sprites only
    FrameData: [
        {SpriteIndex: 0, Duration: 10},
        {SpriteIndex: 0, Duration: 10},  // Reuse S0
        {SpriteIndex: 1, Duration: 10},
        {SpriteIndex: 1, Duration: 10},  // Reuse S1
    ]
}
```

**Recommendation:**  
Add UI option: "Copy Frame (new sprite)" vs "Duplicate Frame (reuse sprite)"

---

## 6. NAMING INCONSISTENCIES

### 6.1 Frame vs Sprite Terminology
**Location:** Multiple files  
**Severity:** LOW

**Problem:**  
Confusing terminology mixing "frame" and "sprite":

```go
getCurrentSprite()           // Returns sprite for current frame
frameIndex                   // Index into FrameData
FrameData[i].SpriteIndex    // Index into Sprites array
```

Timeline UI says "Frame Properties" but shows sprite+framedata properties.

**Clarification Needed:**
- **Sprite**: An image + collision boxes (static data)
- **FrameData**: Timing + game properties + reference to a sprite
- **Frame**: The combination of FrameData + its referenced Sprite

**Recommendation:**  
Use consistent terminology in comments and variable names:
- `currentFrameIndex` instead of `frameIndex`
- `getActiveFrame()` that returns both FrameData and Sprite
- UI: "Frame Properties" tab with "Sprite" and "Frame Data" sections

---

### 6.2 boxActionIndex vs boxEditor.activeBoxType
**Location:** `box.go`, `manager.go`  
**Severity:** LOW

**Problem:**  
Two variables track box type:
```go
type EditorManager struct {
    boxActionIndex int  // Used by UI dropdown
}

type BoxEditor struct {
    activeBoxType collision.BoxType  // Used by editor logic
}

// They're synced in multiple places:
g.editorManager.boxActionIndex = int(selectedBoxType)
g.editorManager.boxEditor.activeBoxType = collision.BoxType(g.editorManager.boxActionIndex)
```

**Recommendation:**  
Remove `boxActionIndex`, use only `boxEditor.activeBoxType`. Convert to int when needed for UI.

---

## 7. POSITIVE OBSERVATIONS

### 7.1 Good Separation of Concerns in Files
The file structure is logical:
- `editorMain.go` - Entry point and main loop
- `manager.go` - State management
- `box.go`, `sprite.go`, `character.go` - Feature-specific logic
- `ui*.go` - UI layout
- `serialization.go` - I/O operations

### 7.2 Good Use of Helper Functions
Functions like `getCurrentSprite()`, `getActiveBox()`, `refreshBoxEditor()` provide clean abstraction.

### 7.3 Consistent Deep Copy Pattern
The serialization code properly deep copies animations and sprites to avoid reference issues during save.

### 7.4 Good Camera Controls
Camera panning and scaling are well-implemented with proper coordinate transformations.

### 7.5 Proper YAML Path Resolution
`makeRelativePath()` and `resolveRelativePath()` ensure portability of saved files.

---

## 8. RECOMMENDED ACTION PLAN

### Phase 1: Critical Fixes (High Priority)
1. ✅ Fix AnimationPlayer synchronization in editor
2. ✅ Add sprite cleanup when removing frames
3. ✅ Make `refreshBoxEditor()` consistently called
4. ✅ Add animation validation before save

### Phase 2: Logic Improvements (Medium Priority)
5. ✅ Simplify `loadCharacter()` initialization
6. ✅ Unify frame/animation reset pattern
7. ✅ Add consistent error handling and logging
8. ✅ Extract camera controller

### Phase 3: Structure Refactoring (Low Priority)
9. ⚠️ Move business logic out of UI files
10. ⚠️ Improve EditorManager encapsulation
11. ⚠️ Clean up naming inconsistencies
12. ⚠️ Add sprite reuse system

### Phase 4: Features (Future)
13. 💡 Add undo/redo system
14. 💡 Add dirty state tracking
15. 💡 Add animation preview in project panel

---

## 9. SPECIFIC CODE RECOMMENDATIONS

### 9.1 Fix AnimationPlayer Sync
Add to `editorMain.go:updateAnimationPlayback()`:

```go
func (g *Game) updateAnimationPlayback() {
    if !g.editorManager.playingAnim || g.editorManager.activeAnimation == nil {
        return
    }
    
    g.editorManager.frameCounter++
    
    // Sync with character's AnimationPlayer
    if g.activeCharacter != nil && g.activeCharacter.AnimationPlayer != nil {
        g.activeCharacter.AnimationPlayer.FrameCounter = g.editorManager.frameCounter
    }
    
    // Rest of existing logic...
}
```

### 9.2 Add Validation Helper
Add to `manager.go`:

```go
func (e *EditorManager) ValidateAnimation() []string {
    if e.activeAnimation == nil {
        return []string{"No active animation"}
    }
    
    var errors []string
    
    if len(e.activeAnimation.FrameData) == 0 {
        errors = append(errors, "Animation has no frames")
    }
    
    for i, fd := range e.activeAnimation.FrameData {
        if fd.Duration < 1 {
            errors = append(errors, fmt.Sprintf("Frame %d: duration must be >= 1", i))
        }
        if fd.SpriteIndex < 0 || fd.SpriteIndex >= len(e.activeAnimation.Sprites) {
            errors = append(errors, fmt.Sprintf("Frame %d: invalid sprite index %d", i, fd.SpriteIndex))
        }
    }
    
    return errors
}
```

### 9.3 Consistent Box Editor Refresh
Replace all instances of setting `boxEditor = nil` with:

```go
func (g *Game) ensureBoxEditorValid() {
    sprite := g.editorManager.getCurrentSprite()
    if sprite == nil {
        g.editorManager.boxEditor = nil
        return
    }
    
    if g.editorManager.boxEditor == nil {
        g.loadBoxEditor(sprite)
    } else {
        g.refreshBoxEditor()
    }
}
```

Call before any box operation.

---

## CONCLUSION

The editor package is functional but has several areas where it deviates from its own patterns and best practices. The most critical issues are:

1. **AnimationPlayer desynchronization** - The editor's preview may not match actual game behavior
2. **Box editor state management** - Can edit wrong frame's boxes
3. **Lack of validation** - Can create invalid animations
4. **Inconsistent state management** - Different patterns for similar operations

The codebase shows good structure overall, but would benefit from consolidation of state management patterns and more defensive programming practices.

**Overall Assessment:** 7/10
- Structure: 8/10 ✓
- Consistency: 6/10 ⚠️
- Logic Correctness: 7/10 ⚠️
- Error Handling: 6/10 ⚠️
- Maintainability: 7/10 ✓

