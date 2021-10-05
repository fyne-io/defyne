package main

import "fyne.io/fyne/v2"

var editorsByFilename = map[string]func(fyne.URI, fyne.Window) editor{
	".gui.json": newGuiEditor,
}

var editorsByMime = map[string]func(fyne.URI) editor{
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
