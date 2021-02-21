package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	editForm    *widget.Form
	widType     *widget.Label
	paletteList *fyne.Container
)

func buildLibrary() fyne.CanvasObject {
	var selected *widgetInfo
	tempNames := []string{}
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

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Search Widgets")
	entry.OnChanged = func(s string) {
		tempNames = []string{}
		for _, n := range widgetNames {
			if strings.Contains(n, s) {
				tempNames = append(tempNames, n)
			}
		}
		list.Refresh()
	}
	entry.SetText(" ")
	entry.SetText("")

	return container.NewBorder(entry, widget.NewButtonWithIcon("Insert", theme.ContentAddIcon(), func() {
		if c, ok := current.(*overlayContainer); ok {
			if selected != nil {
				c.c.Objects = append(c.c.Objects, wrapContent(selected.create()))
				c.c.Refresh()
			}
			return
		}
		log.Println("Please select a container")
	}), nil, nil, list)
}

func buildUI() fyne.CanvasObject {
	content := previewUI().(*fyne.Container)
	overlay := wrapContent(content)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			log.Println("TODO")
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			log.Println("TODO")
		}),
		widget.NewToolbarAction(theme.MailForwardIcon(), func() {
			code := fmt.Sprintf("%#v", overlay)
			fmt.Println(code)
		}))

	widType = widget.NewLabelWithStyle("(None Selected)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	paletteList = container.NewVBox()
	palette := container.NewBorder(widType, nil, nil, nil,
		container.NewGridWithRows(2, widget.NewCard("Properties", "", paletteList),
			widget.NewCard("Component List", "", buildLibrary()),
		))

	split := container.NewHSplit(overlay, palette)
	split.Offset = 0.8
	return container.New(layout.NewBorderLayout(toolbar, nil, nil, nil), toolbar,
		split)
}

func choose(o fyne.CanvasObject) {
	typeName := reflect.TypeOf(o).Elem().Name()
	widName := reflect.TypeOf(o).String()
	l := reflect.ValueOf(o).Elem()
	if typeName == "Entry" {
		if l.FieldByName("Password").Bool() {
			typeName = "PasswordEntry"
		} else if l.FieldByName("MultiLine").Bool() {
			typeName = "MultiLineEntry"
		}
		widName = "*widget." + typeName
	}
	widType.SetText(typeName)

	var items []*widget.FormItem
	if match, ok := widgets[widName]; ok {
		items = match.edit(o)
	}

	editForm = widget.NewForm(items...)
	paletteList.Objects = []fyne.CanvasObject{editForm}
	paletteList.Refresh()
}

func main() {
	a := app.NewWithID("xyz.andy.fynebuilder")
	w := a.NewWindow("Fyne Builder")
	w.SetContent(buildUI())
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}

func previewUI() fyne.CanvasObject {
	return container.New(layout.NewVBoxLayout(),
		widget.NewIcon(theme.ContentAddIcon()),
		widget.NewLabel("label"),
		widget.NewButton("Button", func() {}))
}
