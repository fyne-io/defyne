//go:generate fyne bundle -o bundled.go Icon.png

package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
)

func (d *defyne) setProject(u fyne.URI) {
	d.projectRoot = u

	content := container.NewVSplit(d.makeEditorPanel(), d.makeTerminalPanel())
	content.Offset = 0.8
	d.fileTree = d.makeFilesPanel()
	mainSplit := container.NewHSplit(d.fileTree, content)
	mainSplit.Offset = 0.2

	d.win.SetMainMenu(d.makeMenu())
	d.win.SetContent(container.NewBorder(d.makeToolbar(), nil, nil, nil, mainSplit))
}

func main() {
	a := app.NewWithID("io.fyne.defyne")
	a.SetIcon(resourceIconPng)
	w := a.NewWindow("Defyne")

	pwd, _ := os.Getwd()
	root := storage.NewFileURI(pwd)
	ide := &defyne{win: w}
	ide.setProject(root)
	w.Resize(fyne.NewSize(1024, 768))
	w.ShowAndRun()
}
