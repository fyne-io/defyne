package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Declare conformity with editor interface
var _ editor = (*imageEditor)(nil)

type imageEditor struct {
	uri   fyne.URI
	image *canvas.Image
}

func newImageEditor(u fyne.URI, _ fyne.Window) editor {
	img := canvas.NewImageFromURI(u)
	if u.Extension() == ".svg" {
		img.FillMode = canvas.ImageFillContain
	} else {
		img.FillMode = canvas.ImageFillOriginal
	}
	editor := &imageEditor{uri: u, image: img}

	return editor
}

func (i *imageEditor) changed() bool {
	return false
}

func (i *imageEditor) content() fyne.CanvasObject {
	content := fyne.CanvasObject(i.image)
	if i.uri.Extension() != ".svg" {
		content = container.NewCenter(content)
	}

	return container.NewScroll(content)
}

func (i *imageEditor) close() {
}

func (i *imageEditor) run() {
}

func (i *imageEditor) save() {
}
