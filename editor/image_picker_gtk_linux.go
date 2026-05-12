//go:build linux && cgo

package editor

/*
#cgo pkg-config: gtk+-3.0
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <gtk/gtk.h>
#include <stdlib.h>

static GtkWidget* fg_open_image_dialog(const char* title) {
	return gtk_file_chooser_dialog_new(
		title,
		NULL,
		GTK_FILE_CHOOSER_ACTION_OPEN,
		"Cancel", GTK_RESPONSE_CANCEL,
		"Open", GTK_RESPONSE_ACCEPT,
		NULL
	);
}

static void fg_add_image_filter(GtkFileChooser* chooser) {
	GtkFileFilter* filter = gtk_file_filter_new();
	gtk_file_filter_set_name(filter, "Image files");
	gtk_file_filter_add_pattern(filter, "*.png");
	gtk_file_chooser_add_filter(chooser, filter);
}

*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

var (
	gtkInitOnce sync.Once
	gtkInitOK   bool
)

func ensureGTK() error {
	gtkInitOnce.Do(func() {
		C.XInitThreads()
		gtkInitOK = C.gtk_init_check(nil, nil) == C.TRUE
	})
	if !gtkInitOK {
		return errors.New("gtk initialization failed (is an X/Wayland session available?)")
	}
	return nil
}

func closeGTKDialog(dlg *C.GtkWidget) {
	C.gtk_widget_destroy(dlg)
	for C.gtk_events_pending() != 0 {
		C.gtk_main_iteration()
	}
}

func (ed *CharacterEditor) pickImagesWithDialog() ([]string, error) {
	if err := ensureGTK(); err != nil {
		return nil, err
	}

	ctitle := C.CString("Load Images")
	defer C.free(unsafe.Pointer(ctitle))

	dlg := C.fg_open_image_dialog(ctitle)
	chooser := (*C.GtkFileChooser)(unsafe.Pointer(dlg))
	C.fg_add_image_filter(chooser)
	C.gtk_file_chooser_set_select_multiple(chooser, C.TRUE)

	response := C.gtk_dialog_run((*C.GtkDialog)(unsafe.Pointer(dlg)))
	defer closeGTKDialog(dlg)
	if response != C.GTK_RESPONSE_ACCEPT {
		return nil, nil
	}

	files := C.gtk_file_chooser_get_filenames(chooser)
	if files == nil {
		return nil, nil
	}
	defer C.g_slist_free(files)

	paths := make([]string, 0, 8)
	for node := files; node != nil; node = node.next {
		cpath := (*C.char)(node.data)
		if cpath == nil {
			continue
		}
		paths = append(paths, C.GoString(cpath))
		C.g_free(C.gpointer(node.data))
	}

	return dedupeNonEmpty(paths), nil
}

func (ed *CharacterEditor) pickCharacterWithDialog() (string, error) {
	if err := ensureGTK(); err != nil {
		return "", err
	}

	ctitle := C.CString("Load Character")
	defer C.free(unsafe.Pointer(ctitle))

	dlg := C.fg_open_image_dialog(ctitle)
	chooser := (*C.GtkFileChooser)(unsafe.Pointer(dlg))
	C.gtk_file_chooser_set_select_multiple(chooser, C.TRUE)

	response := C.gtk_dialog_run((*C.GtkDialog)(unsafe.Pointer(dlg)))
	defer closeGTKDialog(dlg)
	if response != C.GTK_RESPONSE_ACCEPT {
		return "", nil
	}

	files := C.gtk_file_chooser_get_filenames(chooser)
	if files == nil {
		return "", nil
	}
	defer C.g_slist_free(files)

	node := files
	if node == nil || node.data == nil {
		return "", nil
	}

	cpath := (*C.char)(node.data)
	defer C.g_free(C.gpointer(node.data))

	return C.GoString(cpath), nil
}
