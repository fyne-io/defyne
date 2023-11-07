package gui

import (
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/defyne/internal/guidefs"
)

func CreateNew(name string) fyne.CanvasObject {
	guidefs.InitOnce()

	if match, ok := guidefs.Widgets[name]; ok {
		return match.Create()
	}

	return nil
}

func EditorFor(o fyne.CanvasObject, props map[string]string) []*widget.FormItem {
	guidefs.InitOnce()

	_, clazz := getTypeOf(o)

	if match, ok := guidefs.Widgets[clazz]; ok {
		return match.Edit(o, props)
	}

	return nil
}
func GoStringFor(o fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string) string {
	guidefs.InitOnce()

	name := reflect.TypeOf(o).String()

	if match, ok := guidefs.Widgets[name]; ok {
		return match.Gostring(o, props)
	}

	return ""
}

func getTypeOf(o fyne.CanvasObject) (string, string) {
	typeName := reflect.TypeOf(o).Elem().Name()
	class := reflect.TypeOf(o).String()
	l := reflect.ValueOf(o).Elem()
	if typeName == "Entry" {
		if l.FieldByName("Password").Bool() {
			typeName = "PasswordEntry"
		} else if l.FieldByName("MultiLine").Bool() {
			typeName = "MultiLineEntry"
		}
		class = "*widget." + typeName
	}

	return typeName, class
}
