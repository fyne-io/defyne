package gui

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/url"
	"reflect"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/defyne/internal/guidefs"
)

const jsonKeyObject = "Object"

type canvObj struct {
	Type   string
	Name   string
	Struct fyne.CanvasObject `json:",omitempty"`
}

type form struct {
	Type   string
	Name   string
	Struct map[string]interface{} `json:",omitempty"`
}

type formItem struct {
	HintText, Text string
	Widget         *canvObj
}

type cont struct {
	canvObj
	Layout     string `json:",omitempty"`
	Name       string
	Objects    []interface{}
	Properties map[string]string `json:",omitempty"`
}

// DecodeJSON returns a tree of `CanvasObject` elements from the provided JSON `Reader` and
// the tree of wrapped elements that describe their metadata.
func DecodeJSON(r io.Reader) (fyne.CanvasObject, map[fyne.CanvasObject]map[string]string, error) {
	guidefs.InitOnce()

	var data interface{}
	err := json.NewDecoder(r).Decode(&data)
	if err != nil || data == nil {
		return nil, nil, err
	}

	meta := make(map[fyne.CanvasObject]map[string]string)
	root := data.(map[string]interface{})
	node, ok := root[jsonKeyObject]
	if !ok {
		return nil, nil, errors.New("cannot parse old format of .gui.json file")
	}
	obj := decodeMap(node.(map[string]interface{}), nil, meta)
	return obj, meta, nil
}

// EncodeJSON writes a JSON stream for the tree of `CanvasObject` elements provided.
// If an error occurs it will be returned, otherwise nil.
func EncodeJSON(obj fyne.CanvasObject, meta map[fyne.CanvasObject]map[string]string, w io.Writer) error {
	guidefs.InitOnce()

	if meta == nil {
		meta = make(map[fyne.CanvasObject]map[string]string)
	}
	tree := encodeObj(obj, meta)

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(map[string]interface{}{jsonKeyObject: tree})
}

func encodeForm(obj *widget.Form, name string) interface{} {
	var items []*formItem
	for _, o := range obj.Items {
		items = append(items,
			&formItem{
				HintText: o.HintText,
				Text:     o.Text,
				Widget:   encodeWidget(o.Widget, ""),
			})
	}

	var node form
	node.Type = "*widget.Form"
	node.Name = name
	node.Struct = map[string]interface{}{
		"Hidden":     obj.Hidden,
		"Items":      items,
		"SubmitText": obj.SubmitText,
		"CancelText": obj.CancelText,
	}

	return &node
}

func encodeObj(obj fyne.CanvasObject, meta map[fyne.CanvasObject]map[string]string) interface{} {
	props := meta[obj]
	name := ""
	if props == nil {
		props = make(map[string]string)
		meta[obj] = props
	} else if props["name"] != "" {
		name = props["name"]
	}

	switch c := obj.(type) {
	case *widget.Button:
		if c.Icon == nil {
			return encodeWidget(c, name)
		}

		ic := c.Icon
		c.Icon = guidefs.WrapResource(c.Icon)
		wid := encodeWidget(c, name)
		go func() { // TODO find a better way to reset this after encoding
			time.Sleep(time.Millisecond * 100)
			c.Icon = ic
		}()
		return wid
	case *widget.Icon:
		if c.Resource == nil {
			return encodeWidget(c, name)
		}

		ic := c.Resource
		c.Resource = guidefs.WrapResource(c.Resource)
		wid := encodeWidget(c, name)
		go func() { // TODO find a better way to reset this after encoding
			time.Sleep(time.Millisecond * 100)
			c.Resource = ic
		}()
		return wid
	case fyne.Widget:
		if form, ok := c.(*widget.Form); ok {
			return encodeForm(form, name)
		}
		return encodeWidget(c, name)
	case *fyne.Container:
		var node cont
		node.Type = "*fyne.Container"
		node.Layout = strings.Split(reflect.TypeOf(c.Layout).String(), ".")[1]
		node.Layout = strings.ToTitle(node.Layout[0:1]) + node.Layout[1:]
		node.Name = name
		p := strings.Index(node.Layout, "Layout")
		if p > 0 {
			node.Layout = node.Layout[:p]
		}
		if node.Layout == "Box" {
			if props["dir"] == "horizontal" {
				node.Layout = "HBox"
			} else {
				node.Layout = "VBox"
			}
		}
		for _, o := range c.Objects {
			node.Objects = append(node.Objects, encodeObj(o, meta))
		}
		node.Properties = meta[c]
		return &node
	}

	return nil
}

func encodeWidget(obj fyne.CanvasObject, name string) *canvObj {
	return &canvObj{Type: reflect.TypeOf(obj).String(), Name: name, Struct: obj}
}

func decodeFormItem(m map[string]interface{}) *widget.FormItem {
	f := &widget.FormItem{}
	if str, ok := m["HintText"]; ok {
		f.HintText = str.(string)
	}
	if str, ok := m["Text"]; ok {
		f.Text = str.(string)
	}
	if wid, ok := m["Widget"]; ok {
		f.Widget = decodeWidget(wid.(map[string]interface{}))
	}
	return f
}

func decodeFromMap(m map[string]interface{}, in interface{}) {
	t := reflect.ValueOf(in).Elem()
	for k, v := range m {
		val := t.FieldByName(k)
		if val.Type().Kind() == reflect.Ptr {
			continue
		}
		val.Set(reflect.ValueOf(v))
	}
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

func decodeMap(m map[string]interface{}, p *fyne.Container, meta map[fyne.CanvasObject]map[string]string) fyne.CanvasObject {
	if m["Type"] == "*fyne.Container" {
		obj := &fyne.Container{}
		name := m["Layout"].(string)

		props := map[string]string{"layout": name}
		if m["Properties"] != nil {
			for k, v := range m["Properties"].(map[string]interface{}) {
				props[k] = v.(string)
			}
		}
		if name == "HBox" {
			props["dir"] = "horizontal"
		} else if name == "VBox" {
			props["dir"] = "vertical"
		}

		if m["Objects"] != nil {
			for _, data := range m["Objects"].([]interface{}) {
				if data == nil {
					// Nil object?
					continue
				}
				child := decodeMap(data.(map[string]interface{}), obj, meta)
				obj.Objects = append(obj.Objects, child)
			}
		}
		obj.Layout = guidefs.Layouts[name].Create(obj, props)
		if name, ok := m["Name"]; ok {
			props["name"] = name.(string)
		}

		meta[obj] = props
		return obj
	}

	obj := decodeWidget(m)
	obj.Refresh()
	props := map[string]string{}
	if name, ok := m["Name"]; ok {
		props["name"] = name.(string)
	}

	meta[obj] = props
	return obj
}

func decodeWidget(m map[string]interface{}) fyne.Widget {
	obj := guidefs.Widgets[m["Type"].(string)].Create().(fyne.Widget)
	e := reflect.ValueOf(obj).Elem()
	for k, v := range m["Struct"].(map[string]interface{}) {
		f := e.FieldByName(k)

		typeName := f.Type().String()
		switch typeName {
		case "fyne.TextAlign", "fyne.TextTruncation", "fyne.TextWrap", "widget.ButtonAlign", "widget.ButtonImportance", "widget.ButtonIconPlacement", "widget.Importance":
			f.SetInt(int64(reflect.ValueOf(v).Float()))
		case "fyne.TextStyle":
			f.Set(reflect.ValueOf(decodeTextStyle(reflect.ValueOf(v).Interface().(map[string]interface{}))))
		case "fyne.Resource":
			res := guidefs.Icons[reflect.ValueOf(v).String()]
			if res != nil {
				f.Set(reflect.ValueOf(res))
			}
		case "[]*widget.FormItem":
			var items []*widget.FormItem
			for _, item := range reflect.ValueOf(v).Interface().([]interface{}) {
				items = append(items, decodeFormItem(item.(map[string]interface{})))
			}
			f.Set(reflect.ValueOf(items))
		case "fyne.CanvasObject":
			log.Println("Unsupported field")
		case "*url.URL":
			u := &url.URL{}
			decodeFromMap(reflect.ValueOf(v).Interface().(map[string]interface{}), u)
			f.Set(reflect.ValueOf(u))
		default:
			if strings.Index(typeName, "int") == 0 {
				f.SetInt(int64(reflect.ValueOf(v).Float()))
			} else if v != nil {
				f.Set(reflect.ValueOf(v))
			}
		}
	}

	return obj
}
