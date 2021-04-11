package main

import (
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type textEditor struct {
	uri    fyne.URI
	entry  *widget.Entry
	edited bool
}

func newTextEditor(u fyne.URI) editor {
	text := widget.NewMultiLineEntry()
	text.TextStyle.Monospace = true
	editor := &textEditor{uri: u, entry: text}

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

func (t *textEditor) save() {
	w, _ := storage.Writer(t.uri)
	_, _ = w.Write([]byte(t.entry.Text))
	_ = w.Close()

	t.edited = false
}
