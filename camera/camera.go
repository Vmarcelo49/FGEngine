package camera

import (
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/types"
)

var (
	Camera types.Rect = types.Rect{X: 0, Y: 0, W: 640, H: 360}
)

func InsideCameraView(renderable graphics.Renderable) bool {
	pos := renderable.GetPosition()
	sprite := renderable.GetSprite()

	// If no sprite, fall back to point check
	if sprite == nil {
		return Camera.Contains(pos.X, pos.Y)
	}
	// Check if the renderable rect overlaps with the camera viewport
	return Camera.IsOverlapping(sprite.Rect)
}

func UpdateCamera(targetPos types.Vector2) {
	Camera.X = targetPos.X - Camera.W/2
	Camera.Y = targetPos.Y - Camera.H/2
	restrictCameraToWorldBounds()
}

func restrictCameraToWorldBounds() {
	if Camera.X < 0 {
		Camera.X = 0
	}
	if Camera.X > constants.WorldWidth-Camera.W {
		Camera.X = constants.WorldWidth - Camera.W
	}
	if Camera.Y < 0 {
		Camera.Y = 0
	}
	if Camera.Y > constants.WorldHeight-Camera.H {
		Camera.Y = constants.WorldHeight - Camera.H
	}
}
