package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (d *defyne) makeMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("New File...", func() {}),
		))
}

func (d *defyne) makeToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.FileIcon(), func() {
			input := widget.NewEntry()
			typeNames := make([]string, len(templates))
			for i, t := range templates {
				typeNames[i] = t.name
			}

			types := widget.NewSelect(typeNames, func(s string) {
				name := strings.Split(input.Text, ".")[0]
				for _, t := range templates {
					if t.name == s {
						name += t.ext
						continue
					}
				}
				input.SetText(name)
			})
			dialog.ShowForm("New file name", "Create", "Cancel",
				[]*widget.FormItem{
					widget.NewFormItem("File type", types),
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
		}),
		widget.NewToolbarAction(theme.MailForwardIcon(), func() {
			if ed, ok := d.openEditors[d.fileTabs.Selected()]; ok {
				ed.run()
			}
		}))
}
