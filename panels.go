package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	xWidget "fyne.io/x/fyne/widget"

	"github.com/fyne-io/terminal"
)

func (d *defyne) makeEditorPanel() fyne.CanvasObject {
	welcome := widget.NewLabel("Welcome to Defyne, the Fyne IDE.\n\nChoose a project file from the list.")
	d.fileTabs = container.NewDocTabs(
		container.NewTabItemWithIcon("Welcome", theme.HomeIcon(),
			container.NewCenter(welcome)))

	d.fileTabs.CloseIntercept = func(t *container.TabItem) {
		ed, ok := d.openEditors[t]
		if !ok { // welcome tab
			return
		}
		if !ed.changed() {
			ed.close()
			d.fileTabs.Remove(t)
			delete(d.openEditors, t)
			return
		}
		dialog.ShowConfirm("File is unsaved", "Are you sure you wish to close?",
			func(ok bool) {
				if !ok {
					return
				}

				ed.close()
				d.fileTabs.Remove(t)
				delete(d.openEditors, t)
			}, d.win)
	}

	return container.NewStack(d.fileTabs)
}

func (d *defyne) makeFilesPanel() *xWidget.FileTree {
	d.openEditors = make(map[*container.TabItem]*fileTab)

	files := xWidget.NewFileTree(d.projectRoot)
	files.Filter = filterHidden()
	files.Sorter = func(u1, u2 fyne.URI) bool {
		return u1.String() < u2.String() // Sort alphabetically
	}

	files.OnSelected = func(uid widget.TreeNodeID) {
		u, err := storage.ParseURI(uid)
		isDir, _ := storage.CanList(u)
		if isDir {
			return
		}

		if err != nil {
			dialog.ShowError(err, d.win)
		}

		d.openEditor(u)
	}
	return files
}

func (d *defyne) makeTerminalPanel() fyne.CanvasObject {
	term := terminal.New()
	term.SetStartDir(d.projectRoot.Path())
	go term.RunLocalShell()

	return term
}

type filter struct{}

func (f *filter) Matches(u fyne.URI) bool {
	name := u.Name()
	if name[0] == '.' {
		return false
	}

	if name == "go.sum" || (len(name) > 7 && strings.Index(name, ".gui.go") == len(u.Name())-7) {
		return false
	}
	return true
}

func filterHidden() storage.FileFilter {
	return &filter{}
}
