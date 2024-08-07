package gui

import (
	"reflect"

	"fyne.io/fyne/v2"
	"github.com/fyne-io/defyne/internal/guidefs"
)

// WidgetClassList returns the list of supported widget & container classes.
// These can be used for passing to `CreateNew` or `EditorFor`.
func WidgetClassList() []string {
	return guidefs.WidgetNames
}

func DropZonesForObject(o fyne.CanvasObject) []fyne.CanvasObject {
	class := reflect.TypeOf(o).String()
	info := guidefs.Widgets[class]

	if !info.IsContainer() {
		return nil
	}

	return info.Children(o)
}
