package guidefs

import (
	"bytes"
	"fmt"
	"reflect"

	"fyne.io/fyne/v2"
)

// GoString generates Go code for the given type and object
func GoString(clazz string, obj fyne.CanvasObject, c DefyneContext, defs map[string]string) string {
	info := Lookup(clazz)
	if info == nil {
		return ""
	}

	if fn := info.Gostring; fn != nil {
		return fn(obj, c, defs)
	}

	buf := bytes.Buffer{}
	fallbackPrint(reflect.ValueOf(obj), &buf)
	return widgetRef(c.Metadata()[obj], defs, buf.String())
}

// fallbackPrint is derived from printValue in the BSD licensed Go source code at: src/fmt/print.go.
// We use it here as a fallback Go printer that handles only exported fields.
func fallbackPrint(value reflect.Value, buf *bytes.Buffer) {
	switch value.Kind() {
	case reflect.Struct:
		t := value.Type()
		buf.WriteString(t.String())
		buf.WriteByte('{')
		j := 0

		visible := reflect.VisibleFields(t)
		hideBase := false
		if len(visible) > 0 && visible[0].Name == "baseObject" {
			hideBase = true
		}
		for i := 0; i < len(visible); i++ {
			vv := visible[i]
			f2 := value.FieldByIndex(vv.Index)
			if !visible[i].IsExported() || (hideBase && len(vv.Index) > 1 && vv.Index[0] == 0) {
				continue
			}

			if j > 0 {
				buf.WriteString(", ")
			}
			j++

			buf.WriteString(visible[i].Name)
			buf.WriteByte(':')
			fallbackPrint(f2, buf)
		}
		buf.WriteByte('}')
	case reflect.Interface:
		vv := value.Elem()
		if !vv.IsValid() {
			buf.WriteString("nil")
		} else {
			fallbackPrint(vv, buf)
		}
	case reflect.Pointer:
		switch a := value.Elem(); a.Kind() {
		case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
			buf.WriteByte('&')
			fallbackPrint(a, buf)
			return
		}
		fallthrough
	default:
		buf.WriteString(fmt.Sprintf("%#v", value))
	}
}
