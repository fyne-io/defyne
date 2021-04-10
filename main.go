//go:generate fyne bundle -o bundled.go Icon.png

package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("io.fyne.defyne")
	a.SetIcon(resourceIconPng)
	w := a.NewWindow("DEFyne")

	pwd, _ := os.Getwd()
	root := storage.NewFileURI(pwd)
	ide := &defyne{win: w, projectRoot: root}

	content := container.NewVSplit(ide.makeEditorPanel(), ide.makeTerminalPanel())
	content.Offset = 0.8
	mainSplit := container.NewHSplit(ide.makeFilesPanel(), content)
	mainSplit.Offset = 0.2

	w.SetContent(container.NewBorder(makeToolbar(), nil, nil, nil, mainSplit))
	w.Resize(fyne.NewSize(1024, 768))
	w.ShowAndRun()
}

func makeToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.FileIcon(), func() {}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {}))
}
