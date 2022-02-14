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
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	editForm    *widget.Form
	widType     *widget.Label
	paletteList *fyne.Container
	once        sync.Once
)

// Builder is a simple type handle for a GUI builder instance.
type Builder struct {
	root, wrapped fyne.CanvasObject
	uri           fyne.URI
	win           fyne.Window
}

// NewBuilder returns an instance of the GUI builder for the specified URI.
// The Window parameter allows presenting dialogs etc.
func NewBuilder(u fyne.URI, win fyne.Window) *Builder {
	initOnce()
	r, err := storage.Reader(u)
	if err != nil {
		dialog.ShowError(err, win)
	}

	var obj, w fyne.CanvasObject
	if r == nil {
		obj = previewUI()
	} else {
		obj, w = DecodeJSON(r)
		_ = r.Close()

		if obj == nil {
			obj = previewUI()
			w = wrapContent(obj, nil)
		}
	}

	return &Builder{root: obj, wrapped: w, uri: u, win: win}
}

// MakeUI builds the UI for the current GUI builder.
func (b *Builder) MakeUI() fyne.CanvasObject {
	return b.buildUI(b.root)
}

// Run generates a go main function and runs it so we can preview the UI in a real app.
func (b *Builder) Run() {
	packagesList := append(packagesRequired(b.wrapped), "app")
	code := exportCode(packagesList, b.wrapped)
	code += `
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	myWindow.SetContent(makeUI())
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
	packagesList := packagesRequired(b.wrapped)
	code := exportCode(packagesList, b.wrapped)
	w, err = storage.Writer(goURI)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(code))
	_ = w.Close()
	return err
}

func (b *Builder) save(w fyne.URIWriteCloser) error {
	err := EncodeJSON(b.wrapped, w)
	_ = w.Close()
	return err
}

func buildLibrary() fyne.CanvasObject {
	var selected *widgetInfo
	tempNames := []string{}
	widgetLowerNames := []string{}
	for _, name := range widgetNames {
		widgetLowerNames = append(widgetLowerNames, strings.ToLower(name))
		tempNames = append(tempNames, name)
	}
	list := widget.NewList(func() int {
		return len(tempNames)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(i widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*widget.Label).SetText(widgets[tempNames[i]].name)
	})
	list.OnSelected = func(i widget.ListItemID) {
		if match, ok := widgets[tempNames[i]]; ok {
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
				tempNames = append(tempNames, widgetNames[i])
			}
		}
		list.Refresh()
		list.Select(0)   // Needed for new selection
		list.Unselect(0) // Without this (and with the above), list is behaving in a weird way
	}

	return container.NewBorder(searchBox, widget.NewButtonWithIcon("Insert", theme.ContentAddIcon(), func() {
		if c, ok := current.(*overlayContainer); ok {
			if selected != nil {
				c.c.Objects = append(c.c.Objects, wrapContent(selected.create(), c.c))
				c.c.Refresh()
			}
			return
		}
		log.Println("Please select a container")
	}), nil, nil, list)
}

func (b *Builder) buildUI(content fyne.CanvasObject) fyne.CanvasObject {
	wrap := container.NewMax(b.wrapped)

	widType = widget.NewLabelWithStyle("(None Selected)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	paletteList = container.NewVBox()
	palette := container.NewBorder(widType, nil, nil, nil,
		container.NewGridWithRows(2, widget.NewCard("Properties", "", paletteList),
			widget.NewCard("Component List", "", buildLibrary()),
		))

	split := container.NewHSplit(wrap, palette)
	split.Offset = 0.8
	return split
}

func packagesRequired(obj fyne.CanvasObject) []string {
	if w, ok := obj.(*overlayWidget); ok {
		return w.Packages()
	}

	ret := []string{"container"}
	var objs []fyne.CanvasObject
	if c, ok := obj.(*fyne.Container); ok {
		objs = c.Objects
	} else if c, ok := obj.(*overlayContainer); ok {
		objs = c.c.Objects
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

func choose(ow fyne.CanvasObject) {
	var o fyne.CanvasObject
	var name string
	o1, o2 := unwrap(ow)
	if o1 != nil {
		o = o1.c
		name = o1.name
	} else {
		o = o2.child
		name = o2.name
	}
	typeName := getTypeOf(o)
	class := "*widget." + typeName
	widType.SetText(typeName)
	widName.OnChanged = func(s string) {
		if o1 != nil {
			o1.name = s
		} else {
			o2.name = s
		}
	}
	widName.SetText(name)

	var items []*widget.FormItem
	if match, ok := widgets[class]; ok {
		items = match.edit(o)
	}

	editForm = widget.NewForm(items...)
	remove := widget.NewButton("Remove", func() {
		var parent *fyne.Container
		var obj fyne.CanvasObject
		if c, ok := current.(*overlayContainer); ok {
			parent = c.parent
			obj = c
		} else if w, ok := current.(*overlayWidget); ok {
			parent = w.parent
			for _, o := range parent.Objects { // match our widget in the container wrapping us
				if c, ok := o.(*fyne.Container); ok && c.Objects[0] == w.child {
					obj = c
					break
				}
			}
		}
		if parent == nil {
			log.Println("Nothing to remove")
			return
		}

		parent.Remove(obj)
		parent.Refresh()
	})
	paletteList.Objects = []fyne.CanvasObject{editForm, remove}
	paletteList.Refresh()
}

func exportCode(pkgs []string, obj fyne.CanvasObject) string {
	for i := 0; i < len(pkgs); i++ {
		pkgs[i] = fmt.Sprintf(`	"fyne.io/fyne/v2/%s"`, pkgs[i])
	}
	code := fmt.Sprintf(`// auto-generated
// Code generated by Defyne GUI builder.

package main

import (
	"fyne.io/fyne/v2"
%s
)

func makeUI() fyne.CanvasObject {
	return %#v
}
`,
		strings.Join(pkgs, "\n"),
		obj)

	formatted, err := format.Source([]byte(code))
	if err != nil {
		fyne.LogError("Failed to encode GUI code", err)
		return ""
	}
	return string(formatted)
}

func initOnce() {
	once.Do(func() {
		initIcons()
		initWidgets()
	})
}

func previewUI() fyne.CanvasObject {
	return container.New(layout.NewVBoxLayout(),
		widget.NewLabel("label"),
		widget.NewButtonWithIcon("Button", theme.HomeIcon(), func() {}))
}
