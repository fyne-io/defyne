package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type layoutInfo struct {
	create func() fyne.Layout
	edit   func() []*widget.FormItem
}

var layouts = map[string]layoutInfo{
	"Center": {
		layout.NewCenterLayout,
		nil,
	},
	"Form": {
		layout.NewFormLayout,
		nil,
	},
	"Grid": {
		func() fyne.Layout {
			return layout.NewGridLayout(2)
		},
		func() []*widget.FormItem {
			return []*widget.FormItem{
				widget.NewFormItem("Columns", widget.NewEntry()),
				widget.NewFormItem("Vertical", widget.NewCheck("", func(bool) {})),
			}
		},
	},
	"GridWrap": {
		func() fyne.Layout {
			return layout.NewGridWrapLayout(fyne.NewSize(100, 100))
		},
		func() []*widget.FormItem {
			return []*widget.FormItem{
				widget.NewFormItem("Item Width", widget.NewEntry()),
				widget.NewFormItem("Item Height", widget.NewEntry()),
			}
		},
	},
	"HBox": {
		layout.NewHBoxLayout,
		nil,
	},
	"Max": {
		layout.NewMaxLayout,
		nil,
	},
	"Padded": {
		layout.NewPaddedLayout,
		nil,
	},
	"VBox": {
		layout.NewVBoxLayout,
		nil,
	},
}

// layoutNames is an array with the list of names of all the layouts
var layoutNames = extractLayoutNames()

// extractLayoutNames returns all the list of names of all the layouts known
func extractLayoutNames() []string {
	var layoutsNamesFromData = make([]string, len(layouts))
	i := 0
	for k := range layouts {
		layoutsNamesFromData[i] = k
		i++
	}
	return layoutsNamesFromData
}
