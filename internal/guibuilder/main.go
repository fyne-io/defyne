package guibuilder

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
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
	widType     *widget.Label
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
		obj, meta = gui.DecodeJSON(r)
		_ = r.Close()

		if obj == nil {
			obj = previewUI()
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
	packagesList := append(packagesRequired(b.root), "app")
	code := exportCode(packagesList, varsRequired(b.root, b.meta[b.root]), b.root)
	code += `
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	gui := newGUI()
	myWindow.SetContent(gui.makeUI())
	myWindow.ShowAndRun()
}
`
	path := filepath.Join(os.TempDir(), "fynebuilder")
	os.MkdirAll(path, 0711)
	path = filepath.Join(path, "main.go")
	_ = ioutil.WriteFile(path, []byte(code), 0600)

	cmd := exec.Command("go", "run", path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Start()
}

// Save will trigger the current state to be written out to the file this was opened from.
func (b *Builder) Save() error {
	w, err := storage.Writer(b.uri)
	if err != nil {
		return err
	}
	err = b.save(w)
	if err != nil {
		return err
	}

	goFile := strings.ReplaceAll(w.URI().Name(), ".gui.json", ".gui.go")
	dir, _ := storage.Parent(w.URI())
	goURI, err := storage.Child(dir, goFile)
	if err != nil {
		return err
	}
	packagesList := packagesRequired(b.root)
	code := exportCode(packagesList, varsRequired(b.root, b.meta[b.root]), b.root)
	w, err = storage.Writer(goURI)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(code))
	_ = w.Close()
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

	widType = widget.NewLabelWithStyle("(None Selected)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	widName = widget.NewEntry()
	widName.Validator = validation.NewRegexp("^$|^[a-zA-Z_][a-zA-Z0-9_]*$", "Invalid variable name")
	paletteList = container.NewVBox()
	palette := container.NewBorder(container.NewVBox(widType,
		widget.NewForm(widget.NewFormItem("Variable", widName))), nil, nil, nil,
		container.NewGridWithRows(2, widget.NewCard("Properties", "",
			container.NewVScroll(paletteList)),
			widget.NewCard("Component List", "", b.buildLibrary()),
		))

	split := container.NewHSplit(wrap, palette)
	split.Offset = 0.8
	return split
}

func packagesRequired(obj fyne.CanvasObject) []string {
	if w, ok := obj.(fyne.Widget); ok {
		return packagesRequiredForWidget(w)
	}

	ret := []string{"container"}
	var objs []fyne.CanvasObject
	if c, ok := obj.(*fyne.Container); ok {
		objs = c.Objects
	} else if c, ok := obj.(*fyne.Container); ok {
		objs = c.Objects
	}
	for _, w := range objs {
		for _, p := range packagesRequired(w) {
			added := false
			for _, exists := range ret {
				if p == exists {
					added = true
					break
				}
			}
			if !added {
				ret = append(ret, p)
			}
		}
	}
	return ret
}

func packagesRequiredForWidget(w fyne.Widget) []string {
	name := reflect.TypeOf(w).String()
	if guidefs.Widgets[name].Packages != nil {
		return guidefs.Widgets[name].Packages(w)
	}

	return []string{"widget"}
}

func varsRequired(obj fyne.CanvasObject, props map[string]string) []string {
	name := props["name"]
	if w, ok := obj.(fyne.Widget); ok {
		if name == "" {
			return []string{}
		}

		_, class := getTypeOf(w)
		return []string{name + " " + class}
	}

	var ret []string
	var objs []fyne.CanvasObject
	if c, ok := obj.(*fyne.Container); ok {
		objs = c.Objects
	} else if c, ok := obj.(*fyne.Container); ok {
		objs = c.Objects

		if name != "" {
			ret = append(ret, name+" "+"*fyne.Container")
		}
	}
	for _, w := range objs {
		ret = append(ret, varsRequired(w, props)...)
	}
	return ret
}

func (b *Builder) choose(o fyne.CanvasObject) {
	b.current = o

	name := b.meta[o]["name"]
	typeName, class := getTypeOf(o)
	widType.SetText(typeName)
	widName.OnChanged = func(s string) {
		props := b.meta[o]
		if props == nil {
			b.meta[o] = make(map[string]string)
		}
		b.meta[o]["name"] = s
	}
	widName.SetText(name)

	var items []*widget.FormItem
	if match, ok := guidefs.Widgets[class]; ok {
		props := b.meta[o]
		items = match.Edit(o, props)
		b.meta[o] = props
	}

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

var defs map[string]string // TODO find a better (non-global, non-race) way...

func exportCode(pkgs, vars []string, obj fyne.CanvasObject) string {
	for i := 0; i < len(pkgs); i++ {
		if pkgs[i] != "net/url" {
			pkgs[i] = "fyne.io/fyne/v2/" + pkgs[i]
		}

		pkgs[i] = fmt.Sprintf(`	"%s"`, pkgs[i])
	}

	defs = make(map[string]string)
	main := fmt.Sprintf("%#v", obj) // start GoString conversion
	setup := ""
	for k, v := range defs {
		setup += "g." + k + " = " + v + "\n"
	}

	code := fmt.Sprintf(`// auto-generated
// Code generated by Defyne GUI builder.

package main

import (
	"fyne.io/fyne/v2"
%s
)

type gui struct {
%s
}

func newGUI() *gui {
	return &gui{}
}

func (g *gui) makeUI() fyne.CanvasObject {
	%s

	return %s}
`,
		strings.Join(pkgs, "\n"),
		strings.Join(vars, "\n"),
		setup, main)

	formatted, err := format.Source([]byte(code))
	if err != nil {
		log.Println(code)
		fyne.LogError("Failed to encode GUI code", err)
		return ""
	}
	return string(formatted)
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

func getTypeOf(o fyne.CanvasObject) (string, string) {
	typeName := reflect.TypeOf(o).Elem().Name()
	class := reflect.TypeOf(o).String()
	l := reflect.ValueOf(o).Elem()
	if typeName == "Entry" {
		if l.FieldByName("Password").Bool() {
			typeName = "PasswordEntry"
		} else if l.FieldByName("MultiLine").Bool() {
			typeName = "MultiLineEntry"
		}
		class = "*widget." + typeName
	}

	return typeName, class
}

func previewUI() fyne.CanvasObject {
	return container.New(layout.NewVBoxLayout(),
		widget.NewLabel("label"),
		widget.NewButtonWithIcon("Button", theme.HomeIcon(), func() {}))
}
