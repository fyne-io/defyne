package gui

import "github.com/fyne-io/defyne/internal/guidefs"

// WidgetClassList returns the list of supported widget & container classes.
// These can be used for passing to `CreateNew` or `EditorFor`.
func WidgetClassList() []string {
	return guidefs.WidgetNames
}
