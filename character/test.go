package character

import (
	"fgengine/animation"
	"fgengine/state"
	"fgengine/types"
)

func MakeTestCharacter() *Character {
	char := &Character{
		ID:         0,
		Name:       "Test Character",
		Friction:   2,
		JumpHeight: 10,
		Animations: map[string]*animation.Animation{
			"idle": {
				Name: "idle",
				Sprites: []*animation.Sprite{
					{
						ImagePath: "assets/common/idle.png",
						Rect: types.Rect{
							X: 0, Y: 0, W: 100, H: 200,
						},
					},
				},
				FrameData: []animation.FrameData{
					{SpriteIndex: 0, Duration: 10},
				},
				TotalDuration: 10,
			},
			"walk": {
				Name: "walk",
				Sprites: []*animation.Sprite{
					{
						ImagePath: "assets/common/walk.png",
						Rect: types.Rect{
							X: 0, Y: 0, W: 100, H: 200,
						},
					},
					{
						ImagePath: "assets/common/walk.png",
						Rect: types.Rect{
							X: 0, Y: 0, W: 100, H: 200,
						},
					},
				},
				FrameData: []animation.FrameData{
					{SpriteIndex: 0, Duration: 10},
					{SpriteIndex: 1, Duration: 10},
				},
				TotalDuration: 20,
			},
		},
		StateMachine:    &state.StateMachine{},
		AnimationPlayer: &animation.AnimationPlayer{},
	}
	char.SetAnimation("idle")
	return char
}
