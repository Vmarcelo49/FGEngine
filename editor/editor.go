package editor

import (
	"fgengine/character"
	"fgengine/types"
	"image/color"

	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	defaultEditorWidth  = 1920
	defaultEditorHeight = 1080
)

// CharacterEditor is an ebiten.Game implementation for editing characters and animations.
type CharacterEditor struct {
	width  int
	height int

	char *character.Character

	activeAnimationName string
	renameCharacterTo   string
	renameAnimationTo   string
	newAnimationName    string

	selectedFrame int
	cancelTypes   string

	selectedBoxType  types.BoxType
	selectedBoxIndex int
	targetBoxType    types.BoxType

	paused       bool
	previewScale int

	showCreateWindow              bool
	showLoadWindow                bool
	showSaveWindow                bool
	showImportWindow              bool
	showExitWindow                bool
	showChangeCharacterNameWindow bool
	showRenameAnimationWindow     bool
	showDeleteAnimationWindow     bool
	ignoreWindowClose             bool

	newCharacterName   string
	loadPath           string
	savePath           string
	pendingImportPaths []string
	exitAfterSave      bool
	dirty              bool

	imagePreviewCache map[string]*imagePreview
	nextTextureID     int

	statusLine string
	exitEditor bool
}

func (ed *CharacterEditor) Update() error {
	if ebiten.IsWindowBeingClosed() {
		if !ed.ignoreWindowClose {
			ed.requestExit()
		}
	}

	if ed.exitEditor {
		return ebiten.Termination
	}

	ebimgui.Update(1.0 / 60.0)

	if ed.char != nil && !ed.paused {
		player := ed.player()
		if player != nil && player.ActiveAnimation != nil {
			player.Update(ed.activeAnimationName) // may cause issues, TODO verify
			ed.selectedFrame = player.FrameIndex
		}
	}

	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()

	ed.drawTopMenuBar()
	ed.drawCharacterWindow()
	ed.drawAnimationPlayerWindow()
	ed.drawFrameDataWindow()
	ed.drawBoxEditorWindow()
	ed.drawImagesWindow()

	ed.drawCreateCharacterWindow()
	ed.drawLoadCharacterWindow()
	ed.drawSaveCharacterWindow()
	ed.drawImportImagesAsAnimationWindow()
	ed.drawUnsavedChangesWindow()
	ed.drawChangeCharacterNameWindow()
	ed.drawRenameActiveAnimationWindow()
	ed.drawDeleteAnimationWindow()

	return nil
}

func (ed *CharacterEditor) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 24, G: 24, B: 26, A: 255})
	ed.drawCharacterPreview(screen)
	ebimgui.Draw(screen)
}

func (ed *CharacterEditor) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth > 0 {
		ed.width = outsideWidth
	}
	if outsideHeight > 0 {
		ed.height = outsideHeight
	}
	ebimgui.SetDisplaySize(float32(ed.width), float32(ed.height))
	return ed.width, ed.height
}

// Entrypoint
func NewCharacterEditor() *CharacterEditor {
	ed := &CharacterEditor{
		width:             defaultEditorWidth,
		height:            defaultEditorHeight,
		newCharacterName:  "NewCharacter",
		renameCharacterTo: "NewCharacter",
		newAnimationName:  "new_animation",
		selectedBoxType:   types.Collision,
		targetBoxType:     types.Hit,
		paused:            true,
		previewScale:      3,
		imagePreviewCache: make(map[string]*imagePreview),
	}
	ebiten.SetWindowClosingHandled(true)
	ed.createNewCharacter("NewCharacter")
	ed.clearDirty()
	return ed
}
