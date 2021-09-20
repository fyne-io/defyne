package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type widgetInfo struct {
	name     string
	create   func() fyne.CanvasObject
	edit     func(fyne.CanvasObject) []*widget.FormItem
	gostring func(fyne.CanvasObject) string
	packages func(object fyne.CanvasObject) []string
}

var widgets map[string]widgetInfo

func initWidgets() {
	widgets = map[string]widgetInfo{
		"*widget.Button": {
			name: "Button",
			create: func() fyne.CanvasObject {
				return widget.NewButton("Button", func() {})
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				b := obj.(*widget.Button)
				entry := widget.NewEntry()
				entry.SetText(b.Text)
				entry.OnChanged = func(text string) {
					b.SetText(text)
				}
				return []*widget.FormItem{
					widget.NewFormItem("Text", entry),
					widget.NewFormItem("Icon", widget.NewSelect(iconNames, func(selected string) {
						b.SetIcon(icons[selected])
					}))}
			},
			gostring: func(obj fyne.CanvasObject) string {
				b := obj.(*widget.Button)
				if b.Icon == nil {
					return fmt.Sprintf("widget.NewButton(\"%s\", func() {})", encodeDoubleQuote(b.Text))
				}

				icon := "theme." + iconReverse[fmt.Sprintf("%p", b.Icon)] + "()"
				return fmt.Sprintf("widget.NewButtonWithIcon(\"%s\", %s, func() {})", encodeDoubleQuote(b.Text), icon)
			},
			packages: func(obj fyne.CanvasObject) []string {
				b := obj.(*widget.Button)
				if b.Icon == nil {
					return []string{"widget"}
				}

				return []string{"widget", "theme"}
			},
		},
		"*widget.Hyperlink": {
			name: "Hyperlink",
			create: func() fyne.CanvasObject {
				fyneURL, _ := url.Parse("https://fyne.io")
				return widget.NewHyperlink("Link Text", fyneURL)
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				link := obj.(*widget.Hyperlink)
				return fmt.Sprintf("widget.NewHyperLink(\"%s\", \"%s\")", encodeDoubleQuote(link.Text), encodeDoubleQuote(link.URL.String()))
			},
		},
		"*widget.Card": {
			name: "Card",
			create: func() fyne.CanvasObject {
				return widget.NewCard("Title", "Subtitle", widget.NewLabel("Content here"))
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				c := obj.(*widget.Card)
				return fmt.Sprintf("widget.NewCard(\"%s\", \"%s\", widget.NewLabel(\"Content here\")",
					encodeDoubleQuote(c.Title), encodeDoubleQuote(c.Subtitle))
			},
		},
		"*widget.Entry": {
			name: "Entry",
			create: func() fyne.CanvasObject {
				e := widget.NewEntry()
				e.SetPlaceHolder("Entry")
				return e
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				l := obj.(*widget.Entry)
				return fmt.Sprintf("&widget.Entry{Text: \"%s\", PlaceHolder: \"%s\"}", encodeDoubleQuote(l.Text), encodeDoubleQuote(l.PlaceHolder))
			},
		},
		"*widget.Icon": {
			name: "Icon",
			create: func() fyne.CanvasObject {
				return widget.NewIcon(theme.HelpIcon())
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				i := obj.(*widget.Icon)
				return []*widget.FormItem{
					widget.NewFormItem("Icon", widget.NewSelect(iconNames, func(selected string) {
						i.SetResource(icons[selected])
					}))}
			},
			gostring: func(obj fyne.CanvasObject) string {
				i := obj.(*widget.Icon)

				res := "theme." + iconReverse[fmt.Sprintf("%p", i.Resource)] + "()"
				return fmt.Sprintf("widget.NewIcon(%s)", res)
			},
			packages: func(obj fyne.CanvasObject) []string {
				return []string{"widget", "theme"}
			},
		},
		"*widget.Label": {
			name: "Label",
			create: func() fyne.CanvasObject {
				return widget.NewLabel("Label")
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				l := obj.(*widget.Label)
				entry := widget.NewEntry()
				entry.SetText(l.Text)
				entry.OnChanged = func(text string) {
					l.SetText(text)
				}
				return []*widget.FormItem{
					widget.NewFormItem("Text", entry)}
			},
			gostring: func(obj fyne.CanvasObject) string {
				l := obj.(*widget.Label)
				return fmt.Sprintf("widget.NewLabel(\"%s\")", encodeDoubleQuote(l.Text))
			},
		},
		"*widget.Check": {
			name: "CheckBox",
			create: func() fyne.CanvasObject {
				return widget.NewCheck("Tick it or don't", func(b bool) {})
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				c := obj.(*widget.Check)
				return fmt.Sprintf("widget.NewCheck(\"%s\", func(b bool) {}", encodeDoubleQuote(c.Text))
			},
		},
		"*widget.RadioGroup": {
			name: "RadioGroup",
			create: func() fyne.CanvasObject {
				return widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(s string) {})
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				r := obj.(*widget.RadioGroup)
				var opts []string
				for _, v := range r.Options {
					opts = append(opts, encodeDoubleQuote(v))
				}
				return fmt.Sprintf("widget.NewRadioGroup([]string{%s}, func(s string) {})", "\""+strings.Join(opts, "\", \"")+"\"")
			},
		},
		"*widget.Select": {
			name: "Select",
			create: func() fyne.CanvasObject {
				return widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {})
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				s := obj.(*widget.Select)
				var opts []string
				for _, v := range s.Options {
					opts = append(opts, encodeDoubleQuote(v))
				}
				return fmt.Sprintf("widget.NewSelect([]string{%s}, func(s string) {})", "\""+strings.Join(opts, "\", \"")+"\"")
			},
		},
		"*widget.Accordion": {
			name: "Accordion",
			create: func() fyne.CanvasObject {
				return widget.NewAccordion(widget.NewAccordionItem("Item 1", widget.NewLabel("The content goes here")), widget.NewAccordionItem("Item 2", widget.NewLabel("Content part 2 goes here")))
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				// TODO: Need to add the properties
				// entry := widget.NewEntry()
				return []*widget.FormItem{}
			},
			gostring: func(obj fyne.CanvasObject) string {
				return "widget.NewAccordion(\"widget.NewAccordionItem(\"Item 1\", widget.NewLabel(\"The content goes here\")), widget.NewAccordionItem(\"Item 2\", widget.NewLabel(\"Content part 2 goes here\")))"
			},
		},
		"*widget.List": {
			name: "List",
			create: func() fyne.CanvasObject {
				myList := []string{"Item 1", "Item 2", "Item 3", "Item 4"}
				// TODO: Need to make the list get adjusted to show the full list of items, currently it has only one item height apprx.
				return widget.NewList(func() int { return len(myList) }, func() fyne.CanvasObject {
					return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
				}, func(id widget.ListItemID, item fyne.CanvasObject) {
					item.(*fyne.Container).Objects[1].(*widget.Label).SetText(myList[id])
				})
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			gostring: func(obj fyne.CanvasObject) string {
				return `widget.NewList(func() int { return len(myList) }, func() fyne.CanvasObject {
				return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
			}, func(id widget.ListItemID, item fyne.CanvasObject) {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(myList[id])
			})`
			},
			packages: func(obj fyne.CanvasObject) []string {
				return []string{"widget", "container"}
			},
		},
		"*widget.Menu": {
			name: "Menu",
			create: func() fyne.CanvasObject {
				myMenu := fyne.NewMenu("Menu Name", fyne.NewMenuItem("Item 1", func() { fmt.Println("From Item 1") }), fyne.NewMenuItem("Item 2", func() { fmt.Println("From Item 2") }), fyne.NewMenuItem("Item 3", func() { fmt.Println("From Item 3") }))
				return widget.NewMenu(myMenu)
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			gostring: func(obj fyne.CanvasObject) string {
				return "widget.NewMenu(fyne.NewMenu(\"Menu Name\", fyne.NewMenuItem(\"Item 1\", func() {}), fyne.NewMenuItem(\"Item 2\", func() {}), fyne.NewMenuItem(\"Item 3\", func() {})))"
			},
		},
		"*widget.Form": {
			name: "Form",
			create: func() fyne.CanvasObject {
				return widget.NewForm(widget.NewFormItem("Username", widget.NewEntry()), widget.NewFormItem("Password", widget.NewPasswordEntry()), widget.NewFormItem("", container.NewGridWithColumns(2, widget.NewButton("Submit", func() { fmt.Println("Form is submitted") }), widget.NewButton("Cancel", func() { fmt.Println("Form is Cancelled") }))))
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			gostring: func(obj fyne.CanvasObject) string {
				return "widget.NewForm(widget.NewFormItem(\"Username\", widget.NewEntry()), widget.NewFormItem(\"Password\", widget.NewPasswordEntry()), widget.NewFormItem(\"\", container.NewGridWithColumns(2, widget.NewButton(\"Submit\", func() {}), widget.NewButton(\"Cancel\", func() {}))))"
			},
		},
		"*widget.MultiLineEntry": {
			name: "Multi Line Entry",
			create: func() fyne.CanvasObject {
				mle := widget.NewMultiLineEntry()
				mle.SetPlaceHolder("Enter Some \nLong text \nHere")
				mle.Wrapping = fyne.TextWrapWord
				return mle
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				mle := obj.(*widget.Entry)
				return fmt.Sprintf("&widget.MultiLineEntry{Text: \"%s\", PlaceHolder: \"%s\"}", encodeDoubleQuote(mle.Text), encodeDoubleQuote(mle.PlaceHolder))
			},
		},
		"*widget.PasswordEntry": {
			name: "Password Entry",
			create: func() fyne.CanvasObject {
				e := widget.NewPasswordEntry()
				e.SetPlaceHolder("Password Entry")
				return e
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				l := obj.(*widget.Entry)
				return fmt.Sprintf("&widget.MultiLineEntry{Text: \"%s\", PlaceHolder: \"%s\"}", encodeDoubleQuote(l.Text), encodeDoubleQuote(l.PlaceHolder))
			},
		},
		"*widget.ProgressBar": {
			name: "Progress Bar",
			create: func() fyne.CanvasObject {
				p := widget.NewProgressBar()
				p.SetValue(0.1)
				return p
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				p := obj.(*widget.ProgressBar)
				return fmt.Sprintf("&widget.ProgressBar{Value: %f}", p.Value)
			},
		},
		"*widget.Separator": {
			// Separator's height(or width as you may call) and color come from the theme, so not sure if we can change the color and height here
			name: "Separator",
			create: func() fyne.CanvasObject {
				return widget.NewSeparator()
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			gostring: func(obj fyne.CanvasObject) string {
				return "widget.NewSeparator()"
			},
		},
		"*widget.Slider": {
			name: "Slider",
			create: func() fyne.CanvasObject {
				s := widget.NewSlider(0, 100)
				s.OnChanged = func(f float64) {
					fmt.Println("Slider changed to", f)
				}
				return s
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
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
			gostring: func(obj fyne.CanvasObject) string {
				slider := obj.(*widget.Slider)
				return fmt.Sprintf("widget.NewSlider(Min:0, Max:100, Value:%f)", slider.Value)
			},
		},
		"*widget.Table": {
			name: "Table",
			create: func() fyne.CanvasObject {
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
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			gostring: func(obj fyne.CanvasObject) string {
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
			name: "Text Grid",
			create: func() fyne.CanvasObject {
				to := widget.NewTextGrid()
				to.SetText("ABCD \nEFGH")
				return to
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				to := obj.(*widget.TextGrid)
				entry := widget.NewEntry()
				entry.SetText(to.Text())
				entry.OnChanged = func(s string) {
					to.SetText(s)
				}
				return []*widget.FormItem{
					widget.NewFormItem("Text", entry)}
			},
			gostring: func(obj fyne.CanvasObject) string {
				to := obj.(*widget.TextGrid)
				return fmt.Sprintf("widget.NewTextGrid(\"%s\")", encodeDoubleQuote(to.Text()))
			},
		},
		"*widget.Toolbar": {
			name: "Toolbar",
			create: func() fyne.CanvasObject {
				return widget.NewToolbar(
					widget.NewToolbarAction(icons["FileIcon"], func() { fmt.Println("Clicked on FileIcon") }),
					widget.NewToolbarSeparator(),
					widget.NewToolbarAction(icons["HomeIcon"], func() { fmt.Println("Clicked on HomeIcon") }),
					widget.NewToolbarSeparator(),
					widget.NewToolbarAction(icons["DownloadIcon"], func() { fmt.Println("Clicked on DownloadIcon") }),
					widget.NewToolbarSeparator(),
					widget.NewToolbarAction(icons["ViewRefreshIcon"], func() { fmt.Println("Clicked on ViewRefreshIcon") }),
					widget.NewToolbarAction(icons["NavigateBackIcon"], func() { fmt.Println("Clicked on NavigateBackIcon") }),
					widget.NewToolbarAction(icons["NavigateNextIcon"], func() { fmt.Println("Clicked on NavigateNextIcon") }),
					widget.NewToolbarAction(icons["MailSendIcon"], func() { fmt.Println("Clicked on MailSendIcon") }),
					widget.NewToolbarSpacer(),
					widget.NewToolbarAction(icons["HelpIcon"], func() { fmt.Println("Clicked on HelpIcon") }),
				)
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				return []*widget.FormItem{}
			},
			gostring: func(obj fyne.CanvasObject) string {
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
			name: "Tree",
			create: func() fyne.CanvasObject {
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
			edit: func(co fyne.CanvasObject) []*widget.FormItem {
				return []*widget.FormItem{}
			},
		},

		"*fyne.Container": {
			name: "Container",
			create: func() fyne.CanvasObject {
				return container.NewMax()
			},
			edit: func(obj fyne.CanvasObject) []*widget.FormItem {
				c := obj.(*fyne.Container)
				props := layoutProps[c]

				var items []*widget.FormItem
				var choose *widget.FormItem
				// TODO figure out how to work Border...
				choose = widget.NewFormItem("Layout", widget.NewSelect(layoutNames, func(l string) {
					lay := layouts[l]
					props["layout"] = l
					c.Layout = lay.create(props)
					c.Refresh()
					choose.Widget.Hide()

					edit := lay.edit
					items = []*widget.FormItem{choose}
					if edit != nil {
						items = append(items, edit(c, props)...)
					}

					editForm = widget.NewForm(items...)
					paletteList.Objects = []fyne.CanvasObject{editForm}
					choose.Widget.Show()
					paletteList.Refresh()
				}))
				choose.Widget.(*widget.Select).SetSelected(props["layout"])
				return items
			},
			gostring: func(obj fyne.CanvasObject) string {
				c := obj.(*fyne.Container)
				l := layoutProps[c]["layout"]
				str := strings.Builder{}
				str.WriteString(fmt.Sprintf("container.New%s(", l))
				for i, o := range c.Objects {
					if _, ok := o.(*overlayContainer); !ok {
						o = o.(*fyne.Container).Objects[1]
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
	}

	widgetNames = extractWidgetNames()
}

// widgetNames is an array with the list of names of all the widgets
var widgetNames []string

// extractWidgetNames returns all the list of names of all the widgets from our data
func extractWidgetNames() []string {
	var widgetNamesFromData = make([]string, len(widgets))
	i := 0
	for k := range widgets {
		widgetNamesFromData[i] = k
		i++
	}
	return widgetNamesFromData
}
