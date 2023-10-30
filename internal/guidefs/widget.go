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
)

type WidgetInfo struct {
	Name     string
	Create   func() fyne.CanvasObject
	Edit     func(fyne.CanvasObject, map[string]string) []*widget.FormItem
	Gostring func(fyne.CanvasObject, map[fyne.CanvasObject]map[string]string) string
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
				return []*widget.FormItem{
					widget.NewFormItem("Text", entry),
					widget.NewFormItem("Icon", widget.NewSelect(IconNames, func(selected string) {
						b.SetIcon(WrapResource(Icons[selected]))
					}))}
			},
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				b := obj.(*widget.Button)
				if b.Icon == nil {
					return fmt.Sprintf("widget.NewButton(\"%s\", func() {})", encodeDoubleQuote(b.Text))
				}

				icon := "theme." + IconReverse[fmt.Sprintf("%p", b.Icon.(*jsonResource).Resource)] + "()"
				return fmt.Sprintf("widget.NewButtonWithIcon(\"%s\", %s, func() {})", encodeDoubleQuote(b.Text), icon)
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				link := obj.(*widget.Hyperlink)
				return fmt.Sprintf(`widget.NewHyperlink("%s", %#v)`, encodeDoubleQuote(link.Text), link.URL)
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				c := obj.(*widget.Card)
				return fmt.Sprintf("widget.NewCard(\"%s\", \"%s\", widget.NewLabel(\"Content here\"))",
					encodeDoubleQuote(c.Title), encodeDoubleQuote(c.Subtitle))
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				l := obj.(*widget.Entry)
				return fmt.Sprintf("&widget.Entry{Text: \"%s\", PlaceHolder: \"%s\"}", encodeDoubleQuote(l.Text), encodeDoubleQuote(l.PlaceHolder))
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
						i.SetResource(WrapResource(Icons[selected]))
					}))}
			},
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				i := obj.(*widget.Icon)

				res := "theme." + IconReverse[fmt.Sprintf("%p", i.Resource.(*jsonResource).Resource)] + "()"
				return fmt.Sprintf("widget.NewIcon(%s)", res)
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
				return []*widget.FormItem{
					widget.NewFormItem("Text", entry)}
			},
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				l := obj.(*widget.Label)
				return fmt.Sprintf("widget.NewLabel(\"%s\")", encodeDoubleQuote(l.Text))
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				c := obj.(*widget.Check)
				return fmt.Sprintf("widget.NewCheck(\"%s\", func(b bool) {})", encodeDoubleQuote(c.Text))
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				r := obj.(*widget.RadioGroup)
				var opts []string
				for _, v := range r.Options {
					opts = append(opts, encodeDoubleQuote(v))
				}
				return fmt.Sprintf("widget.NewRadioGroup([]string{%s}, func(s string) {})", "\""+strings.Join(opts, "\", \"")+"\"")
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				s := obj.(*widget.Select)
				var opts []string
				for _, v := range s.Options {
					opts = append(opts, encodeDoubleQuote(v))
				}
				return fmt.Sprintf("widget.NewSelect([]string{%s}, func(s string) {})", "\""+strings.Join(opts, "\", \"")+"\"")
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				return "widget.NewAccordion(\"widget.NewAccordionItem(\"Item 1\", widget.NewLabel(\"The content goes here\")), widget.NewAccordionItem(\"Item 2\", widget.NewLabel(\"Content part 2 goes here\")))"
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				return `widget.NewList(func() int { return len(myList) }, func() fyne.CanvasObject {
				return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
			}, func(id widget.ListItemID, item fyne.CanvasObject) {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(myList[id])
			})`
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				return "widget.NewMenu(fyne.NewMenu(\"Menu Name\", fyne.NewMenuItem(\"Item 1\", func() {}), fyne.NewMenuItem(\"Item 2\", func() {}), fyne.NewMenuItem(\"Item 3\", func() {})))"
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				return "&widget.Form{Items: []*widget.FormItem{widget.NewFormItem(\"Username\", widget.NewEntry()), widget.NewFormItem(\"Password\", widget.NewPasswordEntry())}, OnSubmit: func() {}, OnCancel: func() {}}"
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
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				mle := obj.(*widget.Entry)
				placeholder := widget.NewMultiLineEntry()
				placeholder.Wrapping = fyne.TextWrapWord
				placeholder.SetText(mle.PlaceHolder)
				placeholder.OnChanged = func(s string) {
					mle.SetPlaceHolder(s)
				}
				value := widget.NewMultiLineEntry()
				value.Wrapping = fyne.TextWrapWord
				value.SetText(mle.Text)
				value.OnChanged = func(s string) {
					mle.SetText(s)
				}
				return []*widget.FormItem{
					widget.NewFormItem("Placeholder", placeholder),
					widget.NewFormItem("Value", value)}
			},
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				mle := obj.(*widget.Entry)
				return fmt.Sprintf("&widget.MultiLineEntry{Text: \"%s\", PlaceHolder: \"%s\"}", encodeDoubleQuote(mle.Text), encodeDoubleQuote(mle.PlaceHolder))
			},
		},
		"*widget.PasswordEntry": {
			Name: "Password Entry",
			Create: func() fyne.CanvasObject {
				e := widget.NewPasswordEntry()
				e.SetPlaceHolder("Password Entry")
				return e
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				l := obj.(*widget.Entry)
				text := widget.NewPasswordEntry()
				text.SetText(l.Text)
				text.OnChanged = func(text string) {
					l.SetText(text)
				}
				placeholder := widget.NewEntry()
				placeholder.SetText(l.PlaceHolder)
				placeholder.OnChanged = func(text string) {
					l.SetPlaceHolder(text)
				}
				// hidePassword := widget.NewCheck("Hide Password", func(b bool) {})
				// hidePassword.SetChecked(l.Hidden)
				// hidePassword.OnChanged = func(b bool) {
				// 	l.Hidden = b
				// }
				return []*widget.FormItem{
					widget.NewFormItem("Text", text),
					// widget.NewFormItem("Hide password", placeholder),
					widget.NewFormItem("PlaceHolder", placeholder)}
			},
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				l := obj.(*widget.Entry)
				return fmt.Sprintf("&widget.MultiLineEntry{Text: \"%s\", PlaceHolder: \"%s\"}", encodeDoubleQuote(l.Text), encodeDoubleQuote(l.PlaceHolder))
			},
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				p := obj.(*widget.ProgressBar)
				return fmt.Sprintf("&widget.ProgressBar{Value: %f}", p.Value)
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				return "widget.NewSeparator()"
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				slider := obj.(*widget.Slider)
				return fmt.Sprintf("widget.NewSlider(Min:0, Max:100, Value:%f)", slider.Value)
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				return `widget.NewTable(func() (int, int) { return 3, 3 }, func() fyne.CanvasObject {
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
			})`
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
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				to := obj.(*widget.TextGrid)
				return fmt.Sprintf("widget.NewTextGrid(\"%s\")", encodeDoubleQuote(to.Text()))
			},
		},
		"*widget.Toolbar": {
			Name: "Toolbar",
			Create: func() fyne.CanvasObject {
				return widget.NewToolbar(
					widget.NewToolbarAction(Icons["FileIcon"], func() { fmt.Println("Clicked on FileIcon") }),
					widget.NewToolbarSeparator(),
					widget.NewToolbarAction(Icons["HomeIcon"], func() { fmt.Println("Clicked on HomeIcon") }),
					widget.NewToolbarSeparator(),
					widget.NewToolbarAction(Icons["DownloadIcon"], func() { fmt.Println("Clicked on DownloadIcon") }),
					widget.NewToolbarSeparator(),
					widget.NewToolbarAction(Icons["ViewRefreshIcon"], func() { fmt.Println("Clicked on ViewRefreshIcon") }),
					widget.NewToolbarAction(Icons["NavigateBackIcon"], func() { fmt.Println("Clicked on NavigateBackIcon") }),
					widget.NewToolbarAction(Icons["NavigateNextIcon"], func() { fmt.Println("Clicked on NavigateNextIcon") }),
					widget.NewToolbarAction(Icons["MailSendIcon"], func() { fmt.Println("Clicked on MailSendIcon") }),
					widget.NewToolbarSpacer(),
					widget.NewToolbarAction(Icons["HelpIcon"], func() { fmt.Println("Clicked on HelpIcon") }),
				)
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			Gostring: func(obj fyne.CanvasObject, _ map[fyne.CanvasObject]map[string]string) string {
				return `widget.NewToolbar(
				widget.NewToolbarAction(theme.FileIcon(), func() {}),
				widget.NewToolbarSeparator(),
				widget.NewToolbarAction(theme.HomeIcon(), func() {}),
				widget.NewToolbarSeparator(),
				widget.NewToolbarAction(theme.DownloadIcon(), func() {}),
				widget.NewToolbarSeparator(),
				widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {}),
				widget.NewToolbarAction(theme.NavigateBackIcon(), func() {}),
				widget.NewToolbarAction(theme.NavigateNextIcon(), func() {}),
				widget.NewToolbarAction(theme.MailSendIcon(), func() {}),
				widget.NewToolbarSpacer(),
				widget.NewToolbarAction(theme.HelpIcon(), func() {}),
			)`
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
			Gostring: func(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string) string {
				c := obj.(*fyne.Container)
				l := props[c]["layout"]
				lay := Layouts[l]
				if lay.goText != nil {
					return lay.goText(c, props)
				}

				str := &strings.Builder{}
				str.WriteString(fmt.Sprintf("container.New%s(", l))
				writeGoString(str, nil, props, c.Objects...)
				str.WriteString(")")
				return str.String()
			},
		},
	}

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

func InitOnce() {
	once.Do(func() {
		initIcons()
		initWidgets()
	})
}
