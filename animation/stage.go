package animation

import (
	"fgengine/constants"
	"fgengine/graphics"
	"fgengine/types"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// StageType defines the visual style of the stage
type StageType int

const (
	StageTypeSolidColor StageType = iota
	StageTypeGrid
	StageTypeImage
)

// Stage implements graphics.Drawable and represents the game background
type Stage struct {
	stageType StageType
	image     *ebiten.Image // cached image for rendering
	dirty     bool          // whether the image needs to be regenerated

	// Solid color stage
	bgColor color.RGBA

	// Grid stage
	gridSize  int
	lineColor color.RGBA

	// Image stage
	sourceImage *ebiten.Image

	// World position (usually 0,0 for background)
	Position types.Vector2
	Size     types.Vector2
}

// NewSolidColorStage creates a stage with a solid background color
func NewSolidColorStage(bgColor color.RGBA) *Stage {
	return &Stage{
		stageType: StageTypeSolidColor,
		bgColor:   bgColor,
		dirty:     true,
		Position:  types.Vector2{X: 0, Y: 0},
		Size:      types.Vector2{X: constants.World.W, Y: constants.World.H},
	}
}

// NewGridStage creates a stage with a grid pattern
func NewGridStage(gridSize int, lineColor, bgColor color.RGBA) *Stage {
	return &Stage{
		stageType: StageTypeGrid,
		gridSize:  gridSize,
		lineColor: lineColor,
		bgColor:   bgColor,
		dirty:     true,
		Position:  types.Vector2{X: 0, Y: 0},
		Size:      types.Vector2{X: constants.World.W, Y: constants.World.H},
	}
}

// NewImageStage creates a stage with a custom image background
func NewImageStage(img *ebiten.Image) *Stage {
	var size types.Vector2
	if img != nil {
		bounds := img.Bounds()
		size = types.Vector2{X: float64(bounds.Dx()), Y: float64(bounds.Dy())}
	} else {
		size = types.Vector2{X: constants.World.W, Y: constants.World.H}
	}

	return &Stage{
		stageType:   StageTypeImage,
		sourceImage: img,
		dirty:       false, // no need to generate, we use source directly
		Position:    types.Vector2{X: 0, Y: 0},
		Size:        size,
	}
}

// Draw implements graphics.Drawable interface
func (s *Stage) Draw(screen *ebiten.Image, camera *graphics.Camera) {
	s.ensureImage()

	if s.image == nil && s.sourceImage == nil {
		return
	}

	imgToDraw := s.image
	if s.stageType == StageTypeImage {
		imgToDraw = s.sourceImage
	}

	if imgToDraw == nil {
		return
	}

	screenPos := camera.WorldToScreen(s.Position)
	options := &ebiten.DrawImageOptions{}

	graphics.CameraTransform(options, camera, types.Vector2{X: 1, Y: 1}, screenPos)

	screen.DrawImage(imgToDraw, options)
}

// ensureImage creates or regenerates the cached image if needed
func (s *Stage) ensureImage() {
	if !s.dirty {
		return
	}

	switch s.stageType {
	case StageTypeSolidColor:
		s.generateSolidColorImage()
	case StageTypeGrid:
		s.generateGridImage()
	case StageTypeImage:
		// nothing to generate, uses sourceImage directly
	}

	s.dirty = false
}

func (s *Stage) generateSolidColorImage() {
	if s.image == nil {
		s.image = ebiten.NewImage(int(s.Size.X), int(s.Size.Y))
	}
	s.image.Fill(s.bgColor)
}

func (s *Stage) generateGridImage() {
	if s.image == nil {
		s.image = ebiten.NewImage(int(s.Size.X), int(s.Size.Y))
	}

	s.image.Fill(s.bgColor)

	width := int(s.Size.X)
	height := int(s.Size.Y)

	// Draw vertical lines
	for x := 0; x <= width; x += s.gridSize {
		vector.StrokeLine(s.image, float32(x), 0, float32(x), float32(height), 1, s.lineColor, false)
	}

	// Draw horizontal lines
	for y := 0; y <= height; y += s.gridSize {
		vector.StrokeLine(s.image, 0, float32(y), float32(width), float32(y), 1, s.lineColor, false)
	}
}

// SetBackgroundColor changes the background color and marks for regeneration
func (s *Stage) SetBackgroundColor(c color.RGBA) {
	s.bgColor = c
	s.dirty = true
}

// SetGridSize changes the grid size and marks for regeneration
func (s *Stage) SetGridSize(size int) {
	if s.stageType != StageTypeGrid {
		return
	}
	s.gridSize = size
	s.dirty = true
}

// SetLineColor changes the grid line color and marks for regeneration
func (s *Stage) SetLineColor(c color.RGBA) {
	if s.stageType != StageTypeGrid {
		return
	}
	s.lineColor = c
	s.dirty = true
}

// SetSourceImage changes the source image for image-type stages
func (s *Stage) SetSourceImage(img *ebiten.Image) {
	s.sourceImage = img
	if img != nil {
		bounds := img.Bounds()
		s.Size = types.Vector2{X: float64(bounds.Dx()), Y: float64(bounds.Dy())}
	}
}

// Invalidate forces the stage to regenerate its image on next draw
func (s *Stage) Invalidate() {
	s.dirty = true
}

// GetType returns the stage type
func (s *Stage) GetType() StageType {
	return s.stageType
}
