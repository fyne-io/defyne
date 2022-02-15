package guibuilder

import (
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type layoutInfo struct {
	create func(map[string]string) fyne.Layout
	edit   func(*fyne.Container, map[string]string) []*widget.FormItem
}

var (
	// layoutNames is an array with the list of names of all the layouts
	layoutNames = extractLayoutNames()
	layoutProps = make(map[*fyne.Container]map[string]string)

	layouts = map[string]layoutInfo{
		"Center": {
			func(map[string]string) fyne.Layout {
				return layout.NewCenterLayout()
			},
			nil,
		},
		"Form": {
			func(map[string]string) fyne.Layout {
				return layout.NewFormLayout()
			},
			nil,
		},
		"Grid": {
			func(props map[string]string) fyne.Layout {
				rowCol := props["grid_type"]
				if rowCol == "" {
					rowCol = "Columns"
				}
				count := props["count"]
				if count == "" {
					count = "2"
				}

				num, err := strconv.ParseInt(count, 0, 0)
				if err != nil {
					num = 2
				}

				if rowCol == "Rows" {
					return layout.NewGridLayoutWithRows(int(num))
				}
				return layout.NewGridLayoutWithColumns(int(num))
			},
			func(c *fyne.Container, props map[string]string) []*widget.FormItem {
				rowCol := props["grid_type"]
				if rowCol == "" {
					rowCol = "Columns"
				}
				count := props["count"]
				if count == "" {
					count = "2"
				}

				cols := widget.NewEntry()
				cols.SetText(count)
				vert := widget.NewSelect([]string{"Columns", "Rows"}, nil)
				vert.SetSelected(rowCol)
				change := func(string) {
					if cols.Text == "" {
						return
					}
					num, err := strconv.ParseInt(cols.Text, 0, 0)
					if err != nil {
						return
					}

					props["grid_type"] = vert.Selected
					props["count"] = cols.Text
					if vert.Selected == "Rows" {
						c.Layout = layout.NewGridLayoutWithRows(int(num))
					} else {
						c.Layout = layout.NewGridLayoutWithColumns(int(num))
					}
					c.Refresh()
				}
				cols.OnChanged = change
				vert.OnChanged = change
				return []*widget.FormItem{
					widget.NewFormItem("Count", cols),
					widget.NewFormItem("Arrange in", vert),
				}
			},
		},
		"GridWrap": {
			func(props map[string]string) fyne.Layout {
				width := props["width"]
				if width == "" {
					width = "100"
				}
				height := props["height"]
				if height == "" {
					height = "100"
				}
				w, err := strconv.ParseInt(width, 0, 0)
				if err != nil {
					w = 100
				}
				h, err := strconv.ParseInt(height, 0, 0)
				if err != nil {
					h = 100
				}

				return layout.NewGridWrapLayout(fyne.NewSize(float32(w), float32(h)))
			},
			func(c *fyne.Container, props map[string]string) []*widget.FormItem {
				width := props["width"]
				if width == "" {
					width = "100"
				}
				height := props["height"]
				if height == "" {
					height = "100"
				}

				widthEnt := widget.NewEntry()
				widthEnt.SetText(width)
				heightEnt := widget.NewEntry()
				heightEnt.SetText(height)
				change := func(string) {
					if widthEnt.Text == "" {
						return
					}
					w, err := strconv.ParseInt(widthEnt.Text, 0, 0)
					if err != nil {
						return
					}
					if widthEnt.Text == "" {
						return
					}
					h, err := strconv.ParseInt(heightEnt.Text, 0, 0)
					if err != nil {
						return
					}

					props["width"] = widthEnt.Text
					props["height"] = heightEnt.Text
					c.Layout = layout.NewGridWrapLayout(fyne.NewSize(float32(w), float32(h)))
					c.Refresh()
				}
				widthEnt.OnChanged = change
				heightEnt.OnChanged = change
				return []*widget.FormItem{
					widget.NewFormItem("Item Width", widthEnt),
					widget.NewFormItem("Item Height", heightEnt),
				}
			},
		},
		"HBox": {
			func(props map[string]string) fyne.Layout {
				props["dir"] = "horizontal"
				return layout.NewHBoxLayout()
			},
			nil,
		},
		"Max": {
			func(props map[string]string) fyne.Layout {
				return layout.NewMaxLayout()
			},
			nil,
		},
		"Padded": {
			func(props map[string]string) fyne.Layout {
				return layout.NewPaddedLayout()
			},
			nil,
		},
		"VBox": {
			func(props map[string]string) fyne.Layout {
				props["dir"] = "vertical"
				return layout.NewVBoxLayout()
			},
			nil,
		},
	}
)

// extractLayoutNames returns all the list of names of all the layouts known
func extractLayoutNames() []string {
	var layoutsNamesFromData = make([]string, len(layouts))
	i := 0
	for k := range layouts {
		layoutsNamesFromData[i] = k
		i++
	}

	sort.Strings(layoutsNamesFromData)
	return layoutsNamesFromData
}
