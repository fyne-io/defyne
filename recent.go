package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

const preferenceKeyRecent = "recentlist"

func addRecent(u fyne.URI, p fyne.Preferences) {
	old := loadRecents(p)
	if len(old) > 0 {
		if old[0].String() == u.String() {
			return
		}

		for i, u2 := range old {
			if u2.String() == u.String() {
				if i < len(old)-1 {
					old = old[:i]
				} else {
					old = append(old[:i], old[i+1:]...)
				}
			}
		}
	}

	all := append([]fyne.URI{u}, old...)
	str := ""
	for _, u := range all {
		str += "|" + u.String()
	}
	p.SetString(preferenceKeyRecent, str[1:])
}

func loadRecents(p fyne.Preferences) []fyne.URI {
	var ret []fyne.URI
	val := p.String(preferenceKeyRecent)
	for _, s := range strings.Split(val, "|") {
		u, err := storage.ParseURI(s)
		if err == nil {
			ret = append(ret, u)
		}
	}
	return ret
}

func makeRecentList(p fyne.Preferences, f func(fyne.URI)) *widget.List {
	recents := loadRecents(p)
	r := widget.NewList(
		func() int {
			return len(recents)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			name := recents[id].Name()
			o.(*widget.Label).SetText(name)
		})
	r.OnSelected = func(id widget.ListItemID) {
		f(recents[id])
	}
	return r
}
