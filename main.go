//go:generate fyne bundle -o bundled.go Icon.png

package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

	w.SetContent(container.NewBorder(ide.makeToolbar(), nil, nil, nil, mainSplit))
	w.Resize(fyne.NewSize(1024, 768))
	w.ShowAndRun()
}

func (d *defyne) makeToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.FileIcon(), func() {
			input := widget.NewEntry()
			dialog.ShowForm("New file name", "Create", "Cancel",
				[]*widget.FormItem{
					widget.NewFormItem("File name", input),
				},
				func(ok bool) {
					if !ok || input.Text == "" {
						return
					}

					uri, err := storage.Child(d.projectRoot, input.Text)
					if err != nil {
						dialog.ShowError(err, d.win)
						return
					}

					w, err := storage.Writer(uri)
					if err != nil {
						dialog.ShowError(err, d.win)
						return
					}
					_, _ = w.Write([]byte{})
					_ = w.Close()

					d.openEditor(uri)
					d.fileTree.Refresh()
				}, d.win)
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			if ed, ok := d.openEditors[d.fileTabs.Selected()]; ok {
				ed.save()
			}
		}))
}
