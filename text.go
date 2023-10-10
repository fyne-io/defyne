package main

import (
	"io/ioutil"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// Declare conformity with editor interface
var _ editor = (*textEditor)(nil)

type codeEntry struct {
	widget.Entry

	editor *textEditor
}

func (e *codeEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if sh, ok := shortcut.(*desktop.CustomShortcut); ok {
		ctrlSuper := (runtime.GOOS == "darwin" && sh.Modifier == fyne.KeyModifierSuper) ||
			(runtime.GOOS != "darwin" && sh.Modifier == fyne.KeyModifierControl)
		if sh.KeyName == "S" && ctrlSuper {
			e.editor.save()
		}
	}
}

type textEditor struct {
	uri    fyne.URI
	entry  *codeEntry
	edited bool
}

func newTextEditor(u fyne.URI, _ fyne.Window) editor {
	text := &codeEntry{}
	text.MultiLine = true
	text.TextStyle.Monospace = true
	text.ExtendBaseWidget(text)
	editor := &textEditor{uri: u, entry: text}
	text.editor = editor

	f, _ := storage.Reader(u)
	b, _ := ioutil.ReadAll(f)
	text.SetText(string(b))
	_ = f.Close()

	text.OnChanged = func(_ string) {
		editor.edited = true
	}
	return editor
}

func (t *textEditor) changed() bool {
	return t.edited
}

func (t *textEditor) content() fyne.CanvasObject {
	return t.entry
}

func (t *textEditor) close() {
}

func (t *textEditor) run() {
}

func (t *textEditor) save() {
	w, _ := storage.Writer(t.uri)
	_, _ = w.Write([]byte(t.entry.Text))
	_ = w.Close()

	t.edited = false
}
