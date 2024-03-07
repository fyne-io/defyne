package guidefs

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	// WidgetNames is an array with the list of names of all the Widgets
	WidgetNames []string

	Widgets map[string]WidgetInfo
	once    sync.Once

	importances = []string{"Medium", "High", "Low", "Danger", "Warning", "Success"}
)

type WidgetInfo struct {
	Name     string
	Create   func() fyne.CanvasObject
	Edit     func(fyne.CanvasObject, map[string]string) []*widget.FormItem
	Gostring func(fyne.CanvasObject, map[fyne.CanvasObject]map[string]string, map[string]string) string
	Packages func(object fyne.CanvasObject) []string
}

func initWidgets() {
	Widgets = map[string]WidgetInfo{
		"*widget.Button": {
			Name: "Button",
			Create: func() fyne.CanvasObject {
				return widget.NewButton("Button", func() {})
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				b := obj.(*widget.Button)
				entry := widget.NewEntry()
				entry.SetText(b.Text)
				entry.OnChanged = func(text string) {
					b.SetText(text)
				}

				icon := widget.NewSelect(append([]string{"(No Icon)"}, IconNames...), func(selected string) {
					b.SetIcon(Icons[selected])
				})
				if b.Icon != nil {
					name := IconName(b.Icon)
					for _, n := range IconNames {
						if n == name {
							icon.SetSelected(n)
							break
						}
					}
				}
				importance := widget.NewSelect(importances, func(s string) {
					var i widget.Importance
					for ii, imp := range importances {
						if imp == s {
							i = widget.Importance(ii)
						}
					}
					b.Importance = i
					b.Refresh()
				})
				importance.SetSelectedIndex(int(b.Importance))
				return []*widget.FormItem{
					widget.NewFormItem("Text", entry),
					widget.NewFormItem("Icon", icon),
					widget.NewFormItem("Importance", importance),
				}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				b := obj.(*widget.Button)
				if b.Icon == nil {
					if b.Importance == widget.MediumImportance {
						return widgetRef(props[obj], defs, fmt.Sprintf("widget.NewButton(\"%s\", func() {})", escapeLabel(b.Text)))
					} else {
						return widgetRef(props[obj], defs, fmt.Sprintf("&widget.Button{Text: \"%s\", Importance: %d, OnTapped: func() {}}",
							escapeLabel(b.Text), b.Importance))
					}
				}

				icon := "theme." + IconName(b.Icon) + "()"
				if b.Importance == widget.MediumImportance {
					return widgetRef(props[obj], defs, fmt.Sprintf("widget.NewButtonWithIcon(\"%s\", %s, func() {})", escapeLabel(b.Text), icon))
				} else {
					return widgetRef(props[obj], defs, fmt.Sprintf("&widget.Button{Text: \"%s\", Importance: %d, Icon: %s, OnTapped: func() {}}",
						escapeLabel(b.Text), b.Importance, icon))
				}
			},
			Packages: func(obj fyne.CanvasObject) []string {
				b := obj.(*widget.Button)
				if b.Icon == nil {
					return []string{"widget"}
				}

				return []string{"widget", "theme"}
			},
		},
		"*widget.Hyperlink": {
			Name: "Hyperlink",
			Create: func() fyne.CanvasObject {
				fyneURL, _ := url.Parse("https://fyne.io")
				return widget.NewHyperlink("Link Text", fyneURL)
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				link := obj.(*widget.Hyperlink)
				title := widget.NewEntry()
				title.SetText(link.Text)
				title.OnChanged = func(text string) {
					link.SetText(text)
				}
				subtitle := widget.NewEntry()
				subtitle.SetText(link.URL.String())
				subtitle.OnChanged = func(text string) {
					fyneURL, _ := url.Parse(text)
					link.SetURL(fyneURL)
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
		},
		"*widget.Card": {
			Name: "Card",
			Create: func() fyne.CanvasObject {
				return widget.NewCard("Title", "Subtitle", widget.NewLabel("Content here"))
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				c := obj.(*widget.Card)
				title := widget.NewEntry()
				title.SetText(c.Title)
				title.OnChanged = func(text string) {
					c.SetTitle(text)
				}
				subtitle := widget.NewEntry()
				subtitle.SetText(c.Subtitle)
				subtitle.OnChanged = func(text string) {
					c.SetSubTitle(text)
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
		},
		"*widget.Entry": {
			Name: "Entry",
			Create: func() fyne.CanvasObject {
				e := widget.NewEntry()
				e.SetPlaceHolder("Entry")
				return e
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				l := obj.(*widget.Entry)
				entry1 := widget.NewEntry()
				entry1.SetText(l.Text)
				entry1.OnChanged = func(text string) {
					l.SetText(text)
				}
				entry2 := widget.NewEntry()
				entry2.SetText(l.PlaceHolder)
				entry2.OnChanged = func(text string) {
					l.SetPlaceHolder(text)
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
		},
		"*widget.Icon": {
			Name: "Icon",
			Create: func() fyne.CanvasObject {
				return widget.NewIcon(theme.HelpIcon())
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				i := obj.(*widget.Icon)
				return []*widget.FormItem{
					widget.NewFormItem("Icon", widget.NewSelect(IconNames, func(selected string) {
						i.SetResource(Icons[selected])
					}))}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				i := obj.(*widget.Icon)

				res := "theme." + IconName(i.Resource) + "()"
				return widgetRef(props[obj], defs, fmt.Sprintf("widget.NewIcon(%s)", res))
			},
			Packages: func(obj fyne.CanvasObject) []string {
				return []string{"widget", "theme"}
			},
		},
		"*widget.Label": {
			Name: "Label",
			Create: func() fyne.CanvasObject {
				return widget.NewLabel("Label")
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				l := obj.(*widget.Label)
				entry := widget.NewEntry()
				entry.SetText(l.Text)
				entry.OnChanged = func(text string) {
					l.SetText(text)
				}

				bold := widget.NewCheck("", func(on bool) {
					l.TextStyle.Bold = on
					l.Refresh()
				})
				bold.Checked = l.TextStyle.Bold
				italic := widget.NewCheck("", func(on bool) {
					l.TextStyle.Italic = on
					l.Refresh()
				})
				italic.Checked = l.TextStyle.Italic
				mono := widget.NewCheck("", func(on bool) {
					l.TextStyle.Monospace = on
					l.Refresh()
				})
				mono.Checked = l.TextStyle.Monospace

				return []*widget.FormItem{
					widget.NewFormItem("Text", entry),
					widget.NewFormItem("Bold", bold),
					widget.NewFormItem("Italic", italic),
					widget.NewFormItem("Monospace", mono)}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				l := obj.(*widget.Label)
				if l.TextStyle.Bold || l.TextStyle.Italic || l.TextStyle.Monospace {
					return widgetRef(props[obj], defs,
						fmt.Sprintf("widget.NewLabelWithStyle(\"%s\", fyne.TextAlignLeading, %#v)", escapeLabel(l.Text), l.TextStyle))
				}
				return widgetRef(props[obj], defs,
					fmt.Sprintf("widget.NewLabel(\"%s\")", escapeLabel(l.Text)))
			},
		},
		"*widget.Check": {
			Name: "Check",
			Create: func() fyne.CanvasObject {
				return widget.NewCheck("Tick it or don't", func(b bool) {})
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				c := obj.(*widget.Check)
				title := widget.NewEntry()
				title.SetText(c.Text)
				title.OnChanged = func(text string) {
					c.Text = text
					c.Refresh()
				}
				isChecked := widget.NewCheck("", func(b bool) { c.SetChecked(b) })
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
		},
		"*widget.RadioGroup": {
			Name: "RadioGroup",
			Create: func() fyne.CanvasObject {
				return widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(s string) {})
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				r := obj.(*widget.RadioGroup)
				initialOption := widget.NewRadioGroup(r.Options, func(s string) { r.SetSelected(s) })
				initialOption.SetSelected(r.Selected)
				entry := widget.NewMultiLineEntry()
				entry.SetText(strings.Join(r.Options, "\n"))
				entry.OnChanged = func(text string) {
					r.Options = strings.Split(text, "\n")
					r.Refresh()
					initialOption.Options = strings.Split(text, "\n")
					initialOption.Refresh()
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
		},
		"*widget.Select": {
			Name: "Select",
			Create: func() fyne.CanvasObject {
				return widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {})
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				s := obj.(*widget.Select)
				initialOption := widget.NewSelect(append([]string{"(Select one)"}, s.Options...), func(opt string) {
					s.SetSelected(opt)
					if opt == "(Select one)" {
						s.ClearSelected()
					}
				})
				initialOption.SetSelected(s.Selected)
				entry := widget.NewMultiLineEntry()
				entry.SetText(strings.Join(s.Options, "\n"))
				entry.OnChanged = func(text string) {
					s.Options = strings.Split(text, "\n")
					s.Refresh()
					initialOption.Options = strings.Split(text, "\n")
					initialOption.Refresh()
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
				return widgetRef(props[obj], defs,
					fmt.Sprintf("widget.NewSelect([]string{%s}, func(s string) {})", "\""+strings.Join(opts, "\", \"")+"\""))
			},
		},
		"*widget.Accordion": {
			Name: "Accordion",
			Create: func() fyne.CanvasObject {
				return widget.NewAccordion(widget.NewAccordionItem("Item 1", widget.NewLabel("The content goes here")), widget.NewAccordionItem("Item 2", widget.NewLabel("Content part 2 goes here")))
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				// TODO: Need to add the properties
				// entry := widget.NewEntry()
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs,
					"widget.NewAccordion(\"widget.NewAccordionItem(\"Item 1\", widget.NewLabel(\"The content goes here\")), widget.NewAccordionItem(\"Item 2\", widget.NewLabel(\"Content part 2 goes here\")))")
			},
		},
		"*widget.List": {
			Name: "List",
			Create: func() fyne.CanvasObject {
				myList := []string{"Item 1", "Item 2", "Item 3", "Item 4"}
				// TODO: Need to make the list get adjusted to show the full list of items, currently it has only one item height apprx.
				return widget.NewList(func() int { return len(myList) }, func() fyne.CanvasObject {
					return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
				}, func(id widget.ListItemID, item fyne.CanvasObject) {
					item.(*fyne.Container).Objects[1].(*widget.Label).SetText(myList[id])
				})
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs,
					`widget.NewList(func() int { return len(myList) }, func() fyne.CanvasObject {
				return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
			}, func(id widget.ListItemID, item fyne.CanvasObject) {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(myList[id])
			})`)
			},
			Packages: func(obj fyne.CanvasObject) []string {
				return []string{"widget", "container"}
			},
		},
		"*widget.Menu": {
			Name: "Menu",
			Create: func() fyne.CanvasObject {
				myMenu := fyne.NewMenu("Menu Name", fyne.NewMenuItem("Item 1", func() { fmt.Println("From Item 1") }), fyne.NewMenuItem("Item 2", func() { fmt.Println("From Item 2") }), fyne.NewMenuItem("Item 3", func() { fmt.Println("From Item 3") }))
				return widget.NewMenu(myMenu)
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs,
					"widget.NewMenu(fyne.NewMenu(\"Menu Name\", fyne.NewMenuItem(\"Item 1\", func() {}), fyne.NewMenuItem(\"Item 2\", func() {}), fyne.NewMenuItem(\"Item 3\", func() {})))")
			},
		},
		"*widget.Form": {
			Name: "Form",
			Create: func() fyne.CanvasObject {
				f := widget.NewForm(widget.NewFormItem("Username", widget.NewEntry()), widget.NewFormItem("Password", widget.NewPasswordEntry()))
				f.OnSubmit = func() {}
				f.OnCancel = func() {}
				return f
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs,
					"&widget.Form{Items: []*widget.FormItem{widget.NewFormItem(\"Username\", widget.NewEntry()), widget.NewFormItem(\"Password\", widget.NewPasswordEntry())}, OnSubmit: func() {}, OnCancel: func() {}}")
			},
		},
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
		"*widget.ProgressBar": {
			Name: "Progress Bar",
			Create: func() fyne.CanvasObject {
				p := widget.NewProgressBar()
				p.SetValue(0.1)
				return p
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				p := obj.(*widget.ProgressBar)
				value := widget.NewEntry()
				value.SetText(fmt.Sprintf("%f", p.Value))
				value.OnChanged = func(s string) {
					if f, err := strconv.ParseFloat(s, 64); err == nil {
						p.SetValue(f)
					}
				}
				return []*widget.FormItem{
					widget.NewFormItem("Value", value)}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				p := obj.(*widget.ProgressBar)
				return widgetRef(props[obj], defs,
					fmt.Sprintf("&widget.ProgressBar{Value: %f}", p.Value))
			},
		},
		"*widget.Separator": {
			// Separator's height(or width as you may call) and color come from the theme, so not sure if we can change the color and height here
			Name: "Separator",
			Create: func() fyne.CanvasObject {
				return widget.NewSeparator()
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs, "widget.NewSeparator()")
			},
		},
		"*widget.Slider": {
			Name: "Slider",
			Create: func() fyne.CanvasObject {
				s := widget.NewSlider(0, 100)
				s.OnChanged = func(f float64) {
					fmt.Println("Slider changed to", f)
				}
				return s
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				slider := obj.(*widget.Slider)
				val := widget.NewEntry()
				val.SetText(fmt.Sprintf("%f", slider.Value))
				val.OnChanged = func(s string) {
					if f, err := strconv.ParseFloat(s, 64); err == nil {
						slider.SetValue(f)
					}
				}
				return []*widget.FormItem{
					widget.NewFormItem("Value", val)}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				slider := obj.(*widget.Slider)
				return widgetRef(props[obj], defs, fmt.Sprintf("widget.NewSlider(Min:0, Max:100, Value:%f)", slider.Value))
			},
		},
		"*widget.Table": {
			Name: "Table",
			Create: func() fyne.CanvasObject {
				return widget.NewTable(func() (int, int) { return 3, 3 }, func() fyne.CanvasObject {
					return widget.NewLabel("Cell 000, 000")
				}, func(id widget.TableCellID, cell fyne.CanvasObject) {
					label := cell.(*widget.Label)
					switch id.Col {
					case 0:
						label.SetText(fmt.Sprintf("%d", id.Row+1))
					case 1:
						label.SetText("A longer cell")
					default:
						label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
					}
				})
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs,
					`widget.NewTable(func() (int, int) { return 3, 3 }, func() fyne.CanvasObject {
				return widget.NewLabel("Cell 000, 000")
			}, func(id widget.TableCellID, cell fyne.CanvasObject) {
				label := cell.(*widget.Label)
				switch id.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", id.Row+1))
				case 1:
					label.SetText("A longer cell")
				default:
					label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
				}
			})`)
			},
		},
		"*widget.TextGrid": {
			Name: "Text Grid",
			Create: func() fyne.CanvasObject {
				to := widget.NewTextGrid()
				to.SetText("ABCD \nEFGH")
				return to
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				to := obj.(*widget.TextGrid)
				entry := widget.NewEntry()
				entry.SetText(to.Text())
				entry.OnChanged = func(s string) {
					to.SetText(s)
				}
				return []*widget.FormItem{
					widget.NewFormItem("Text", entry)}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				to := obj.(*widget.TextGrid)
				return widgetRef(props[obj], defs,
					fmt.Sprintf("widget.NewTextGrid(\"%s\")", escapeLabel(to.Text())))
			},
		},
		"*widget.Toolbar": {
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
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
				return widgetRef(props[obj], defs, `widget.NewToolbar(
				widget.NewToolbarAction(theme.FileIcon(), func() {}),
				widget.NewToolbarSeparator(),
				widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {}),
				widget.NewToolbarAction(theme.NavigateBackIcon(), func() {}),
				widget.NewToolbarAction(theme.NavigateNextIcon(), func() {}),
				widget.NewToolbarSpacer(),
				widget.NewToolbarAction(theme.HelpIcon(), func() {}),
			)`)
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
			Edit: func(co fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			//GoString: // TODO
		},

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
			Create: func() fyne.CanvasObject {
				s := container.NewScroll(nil)
				s.Content = container.NewStack()
				return s
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
		},
	}

	Widgets["*widget.Scroll"] = Widgets["*container.Scroll"]
	WidgetNames = extractWidgetNames()
}

// extractWidgetNames returns all the list of names of all the Widgets from our data
func extractWidgetNames() []string {
	var widgetNamesFromData = make([]string, len(Widgets))
	i := 0
	for k := range Widgets {
		widgetNamesFromData[i] = k
		i++
	}

	sort.Strings(widgetNamesFromData)
	return widgetNamesFromData
}

func widgetRef(props map[string]string, defs map[string]string, code string) string {
	if name, ok := props["name"]; ok && name != "" {
		defs[name] = code
		return "g." + name
	}

	return code
}

func InitOnce() {
	once.Do(func() {
		initIcons()
		initWidgets()
	})
}
