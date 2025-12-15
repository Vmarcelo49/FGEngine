package graphics

import (
	"fgengine/config"
	"fgengine/constants"
	"fgengine/types"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	Viewport        types.Rect
	LockWorldBounds bool
	Scaling         float64
}

// Makes a camera centered in the world
func NewCamera() *Camera {
	viewport := constants.Camera // Start with camera dimensions
	viewport.AlignCenter(constants.World)
	// Position camera at bottom of world instead of center vertically
	viewport.Y = constants.World.Bottom() - constants.Camera.H

	return &Camera{
		Viewport:        viewport,
		LockWorldBounds: false,
		Scaling:         1,
	}
}

func (c *Camera) UpdatePosition(targetPos types.Vector2) {
	// Center viewport around target position
	c.Viewport.X = targetPos.X - c.Viewport.W/2
	c.Viewport.Y = targetPos.Y - c.Viewport.H/2
	if c.LockWorldBounds {
		c.lockToWorldBounds()
	}
}

func (c *Camera) SetPosition(pos types.Vector2) {
	c.Viewport.X = pos.X
	c.Viewport.Y = pos.Y
	if c.LockWorldBounds {
		c.lockToWorldBounds()
	}
}

func (c *Camera) lockToWorldBounds() {
	if c.Viewport.X < constants.World.X {
		c.Viewport.X = constants.World.X
	}
	if c.Viewport.X > constants.World.Right()-c.Viewport.W {
		c.Viewport.X = constants.World.Right() - c.Viewport.W
	}
	if c.Viewport.Y < constants.World.Y {
		c.Viewport.Y = constants.World.Y
	}
	if c.Viewport.Y > constants.World.Bottom()-c.Viewport.H {
		c.Viewport.Y = constants.World.Bottom() - c.Viewport.H
	}
}

func (c *Camera) WorldToScreen(worldPos types.Vector2) types.Vector2 {
	screenX := worldPos.X - c.Viewport.X
	screenY := worldPos.Y - c.Viewport.Y

	return types.Vector2{
		X: screenX,
		Y: screenY,
	}
}

func (c *Camera) ScreenToWorld(screenPos types.Vector2) types.Vector2 {
	worldX := screenPos.X
	worldY := screenPos.Y

	return types.Vector2{
		X: worldX + c.Viewport.X,
		Y: worldY + c.Viewport.Y,
	}
}

func layoutMatchesCamSize(camera *Camera) bool {
	return (float64(config.LayoutSizeW) == camera.Viewport.W && float64(config.LayoutSizeH) == camera.Viewport.H)
}

func CameraTransform(options *ebiten.DrawImageOptions, camera *Camera, entityScale types.Vector2, screenPos types.Vector2) {
	options.GeoM.Scale(entityScale.X, entityScale.Y)

	centerX := camera.Viewport.W / 2
	centerY := camera.Viewport.H / 2

	if camera.Scaling != 0 && camera.Scaling != 1 {
		options.GeoM.Translate(-centerX, -centerY)
		options.GeoM.Scale(camera.Scaling, camera.Scaling)
		options.GeoM.Translate(centerX, centerY)

		// scale again?
		if !layoutMatchesCamSize(camera) {
			options.GeoM.Scale(entityScale.X, entityScale.Y)
		}
	}

	options.GeoM.Translate(screenPos.X, screenPos.Y)
}
