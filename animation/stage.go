package animation

import (
	"fgengine/constants"
	"fgengine/types"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TODO, remove the global var and place stage logic in graphics to use the image cache and the camera instead of just the screen position

var StageImage *ebiten.Image

func DrawStaticColorStage(color color.RGBA, screen *ebiten.Image, screenPos types.Vector2, scaling float64) {
	if StageImage == nil {
		StageImage = ebiten.NewImage(int(constants.World.W), int(constants.World.H))
	}
	StageImage.Fill(color)

	options := &ebiten.DrawImageOptions{}

	// Apply scaling around center of viewport instead of (0,0)
	if scaling != 0 && scaling != 1 {
		// Calculate center of viewport (in screen space)
		centerX := constants.CameraWidth / 2
		centerY := constants.CameraHeight / 2

		// Calculate relative position from center
		relativeX := screenPos.X - centerX
		relativeY := screenPos.Y - centerY

		// Scale the relative position
		scaledRelativeX := relativeX * scaling
		scaledRelativeY := relativeY * scaling

		// Calculate final position (back from center)
		finalX := scaledRelativeX + centerX
		finalY := scaledRelativeY + centerY

		options.GeoM.Scale(scaling, scaling)
		options.GeoM.Translate(finalX, finalY)
	} else {
		options.GeoM.Translate(screenPos.X, screenPos.Y)
	}

	screen.DrawImage(StageImage, options)
}

func DrawGridStage(gridPixels int, lineColor, bgColor color.RGBA, screen *ebiten.Image, screenPos types.Vector2, scaling float64) {
	if StageImage == nil {
		StageImage = ebiten.NewImage(int(constants.World.W), int(constants.World.H))
		StageImage.Fill(bgColor)
	}
	for x := 0; x < int(constants.World.W); x += gridPixels {
		for y := 0; y < int(constants.World.H); y += gridPixels {
			vector.StrokeLine(StageImage, float32(x), 0, float32(x), float32(constants.World.H), 1, lineColor, false)
			vector.StrokeLine(StageImage, 0, float32(y), float32(constants.World.W), float32(y), 1, lineColor, false)
		}
	}

	options := &ebiten.DrawImageOptions{}

	// Apply scaling around center of viewport instead of (0,0)
	if scaling != 0 && scaling != 1 {
		// Calculate center of viewport (in screen space)
		centerX := constants.CameraWidth / 2
		centerY := constants.CameraHeight / 2

		// Calculate relative position from center
		relativeX := screenPos.X - centerX
		relativeY := screenPos.Y - centerY

		// Scale the relative position
		scaledRelativeX := relativeX * scaling
		scaledRelativeY := relativeY * scaling

		// Calculate final position (back from center)
		finalX := scaledRelativeX + centerX
		finalY := scaledRelativeY + centerY

		options.GeoM.Scale(scaling, scaling)
		options.GeoM.Translate(finalX, finalY)
	} else {
		options.GeoM.Translate(screenPos.X, screenPos.Y)
	}

	screen.DrawImage(StageImage, options)
}

func DrawStaticImageStage(img *ebiten.Image, screen *ebiten.Image, screenPos types.Vector2, scaling float64) {
	if img == nil {
		return
	}
	options := &ebiten.DrawImageOptions{}

	// Apply scaling around center of viewport instead of (0,0)
	if scaling != 0 && scaling != 1 {
		// Calculate center of viewport (in screen space)
		centerX := constants.CameraWidth / 2
		centerY := constants.CameraHeight / 2

		// Calculate relative position from center
		relativeX := screenPos.X - centerX
		relativeY := screenPos.Y - centerY

		// Scale the relative position
		scaledRelativeX := relativeX * scaling
		scaledRelativeY := relativeY * scaling

		// Calculate final position (back from center)
		finalX := scaledRelativeX + centerX
		finalY := scaledRelativeY + centerY

		options.GeoM.Scale(scaling, scaling)
		options.GeoM.Translate(finalX, finalY)
	} else {
		options.GeoM.Translate(screenPos.X, screenPos.Y)
	}

	screen.DrawImage(img, options)
}
