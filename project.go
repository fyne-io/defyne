package main

import (
	"io"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func (d *defyne) showNewProjectDialog(w fyne.Window) {
	var dir fyne.ListableURI
	parent := widget.NewButton("Choose directory", nil)
	parent.OnTapped = func() {
		dialog.ShowFolderOpen(func(u fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, d.win)
				return
			}
			if u == nil {
				return
			}

			dir = u
			parent.SetText(u.Name())
		}, w)
	}

	name := widget.NewEntry()
	dialog.ShowForm("Create project", "Create", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Parent directory", parent),
		widget.NewFormItem("Project name", name),
	}, func(ok bool) {
		if !ok {
			return
		}

		dir, err := createProject(dir, name.Text)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			d.setProject(dir)
			d.win.Show()
			w.Close()
		}
	}, w)
}

func (d *defyne) showOpenProjectDialog(w fyne.Window) {
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
	openProject := widget.NewButton("Open Project", func() {
		d.showOpenProjectDialog(w)
	})
	openProject.Importance = widget.HighImportance
	open := container.NewBorder(widget.NewLabelWithStyle("Defyne", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		openProject, nil, nil, container.NewCenter(img))

	create := widget.NewButton("New Project", func() {
		d.showNewProjectDialog(w)
	})
	w.SetContent(container.NewGridWithColumns(2,
		container.NewBorder(widget.NewLabel("Recent projects"), create, nil, nil, recent),
		open))
	w.Resize(fyne.NewSize(620, 440))
	w.Show()
}

func createProject(parent fyne.URI, name string) (fyne.URI, error) {
	dir, err := storage.Child(parent, name)
	if err != nil {
		return nil, err
	}

	err = storage.CreateListable(dir)
	if err != nil {
		return nil, err
	}

	err = writeFile(dir, "main.go", `package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("`+name+`")

	w.SetContent(widget.NewLabel("Hello `+name+`"))
	w.ShowAndRun()
}
`)
	if err != nil {
		fyne.LogError("Failed to write main.go", err) // we can just return the partial project
		return dir, nil
	}

	err = writeFile(dir, "go.mod", `module `+name+`

require fyne.io/fyne/v2 v2.1.1
`)
	if err != nil {
		fyne.LogError("Failed to write go.mod", err) // we can just return the partial project
		return dir, nil
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir.Path()
	err = cmd.Start() // run in background - may take a little while but should not block file editing
	if err != nil {   // just print, can just continue to open project
		fyne.LogError("Could not run go mod tidy", err)
	}
	return dir, nil
}

func writeFile(dir fyne.URI, name, content string) error {
	modURI, _ := storage.Child(dir, name)

	w, err := storage.Writer(modURI)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, content)
	if err != nil {
		fyne.LogError("Failed to write go.mod", err)
		return err
	}
	return w.Close()
}
