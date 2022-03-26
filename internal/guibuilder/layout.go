package guibuilder

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type layoutInfo struct {
	create func(*fyne.Container, map[string]string) fyne.Layout
	edit   func(*fyne.Container, map[string]string) []*widget.FormItem
	goText func(*fyne.Container, map[string]string) string
}

var (
	// layoutNames is an array with the list of names of all the layouts
	layoutNames = extractLayoutNames()
	layoutProps = make(map[*fyne.Container]map[string]string)

	layouts = map[string]layoutInfo{
		"Border": {
			func(c *fyne.Container, props map[string]string) fyne.Layout {
				topNum := props["top"]
				topID, _ := strconv.Atoi(topNum)
				bottomNum := props["bottom"]
				bottomID, _ := strconv.Atoi(bottomNum)
				leftNum := props["left"]
				leftID, _ := strconv.Atoi(leftNum)
				rightNum := props["right"]
				rightID, _ := strconv.Atoi(rightNum)

				var t, b, l, r fyne.CanvasObject
				if topNum != "" && topID < len(c.Objects) {
					t = c.Objects[topID]
				}
				if bottomNum != "" && bottomID < len(c.Objects) {
					b = c.Objects[bottomID]
				}
				if leftNum != "" && leftID < len(c.Objects) {
					l = c.Objects[leftID]
				}
				if rightNum != "" && rightID < len(c.Objects) {
					r = c.Objects[rightID]
				}

				return layout.NewBorderLayout(t, b, l, r)
			},
			func(c *fyne.Container, props map[string]string) []*widget.FormItem {
				topNum := props["top"]
				topID, _ := strconv.Atoi(topNum)
				bottomNum := props["bottom"]
				bottomID, _ := strconv.Atoi(bottomNum)
				leftNum := props["left"]
				leftID, _ := strconv.Atoi(leftNum)
				rightNum := props["right"]
				rightID, _ := strconv.Atoi(rightNum)

				var t, b, l, r fyne.CanvasObject
				list := []string{"(Empty)"}
				for _, w := range c.Objects {
					label := ""
					if c, ok := w.(*overlayContainer); ok {
						name := c.name
						if name == "" {
							name = fmt.Sprintf("%p", c.c)
						}
						label = fmt.Sprintf("Container (%s)", name)
					} else {
						wid := w.(*fyne.Container).Objects[0]
						name := w.(*fyne.Container).Objects[1].(*overlayWidget).name
						if name == "" {
							name = fmt.Sprintf("%p", wid)
						}
						label = fmt.Sprintf("%s (%s)", reflect.TypeOf(wid).Elem().Name(), name)
					}
					list = append(list, label)
				}
				top := widget.NewSelect(list, nil)
				if topNum != "" && topID < len(c.Objects) {
					top.SetSelectedIndex(topID + 1)
					t = c.Objects[topID]
				}
				bottom := widget.NewSelect(list, nil)
				if bottomNum != "" && bottomID < len(c.Objects) {
					bottom.SetSelectedIndex(bottomID + 1)
					b = c.Objects[bottomID]
				}
				left := widget.NewSelect(list, nil)
				if leftNum != "" && leftID < len(c.Objects) {
					left.SetSelectedIndex(leftID + 1)
					l = c.Objects[leftID]
				}
				right := widget.NewSelect(list, nil)
				if rightNum != "" && rightID < len(c.Objects) {
					right.SetSelectedIndex(rightID + 1)
					r = c.Objects[rightID]
				}
				change := func(string) {
					t, b, l, r = nil, nil, nil, nil
					props["top"] = ""
					props["bottom"] = ""
					props["left"] = ""
					props["right"] = ""
					if top.SelectedIndex() > 0 {
						props["top"] = strconv.Itoa(top.SelectedIndex() - 1)
						t = c.Objects[top.SelectedIndex()-1]
					}
					if bottom.SelectedIndex() > 0 {
						props["bottom"] = strconv.Itoa(bottom.SelectedIndex() - 1)
						b = c.Objects[bottom.SelectedIndex()-1]
					}
					if left.SelectedIndex() > 0 {
						props["left"] = strconv.Itoa(left.SelectedIndex() - 1)
						l = c.Objects[left.SelectedIndex()-1]
					}
					if right.SelectedIndex() > 0 {
						props["right"] = strconv.Itoa(right.SelectedIndex() - 1)
						r = c.Objects[right.SelectedIndex()-1]
					}

					c.Layout = layout.NewBorderLayout(t, b, l, r)
					c.Refresh()
				}
				top.OnChanged = change
				bottom.OnChanged = change
				left.OnChanged = change
				right.OnChanged = change
				c.Layout = layout.NewBorderLayout(t, b, l, r)

				return []*widget.FormItem{
					widget.NewFormItem("Top", top),
					widget.NewFormItem("Bottom", bottom),
					widget.NewFormItem("Left", left),
					widget.NewFormItem("Right", right),
				}
			},
			func(c *fyne.Container, props map[string]string) string {
				topNum := props["top"]
				topID, _ := strconv.Atoi(topNum)
				bottomNum := props["bottom"]
				bottomID, _ := strconv.Atoi(bottomNum)
				leftNum := props["left"]
				leftID, _ := strconv.Atoi(leftNum)
				rightNum := props["right"]
				rightID, _ := strconv.Atoi(rightNum)

				var t, b, l, r fyne.CanvasObject
				if topNum != "" && topID < len(c.Objects) {
					t = c.Objects[topID]
					if _, ok := t.(*overlayContainer); !ok {
						t = t.(*fyne.Container).Objects[1]
					}
				}
				if bottomNum != "" && bottomID < len(c.Objects) {
					b = c.Objects[bottomID]
					if _, ok := b.(*overlayContainer); !ok {
						b = b.(*fyne.Container).Objects[1]
					}
				}
				if leftNum != "" && leftID < len(c.Objects) {
					l = c.Objects[leftID]
					if _, ok := l.(*overlayContainer); !ok {
						l = l.(*fyne.Container).Objects[1]
					}
				}
				if rightNum != "" && rightID < len(c.Objects) {
					r = c.Objects[rightID]
					if _, ok := r.(*overlayContainer); !ok {
						r = r.(*fyne.Container).Objects[1]
					}
				}

				str := strings.Builder{}
				str.WriteString(fmt.Sprintf("container.NewBorder(\n\t\t%s, \n\t\t%s, \n\t\t%s, \n\t\t%s, ",
					goStringOrNil(t), goStringOrNil(b), goStringOrNil(l), goStringOrNil(r)))
				for i, o := range c.Objects {
					if _, ok := o.(*overlayContainer); !ok {
						o = o.(*fyne.Container).Objects[1]
					}
					if o == t || o == b || o == l || o == r {
						continue
					}

					str.WriteString(fmt.Sprintf("\n\t\t%#v", o))
					if i < len(c.Objects)-1 {
						str.WriteRune(',')
					}
				}
				str.WriteString(")\n")
				return str.String()
			},
		},
		"Center": {
			func(*fyne.Container, map[string]string) fyne.Layout {
				return layout.NewCenterLayout()
			},
			nil,
			nil,
		},
		"Form": {
			func(*fyne.Container, map[string]string) fyne.Layout {
				return layout.NewFormLayout()
			},
			nil,
			nil,
		},
		"Grid": {
			func(c *fyne.Container, props map[string]string) fyne.Layout {
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
			nil,
		},
		"GridWrap": {
			func(c *fyne.Container, props map[string]string) fyne.Layout {
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
			nil,
		},
		"HBox": {
			func(c *fyne.Container, props map[string]string) fyne.Layout {
				props["dir"] = "horizontal"
				return layout.NewHBoxLayout()
			},
			nil,
			nil,
		},
		"Max": {
			func(c *fyne.Container, props map[string]string) fyne.Layout {
				return layout.NewMaxLayout()
			},
			nil,
			nil,
		},
		"Padded": {
			func(c *fyne.Container, props map[string]string) fyne.Layout {
				return layout.NewPaddedLayout()
			},
			nil,
			nil,
		},
		"VBox": {
			func(c *fyne.Container, props map[string]string) fyne.Layout {
				props["dir"] = "vertical"
				return layout.NewVBoxLayout()
			},
			nil,
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

func goStringOrNil(o fyne.CanvasObject) string {
	if o == nil {
		return "nil"
	}

	return fmt.Sprintf("%#v", o)
}
