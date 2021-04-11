package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	xWidget "fyne.io/x/fyne/widget"
)

type fileTab struct {
	editor
	uri fyne.URI
}

type defyne struct {
	win         fyne.Window
	projectRoot fyne.URI
	fileTabs    *container.DocTabs
	fileTree    *xWidget.FileTree
	openEditors map[*container.TabItem]*fileTab
}

func (d *defyne) openEditor(u fyne.URI) {
	if _, ok := editors[u.MimeType()]; !ok {
		dialog.ShowInformation("No registered editor",
			"No known editor for mime "+u.MimeType(), d.win)
		return
	}

	for tab, item := range d.openEditors {
		if item.uri.String() == u.String() {
			d.fileTabs.Select(tab)
			return
		}
	}

	ed := editors[u.MimeType()](u)
	newTab := container.NewTabItemWithIcon(u.Name(), theme.FileTextIcon(), ed.content())
	d.openEditors[newTab] = &fileTab{ed, u}

	d.fileTabs.Append(newTab)
	d.fileTabs.Select(newTab)
}
