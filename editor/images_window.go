package editor

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"

	"fgengine/animation"

	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

type imagePreview struct {
	textureRef *int
	width      float32
	height     float32
}

func (ed *CharacterEditor) drawImagesWindow() {
	open := true
	if !imgui.BeginV("Images", &open, imgui.WindowFlags(0)) {
		imgui.End()
		return
	}
	defer imgui.End()

	anim := ed.activeAnimation()
	if anim == nil {
		imgui.Text("No active animation.")
		return
	}

	if imgui.Button("Load Images") {
		loaded, err := ed.pickImagesWithDialog()
		if err != nil {
			ed.statusLine = "Image load failed: " + err.Error()
		} else if len(loaded) > 0 {
			anim := ed.activeAnimation()
			if anim != nil && len(anim.Sprites) == 0 && len(loaded) > 1 {
				ed.pendingImportPaths = loaded
				ed.showImportWindow = true
			} else {
				ed.addImagesToActiveAnimation(loaded)
				ed.statusLine = fmt.Sprintf("Added %d image(s)", len(loaded))
			}
		}
	}

	imgui.Separator()

	if len(anim.Sprites) == 0 {
		imgui.Text("No images in active animation.")
		return
	}

	for i, spr := range anim.Sprites {
		imgui.PushIDInt(int32(i))

		name := fmt.Sprintf("IMG%d", i+1)
		imgPath := ""
		if spr != nil && spr.ImagePath != "" {
			imgPath = spr.ImagePath
			name = filepath.Base(spr.ImagePath)
		}

		preview, err := ed.ensureImagePreview(imgPath)
		if err != nil {
			imgui.Text("Preview load failed")
			imgui.Text(err.Error())
		} else if preview != nil {
			pw, ph := fitPreviewSize(preview.width, preview.height, 196, 128)
			imgui.Image(imgui.TextureID(preview.textureRef), imgui.Vec2{X: pw, Y: ph})
		} else {
			imgui.Button("No Image")
		}

		if imgui.BeginPopupContextItemV("image_context", 1) {
			if imgui.Button("Delete") {
				ed.deleteImageFromActiveAnimation(i)
			}
			if imgui.Button("Set Image On The Current Frame") {
				ed.setImageOnCurrentFrame(i)
			}
			imgui.EndPopup()
		}

		imgui.Text(name)
		imgui.PopID()
	}
}

func (ed *CharacterEditor) ensureImagePreview(path string) (*imagePreview, error) {
	trimPath := strings.TrimSpace(path)
	if trimPath == "" {
		return nil, nil
	}

	if preview, ok := ed.imagePreviewCache[trimPath]; ok {
		return preview, nil
	}

	texture, rawImage, err := ebitenutil.NewImageFromFile(trimPath)
	if err != nil {
		return nil, err
	}

	bounds := rawImage.Bounds()
	ed.nextTextureID++
	textureRef := new(int)
	*textureRef = ed.nextTextureID + 10_000
	tid := imgui.TextureID(textureRef)
	ebimgui.GlobalManager().Cache.SetTexture(tid, texture)

	preview := &imagePreview{
		textureRef: textureRef,
		width:      float32(bounds.Dx()),
		height:     float32(bounds.Dy()),
	}
	ed.imagePreviewCache[trimPath] = preview

	return preview, nil
}

func fitPreviewSize(width, height, maxWidth, maxHeight float32) (float32, float32) {
	if width <= 0 || height <= 0 {
		return maxWidth, maxHeight
	}
	scale := maxWidth / width
	hScale := maxHeight / height
	if hScale < scale {
		scale = hScale
	}
	if scale > 1 {
		scale = 1
	}
	return width * scale, height * scale
}

func (ed *CharacterEditor) addImagesToActiveAnimation(paths []string) {
	anim := ed.activeAnimation()
	if anim == nil {
		return
	}

	for _, p := range paths {
		if strings.TrimSpace(p) == "" {
			continue
		}
		anim.Sprites = append(anim.Sprites, &animation.Sprite{ImagePath: p})
	}
	ed.applyDefaultIdleAnchorToAnimationSprites(anim)

	ed.markDirty()
}

func (ed *CharacterEditor) deleteImageFromActiveAnimation(deleteIndex int) {
	anim := ed.activeAnimation()
	if anim == nil || deleteIndex < 0 || deleteIndex >= len(anim.Sprites) {
		return
	}

	deletedPath := strings.TrimSpace(anim.Sprites[deleteIndex].ImagePath)
	anim.Sprites = append(anim.Sprites[:deleteIndex], anim.Sprites[deleteIndex+1:]...)

	for i := range anim.FrameData {
		frame := &anim.FrameData[i]
		spriteIndex := frame.SpriteIndex
		switch {
		case spriteIndex == deleteIndex:
			frame.SpriteIndex = 0
		case spriteIndex > deleteIndex:
			frame.SpriteIndex = spriteIndex - 1
		}
	}

	if deletedPath != "" && !isPathUsedByAnimation(anim, deletedPath) {
		ed.removeImagePreview(deletedPath)
	}

	ed.statusLine = "Image deleted"
	ed.markDirty()
}

func isPathUsedByAnimation(anim *animation.Animation, path string) bool {
	for _, spr := range anim.Sprites {
		if spr != nil && strings.TrimSpace(spr.ImagePath) == path {
			return true
		}
	}
	return false
}

func (ed *CharacterEditor) removeImagePreview(path string) {
	preview, ok := ed.imagePreviewCache[path]
	if !ok || preview == nil || preview.textureRef == nil {
		return
	}
	ebimgui.GlobalManager().Cache.RemoveTexture(imgui.TextureID(preview.textureRef))
	delete(ed.imagePreviewCache, path)
}

func (ed *CharacterEditor) setImageOnCurrentFrame(imageIndex int) {
	fd := ed.currentFrameData()
	anim := ed.activeAnimation()
	if fd == nil || anim == nil {
		return
	}
	if imageIndex < 0 || imageIndex >= len(anim.Sprites) {
		return
	}

	fd.SpriteIndex = imageIndex
	ed.statusLine = fmt.Sprintf("Set current frame image to %d", imageIndex)
	ed.markDirty()
}

func dedupeNonEmpty(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		trim := strings.TrimSpace(p)
		if trim == "" {
			continue
		}
		if _, ok := seen[trim]; ok {
			continue
		}
		seen[trim] = struct{}{}
		out = append(out, trim)
	}
	return out
}
