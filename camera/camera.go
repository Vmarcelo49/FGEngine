package camera

var (
	Width  int = 640
	Height int = 360
)

// SetDimensions updates the camera viewport size
func SetDimensions(width, height int) {
	Width = width
	Height = height
}

// GetDimensions returns the current camera viewport size
func GetDimensions() (int, int) {
	return Width, Height
}

// TODO: Future camera system will include:
// - Position/transform management
// - Zoom/scale controls
// - Smooth following/tracking
// - Screen-to-world coordinate conversion
