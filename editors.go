package main

import "fyne.io/fyne/v2"

var editorsByFilename = map[string]func(fyne.URI, fyne.Window) editor{
	".gui.json": newGuiEditor,
}

var editorsByMime = map[string]func(fyne.URI) editor{
	"application/xml-tdd": newTextEditor, // Go mime on older systems?
	"image/jpeg":          newImageEditor,
	"image/png":           newImageEditor,
	"image/svg+xml":       newImageEditor,
	"text/markdown":       newTextEditor,
	"text/plain":          newTextEditor,
	"text/x-go":           newTextEditor,
}

type editor interface {
	changed() bool
	close()
	content() fyne.CanvasObject
	run()
	save()
}
