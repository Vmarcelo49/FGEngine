# Editor Integration Plan

## Overview
This document outlines the plan to integrate the animation/character editor from the older version into the current FGEngine project. The integration will preserve FGEngine's type system while incorporating the editor's functionality.

## Current State Analysis

### FGEngine Architecture (Target System)
- **Structure**: Modular packages (animation, collision, input, player, types)
- **Data Format**: YAML-based character/animation definitions
- **Types**: Clean separation with `types.Rect`, `types.Vector2`, `animation.Sprite`, etc.

### Editor Architecture (Source System)
- **Structure**: Single package with all functionality
- **UI**: Custom debugui-based interface for editing
- **Types**: Similar but extended types (`SpriteEx`, `Properties`, etc.)

## Type Consolidation Strategy

### ✅ Confirmed FGEngine Types to Preserve
| Type | Location | Status | Reason |
|------|----------|--------|---------|
| `Rect` | `types.Rect` | **PRESERVE** | Has helper methods (Right(), Bottom(), Center(), etc.) |
| `Vector2` | `types.Vector2` | **PRESERVE** | Clean, separate package |
| `Sprite` | `animation.Sprite` | **PRESERVE** | FGEngine canonical version |
| `Character` | `animation.Character` | **PRESERVE** | Clean separation of concerns |
| `Animation` | `animation.Animation` | **PRESERVE** | Proper modular design |
| `FrameProperties` | `animation.FrameProperties` | **PRESERVE** | More complete than editor's Properties |

### 🔄 Types Requiring Enhancement/Migration
| Editor Type | FGEngine Equivalent | Action Required |
|-------------|-------------------|-----------------|
| `SpriteEx` | `animation.Sprite` | ✅ **Use FGEngine version** - already compatible |
| `Properties` | `animation.FrameProperties` | Adapt editor UI to use FGEngine version |
| `AnimationState` | `animation.AnimationManager` | Adapt editor logic to FGEngine manager |
| `CharacterState` | `player.PlayerState` | Map editor state to player state |
| `StateMachine` | `player.StateMachine` | Evaluate consolidation needs |

## Integration Phases

### Phase 1: Dependency Resolution
Complete.
All dependencies already updated, and the editor folder was also created.

### Phase 2: Type Adapters and Migration Utilities
**Goal**: Create seamless conversion between editor's legacy types and FGEngine types

#### Actions:
1. **Create `editor/adapters/types.go`**
   - Convert `SpriteEx` → `animation.Sprite` (no-op, already compatible)
   - Convert `Properties` → `animation.FrameProperties`
   - Convert editor's `Character` → `animation.Character`
   - Convert editor's `Animation` → `animation.Animation`

2. **Create `editor/adapters/serialization.go`**
   - Adapt editor's YAML serialization to use FGEngine types
   - Maintain backward compatibility with existing editor files
   - Add migration utilities for old format → new format

### Phase 3: Editor Core Functionality Porting
**Goal**: Port editor's UI and editing functionality to work with FGEngine types

#### Actions:
1. **Port Box Rendering System**
   - Adapt `boxes.go` to work with `types.Rect` instead of editor's `Rect`
   - Update collision box editing to use FGEngine's collision system
   - Integrate with FGEngine's `collision.Box` and `collision.BoxType`

2. **Port Animation Editing**
   - Adapt `animation.go` to work with `animation.AnimationManager`
   - Update frame editing to use `animation.Sprite`
   - Integrate with FGEngine's animation system

3. **Port Character Editing**
   - Adapt `character.go` to work with `animation.Character`
   - Update serialization to match FGEngine's YAML format
   - Maintain editor's save/load functionality

4. **Port Properties Editing**
   - Adapt `properties.go` to work with `animation.FrameProperties`
   - Map editor's property types to FGEngine's property types
   - Update UI to reflect FGEngine's property structure

### Phase 4: UI Integration and Testing
**Goal**: Ensure the editor works seamlessly with FGEngine's data structures

#### Actions:
1. **Debugui Integration**
   - Create editor UI that can edit FGEngine characters
   - Implement real-time preview using FGEngine's rendering system
   - Add validation using FGEngine's type constraints

2. **State Machine Integration**
   - Evaluate whether to consolidate state machines
   - Create mapping between editor states and player states
   - Ensure state consistency between editor and engine

3. **Testing and Validation**
   - Test editor can load existing FGEngine character files
   - Test editor can save files compatible with FGEngine
   - Verify no data loss during edit cycles
   - Test animation playback matches FGEngine behavior

## File Structure After Integration

```
FGEngine/
├── animation/           # Core FGEngine animation system (preserved)
├── collision/          # Core FGEngine collision system (preserved)
├── input/              # Core FGEngine input system (preserved)
├── player/             # Core FGEngine player system (preserved)
├── types/              # Core FGEngine types (preserved)
├── editor/             # New editor integration
│   ├── go.mod         # Editor module depending on fgengine
│   ├── main.go        # Editor entry point
│   ├── ui/            # Debugui interface
│   │   ├── animation_ui.go
│   │   ├── character_ui.go
│   │   ├── properties_ui.go
│   │   └── boxes_ui.go
│   ├── adapters/      # Type conversion utilities
│   │   ├── types.go
│   │   └── serialization.go
│   └── assets/        # Editor-specific assets
├── assets/             # Shared game assets
└── main.go            # FGEngine main (preserved)
```

## Migration Path for Existing Data

### Editor Files → FGEngine Format
- **Character files**: Direct compatibility (same YAML structure)
- **Animation files**: Direct compatibility (same YAML structure)  
- **Sprite references**: Path resolution may need updates

### Backward Compatibility
- Editor will be able to load old editor-format files
- Editor will save in FGEngine-compatible format
- Migration utility for bulk conversion of old files

## Success Criteria

### ✅ Technical Requirements
- [ ] Editor builds and runs without dependency conflicts
- [ ] Editor can load existing FGEngine character files (e.g., `helmet.yaml`)
- [ ] Editor can save files that work with FGEngine's runtime
- [ ] No breaking changes to core FGEngine types
- [ ] All editor functionality (box editing, animation editing, properties) works

### ✅ User Experience Requirements
- [ ] Editor provides visual feedback using FGEngine's rendering
- [ ] Editor preserves all existing editing capabilities
- [ ] Editor workflow remains intuitive
- [ ] Real-time preview shows how animations will look in-engine

### ✅ Integration Requirements
- [ ] FGEngine can load editor-created files without modification
- [ ] Character/animation data flows seamlessly between editor and engine
- [ ] Type system remains clean and maintainable
- [ ] No duplicate code between editor and engine

## Risk Mitigation

### Version Conflicts
- **Risk**: Ebiten version differences cause compatibility issues
- **Mitigation**: Upgrade FGEngine to match editor's Ebiten version

### Type Conversion Issues
- **Risk**: Data loss during type conversion
- **Mitigation**: Comprehensive testing and validation of all conversion paths

### Performance Impact
- **Risk**: Editor adds overhead to FGEngine
- **Mitigation**: Editor is separate module, no impact on FGEngine runtime

### Maintenance Burden
- **Risk**: Two systems to maintain
- **Mitigation**: Editor depends on FGEngine, single source of truth for types

## Timeline Estimate

- **Phase 1**: 1-2 days (dependency resolution)
- **Phase 2**: 2-3 days (adapter creation)
- **Phase 3**: 5-7 days (core functionality porting)
- **Phase 4**: 2-3 days (testing and validation)

**Total**: ~2 weeks for complete integration

## Next Steps

1. **Start with Phase 1**: Update FGEngine dependencies to match editor
2. **Create basic adapter structure**: Set up editor subdirectory
3. **Test compatibility**: Ensure editor types can map to FGEngine types
4. **Iterative porting**: Port one editor feature at a time
5. **Continuous testing**: Validate each phase before moving to next
