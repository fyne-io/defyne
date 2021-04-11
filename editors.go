package main

import "fyne.io/fyne/v2"

var editors = map[string]func(fyne.URI) editor{
	"application/xml-tdd": newTextEditor,
	"image/jpeg":          newImageEditor,
	"image/png":           newImageEditor,
	"image/svg+xml":       newImageEditor,
	"text/markdown":       newTextEditor,
	"text/plain":          newTextEditor,
}

type editor interface {
	changed() bool
	close()
	content() fyne.CanvasObject
	save()
}
