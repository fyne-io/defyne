package guibuilder

import (
	"encoding/json"
	"io"
	"log"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
)

type canvObj struct {
	Type   string
	Name   string
	Struct fyne.CanvasObject `json:",omitempty"`
}

type cont struct {
	canvObj
	Layout     string `json:",omitempty"`
	Name       string
	Objects    []interface{}
	Properties map[string]string `json:",omitempty"`
}

func encodeObj(obj fyne.CanvasObject) interface{} {
	c, w := unwrap(obj)
	if w != nil {
		return encodeWidget(w.child, w.name)
	} else if c != nil {
		var node cont
		node.Type = "*fyne.Container"
		node.Layout = strings.Split(reflect.TypeOf(c.c.Layout).String(), ".")[1]
		node.Layout = strings.ToTitle(node.Layout[0:1]) + node.Layout[1:]
		node.Name = c.name
		p := strings.Index(node.Layout, "Layout")
		if p > 0 {
			node.Layout = node.Layout[:p]
		}
		if node.Layout == "Box" {
			props := layoutProps[c.c]
			if props["dir"] == "horizontal" {
				node.Layout = "HBox"
			} else {
				node.Layout = "VBox"
			}
		}
		for _, o := range c.c.Objects {
			node.Objects = append(node.Objects, encodeObj(o)) // what are these? TODO
		}
		node.Properties = layoutProps[c.c]
		return &node
	}

	return nil
}

func encodeWidget(obj fyne.CanvasObject, name string) interface{} {
	return &canvObj{Type: reflect.TypeOf(obj).String(), Name: name, Struct: obj}
}

// DecodeJSON returns a tree of `CanvasObject` elements from the provided JSON `Reader` and
// the tree of wrapped elements that describe their metadata.
func DecodeJSON(r io.Reader) (fyne.CanvasObject, fyne.CanvasObject) {
	var data interface{}
	_ = json.NewDecoder(r).Decode(&data)
	if data == nil {
		return nil, nil
	}

	return decodeMap(data.(map[string]interface{}), nil)
}

func decodeTextStyle(m map[string]interface{}) (s fyne.TextStyle) {
	if m["Bold"] == true {
		s.Bold = true
	}
	if m["Italic"] == true {
		s.Italic = true
	}
	if m["Monospace"] == true {
		s.Monospace = true
	}

	if m["TabWidth"] != 0 {
		s.TabWidth = int(m["TabWidth"].(float64))
	}
	return
}

func decodeMap(m map[string]interface{}, p *fyne.Container) (fyne.CanvasObject, fyne.CanvasObject) {
	if m["Type"] == "*fyne.Container" {
		obj := &fyne.Container{}
		name := m["Layout"].(string)

		props := map[string]string{"layout": name}
		if m["Properties"] != nil {
			for k, v := range m["Properties"].(map[string]interface{}) {
				props[k] = v.(string)
			}
		}
		layoutProps[obj] = props
		if name == "HBox" {
			layoutProps[obj]["dir"] = "horizontal"
		} else if name == "VBox" {
			layoutProps[obj]["dir"] = "vertical"
		}

		wrap := wrapContent(obj, p).(*overlayContainer)
		if m["Objects"] != nil {
			for _, data := range m["Objects"].([]interface{}) {
				child, childWrap := decodeMap(data.(map[string]interface{}), wrap.c)
				obj.Objects = append(obj.Objects, child)
				wrap.c.Objects = append(wrap.c.Objects, childWrap)
			}
		}
		obj.Layout = layouts[name].create(wrap.c, layoutProps[obj])
		wrap.c.Layout = obj.Layout
		if name, ok := m["Name"]; ok {
			wrap.name = name.(string)
		}
		return obj, wrap
	}

	obj := widgets[m["Type"].(string)].create().(fyne.Widget)
	e := reflect.ValueOf(obj).Elem()
	for k, v := range m["Struct"].(map[string]interface{}) {
		f := e.FieldByName(k)

		if f.Type().String() == "fyne.TextAlign" || f.Type().String() == "fyne.TextWrap" ||
			f.Type().String() == "widget.ButtonAlign" || f.Type().String() == "widget.ButtonImportance" || f.Type().String() == "widget.ButtonIconPlacement" {
			f.SetInt(int64(reflect.ValueOf(v).Float()))
		} else if f.Type().String() == "fyne.TextStyle" {
			f.Set(reflect.ValueOf(decodeTextStyle(reflect.ValueOf(v).Interface().(map[string]interface{}))))
		} else if f.Type().String() == "fyne.Resource" {
			res := icons[reflect.ValueOf(v).String()]
			if res != nil {
				f.Set(reflect.ValueOf(wrapResource(res)))
			}
		} else if f.Type().String() == "fyne.CanvasObject" {
			log.Println("Unsupported field")
		} else {
			if strings.Index(f.Type().String(), "int") == 0 {
				f.SetInt(int64(reflect.ValueOf(v).Float()))
			} else {
				f.Set(reflect.ValueOf(v))
			}
		}
	}

	obj.Refresh()
	w := wrapWidget(obj, p)
	if name, ok := m["Name"]; ok {
		w.(*fyne.Container).Objects[1].(*overlayWidget).name = name.(string)
	}
	return obj, w
}

// EncodeJSON writes a JSON stream for the tree of `CanvasObject` elements provided.
// If an error occurs it will be returned, otherwise nil.
func EncodeJSON(obj fyne.CanvasObject, w io.Writer) error {
	tree := encodeObj(obj)

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(tree)
}
