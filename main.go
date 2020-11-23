package main

import (
	"fmt"
	"log"
	"reflect"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var (
	widType     *widget.Label
	paletteList *fyne.Container
)

func buildLibrary() fyne.CanvasObject {
	var selected fyne.CanvasObject
	list := widget.NewList(func() int {
		return len(widgetNames)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(i widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*widget.Label).SetText(widgets[widgetNames[i]].name)
	})
	list.OnSelected = func(i widget.ListItemID) {
		if match, ok := widgets[widgetNames[i]]; ok {
			selected = match.create()
		}
	}

	return container.NewBorder(nil, widget.NewButtonWithIcon("Insert", theme.ContentAddIcon(), func() {
		if c, ok := current.(*overlayContainer); ok {
			if selected != nil {
				c.c.Objects = append(c.c.Objects, wrapContent(selected))
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
			log.Println(code)
		}))

	widType = widget.NewLabelWithStyle("(None Selected)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	paletteList = container.NewVBox()
	palette := container.NewBorder(widType, nil, nil, nil,
		container.NewGridWithRows(2, widget.NewCard("Properties", "", paletteList),
			widget.NewCard("Component List", "", buildLibrary()),
		))

	split := container.NewHSplit(overlay, palette)
	split.Offset = 0.8
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, nil, nil), toolbar,
		split)
}

func choose(o fyne.CanvasObject) {
	typeName := reflect.TypeOf(o).Elem().Name()
	widType.SetText(typeName)

	var items []*widget.FormItem
	if match, ok := widgets[reflect.TypeOf(o).String()]; ok {
		items = match.edit(o)
	}

	paletteList.Objects = []fyne.CanvasObject{widget.NewForm(items...)}
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
	return fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		widget.NewIcon(theme.ContentAddIcon()),
		widget.NewLabel("label"),
		widget.NewButton("Button", func() {}))
}
