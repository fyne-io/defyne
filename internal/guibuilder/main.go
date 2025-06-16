package guibuilder

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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
	gui.DefyneContext

	root, current fyne.CanvasObject
	uri           fyne.URI
	win           fyne.Window
	meta          map[fyne.CanvasObject]map[string]string
	th            fyne.Theme
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
	builder := &Builder{uri: u, win: win, meta: meta}
	var obj fyne.CanvasObject
	if r == nil {
		obj = previewUI()
	} else {
		obj, meta, err = gui.DecodeObject(r, builder)
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

	builder.root = obj
	return builder
}

func (b *Builder) Metadata() map[fyne.CanvasObject]map[string]string {
	return b.meta
}

func (b *Builder) Theme() fyne.Theme {
	return b.th
}

// MakeUI builds the UI for the current GUI builder.
func (b *Builder) MakeUI() fyne.CanvasObject {
	return b.buildUI(b.root)
}

// Run generates a go main function and runs it so we can preview the UI in a real app.
func (b *Builder) Run() {
	path := filepath.Join(os.TempDir(), "fynebuilder")
	dir := storage.NewFileURI(path)
	storage.CreateListable(dir)
	goURI, err := storage.Child(dir, "main.go")
	if err != nil {
		fyne.LogError("Failed to write temporary code", err)
		return
	}

	w, err := storage.Writer(goURI)
	if err != nil {
		fyne.LogError("Failed get storage writer", err)
		return
	}
	err = gui.ExportGoPreview(b.root, b, w)
	if err != nil {
		fyne.LogError("Failed to export go preview", err)
		return
	}

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
	name := strings.ReplaceAll(b.uri.Name(), ".gui.json", "")
	goFile := name + ".gui.go"
	dir, _ := storage.Parent(b.uri)

	goURI, err := storage.Child(dir, goFile)
	if err != nil {
		return err
	}

	w, err := storage.Writer(goURI)
	if err != nil {
		return err
	}
	err = gui.ExportGo(b.root, b, name, w)
	if err != nil {
		return err
	}

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
	err := gui.EncodeObject(b.root, b, w)
	_ = w.Close()
	return err
}

func (b *Builder) buildLibrary() fyne.CanvasObject {
	var selected *guidefs.WidgetInfo
	tempNames := []string{}
	widgetNames := []string{}
	addClass := func(name string) {
		widgetNames = append(widgetNames, name)
		tempNames = append(tempNames, name)
	}
	for _, name := range guidefs.WidgetNames {
		addClass(name)
	}
	for _, name := range guidefs.ContainerNames {
		addClass(name)
	}
	for _, name := range guidefs.CollectionNames {
		addClass(name)
	}
	for _, name := range guidefs.GraphicsNames {
		addClass(name)
	}
	list := widget.NewList(func() int {
		return len(tempNames)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(i widget.ListItemID, obj fyne.CanvasObject) {
		if i >= len(tempNames) {
			return
		}
		obj.(*widget.Label).SetText(guidefs.Lookup(tempNames[i]).Name)
	})
	list.OnSelected = func(i widget.ListItemID) {
		match := guidefs.Lookup(tempNames[i])
		if match != nil {
			selected = match
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
		for i := 0; i < len(widgetNames); i++ {
			test := strings.ToLower(guidefs.Lookup(widgetNames[i]).Name)
			if strings.Contains(test, s) {
				tempNames = append(tempNames, widgetNames[i])
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

		class := reflect.TypeOf(b.current).String()
		if wid := guidefs.Lookup(class); wid != nil && wid.IsContainer() {
			wid.AddChild(b.current, selected.Create())

			return
		}

		dialog.ShowInformation("Selected not a container", "Please select a container to add items", b.win)
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
	nameItem := widget.NewFormItem("Type", widget.NewLabel(gui.NameOf(o)))
	editForm = widget.NewForm()
	items := gui.EditorFor(o, b, func(items []*widget.FormItem) {
		editForm.Items = nil
		editForm.Refresh()
		editForm.Items = append([]*widget.FormItem{nameItem}, items...)
		editForm.Refresh()
	}, nil)

	items = append([]*widget.FormItem{nameItem}, items...)
	b.meta[o] = props

	editForm.Items = items
	remove := widget.NewButton("Remove", func() {
		parent := findParent(b.current, b.root)
		if parent == nil {
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
