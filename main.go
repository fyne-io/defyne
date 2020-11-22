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
	paletteList *widget.Box
)

func buildLibrary() fyne.CanvasObject {
	widgets := []string{"Button", "Icon", "Label", "Entry"}

	var selected fyne.CanvasObject
	list := widget.NewList(func() int {
		return len(widgets)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(i widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*widget.Label).SetText(widgets[i])
	})
	list.OnSelected = func(i widget.ListItemID) {
		switch widgets[i] {
		case "Button":
			selected = widget.NewButton("Button", func() {})
		case "Icon":
			selected = widget.NewIcon(theme.HelpIcon())
		case "Label":
			selected = widget.NewLabel("Label")
		case "Entry":
			selected = widget.NewEntry()
			selected.(*widget.Entry).SetPlaceHolder("Entry")
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
	paletteList = widget.NewVBox()
	palette := container.NewBorder(widType, nil, nil, nil,
		container.NewGridWithRows(2, widget.NewCard("Properties", "", paletteList),
			widget.NewCard("Component List", "", buildLibrary()),
		))

	split := widget.NewHSplitContainer(overlay, palette)
	split.Offset = 0.8
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, nil, nil), toolbar,
		split)
}

func choose(o fyne.CanvasObject) {
	typeName := reflect.TypeOf(o).Elem().Name()
	widType.SetText(typeName)

	var items []fyne.CanvasObject
	switch obj := o.(type) {
	case *fyne.Container:
		items = []fyne.CanvasObject{
			widget.NewForm(widget.NewFormItem("Layout", widget.NewSelect([]string{"Center", "Grid", "GridWrap", "Max"}, func(l string) {
				switch l {
				case "Center":
					obj.Layout = layout.NewCenterLayout()
					obj.Refresh()
				case "Grid":
					obj.Layout = layout.NewGridLayout(2)
					obj.Refresh()
				case "GridWrap":
					obj.Layout = layout.NewGridWrapLayout(fyne.NewSize(100, 100))
					obj.Refresh()
				case "Max":
					obj.Layout = layout.NewMaxLayout()
					obj.Refresh()
				}
			}))),
		}
	case *widget.Label:
		entry := widget.NewEntry()
		entry.SetText(obj.Text)
		entry.OnChanged = func(text string) {
			obj.SetText(text)
		}
		items = []fyne.CanvasObject{widget.NewForm(widget.NewFormItem("Text", entry))}
	case *widget.Button:
		entry := widget.NewEntry()
		entry.SetText(obj.Text)
		entry.OnChanged = func(text string) {
			obj.SetText(text)
		}
		items = []fyne.CanvasObject{widget.NewForm(widget.NewFormItem("Text", entry),
			widget.NewFormItem("Icon", widget.NewSelect([]string{}, func(selected string) {})))}
	}
	paletteList.Children = items
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
