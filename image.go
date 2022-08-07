package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Declare conformity with editor interface
var _ editor = (*imageEditor)(nil)

type imageEditor struct {
	uri   fyne.URI
	image *canvas.Image
}

func newImageEditor(u fyne.URI, _ fyne.Window) editor {
	img := canvas.NewImageFromURI(u)
	img.FillMode = canvas.ImageFillContain
	editor := &imageEditor{uri: u, image: img}

	return editor
}

func (i *imageEditor) changed() bool {
	return false
}

func (i *imageEditor) content() fyne.CanvasObject {
	return i.image
}

func (i *imageEditor) close() {
}

func (i *imageEditor) run() {
}

func (i *imageEditor) save() {
}
