package main

import (
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fyne-io/defyne/internal/envcheck"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (d *defyne) menuActionNew() {
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
}

func (d *defyne) menuActionRunProject() {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = d.projectRoot.Path()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Start()
}

func (d *defyne) menuActionRun() {
	if ed, ok := d.openEditors[d.fileTabs.Selected()]; ok {
		ed.run()
	}
}

func (d *defyne) menuActionSave() {
	if ed, ok := d.openEditors[d.fileTabs.Selected()]; ok {
		ed.save()
	}
}

func (d *defyne) menuActionFullScreenToggle() {
	d.win.SetFullScreen(!d.win.FullScreen())
}

func (d *defyne) makeHelpMenu() *fyne.Menu {
	return fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://fyne.io/defyne")
			_ = fyne.CurrentApp().OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support")
			_ = fyne.CurrentApp().OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Check Environment...", func() {
			envcheck.ShowEnvCheckDialog(d.win)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://fyne.io/sponsor")
			_ = fyne.CurrentApp().OpenURL(u)
		}),
	)
}

func (d *defyne) makeMenu() *fyne.MainMenu {
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Open Project...", d.showProjectSelect),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("New File...", d.menuActionNew),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Save", d.menuActionSave),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Run", d.menuActionRun),
			fyne.NewMenuItem("Run Project", d.menuActionRunProject),
		))
	if runtime.GOOS != "darwin" {
		menu.Items = append(menu.Items,
			fyne.NewMenu("Window",
				fyne.NewMenuItem("Full Screen", d.menuActionFullScreenToggle),
			))
	}
	menu.Items = append(menu.Items,
		d.makeHelpMenu(),
	)
	return menu
}

func (d *defyne) makeToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.FileIcon(), d.menuActionNew),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), d.menuActionSave),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.MediaPlayIcon(), d.menuActionRun),
		widget.NewToolbarAction(theme.NewThemedResource(resourceFolderPlaySvg), d.menuActionRunProject))
}
