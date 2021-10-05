package main

import (
	"fyne.io/fyne/v2"

	"github.com/fyne-io/defyne/internal/guibuilder"
)

type guiEditor struct {
	uri     fyne.URI
	builder fyne.CanvasObject
	edited  bool
}

func newGuiEditor(u fyne.URI, win fyne.Window) editor {
	builder := guibuilder.ShowBuilder(u, win)
	return &guiEditor{uri: u, builder: builder}
}

func (g *guiEditor) changed() bool {
	return g.edited
}

func (g *guiEditor) content() fyne.CanvasObject {
	return g.builder
}

func (g *guiEditor) close() {
}

func (g *guiEditor) save() {
	// TODO

	g.edited = false
}
