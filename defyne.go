package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	xWidget "fyne.io/x/fyne/widget"
)

type defyne struct {
	win         fyne.Window
	projectRoot fyne.URI
	fileTabs    *container.DocTabs
	fileTree    *xWidget.FileTree
	openEditors map[*container.TabItem]editor
}

func (d *defyne) openEditor(u fyne.URI) {
	if u.Extension() != ".go" && u.Extension() != ".mod" && !strings.Contains(u.MimeType(), "text/") {
		return // TODO let's add pluggable editor providers
	}

	editor := newTextEditor(u)
	newTab := container.NewTabItemWithIcon(u.Name(), theme.FileTextIcon(), editor.content())
	d.openEditors[newTab] = editor

	d.fileTabs.Append(newTab)
	d.fileTabs.Select(newTab)
}

type editor interface {
	changed() bool
	close()
	content() fyne.CanvasObject
	save()
}
