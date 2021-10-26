package main

import (
	"strings"

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
	for tab, item := range d.openEditors {
		if item.uri.String() == u.String() {
			d.fileTabs.Select(tab)
			return
		}
	}

	var ed editor
	for name, e := range editorsByFilename {
		if strings.Index(u.Name(), name) > -1 {
			ed = e(u, d.win)
		}
	}
	if ed == nil {
		if _, ok := editorsByMime[u.MimeType()]; !ok {
			dialog.ShowInformation("No registered editor",
				"No known editor for mime "+u.MimeType(), d.win)
			return
		}

		ed = editorsByMime[u.MimeType()](u, d.win)
	}

	newTab := container.NewTabItemWithIcon(u.Name(), theme.FileTextIcon(), ed.content())
	d.openEditors[newTab] = &fileTab{ed, u}

	d.fileTabs.Append(newTab)
	d.fileTabs.Select(newTab)
}
