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

func (c *Camera) IsVisible(renderable Renderable) bool {
	pos := renderable.GetPosition()
	sprite := renderable.GetSprite()

	// If no sprite, fall back to point check
	if sprite == nil {
		return c.Viewport.Contains(pos.X, pos.Y)
	}
	renderableRect := types.Rect{
		X: pos.X,
		Y: pos.Y,
		W: sprite.Rect.W,
		H: sprite.Rect.H,
	}
	return c.Viewport.IsOverlapping(renderableRect)
}

func layoutMatchesCamSize(camera *Camera) bool {
	return (float64(config.LayoutSizeW) == camera.Viewport.W && float64(config.LayoutSizeH) == camera.Viewport.H)
}

func zoomAroundCenterOption(options *ebiten.DrawImageOptions, camera *Camera, renderable Renderable, screenPos types.Vector2) {
	centerViewportX := camera.Viewport.W / 2
	centerViewportY := camera.Viewport.H / 2

	relativeX := screenPos.X - centerViewportX
	relativeY := screenPos.Y - centerViewportY

	scaledRelativeX := relativeX * camera.Scaling
	scaledRelativeY := relativeY * camera.Scaling

	finalX := scaledRelativeX + centerViewportX
	finalY := scaledRelativeY + centerViewportY

	options.GeoM.Scale(renderable.GetRenderProperties().Scale.X*camera.Scaling, renderable.GetRenderProperties().Scale.Y*camera.Scaling)
	options.GeoM.Translate(finalX, finalY)
}
