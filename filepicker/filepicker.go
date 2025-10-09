package filepicker

// FileFilter represents a file filter for the file picker
type FileFilter struct {
	Description string   // e.g. "Image files"
	Extensions  []string // e.g. ["png", "jpg", "jpeg"]
}

// FilePicker defines the interface for file selection
type FilePicker interface {
	// LoadFile opens a dialog to load a file
	LoadFile(filter FileFilter) (string, error)

	// SaveFile opens a dialog to save a file
	SaveFile(filter FileFilter) (string, error)
}

// GetFilePicker returns the appropriate FilePicker implementation for the platform
func GetFilePicker() FilePicker {
	return newPlatformFilePicker()
}
