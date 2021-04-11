package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	if _, ok := editors[u.MimeType()]; !ok {
		dialog.ShowInformation("No registered editor",
			"No known editor for mime "+u.MimeType(), d.win)
		return
	}

	ed := editors[u.MimeType()](u)
	newTab := container.NewTabItemWithIcon(u.Name(), theme.FileTextIcon(), ed.content())
	d.openEditors[newTab] = ed

	d.fileTabs.Append(newTab)
	d.fileTabs.Select(newTab)
}
