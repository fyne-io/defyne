package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (d *defyne) showProjectOpenDialog(w fyne.Window) {
	dialog.ShowFolderOpen(func(u fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, d.win)
			return
		}
		if u == nil {
			return
		}

		d.setProject(u)
		d.win.Show()
		w.Close()
	}, w)
}

func (d *defyne) showProjectSelect() {
	a := fyne.CurrentApp()
	w := a.NewWindow("Defyne : Open Project")

	recent := widget.NewList(
		func() int {
			return 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText("TODO")
		})

	img := canvas.NewImageFromResource(resourceIconPng)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(192, 192))
	open := container.NewBorder(widget.NewLabelWithStyle("Defyne", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Open Project", func() {
			d.showProjectOpenDialog(w)
		}), nil, nil, container.NewCenter(img))

	w.SetContent(container.NewGridWithColumns(2,
		container.NewBorder(widget.NewLabel("Recent projects"), nil, nil, nil, recent),
		open))
	w.Resize(fyne.NewSize(620, 440))
	w.Show()
}
