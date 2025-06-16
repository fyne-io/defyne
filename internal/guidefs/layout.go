package guidefs

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type layoutInfo struct {
	Create func(*fyne.Container, DefyneContext) fyne.Layout
	Edit   func(*fyne.Container, DefyneContext) []*widget.FormItem
	goText func(*fyne.Container, DefyneContext, map[string]string) string
}

var (
	// layoutNames is an array with the list of names of all the Layouts
	layoutNames = extractLayoutNames()

	// Layouts maps container names to layout information to create and edit containers, and generate code
	Layouts = map[string]layoutInfo{
		"Border": {
			func(c *fyne.Container, d DefyneContext) fyne.Layout {
				props := d.Metadata()[c]
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
			func(c *fyne.Container, d DefyneContext) []*widget.FormItem {
				props := d.Metadata()[c]
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
					if c, ok := w.(*fyne.Container); ok {
						name := d.Metadata()[w]["name"]
						if name == "" {
							name = fmt.Sprintf("%p", c)
						}
						label = fmt.Sprintf("Container (%s)", name)
					} else {
						name := d.Metadata()[w]["name"]
						if name == "" {
							name = widgetName(w)
						}
						label = fmt.Sprintf("%s (%s)", reflect.TypeOf(w).Elem().Name(), name)
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
					widget.NewFormItem("Middle", widget.NewLabel("(all other widgets)")),
				}
			},
			func(c *fyne.Container, ctx DefyneContext, defs map[string]string) string {
				props := ctx.Metadata()[c]
				topNum := props["top"]
				topID, _ := strconv.Atoi(topNum)
				bottomNum := props["bottom"]
				bottomID, _ := strconv.Atoi(bottomNum)
				leftNum := props["left"]
				leftID, _ := strconv.Atoi(leftNum)
				rightNum := props["right"]
				rightID, _ := strconv.Atoi(rightNum)

				ignored := 0
				var t, b, l, r fyne.CanvasObject
				if topNum != "" && topID < len(c.Objects) {
					t = c.Objects[topID]
					ignored++
				}
				if bottomNum != "" && bottomID < len(c.Objects) {
					b = c.Objects[bottomID]
					ignored++
				}
				if leftNum != "" && leftID < len(c.Objects) {
					l = c.Objects[leftID]
					ignored++
				}
				if rightNum != "" && rightID < len(c.Objects) {
					r = c.Objects[rightID]
					ignored++
				}

				str := &strings.Builder{}
				str.WriteString("container.NewBorder(\n\t\t")
				writeGoStringOrNil(str, ctx, defs, t)
				str.WriteString(", \n\t\t")
				writeGoStringOrNil(str, ctx, defs, b)
				str.WriteString(", \n\t\t")
				writeGoStringOrNil(str, ctx, defs, l)
				str.WriteString(", \n\t\t")
				writeGoStringOrNil(str, ctx, defs, r)
				if len(c.Objects) > ignored {
					str.WriteString(", ")
					writeGoStringExcluding(str, func(o fyne.CanvasObject) bool {
						return o == t || o == b || o == l || o == r
					}, ctx, defs, c.Objects...)
				}
				str.WriteString(")")
				return str.String()
			},
		},
		"Center": {
			func(*fyne.Container, DefyneContext) fyne.Layout {
				return layout.NewCenterLayout()
			},
			nil,
			nil,
		},
		"Form": {
			func(*fyne.Container, DefyneContext) fyne.Layout {
				return layout.NewFormLayout()
			},
			nil,
			nil,
		},
		"Grid": {
			func(c *fyne.Container, d DefyneContext) fyne.Layout {
				props := d.Metadata()[c]
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
			func(c *fyne.Container, d DefyneContext) []*widget.FormItem {
				props := d.Metadata()[c]
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
			func(c *fyne.Container, ctx DefyneContext, defs map[string]string) string {
				props := ctx.Metadata()[c]
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

				str := &strings.Builder{}
				if rowCol == "Rows" {
					str.WriteString(fmt.Sprintf("container.NewGridWithRows(%d, ", num))
				} else {
					str.WriteString(fmt.Sprintf("container.NewGridWithColumns(%d, ", num))
				}
				writeGoStringExcluding(str, nil, ctx, defs, c.Objects...)
				str.WriteString(")")
				return str.String()
			},
		},
		"GridWrap": {
			func(c *fyne.Container, d DefyneContext) fyne.Layout {
				props := d.Metadata()[c]
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
			func(c *fyne.Container, d DefyneContext) []*widget.FormItem {
				props := d.Metadata()[c]
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
			func(c *fyne.Container, d DefyneContext) fyne.Layout {
				d.Metadata()[c]["dir"] = "horizontal"
				return layout.NewHBoxLayout()
			},
			nil,
			nil,
		},
		"Max": {
			func(_ *fyne.Container, _ DefyneContext) fyne.Layout {
				return layout.NewStackLayout()
			},
			nil,
			nil,
		},
		"Padded": {
			func(_ *fyne.Container, _ DefyneContext) fyne.Layout {
				return layout.NewPaddedLayout()
			},
			nil,
			nil,
		},
		"CustomPadded": {
			func(c *fyne.Container, d DefyneContext) fyne.Layout {
				props := d.Metadata()[c]
				pad := theme.Padding()
				padStr := strconv.FormatFloat(float64(pad), 'f', -2, 64)

				top := props["top"]
				if top == "" {
					top = padStr
				}
				bottom := props["bottom"]
				if bottom == "" {
					bottom = padStr
				}
				left := props["left"]
				if left == "" {
					left = padStr
				}
				right := props["right"]
				if right == "" {
					right = padStr
				}
				tt, err := strconv.ParseFloat(top, 64)
				t := float32(tt)
				if err != nil {
					t = pad
				}
				bb, err := strconv.ParseFloat(bottom, 64)
				b := float32(bb)
				if err != nil {
					b = pad
				}
				ll, err := strconv.ParseFloat(left, 64)
				l := float32(ll)
				if err != nil {
					l = pad
				}
				rr, err := strconv.ParseFloat(right, 64)
				r := float32(rr)
				if err != nil {
					r = pad
				}
				return layout.NewCustomPaddedLayout(t, b, l, r)
			},
			func(c *fyne.Container, d DefyneContext) []*widget.FormItem {
				props := d.Metadata()[c]
				pad := theme.Padding()
				padStr := strconv.FormatFloat(float64(pad), 'f', -2, 64)
				top := props["top"]
				if top == "" {

				}
				bottom := props["bottom"]
				if bottom == "" {
					bottom = padStr
				}
				left := props["left"]
				if left == "" {
					left = padStr
				}
				right := props["right"]
				if right == "" {
					right = padStr
				}

				topEnt := widget.NewEntry()
				topEnt.SetText(top)
				bottomEnt := widget.NewEntry()
				bottomEnt.SetText(bottom)
				leftEnt := widget.NewEntry()
				leftEnt.SetText(left)
				rightEnt := widget.NewEntry()
				rightEnt.SetText(right)
				change := func(string) {
					if topEnt.Text == "" {
						return
					}
					tt, err := strconv.ParseFloat(topEnt.Text, 64)
					if err != nil {
						return
					}
					t := float32(tt)
					if bottomEnt.Text == "" {
						return
					}
					bb, err := strconv.ParseFloat(bottomEnt.Text, 64)
					if err != nil {
						return
					}
					b := float32(bb)
					if leftEnt.Text == "" {
						return
					}
					ll, err := strconv.ParseFloat(leftEnt.Text, 64)
					if err != nil {
						return
					}
					l := float32(ll)
					if rightEnt.Text == "" {
						return
					}
					rr, err := strconv.ParseFloat(rightEnt.Text, 64)
					if err != nil {
						return
					}
					r := float32(rr)

					props["top"] = topEnt.Text
					props["bottom"] = bottomEnt.Text
					props["left"] = leftEnt.Text
					props["right"] = rightEnt.Text
					c.Layout = layout.NewCustomPaddedLayout(t, b, l, r)
					c.Refresh()
				}
				topEnt.OnChanged = change
				bottomEnt.OnChanged = change
				leftEnt.OnChanged = change
				rightEnt.OnChanged = change
				return []*widget.FormItem{
					widget.NewFormItem("Top", topEnt),
					widget.NewFormItem("Bottom", bottomEnt),
					widget.NewFormItem("Left", leftEnt),
					widget.NewFormItem("Right", rightEnt),
				}
			},
			func(c *fyne.Container, d DefyneContext, defs map[string]string) string {
				props := d.Metadata()
				pad := theme.Padding()
				padStr := strconv.FormatFloat(float64(pad), 'f', -2, 64)

				topNum := props[c]["top"]
				if topNum == "" {
					topNum = padStr
				}
				bottomNum := props[c]["bottom"]
				if bottomNum == "" {
					bottomNum = padStr
				}
				leftNum := props[c]["left"]
				if leftNum == "" {
					leftNum = padStr
				}
				rightNum := props[c]["right"]
				if rightNum == "" {
					rightNum = padStr
				}

				str := &strings.Builder{}
				str.WriteString("container.New(layout.NewCustomPaddedLayout(\n\t\t")
				str.WriteString(topNum)
				str.WriteString(", \n\t\t")
				str.WriteString(bottomNum)
				str.WriteString(", \n\t\t")
				str.WriteString(leftNum)
				str.WriteString(", \n\t\t")
				str.WriteString(rightNum)
				str.WriteString("), ")
				writeGoStringExcluding(str, func(o fyne.CanvasObject) bool {
					return false
				}, d, defs, c.Objects...)
				str.WriteString(")")
				return str.String()
			},
		},
		"Stack": {
			func(c *fyne.Container, d DefyneContext) fyne.Layout {
				return layout.NewStackLayout()
			},
			nil,
			nil,
		},
		"VBox": {
			func(c *fyne.Container, d DefyneContext) fyne.Layout {
				d.Metadata()[c]["dir"] = "vertical"
				return layout.NewVBoxLayout()
			},
			nil,
			nil,
		},
	}
)

// extractLayoutNames returns all the list of names of all the Layouts known
func extractLayoutNames() []string {
	var layoutsNamesFromData = make([]string, len(Layouts))
	i := 0
	for k := range Layouts {
		layoutsNamesFromData[i] = k
		i++
	}

	sort.Strings(layoutsNamesFromData)
	return layoutsNamesFromData
}

func trim(in string, count int) string {
	if len(in) > count {
		return in[:count] + "…"
	}

	return in
}

func widgetName(o fyne.CanvasObject) string {
	switch w := o.(type) {
	case *widget.Button:
		return trim(w.Text, 10)
	case *widget.Label:
		return trim(w.Text, 10)
	case *widget.Select:
		if len(w.Options) == 0 {
			return "No options"
		}

		return fmt.Sprintf("[%s, …]", w.Options[0])
	default:
		return fmt.Sprintf("%p", o)
	}
}

func writeGoString(str *strings.Builder, c DefyneContext,
	defs map[string]string, o fyne.CanvasObject) error {
	clazz := reflect.TypeOf(o).String()

	if match := Lookup(clazz); match != nil {
		code := GoString(clazz, o, c, defs)
		str.WriteString(fmt.Sprintf("\n\t\t%s", code))
	} else {
		return errors.New("failed to find go string for type" + clazz)
	}

	return nil
}

func writeGoStringOrNil(str *strings.Builder, c DefyneContext,
	defs map[string]string, o fyne.CanvasObject) {
	if o == nil {
		str.WriteString("nil")
		return
	}

	_ = writeGoString(str, c, defs, o)
}

func writeGoStringExcluding(str *strings.Builder, skip func(object fyne.CanvasObject) bool, c DefyneContext,
	defs map[string]string, objs ...fyne.CanvasObject) {
	for i, o := range objs {
		if skip != nil && skip(o) {
			continue
		}

		err := writeGoString(str, c, defs, o)
		if err != nil {
			fyne.LogError("Error writing Go string", err)
		} else if i < len(objs)-1 {
			str.WriteString(", ")
		}
	}
}
