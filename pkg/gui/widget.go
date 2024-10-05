package gui

import (
	"reflect"

	"fyne.io/fyne/v2"
	"github.com/fyne-io/defyne/internal/guidefs"
)

// CollectionClassList returns the list of supported collection widget classes.
// These can be used for passing to `CreateNew` or `EditorFor`.
func CollectionClassList() []string {
	return guidefs.CollectionNames
}

// ContainerClassList returns the list of supported container classes.
// These can be used for passing to `CreateNew` or `EditorFor`.
func ContainerClassList() []string {
	return guidefs.ContainerNames
}

// WidgetClassList returns the list of supported widget classes.
// These can be used for passing to `CreateNew` or `EditorFor`.
func WidgetClassList() []string {
	return guidefs.WidgetNames
}

func DropZonesForObject(o fyne.CanvasObject) []fyne.CanvasObject {
	class := reflect.TypeOf(o).String()
	info := guidefs.Lookup(class)

	if !info.IsContainer() {
		return nil
	}

	return info.Children(o)
}
