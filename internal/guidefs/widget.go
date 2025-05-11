//go:generate fyne bundle -o bundled.go -package guidefs ../../assets

package guidefs

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const noIconLabel = "(No Icon)"

var (
	// WidgetNames is an array with the list of names of all the Widgets
	WidgetNames []string
	// CollectionNames is an array with collection type names
	CollectionNames []string
	// ContainerNames is an array of known container type names
	ContainerNames []string

	// Widgets maps names to widget information and corresponding functions
	Widgets map[string]WidgetInfo
	// Collections maps names to widget information and corresponding functions
	Collections map[string]WidgetInfo
	// Containers maps names to widget information and corresponding functions
	Containers map[string]WidgetInfo
	once       sync.Once

	importances = []string{"Medium", "High", "Low", "Danger", "Warning", "Success"}
)

// WidgetInfo contains the name and corresponding functions for the widget type
type WidgetInfo struct {
	Name     string
	Children func(o fyne.CanvasObject) []fyne.CanvasObject
	AddChild func(parent, child fyne.CanvasObject)
	Create   func() fyne.CanvasObject
	Edit     func(fyne.CanvasObject, map[string]string, func([]*widget.FormItem), func()) []*widget.FormItem
	Gostring func(fyne.CanvasObject, map[fyne.CanvasObject]map[string]string, map[string]string) string
	Packages func(object fyne.CanvasObject) []string
}

// IsContainer indicates wether a widget children or not
func (w WidgetInfo) IsContainer() bool {
	return w.Children != nil
}

func initWidgets() {
	Widgets = map[string]WidgetInfo{
		"*widget.Button":     initButtonWidget(),
		"*widget.Hyperlink":  initHyperlinkWidget(),
		"*widget.Card":       initCardWidget(),
		"*widget.Entry":      initEntryWidget(),
		"*widget.Icon":       initIconWidget(),
		"*widget.Label":      initLabelWidget(),
		"*widget.RichText":   initRichTextWidget(),
		"*widget.Check":      initCheckWidget(),
		"*widget.RadioGroup": initRadioGroupWidget(),
		"*widget.Select":     initSelectWidget(),
		"*layout.Spacer": {
			Name: "Spacer",
			Create: func() fyne.CanvasObject {
				return layout.NewSpacer()
			},
			Edit: func(_ fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), _ func()) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs, "layout.NewSpacer()")
			},
			Packages: func(_ fyne.CanvasObject) []string {
				return []string{"layout"}
			},
		},
		"*widget.DateEntry": initDateEntryWidget(),
		"*widget.Accordion": initAccordionWidget(),
		"*widget.Form":      initFormWidget(),
		"*widget.MultiLineEntry": {
			Name: "Multi Line Entry",
			Create: func() fyne.CanvasObject {
				mle := widget.NewMultiLineEntry()
				mle.SetPlaceHolder("Enter Some \nLong text \nHere")
				mle.Wrapping = fyne.TextWrapWord
				return mle
			},
			// The rest inherits from Entry
		},
		"*widget.PasswordEntry": {
			Name: "Password Entry",
			Create: func() fyne.CanvasObject {
				e := widget.NewPasswordEntry()
				e.SetPlaceHolder("Password Entry")
				return e
			},
			// The rest inherits from Entry
		},
		"*widget.ProgressBar": initProgressBarWidget(),
		"*widget.Separator": {
			// Separator's height(or width as you may call) and color come from the theme, so not sure if we can change the color and height here
			Name: "Separator",
			Create: func() fyne.CanvasObject {
				return widget.NewSeparator()
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), _ func()) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs, "widget.NewSeparator()")
			},
		},
		"*widget.Slider":   initSliderWidget(),
		"*widget.TextGrid": initTextGridWidget(),
		"*widget.Toolbar":  initToolbarWidget(),
	}

	Collections = map[string]WidgetInfo{
		"*widget.List": {
			Name: "List",
			Create: func() fyne.CanvasObject {
				return widget.NewList(func() int {
					return 5
				}, func() fyne.CanvasObject {
					return widget.NewLabel("Template Object")
				}, func(id widget.ListItemID, item fyne.CanvasObject) {
					item.(*widget.Label).SetText(fmt.Sprintf("Item %d", id+1))
				})
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), _ func()) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs,
					`widget.NewList(func() int {
				return 5
			}, func() fyne.CanvasObject {
				return widget.NewLabel("Template Object")
			}, func(id widget.ListItemID, item fyne.CanvasObject) {
				item.(*widget.Label).SetText(fmt.Sprintf("Item %d", id+1))
			})`)
			},
			Packages: func(obj fyne.CanvasObject) []string {
				return []string{"widget", "fmt"}
			},
		},
		"*widget.Table": {
			Name: "Table",
			Create: func() fyne.CanvasObject {
				return widget.NewTable(func() (int, int) {
					return 3, 3
				}, func() fyne.CanvasObject {
					return widget.NewLabel("Cell 000, 000")
				}, func(id widget.TableCellID, cell fyne.CanvasObject) {
					label := cell.(*widget.Label)
					label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
				})
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), _ func()) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs,
					`widget.NewTable(func() (int, int) {
				return 3, 3
			}, func() fyne.CanvasObject {
				return widget.NewLabel("Cell 000, 000")
			}, func(id widget.TableCellID, cell fyne.CanvasObject) {
				label := cell.(*widget.Label)
				label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			})`)
			},
			Packages: func(obj fyne.CanvasObject) []string {
				return []string{"widget", "fmt"}
			},
		},
		"*widget.Tree": {
			Name: "Tree",
			Create: func() fyne.CanvasObject {
				data := map[string][]string{
					"":  {"A"},
					"A": {"B", "D", "H", "J", "L", "O", "P", "S", "V"},
					"B": {"C"},
					"C": {"abc"},
					"D": {"E"},
					"E": {"F", "G"},
					"F": {"adef"},
					"G": {"adeg"},
					"H": {"I"},
					"I": {"ahi"},
					"O": {"ao"},
					"P": {"Q"},
					"Q": {"R"},
					"R": {"apqr"},
					"S": {"T"},
					"T": {"U"},
					"U": {"astu"},
					"V": {"W"},
					"W": {"X"},
					"X": {"Y"},
					"Y": {"Z"},
					"Z": {"avwxyz"},
				}

				tree := widget.NewTreeWithStrings(data)
				tree.OnSelected = func(id string) {
					fmt.Println("Tree node selected:", id)
				}
				tree.OnUnselected = func(id string) {
					fmt.Println("Tree node unselected:", id)
				}
				tree.OpenBranch("A")
				tree.OpenBranch("D")
				tree.OpenBranch("E")
				tree.OpenBranch("L")
				tree.OpenBranch("M")
				return tree
			},
			Edit: func(co fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), _ func()) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs,
					`widget.NewTreeWithStrings(
map[string][]string{
	"":  {"A"},
	"A": {"B", "D", "H", "J", "L", "O", "P", "S", "V"},
	"B": {"C"},
	"C": {"abc"},
	"D": {"E"},
	"E": {"F", "G"},
	"F": {"adef"},
	"G": {"adeg"},
	"H": {"I"},
	"I": {"ahi"},
	"O": {"ao"},
	"P": {"Q"},
	"Q": {"R"},
	"R": {"apqr"},
	"S": {"T"},
	"T": {"U"},
	"U": {"astu"},
	"V": {"W"},
	"W": {"X"},
	"X": {"Y"},
	"Y": {"Z"},
	"Z": {"avwxyz"},
})`)
			},
		},
	}

	WidgetNames = extractNames(Widgets)
	CollectionNames = extractNames(Collections)
}

func initAccordionWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Accordion",
		Create: func() fyne.CanvasObject {
			return widget.NewAccordion(widget.NewAccordionItem("Item 1", widget.NewLabel("The content goes here")), widget.NewAccordionItem("Item 2", widget.NewLabel("Content part 2 goes here")))
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			acc := obj.(*widget.Accordion)
			multi := widget.NewCheck("", func(on bool) {
				acc.MultiOpen = on
				acc.Refresh()
				onchanged()
			})
			multi.Checked = acc.MultiOpen

			// TODO: Need to add the properties
			// entry := widget.NewEntry()
			return []*widget.FormItem{widget.NewFormItem("Multiple Open", multi)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			items := "widget.NewAccordionItem(\"Item 1\", widget.NewLabel(\"The content goes here\")), widget.NewAccordionItem(\"Item 2\", widget.NewLabel(\"Content part 2 goes here\"))"
			acc := obj.(*widget.Accordion)
			if acc.MultiOpen {
				return widgetRef(props[obj], defs,
					fmt.Sprintf("&widget.Accordion{Items: []*widget.AccordionItem{%s}, MultiOpen: true}", items))
			}

			return widgetRef(props[obj], defs,
				fmt.Sprintf("widget.NewAccordion(%s)", items))
		},
	}
}

func initButtonWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Button",
		Create: func() fyne.CanvasObject {
			return widget.NewButton("Button", func() {})
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			b := obj.(*widget.Button)
			entry := widget.NewEntry()
			entry.SetText(b.Text)
			entry.OnChanged = func(text string) {
				b.SetText(text)
				onchanged()
			}

			ready := false
			importance := widget.NewSelect(importances, func(s string) {
				var i widget.Importance
				for ii, imp := range importances {
					if imp == s {
						i = widget.Importance(ii)
					}
				}
				b.Importance = i
				b.Refresh()
				if ready {
					onchanged()
				}
			})
			importance.SetSelectedIndex(int(b.Importance))

			var left, center, right *widget.Button
			setAlign := func(a widget.ButtonAlign) {
				b.Alignment = a
				b.Refresh()

				setState := func(tb *widget.Button, a widget.ButtonAlign) {
					if b.Alignment == a {
						tb.Importance = widget.HighImportance
					} else {
						tb.Importance = widget.MediumImportance
					}
					tb.Refresh()
				}

				setState(left, widget.ButtonAlignLeading)
				setState(center, widget.ButtonAlignCenter)
				setState(right, widget.ButtonAlignTrailing)
				if ready {
					onchanged()
				}
			}
			left = widget.NewButtonWithIcon("", resourceFormatalignleftSvg, func() {
				setAlign(widget.ButtonAlignLeading)
			})
			center = widget.NewButtonWithIcon("", resourceFormataligncenterSvg, func() {
				setAlign(widget.ButtonAlignCenter)
			})
			right = widget.NewButtonWithIcon("", resourceFormatalignrightSvg, func() {
				setAlign(widget.ButtonAlignTrailing)
			})
			aligns := container.NewHBox(left, center, right)
			setAlign(b.Alignment)

			ready = true
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry),
				widget.NewFormItem("Icon", newIconSelectorButton(b.Icon, b.SetIcon, true)),
				widget.NewFormItem("Importance", importance),
				widget.NewFormItem("Alignment", aligns),
			}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			b := obj.(*widget.Button)
			action := props[obj]["OnTapped"]
			if action == "" {
				action = "func() {}"
			}
			if b.Icon == nil {
				if b.Importance == widget.MediumImportance && b.Alignment == widget.ButtonAlignCenter {
					return widgetRef(props[obj], defs, fmt.Sprintf("widget.NewButton(\"%s\", %s)", escapeLabel(b.Text), action))
				}

				return widgetRef(props[obj], defs, fmt.Sprintf("&widget.Button{Text: \"%s\", Importance: %d, Alignment: %d, OnTapped: %s}",
					escapeLabel(b.Text), b.Importance, b.Alignment, action))
			}

			icon := "theme." + IconName(b.Icon) + "()"
			if b.Importance == widget.MediumImportance && b.Alignment == widget.ButtonAlignCenter {
				return widgetRef(props[obj], defs, fmt.Sprintf("widget.NewButtonWithIcon(\"%s\", %s, %s)", escapeLabel(b.Text), icon, action))
			}

			return widgetRef(props[obj], defs, fmt.Sprintf("&widget.Button{Text: \"%s\", Importance: %d, Icon: %s, Alignment: %d, OnTapped: %s}",
				escapeLabel(b.Text), b.Importance, icon, b.Alignment, action))
		},
		Packages: func(obj fyne.CanvasObject) []string {
			b := obj.(*widget.Button)
			if b.Icon == nil {
				return []string{"widget"}
			}

			return []string{"widget", "theme"}
		},
	}
}

func initCardWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Card",
		Create: func() fyne.CanvasObject {
			return widget.NewCard("Title", "Subtitle", widget.NewLabel("Content here"))
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			c := obj.(*widget.Card)
			title := widget.NewEntry()
			title.SetText(c.Title)
			title.OnChanged = func(text string) {
				c.SetTitle(text)
				onchanged()
			}
			subtitle := widget.NewEntry()
			subtitle.SetText(c.Subtitle)
			subtitle.OnChanged = func(text string) {
				c.SetSubTitle(text)
				onchanged()
			}
			return []*widget.FormItem{
				widget.NewFormItem("Title", title),
				widget.NewFormItem("Subtitle", subtitle)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			c := obj.(*widget.Card)
			return widgetRef(props[obj], defs, fmt.Sprintf("widget.NewCard(\"%s\", \"%s\", widget.NewLabel(\"Content here\"))",
				escapeLabel(c.Title), escapeLabel(c.Subtitle)))
		},
	}
}

func initCheckWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Check",
		Create: func() fyne.CanvasObject {
			return widget.NewCheck("Tick it or don't", func(b bool) {})
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			c := obj.(*widget.Check)
			title := widget.NewEntry()
			title.SetText(c.Text)
			title.OnChanged = func(text string) {
				c.Text = text
				c.Refresh()
				onchanged()
			}
			ready := false
			isChecked := widget.NewCheck("", func(b bool) {
				c.SetChecked(b)
				if ready {
					onchanged()
				}
			})
			isChecked.SetChecked(c.Checked)
			return []*widget.FormItem{
				widget.NewFormItem("Title", title),
				widget.NewFormItem("isChecked", isChecked)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			c := obj.(*widget.Check)
			return widgetRef(props[obj], defs,
				fmt.Sprintf("widget.NewCheck(\"%s\", func(b bool) {})", escapeLabel(c.Text)))
		},
	}
}

func initDateEntryWidget() WidgetInfo {
	return WidgetInfo{
		Name: "DateEntry",
		Create: func() fyne.CanvasObject {
			return widget.NewDateEntry()
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), _ func()) []*widget.FormItem {
			_ = obj.(*widget.DateEntry)

			return []*widget.FormItem{}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			_ = obj.(*widget.DateEntry)

			return widgetRef(props[obj], defs, "widget.NewDateEntry()")
		},
	}
}

func initEntryWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Entry",
		Create: func() fyne.CanvasObject {
			e := widget.NewEntry()
			e.SetPlaceHolder("Entry")
			return e
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			l := obj.(*widget.Entry)
			entry1 := widget.NewEntry()
			entry1.SetText(l.Text)
			entry1.OnChanged = func(text string) {
				l.SetText(text)
				onchanged()
			}
			entry2 := widget.NewEntry()
			entry2.SetText(l.PlaceHolder)
			entry2.OnChanged = func(text string) {
				l.SetPlaceHolder(text)
				onchanged()
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry1),
				widget.NewFormItem("PlaceHolder", entry2)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			l := obj.(*widget.Entry)
			return widgetRef(props[obj], defs,
				fmt.Sprintf("&widget.Entry{Text: \"%s\", PlaceHolder: \"%s\", MultiLine: %t, Password: %t}",
					escapeLabel(l.Text), escapeLabel(l.PlaceHolder), l.MultiLine, l.Password))
		},
	}
}

func initFormWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Form",
		Create: func() fyne.CanvasObject {
			f := widget.NewForm(widget.NewFormItem("Username", widget.NewEntry()), widget.NewFormItem("Password", widget.NewPasswordEntry()), widget.NewFormItem("Remember", widget.NewCheck("", func(bool) {})))
			f.OnSubmit = func() {}
			f.OnCancel = func() {}
			return f
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, refresh func([]*widget.FormItem), _ func()) []*widget.FormItem {
			items := []*widget.FormItem{}
			formItems := obj.(*widget.Form).Items

			tidyWidget := func(wid, src fyne.CanvasObject) {
				switch t := wid.(type) {
				case *widget.Entry:
					t.PlaceHolder = ""
					t.Password = src.(*widget.Entry).Password
					t.MultiLine = src.(*widget.Entry).MultiLine
				case *widget.Check:
					t.Text = ""
				}
			}

			// TODO get the window passed in somehow
			w := fyne.CurrentApp().Driver().AllWindows()[0]

			var editItem, removeItem func(int)
			add := widget.NewButtonWithIcon("Add...", theme.ContentAddIcon(), nil)
			addLine := widget.NewFormItem("", add)
			add.OnTapped = func() {

				insertChose := ""
				name := widget.NewEntry()
				options := widget.NewSelect([]string{"Check", "DateEntry", "Entry", "MultiLineEntry", "PasswordEntry", "Select", "Slider"}, func(s string) {
					insertChose = s
				})

				items := []*widget.FormItem{
					widget.NewFormItem("Name", name),
					widget.NewFormItem("Type", options),
				}

				dialog.ShowForm("Add form item", "Add", "Cancel", items, func(ok bool) {
					if !ok {
						return
					}

					class := "*widget." + insertChose
					wid1 := Lookup(class).Create()
					tidyWidget(wid1, wid1)
					wid2 := Lookup(class).Create()
					tidyWidget(wid2, wid2)

					obj.(*widget.Form).Items = nil // TODO fix this too!
					obj.Refresh()
					formItems = append(formItems, widget.NewFormItem(name.Text, wid1))
					obj.(*widget.Form).Items = formItems
					obj.Refresh()

					var row *fyne.Container
					edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
						index := getFormIndex(row, items)
						editItem(index)
					})
					remove := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
						index := getFormIndex(row, items)
						removeItem(index)
					})
					row = container.NewBorder(nil, nil, edit, remove, wid2)
					items = append(items, widget.NewFormItem(insertChose, row))
					refresh(append(items, addLine))
				}, w)
			}

			editItem = func(i int) {
				editText := widget.NewEntry()
				editText.SetText(formItems[i].Text)

				dialog.ShowForm("Edit form item", "Save", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", editText)}, func(ok bool) {
					if !ok {
						return
					}

					formItems[i].Text = editText.Text
					obj.Refresh()

					items[i].Text = editText.Text
					refresh(append(items, addLine))
				}, w)
			}
			removeItem = func(i int) {
				obj.(*widget.Form).Items = nil
				obj.Refresh() // working around Form refresh bug

				formItems = removeFormItem(i, formItems)
				obj.(*widget.Form).Items = formItems
				obj.Refresh()
				items = removeFormItem(i, items)
				refresh(append(items, addLine))
			}

			for _, o := range formItems {
				// copy items so the widgets can be in both places
				class := reflect.TypeOf(o.Widget).String()
				wid := Lookup(class).Create()
				tidyWidget(wid, o.Widget)

				var row *fyne.Container
				edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
					index := getFormIndex(row, items)
					editItem(index)
				})
				remove := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
					index := getFormIndex(row, items)
					removeItem(index)
				})
				row = container.NewBorder(nil, nil, edit, remove, wid)
				items = append(items, widget.NewFormItem(o.Text, row))
			}

			return append(items, addLine)
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			str := &strings.Builder{}
			str.WriteString("&widget.Form{Items: []*widget.FormItem{")
			for _, i := range obj.(*widget.Form).Items {
				str.WriteString("widget.NewFormItem(\"" + i.Text + "\", ")
				writeGoStringExcluding(str, nil, props, defs, i.Widget)
				str.WriteString("),")
			}
			str.WriteString("}, OnSubmit: func() {}, OnCancel: func() {}}")
			return widgetRef(props[obj], defs, str.String())
		},
	}
}

func initHyperlinkWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Hyperlink",
		Create: func() fyne.CanvasObject {
			fyneURL, _ := url.Parse("https://fyne.io")
			return widget.NewHyperlink("Link Text", fyneURL)
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			link := obj.(*widget.Hyperlink)
			title := widget.NewEntry()
			title.SetText(link.Text)
			title.OnChanged = func(text string) {
				link.SetText(text)
				onchanged()
			}
			subtitle := widget.NewEntry()
			subtitle.SetText(link.URL.String())
			subtitle.OnChanged = func(text string) {
				fyneURL, _ := url.Parse(text)
				link.SetURL(fyneURL)
				onchanged()
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", title),
				widget.NewFormItem("URL", subtitle)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			link := obj.(*widget.Hyperlink)
			return widgetRef(props[obj], defs, fmt.Sprintf(`widget.NewHyperlink("%s", %#v)`, escapeLabel(link.Text), link.URL))
		},
		Packages: func(_ fyne.CanvasObject) []string {
			return []string{"net/url"}
		},
	}
}

func initIconWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Icon",
		Create: func() fyne.CanvasObject {
			return widget.NewIcon(theme.HelpIcon())
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			i := obj.(*widget.Icon)
			return []*widget.FormItem{
				widget.NewFormItem("Icon", newIconSelectorButton(i.Resource, func(res fyne.Resource) {
					i.SetResource(res)
					onchanged()
				}, true))}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			i := obj.(*widget.Icon)

			res := "theme." + IconName(i.Resource) + "()"
			return widgetRef(props[obj], defs, fmt.Sprintf("widget.NewIcon(%s)", res))
		},
		Packages: func(obj fyne.CanvasObject) []string {
			return []string{"widget", "theme"}
		},
	}
}

func initLabelWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Label",
		Create: func() fyne.CanvasObject {
			return widget.NewLabel("Label")
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			l := obj.(*widget.Label)
			entry := widget.NewEntry()
			entry.SetText(l.Text)
			entry.OnChanged = func(text string) {
				l.SetText(text)
				onchanged()
			}

			wrap := widget.NewCheck("", func(on bool) {
				if on {
					l.Wrapping = fyne.TextWrapWord
				} else {
					l.Wrapping = fyne.TextWrapOff
				}
				l.Refresh()
				onchanged()
			})
			wrap.Checked = l.Wrapping == fyne.TextWrapWord

			bold := widget.NewCheck("", func(on bool) {
				l.TextStyle.Bold = on
				l.Refresh()
				onchanged()
			})
			bold.Checked = l.TextStyle.Bold
			italic := widget.NewCheck("", func(on bool) {
				l.TextStyle.Italic = on
				l.Refresh()
				onchanged()
			})
			italic.Checked = l.TextStyle.Italic
			mono := widget.NewCheck("", func(on bool) {
				l.TextStyle.Monospace = on
				l.Refresh()
				onchanged()
			})
			mono.Checked = l.TextStyle.Monospace

			var left, center, right *widget.Button
			ready := false
			setAlign := func(a fyne.TextAlign) {
				l.Alignment = a
				l.Refresh()

				setState := func(b *widget.Button, a fyne.TextAlign) {
					if l.Alignment == a {
						b.Importance = widget.HighImportance
					} else {
						b.Importance = widget.MediumImportance
					}
					b.Refresh()
				}

				setState(left, fyne.TextAlignLeading)
				setState(center, fyne.TextAlignCenter)
				setState(right, fyne.TextAlignTrailing)
				if ready {
					onchanged()
				}
			}
			left = widget.NewButtonWithIcon("", resourceFormatalignleftSvg, func() {
				setAlign(fyne.TextAlignLeading)
			})
			center = widget.NewButtonWithIcon("", resourceFormataligncenterSvg, func() {
				setAlign(fyne.TextAlignCenter)
			})
			right = widget.NewButtonWithIcon("", resourceFormatalignrightSvg, func() {
				setAlign(fyne.TextAlignTrailing)
			})
			aligns := container.NewHBox(left, center, right)
			setAlign(l.Alignment)

			ready = true
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry),
				widget.NewFormItem("Word Wrap", wrap),
				widget.NewFormItem("Bold", bold),
				widget.NewFormItem("Italic", italic),
				widget.NewFormItem("Monospace", mono),
				widget.NewFormItem("Alignment", aligns)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			l := obj.(*widget.Label)
			if l.Alignment != fyne.TextAlignLeading || l.Wrapping != fyne.TextWrapOff {
				style := ""
				if l.TextStyle.Bold || l.TextStyle.Italic || l.TextStyle.Monospace {
					style = fmt.Sprintf(", TextStyle: %#v", l.TextStyle)
				}

				return widgetRef(props[obj], defs,
					fmt.Sprintf("&widget.Label{Text: \"%s\"%s, Alignment: %d, Wrapping: %d}", escapeLabel(l.Text), style, l.Alignment, l.Wrapping))
			}

			if l.TextStyle.Bold || l.TextStyle.Italic || l.TextStyle.Monospace {
				return widgetRef(props[obj], defs,
					fmt.Sprintf("widget.NewLabelWithStyle(\"%s\", %d, %#v)", escapeLabel(l.Text), l.Alignment, l.TextStyle))
			}
			return widgetRef(props[obj], defs,
				fmt.Sprintf("widget.NewLabel(\"%s\")", escapeLabel(l.Text)))
		},
	}
}

func initProgressBarWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Progress Bar",
		Create: func() fyne.CanvasObject {
			p := widget.NewProgressBar()
			p.SetValue(0.1)
			return p
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			p := obj.(*widget.ProgressBar)
			value := widget.NewEntry()
			value.SetText(fmt.Sprintf("%f", p.Value))
			value.OnChanged = func(s string) {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					p.SetValue(f)
				}
				onchanged()
			}
			return []*widget.FormItem{
				widget.NewFormItem("Value", value)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			p := obj.(*widget.ProgressBar)
			return widgetRef(props[obj], defs,
				fmt.Sprintf("&widget.ProgressBar{Value: %f}", p.Value))
		},
	}
}

func initRadioGroupWidget() WidgetInfo {
	return WidgetInfo{
		Name: "RadioGroup",
		Create: func() fyne.CanvasObject {
			return widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(s string) {})
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			r := obj.(*widget.RadioGroup)
			initialOption := widget.NewRadioGroup(r.Options, func(s string) {
				r.SetSelected(s)
				onchanged()
			})
			initialOption.SetSelected(r.Selected)
			entry := widget.NewMultiLineEntry()
			entry.SetText(strings.Join(r.Options, "\n"))
			entry.OnChanged = func(text string) {
				r.Options = strings.Split(text, "\n")
				r.Refresh()
				initialOption.Options = strings.Split(text, "\n")
				initialOption.Refresh()
				onchanged()
			}
			return []*widget.FormItem{
				widget.NewFormItem("Options", entry),
				widget.NewFormItem("Initial Option", initialOption)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			r := obj.(*widget.RadioGroup)
			var opts []string
			for _, v := range r.Options {
				opts = append(opts, escapeLabel(v))
			}
			return widgetRef(props[obj], defs,
				fmt.Sprintf("widget.NewRadioGroup([]string{%s}, func(s string) {})", "\""+strings.Join(opts, "\", \"")+"\""))
		},
	}
}

func initRichTextWidget() WidgetInfo {
	return WidgetInfo{
		Name: "RichText",
		Create: func() fyne.CanvasObject {
			return widget.NewRichTextFromMarkdown("## Rich Text")
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			r := obj.(*widget.RichText)
			entry := widget.NewMultiLineEntry()
			entry.SetText(r.String()) // TODO re-assemble the markdown !?
			entry.OnChanged = func(text string) {
				r.ParseMarkdown(text)
				onchanged()
			}

			wraps := map[string]fyne.TextWrap{
				"Off":   fyne.TextWrapOff,
				"Word":  fyne.TextWrapWord,
				"Break": fyne.TextWrapBreak,
			}
			wrap := widget.NewSelect([]string{"Off", "Word", "Break"}, func(w string) {
				r.Wrapping = wraps[w]
				r.Refresh()
				onchanged()
			})
			wrap.Selected = "Off"
			if r.Wrapping == fyne.TextWrapWord {
				wrap.Selected = "Word"
			} else if r.Wrapping == fyne.TextWrapBreak {
				wrap.Selected = "Break"
			}

			return []*widget.FormItem{
				widget.NewFormItem("Text", entry),
				widget.NewFormItem("Wrapping", wrap)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			l := obj.(*widget.RichText)
			// TODO wrap
			return widgetRef(props[obj], defs,
				fmt.Sprintf("widget.NewRichTextFromMarkdown(`%s`)", l.String())) // TODO re-assemble the markdown !?
		},
	}
}

func initSelectWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Select",
		Create: func() fyne.CanvasObject {
			return widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {})
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			s := obj.(*widget.Select)
			initialOption := widget.NewSelect(append([]string{"(Select one)"}, s.Options...), nil)
			initialOption.SetSelected(s.Selected)
			initialOption.OnChanged = func(opt string) {
				s.SetSelected(opt)
				if opt == "(Select one)" {
					s.ClearSelected()
				}
				onchanged()
			}
			entry := widget.NewMultiLineEntry()
			entry.SetText(strings.Join(s.Options, "\n"))
			entry.OnChanged = func(text string) {
				s.Options = strings.Split(text, "\n")
				s.Refresh()
				initialOption.Options = strings.Split(text, "\n")
				initialOption.Refresh()
				onchanged()
			}
			return []*widget.FormItem{
				widget.NewFormItem("Options", entry),
				widget.NewFormItem("Initial Option", initialOption)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			s := obj.(*widget.Select)
			var opts []string
			for _, v := range s.Options {
				opts = append(opts, escapeLabel(v))
			}

			optionString := "\"" + strings.Join(opts, "\", \"") + "\""
			if s.Selected == "" {
				return widgetRef(props[obj], defs,
					fmt.Sprintf("widget.NewSelect([]string{%s}, func(s string) {})", optionString))
			}

			format := "&widget.Select{Options: []string{%s}, Selected: \"%s\", OnChanged: func(s string) {}}"
			return widgetRef(props[obj], defs, fmt.Sprintf(format, optionString, s.Selected))
		},
	}
}

func initSliderWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Slider",
		Create: func() fyne.CanvasObject {
			s := widget.NewSlider(0, 100)
			s.OnChanged = func(f float64) {
				fmt.Println("Slider changed to", f)
			}
			return s
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			slider := obj.(*widget.Slider)
			val := widget.NewEntry()
			val.SetText(fmt.Sprintf("%f", slider.Value))
			val.OnChanged = func(s string) {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					slider.SetValue(f)
				}
				onchanged()
			}
			vert := widget.NewCheck("", func(on bool) {
				if on {
					slider.Orientation = widget.Vertical
				} else {
					slider.Orientation = widget.Horizontal
				}
				slider.Refresh()
				onchanged()
			})
			vert.Checked = slider.Orientation == widget.Vertical
			return []*widget.FormItem{
				widget.NewFormItem("Value", val),
				widget.NewFormItem("Vertical", vert),
			}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			slider := obj.(*widget.Slider)
			orient := "widget.Horizontal"
			if slider.Orientation == widget.Vertical {
				orient = "widget.Vertical"
			}
			return widgetRef(props[obj], defs, fmt.Sprintf("&widget.Slider{Min:0, Max:100, Value:%f, Orientation: %s}", slider.Value, orient))
		},
	}
}

func initTextGridWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Text Grid",
		Create: func() fyne.CanvasObject {
			to := widget.NewTextGrid()
			to.SetText("ABCD \nEFGH")
			return to
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, _ func([]*widget.FormItem), onchanged func()) []*widget.FormItem {
			to := obj.(*widget.TextGrid)
			entry := widget.NewEntry()
			entry.SetText(to.Text())
			entry.OnChanged = func(s string) {
				to.SetText(s)
				onchanged()
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry)}
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			to := obj.(*widget.TextGrid)
			return widgetRef(props[obj], defs,
				fmt.Sprintf("widget.NewTextGrid(\"%s\")", escapeLabel(to.Text())))
		},
	}
}

func initToolbarWidget() WidgetInfo {
	return WidgetInfo{
		Name: "Toolbar",
		Create: func() fyne.CanvasObject {
			return widget.NewToolbar(
				widget.NewToolbarAction(Icons["FileIcon"], func() { fmt.Println("Clicked on FileIcon") }),
				widget.NewToolbarSeparator(),
				widget.NewToolbarAction(Icons["ViewRefreshIcon"], func() { fmt.Println("Clicked on ViewRefreshIcon") }),
				widget.NewToolbarAction(Icons["NavigateBackIcon"], func() { fmt.Println("Clicked on NavigateBackIcon") }),
				widget.NewToolbarAction(Icons["NavigateNextIcon"], func() { fmt.Println("Clicked on NavigateNextIcon") }),
				widget.NewToolbarSpacer(),
				widget.NewToolbarAction(Icons["HelpIcon"], func() { fmt.Println("Clicked on HelpIcon") }),
			)
		},
		Edit: func(obj fyne.CanvasObject, _ map[string]string, refresh func([]*widget.FormItem), _ func()) []*widget.FormItem {
			items := []*widget.FormItem{}
			toolItems := obj.(*widget.Toolbar).Items

			add := widget.NewButtonWithIcon("Add...", theme.ContentAddIcon(), nil)
			addLine := widget.NewFormItem("", add)
			options := []string{"Action", "Separator", "Spacer"}

			removeItem := func(i int) {
				toolItems = removeToolbarItem(i, toolItems)
				obj.(*widget.Toolbar).Items = toolItems
				obj.Refresh()
				items = removeFormItem(i, items)
				refresh(append(items, addLine))
			}

			newToolEdit := func(id int, o widget.ToolbarItem) *widget.FormItem {
				chosen := ""
				holder := container.NewStack()
				var wid fyne.CanvasObject

				switch t := o.(type) {
				case *widget.ToolbarSeparator:
					chosen = options[1]
				case *widget.ToolbarSpacer:
					chosen = options[2]
				case *widget.ToolbarAction:
					chosen = options[0]
					wid = newIconSelectorButton(t.Icon, t.SetIcon, false)
					holder.Objects = []fyne.CanvasObject{wid}
				}

				chooser := widget.NewSelect(options, func(s string) {
					switch s {
					case "Separator":
						toolItems[id] = widget.NewToolbarSeparator()
						items[id].Text = "Separator"
						holder.Objects = nil
					case "Spacer":
						toolItems[id] = widget.NewToolbarSpacer()
						items[id].Text = "Spacer"
						holder.Objects = nil
					default:
						act := widget.NewToolbarAction(theme.QuestionIcon(), nil)
						toolItems[id] = act
						items[id].Text = "Action"

						holder.Objects = []fyne.CanvasObject{newIconSelectorButton(act.Icon, act.SetIcon, false)}
					}

					obj.Refresh()
					refresh(append(items, addLine))
				})
				chooser.Selected = chosen

				var row *fyne.Container
				remove := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
					index := getFormIndex(row, items)
					removeItem(index)
				})
				row = container.NewBorder(nil, nil, chooser, remove, holder)
				return widget.NewFormItem(chosen, row)
			}

			add.OnTapped = func() {
				newID := len(toolItems)
				newAction := widget.NewToolbarAction(theme.QuestionIcon(), nil)
				toolItems = append(toolItems, newAction)

				obj.(*widget.Toolbar).Items = toolItems
				obj.Refresh()
				items = append(items, newToolEdit(newID, newAction))

				obj.Refresh()
				refresh(append(items, addLine))
			}

			for i, o := range toolItems {
				items = append(items, newToolEdit(i, o))
			}

			return append(items, addLine)
		},
		Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
			str := &strings.Builder{}
			str.WriteString("widget.NewToolbar(\n")
			for _, i := range obj.(*widget.Toolbar).Items {
				switch t := i.(type) {
				case *widget.ToolbarSeparator:
					str.WriteString("\t\t\t\twidget.NewToolbarSeparator(),\n")
				case *widget.ToolbarSpacer:
					str.WriteString("\t\t\t\twidget.NewToolbarSpacer(),\n")
				case *widget.ToolbarAction:
					res := "theme." + IconName(t.Icon) + "()"
					str.WriteString(fmt.Sprintf("\t\t\t\twidget.NewToolbarAction(%s, func() {}),\n", res))
					// TODO action handler
				}
			}
			str.WriteString(")")
			return widgetRef(props[obj], defs, str.String())
		},
	}
}

// extractWidgetNames returns all the list of names of all the Widgets from our data
func extractNames(in map[string]WidgetInfo) []string {
	var widgetNamesFromData = make(widgetNames, len(in))
	dupe := false

	i := 0
	for k := range in {
		if k == "*widget.Scroll" { // do not duplicate scroll
			dupe = true
			continue
		}
		widgetNamesFromData[i] = k
		i++
	}

	if dupe {
		widgetNamesFromData = widgetNamesFromData[:len(widgetNamesFromData)-1]
	}

	sort.Sort(widgetNamesFromData)
	return widgetNamesFromData
}

func widgetRef(props map[string]string, defs map[string]string, code string) string {
	if name, ok := props["name"]; ok && name != "" {
		defs[name] = code
		return "g." + name
	}

	return code
}

// InitOnce initializes icons, graphics, containers, and widgets one time
func InitOnce() {
	once.Do(func() {
		initIcons()
		initGraphics()
		initContainers()
		initWidgets()
	})
}

// Lookup returns the [WidgetInfo] for the given widget type
func Lookup(clazz string) *WidgetInfo {
	if match, ok := Widgets[clazz]; ok {
		return &match
	}
	if match, ok := Collections[clazz]; ok {
		return &match
	}
	if match, ok := Containers[clazz]; ok {
		return &match
	}
	if match, ok := Graphics[clazz]; ok {
		return &match
	}

	return nil
}

type widgetNames []string

func (w widgetNames) Len() int {
	return len(w)
}

func (w widgetNames) Less(i, j int) bool {
	partsI := strings.Split(w[i], ".")
	partsJ := strings.Split(w[j], ".")

	return partsI[1] < partsJ[1]
}

func (w widgetNames) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}
