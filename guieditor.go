package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"github.com/fyne-io/defyne/internal/guibuilder"
)

type guiEditor struct {
	uri     fyne.URI
	builder *guibuilder.Builder
	edited  bool
	win     fyne.Window
}

func newGuiEditor(u fyne.URI, win fyne.Window) editor {
	builder := guibuilder.NewBuilder(u, win)
	return &guiEditor{uri: u, builder: builder, win: win}
}

func (g *guiEditor) changed() bool {
	return g.edited
}

func (g *guiEditor) content() fyne.CanvasObject {
	return g.builder.MakeUI()
}

func (g *guiEditor) close() {
}

func (g *guiEditor) save() {
	err := g.builder.Save()
	if err != nil {
		dialog.ShowError(err, g.win)
		return
	}

	g.edited = false
}
