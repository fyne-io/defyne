package guibuilder

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fyne-io/defyne/internal/guidefs"
	"github.com/fyne-io/defyne/pkg/gui"
)

var (
	editForm    *widget.Form
	widName     *widget.Entry
	paletteList *fyne.Container
)

// Builder is a simple type handle for a GUI builder instance.
type Builder struct {
	root, current fyne.CanvasObject
	uri           fyne.URI
	win           fyne.Window
	meta          map[fyne.CanvasObject]map[string]string
}

// NewBuilder returns an instance of the GUI builder for the specified URI.
// The Window parameter allows presenting dialogs etc.
func NewBuilder(u fyne.URI, win fyne.Window) *Builder {
	guidefs.InitOnce()
	r, err := storage.Reader(u)
	if err != nil {
		dialog.ShowError(err, win)
	}

	meta := make(map[fyne.CanvasObject]map[string]string)
	var obj fyne.CanvasObject
	if r == nil {
		obj = previewUI()
	} else {
		obj, meta, err = gui.DecodeJSON(r)
		if err != nil {
			dialog.ShowError(err, win)
		}
		_ = r.Close()

		if obj == nil {
			obj = previewUI()
		}
		if meta == nil {
			meta = make(map[fyne.CanvasObject]map[string]string)
		}
	}

	return &Builder{root: obj, uri: u, win: win, meta: meta}
}

// MakeUI builds the UI for the current GUI builder.
func (b *Builder) MakeUI() fyne.CanvasObject {
	return b.buildUI(b.root)
}

// Run generates a go main function and runs it so we can preview the UI in a real app.
func (b *Builder) Run() {
	path := filepath.Join(os.TempDir(), "fynebuilder")
	goURI, err := storage.Child(storage.NewFileURI(path), "main.go")
	if err != nil {
		fyne.LogError("Failed to write temporary code", err)
		return
	}

	w, err := storage.Writer(goURI)
	err = gui.ExportGoPreview(b.root, b.meta, w)

	pwd, _ := os.Getwd()
	os.Chdir(path)
	cmd := exec.Command("go", "mod", "init", "temp")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
	cmd = exec.Command("go", "run", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Start()
	os.Chdir(pwd)
}

// Save will trigger the current state to be written out to the file this was opened from.
func (b *Builder) Save() error {
	goFile := strings.ReplaceAll(b.uri.Name(), ".gui.json", ".gui.go")
	dir, _ := storage.Parent(b.uri)

	goURI, err := storage.Child(dir, goFile)
	if err != nil {
		return err
	}

	w, err := storage.Writer(goURI)
	if err != nil {
		return err
	}
	err = gui.ExportGo(b.root, b.meta, w)

	_ = w.Close()

	w, err = storage.Writer(b.uri)
	if err != nil {
		return err
	}
	err = b.save(w)
	if err != nil {
		return err
	}

	return err
}

func (b *Builder) save(w fyne.URIWriteCloser) error {
	err := gui.EncodeJSON(b.root, b.meta, w)
	_ = w.Close()
	return err
}

func (b *Builder) buildLibrary() fyne.CanvasObject {
	var selected *guidefs.WidgetInfo
	tempNames := []string{}
	widgetLowerNames := []string{}
	for _, name := range guidefs.WidgetNames {
		widgetLowerNames = append(widgetLowerNames, strings.ToLower(name))
		tempNames = append(tempNames, name)
	}
	list := widget.NewList(func() int {
		return len(tempNames)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(i widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*widget.Label).SetText(guidefs.Widgets[tempNames[i]].Name)
	})
	list.OnSelected = func(i widget.ListItemID) {
		if match, ok := guidefs.Widgets[tempNames[i]]; ok {
			selected = &match
		}
	}
	list.OnUnselected = func(widget.ListItemID) {
		selected = nil
	}

	searchBox := widget.NewEntry()
	searchBox.SetPlaceHolder("Search Widgets")
	searchBox.OnChanged = func(s string) {
		s = strings.ToLower(s)
		tempNames = []string{}
		for i := 0; i < len(widgetLowerNames); i++ {
			if strings.Contains(widgetLowerNames[i], s) {
				tempNames = append(tempNames, guidefs.WidgetNames[i])
			}
		}
		list.Refresh()
		list.Select(0)   // Needed for new selection
		list.Unselect(0) // Without this (and with the above), list is behaving in a weird way
	}

	return container.NewBorder(searchBox, widget.NewButtonWithIcon("Insert", theme.ContentAddIcon(), func() {
		if c, ok := b.current.(*fyne.Container); ok {
			if selected != nil {
				c.Objects = append(c.Objects, selected.Create())
				c.Refresh()
				// cause property editor to refresh
				b.choose(c)
			}
			return
		}
		log.Println("Please select a container")
	}), nil, nil, list)
}

func (b *Builder) buildUI(content fyne.CanvasObject) fyne.CanvasObject {
	wrap := container.NewStack(b.root, newOverlay(b))

	widName = widget.NewEntry()
	widName.Validator = validation.NewRegexp("^$|^[a-zA-Z_][a-zA-Z0-9_]*$", "Invalid variable name")
	paletteList = container.NewVBox()
	palette := container.NewBorder(
		widget.NewForm(widget.NewFormItem("Variable", widName)), nil, nil, nil,
		container.NewGridWithRows(2, widget.NewCard("Properties", "",
			container.NewVScroll(paletteList)),
			widget.NewCard("Component List", "", b.buildLibrary()),
		))

	split := container.NewHSplit(wrap, palette)
	split.Offset = 0.8
	return split
}

func (b *Builder) choose(o fyne.CanvasObject) {
	b.current = o

	name := b.meta[o]["name"]
	widName.OnChanged = func(s string) {
		props := b.meta[o]
		if props == nil {
			b.meta[o] = make(map[string]string)
		}
		b.meta[o]["name"] = s
	}
	widName.SetText(name)

	props := b.meta[o]
	if props == nil {
		props = make(map[string]string)
	}
	items := gui.EditorFor(o, props)

	nameItem := widget.NewFormItem("Type", widget.NewLabel(gui.NameOf(o)))
	items = append([]*widget.FormItem{nameItem}, items...)
	b.meta[o] = props

	editForm = widget.NewForm(items...)
	remove := widget.NewButton("Remove", func() {
		var parent *fyne.Container
		if c, ok := b.current.(*fyne.Container); ok {
			parent = findParent(c, b.root)
		} else if w, ok := b.current.(fyne.Widget); ok {
			parent = findParent(w, b.root)
		}
		if parent == nil {
			log.Println("Nothing to remove")
			return
		}

		parent.Remove(b.current)
		parent.Refresh()
		b.choose(parent)
	})
	paletteList.Objects = []fyne.CanvasObject{editForm, remove}
	paletteList.Refresh()
}

func findParent(o fyne.CanvasObject, parent fyne.CanvasObject) *fyne.Container {
	switch w := parent.(type) {
	case *fyne.Container:
		for _, child := range w.Objects {
			if child == o {
				return w
			}
			if found := findParent(o, child); found != nil {
				return found
			}
		}
		return nil
	}

	return nil
}

func previewUI() fyne.CanvasObject {
	return container.New(layout.NewVBoxLayout(),
		widget.NewLabel("label"),
		widget.NewButtonWithIcon("Button", theme.HomeIcon(), func() {}))
}
