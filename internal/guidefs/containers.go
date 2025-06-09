package guidefs

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func initContainers() {

	Containers = map[string]WidgetInfo{
		"*fyne.Container": {
			Name: "Container",
			Create: func() fyne.CanvasObject {
				return container.NewVBox()
			},
			Edit: func(obj fyne.CanvasObject, props map[string]string, refresh func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
				c := obj.(*fyne.Container)

				choose := widget.NewFormItem("Layout", widget.NewSelect(layoutNames, nil))
				items := []*widget.FormItem{choose}
				ready := false
				choose.Widget.(*widget.Select).OnChanged = func(l string) {
					lay := Layouts[l]
					props["layout"] = l
					c.Layout = lay.Create(c, props)
					c.Refresh()

					edit := lay.Edit
					items = []*widget.FormItem{choose}
					if edit != nil {
						items = append(items, edit(c, props)...)
					}

					refresh(items)
					if ready {
						onchanged()
					}
				}
				choose.Widget.(*widget.Select).SetSelected(props["layout"])
				ready = true
				return items
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				c := obj.(*fyne.Container)
				l := props[c]["layout"]
				if l == "" {
					l = "VBox"
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
		"*container.AppTabs": {
			Name: "App Tabs",
			Children: func(o fyne.CanvasObject) []fyne.CanvasObject {
				tabs := o.(*container.AppTabs)

				children := make([]fyne.CanvasObject, len(tabs.Items))
				for i, c := range tabs.Items {
					children[i] = c.Content
				}
				return children
			},
			AddChild: func(parent, o fyne.CanvasObject) {
				tabs := o.(*container.AppTabs)

				item := container.NewTabItem("Untitled", o)
				tabs.Append(item)
			},
			Create: func() fyne.CanvasObject {
				return container.NewAppTabs(container.NewTabItem("Untitled", container.NewStack()))
			},
			Edit: func(obj fyne.CanvasObject, props map[string]string, setItems func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
				tabs := obj.(*container.AppTabs)
				items := make([]*widget.FormItem, len(tabs.Items)+2)
				itemNames := make([]string, len(tabs.Items))

				newRow := func(item *container.TabItem, i int) *widget.FormItem {
					icon := newIconSelectorButton(item.Icon, func(i fyne.Resource) {
						item.Icon = i
						tabs.Refresh()
						onchanged()
					}, false)
					edit := widget.NewEntry()
					edit.SetText(item.Text)
					edit.OnChanged = func(s string) {
						item.Text = s
						tabs.Refresh()
						onchanged()
					}
					del := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
						if i == len(tabs.Items)-1 {
							tabs.Items = tabs.Items[:i]
							items = items[:i]
							itemNames = itemNames[:i]
						} else {
							tabs.Items = append(tabs.Items[:i], tabs.Items[i+1:]...)
							items = append(items[:i], items[i+1:]...)
							itemNames = append(itemNames[:i], itemNames[i+1:]...)
						}
						tabs.Refresh()
						setItems(items)
						onchanged()
					})
					del.Importance = widget.DangerImportance

					tools := container.NewBorder(nil, nil, icon, del, edit)
					return widget.NewFormItem(fmt.Sprintf("Tab %d", i+1), tools)
				}
				for i, c := range tabs.Items {
					items[i] = newRow(c, i)
					itemNames[i] = fmt.Sprintf("%s (%d)", c.Text, i+1)
				}

				items[len(items)-2] = widget.NewFormItem("",
					widget.NewButton("Add Tab", func() {
						title := fmt.Sprintf("Tab %d", len(tabs.Items)+1)
						item := container.NewTabItem(title, container.NewStack())

						add := items[len(items)-2]
						sel := items[len(items)-1]
						newItem := newRow(item, len(tabs.Items))
						items = append(items[:len(items)-2], newItem, add, sel)
						itemNames = append(itemNames, title)

						tabs.Append(item)
						setItems(items)
						onchanged()
					}))
				ready := false
				selected := widget.NewSelect(itemNames, nil)
				selected.OnChanged = func(_ string) {
					tabs.SelectIndex(selected.SelectedIndex())
					if ready {
						onchanged()
					}
				}
				selected.SetSelectedIndex(tabs.SelectedIndex())
				ready = true
				items[len(items)-1] = widget.NewFormItem("Selected", selected)
				return items
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				tabs := obj.(*container.AppTabs)
				str := &strings.Builder{}
				str.WriteString("container.NewAppTabs(")

				for i, c := range tabs.Items {
					if i > 0 {
						str.WriteString(",\n")
					}

					hasIcon := c.Icon != nil
					constr := "NewTabItem"
					if hasIcon {
						constr = "NewTabItemWithIcon"
					}
					str.WriteString(fmt.Sprintf("container.%s(\"%s\", ", constr, c.Text))
					if hasIcon {
						str.WriteString("theme." + IconName(c.Icon) + "(), ")
					}
					writeGoStringExcluding(str, nil, props, defs, c.Content)
					str.WriteString(")")
				}
				str.WriteString(")")
				return widgetRef(props[obj], defs, str.String())
			},
			Packages: func(obj fyne.CanvasObject) []string {
				tabs := obj.(*container.AppTabs)
				for _, c := range tabs.Items {
					if c.Icon != nil {
						return []string{"container", "theme"}
					}
				}
				return []string{"container"}
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
			Edit: func(obj fyne.CanvasObject, props map[string]string, _ func([]*widget.FormItem), _ func()) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				s := obj.(*container.Scroll)
				str := &strings.Builder{}
				str.WriteString("container.NewScroll(")
				writeGoStringExcluding(str, nil, props, defs, s.Content)
				str.WriteString(")")
				return widgetRef(props[obj], defs, str.String())
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
			Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
				split := obj.(*container.Split)
				offset := widget.NewEntry()
				offset.SetText(fmt.Sprintf("%f", split.Offset))
				offset.OnChanged = func(s string) {
					if f, err := strconv.ParseFloat(s, 64); err == nil {
						split.SetOffset(f)
					}
					onchanged()
				}
				// TODO - add Fyne split.OnChanged
				vert := widget.NewCheck("", func(on bool) {
					split.Horizontal = !on
					split.Refresh()
					onchanged()
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
				return widgetRef(props[obj], defs, str.String())
			},
			Packages: func(_ fyne.CanvasObject) []string {
				return []string{"container"}
			},
		},
	}

	Containers["*widget.Scroll"] = Containers["*container.Scroll"] // internal widget name may be used

	ContainerNames = extractNames(Containers)
}
