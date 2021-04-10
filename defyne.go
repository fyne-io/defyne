package main

import (
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type defyne struct {
	win         fyne.Window
	projectRoot fyne.URI
	fileTabs    *container.DocTabs
}

func (d *defyne) openEditor(u fyne.URI) {
	if u.Extension() != ".go" && u.Extension() != ".mod" && !strings.Contains(u.MimeType(), "text/") {
		return // TODO let's add pluggable editor providers
	}

	text := widget.NewMultiLineEntry()
	text.TextStyle.Monospace = true
	f, _ := storage.Reader(u)
	b, _ := ioutil.ReadAll(f)
	text.SetText(string(b))
	_ = f.Close()

	newTab := container.NewTabItemWithIcon(u.Name(), theme.FileTextIcon(), text)
	d.fileTabs.Append(newTab)
	d.fileTabs.Select(newTab)
}
