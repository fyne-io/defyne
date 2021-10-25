//go:generate fyne bundle -o bundled.go Icon.png

package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
)

func main() {
	a := app.NewWithID("io.fyne.defyne")
	a.SetIcon(resourceIconPng)
	w := a.NewWindow("Defyne")

	pwd, _ := os.Getwd()
	root := storage.NewFileURI(pwd)
	ide := &defyne{win: w, projectRoot: root}

	content := container.NewVSplit(ide.makeEditorPanel(), ide.makeTerminalPanel())
	content.Offset = 0.8
	ide.fileTree = ide.makeFilesPanel()
	mainSplit := container.NewHSplit(ide.fileTree, content)
	mainSplit.Offset = 0.2

	w.SetMainMenu(ide.makeMenu())
	w.SetContent(container.NewBorder(ide.makeToolbar(), nil, nil, nil, mainSplit))
	w.Resize(fyne.NewSize(1024, 768))
	w.ShowAndRun()
}
