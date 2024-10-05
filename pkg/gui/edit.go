package gui

import (
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/defyne/internal/guidefs"
)

func CreateNew(name string) fyne.CanvasObject {
	guidefs.InitOnce()

	if match := guidefs.Lookup(name); match != nil {
		return match.Create()
	}

	return nil
}

func EditorFor(o fyne.CanvasObject, props map[string]string) []*widget.FormItem {
	guidefs.InitOnce()

	_, clazz := getTypeOf(o)

	if match := guidefs.Lookup(clazz); match != nil {
		return match.Edit(o, props)
	}

	return nil
}

func GoStringFor(o fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string, defs map[string]string) string {
	guidefs.InitOnce()

	name := reflect.TypeOf(o).String()

	if match := guidefs.Lookup(name); match != nil {
		return match.Gostring(o, props, defs)
	}

	return ""
}

func getTypeOf(o fyne.CanvasObject) (string, string) {
	class := reflect.TypeOf(o).String()
	name := NameOf(o)

	return name, class
}

func NameOf(o fyne.CanvasObject) string {
	typeName := reflect.TypeOf(o).Elem().Name()
	l := reflect.ValueOf(o).Elem()
	if typeName == "Entry" {
		if l.FieldByName("Password").Bool() {
			typeName = "PasswordEntry"
		} else if l.FieldByName("MultiLine").Bool() {
			typeName = "MultiLineEntry"
		}
	}

	return typeName
}
