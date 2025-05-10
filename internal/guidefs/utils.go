package guidefs

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func escapeLabel(inStr string) (outStr string) {
	outStr = strings.ReplaceAll(inStr, "\"", "\\\"")
	outStr = strings.ReplaceAll(outStr, "\n", "\\n")
	return
}

func newIconSelectorButton(ic fyne.Resource, fn func(fyne.Resource), showName bool) (iconSel *widget.Button) {
	items := make([]*fyne.MenuItem, len(IconNames)+1)

	items[0] = &fyne.MenuItem{
		Label: noIconLabel,
		Icon:  nil,
		Action: func() {
			iconSel.SetText(noIconLabel)
			iconSel.SetIcon(nil)
			fn(nil)
		},
	}
	for i, n := range IconNames {
		name := n
		items[i+1] = &fyne.MenuItem{
			Label: n,
			Icon:  Icons[n],
			Action: func() {
				if showName {
					iconSel.SetText(name)
				} else {
					iconSel.SetText("")
				}
				iconSel.SetIcon(Icons[name])
				fn(Icons[name])
			},
		}
	}
	iconSel = widget.NewButton(noIconLabel, func() {
		d := fyne.CurrentApp().Driver()
		c := d.CanvasForObject(iconSel)
		p := d.AbsolutePositionForObject(iconSel).AddXY(0, iconSel.Size().Height)
		widget.NewPopUpMenu(fyne.NewMenu("", items...), c).ShowAtPosition(p)
	})
	if ic != nil {
		name := IconName(ic)
		for _, n := range IconNames {
			if n == name {
				if showName {
					iconSel.SetText(n)
				} else {
					iconSel.SetText("")
				}
				iconSel.SetIcon(Icons[n])
				break
			}
		}
	}

	return iconSel
}

func getFormIndex(obj fyne.CanvasObject, list []*widget.FormItem) int {
	for i, item := range list {
		if item.Widget == obj {
			return i
		}
	}

	return 0
}

func removeFormItem(i int, l []*widget.FormItem) []*widget.FormItem {
	copy(l[i:], l[i+1:])
	l[len(l)-1] = nil
	return l[:len(l)-1]
}

func removeToolbarItem(i int, l []widget.ToolbarItem) []widget.ToolbarItem {
	copy(l[i:], l[i+1:])
	l[len(l)-1] = nil
	return l[:len(l)-1]
}
