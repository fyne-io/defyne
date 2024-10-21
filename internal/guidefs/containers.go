package guidefs

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func initContainers() {

	Containers = map[string]WidgetInfo{
		"*fyne.Container": {
			Name: "Container",
			Create: func() fyne.CanvasObject {
				return container.NewStack()
			},
			Edit: func(obj fyne.CanvasObject, props map[string]string) []*widget.FormItem {
				c := obj.(*fyne.Container)

				choose := widget.NewFormItem("Layout", widget.NewSelect(layoutNames, nil))
				items := []*widget.FormItem{choose}
				choose.Widget.(*widget.Select).OnChanged = func(l string) {
					lay := Layouts[l]
					props["layout"] = l
					c.Layout = lay.Create(c, props)
					c.Refresh()
					choose.Widget.Hide()

					edit := lay.Edit
					items = []*widget.FormItem{choose}
					if edit != nil {
						items = append(items, edit(c, props)...)
					}

					// TODO wtf?					editForm = widget.NewForm(items...)
					//					paletteList.Objects = []fyne.CanvasObject{editForm}
					choose.Widget.Show()
					//					paletteList.Refresh()
				}
				choose.Widget.(*widget.Select).SetSelected(props["layout"])
				return items
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				c := obj.(*fyne.Container)
				l := props[c]["layout"]
				if l == "" {
					l = "Stack"
				}
				lay := Layouts[l]
				if lay.goText != nil {
					return lay.goText(c, props, defs)
				}

				str := &strings.Builder{}
				if l == "Form" {
					str.WriteString("container.New(layout.NewFormLayout(), ")
				} else {
					str.WriteString(fmt.Sprintf("container.New%s(", l))
				}
				writeGoStringExcluding(str, nil, props, defs, c.Objects...)
				str.WriteString(")")
				return widgetRef(props[obj], defs, str.String())
			},
		},
		"*container.Scroll": {
			Name: "Scroll",
			Children: func(o fyne.CanvasObject) []fyne.CanvasObject {
				scr := o.(*container.Scroll)
				return []fyne.CanvasObject{scr.Content}
			},
			AddChild: func(parent, o fyne.CanvasObject) {
				scr := parent.(*container.Scroll)
				scr.Content = o
				scr.Refresh()
			},
			Create: func() fyne.CanvasObject {
				return container.NewScroll(container.NewStack())
			},
			Edit: func(obj fyne.CanvasObject, props map[string]string) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				s := obj.(*container.Scroll)
				str := &strings.Builder{}
				str.WriteString("container.NewScroll(")
				writeGoStringExcluding(str, nil, props, defs, s.Content)
				str.WriteString(")")
				return str.String()
			},
			Packages: func(_ fyne.CanvasObject) []string {
				return []string{"container"}
			},
		},
		"*container.Split": {
			Name: "Split",
			Children: func(o fyne.CanvasObject) []fyne.CanvasObject {
				split := o.(*container.Split)
				return []fyne.CanvasObject{split.Leading, split.Trailing}
			},
			AddChild: func(parent, o fyne.CanvasObject) {
				split := parent.(*container.Split)
				if split.Leading == nil {
					split.Leading = o
				} else {
					split.Trailing = o
				}
				split.Refresh()
			},
			Create: func() fyne.CanvasObject {
				return container.NewHSplit(container.NewStack(), container.NewStack())
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				split := obj.(*container.Split)
				offset := widget.NewEntry()
				offset.SetText(fmt.Sprintf("%f", split.Offset))
				offset.OnChanged = func(s string) {
					if f, err := strconv.ParseFloat(s, 64); err == nil {
						split.SetOffset(f)
					}
				}
				// TODO - add Fyne split.OnChanged
				vert := widget.NewCheck("", func(on bool) {
					split.Horizontal = !on
					split.Refresh()
				})
				vert.Checked = !split.Horizontal
				return []*widget.FormItem{
					widget.NewFormItem("Offset", offset),
					widget.NewFormItem("Vertical", vert),
				}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				s := obj.(*container.Split)
				str := &strings.Builder{}
				str.WriteString(fmt.Sprintf("&container.Split{Horizontal: %t, Offset: %f, Leading: ", s.Horizontal, s.Offset))
				writeGoStringExcluding(str, nil, props, defs, s.Leading)
				str.WriteString(", Trailing: ")
				writeGoStringExcluding(str, nil, props, defs, s.Trailing)
				str.WriteString("}")
				return str.String()
			},
			Packages: func(_ fyne.CanvasObject) []string {
				return []string{"container"}
			},
		},
	}

	Containers["*widget.Scroll"] = Containers["*container.Scroll"] // internal widget name may be used

	ContainerNames = extractNames(Containers)
}
