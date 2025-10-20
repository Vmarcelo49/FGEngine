package filepicker

type FileFilter struct {
	Description string   // e.g. "Image files"
	Extensions  []string // e.g. ["png", "jpg", "jpeg"]
}

type FilePicker interface {
	LoadFile(filter FileFilter) (string, error)
	SaveFile(filter FileFilter) (string, error)
}

// GetFilePicker returns the appropriate FilePicker implementation for the platform
func GetFilePicker() FilePicker {
	return newPlatformFilePicker()
}

// Filepicker needs to be refactored to return files instead of paths for better cross-platform support, browser environments do not have direct file system access.
