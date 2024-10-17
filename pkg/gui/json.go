package gui

import (
	"encoding/json"
	"errors"
	"image/color"
	"io"
	"net/url"
	"reflect"
	"strings"
	"time"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/fyne-io/defyne/internal/guidefs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const jsonKeyObject = "Object"

type canvObj struct {
	Type    string
	Name    string            `json:",omitempty"`
	Actions map[string]string `json:",omitempty"`
	Struct  fyne.CanvasObject `json:",omitempty"`
}

type cntObj struct {
	canvObj
	Struct map[string]interface{}
}

type form struct {
	Type   string
	Name   string                 `json:",omitempty"`
	Struct map[string]interface{} `json:",omitempty"`
}

type formItem struct {
	HintText, Text string
	Widget         *canvObj
}

type cont struct {
	canvObj
	Layout     string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Objects    []interface{}
	Properties map[string]string `json:",omitempty"`
}

// DecodeObject returns a tree of `CanvasObject` elements from the provided JSON `Reader` and
// updates the metadata map to include any additional information.
func DecodeObject(r io.Reader) (fyne.CanvasObject, map[fyne.CanvasObject]map[string]string, error) {
	guidefs.InitOnce()

	var data interface{}
	err := json.NewDecoder(r).Decode(&data)
	if err != nil || data == nil {
		return nil, nil, err
	}

	meta := make(map[fyne.CanvasObject]map[string]string)
	root := data.(map[string]interface{})

	obj, err := DecodeMap(root, meta)
	return obj, meta, err
}

// DecodeMap returns a tree of `CanvasObject` elements from the provided JSON map and
// updates the metadata map to include any additional information.
func DecodeMap(m map[string]interface{}, meta map[fyne.CanvasObject]map[string]string) (fyne.CanvasObject, error) {
	guidefs.InitOnce()

	switch m["Type"] {
	case "*fyne.Container":
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
				child, _ := DecodeMap(data.(map[string]interface{}), meta)
				obj.Objects = append(obj.Objects, child)
			}
		}
		obj.Layout = guidefs.Layouts[name].Create(obj, props)
		if name, ok := m["Name"]; ok {
			props["name"] = name.(string)
		}

		meta[obj] = props
		return obj, nil
	case "*container.Scroll":
		obj := &container.Scroll{}
		info := m["Struct"].(map[string]interface{})
		if off, ok := info["Direction"]; ok {
			obj.Direction = container.ScrollDirection(off.(float64))
		}
		if info["Content"] != nil {
			child, _ := DecodeMap(info["Content"].(map[string]interface{}), meta)
			obj.Content = child
		}

		props := map[string]string{}
		if name, ok := m["Name"]; ok {
			props["name"] = name.(string)
		}

		meta[obj] = props
		return obj, nil
	case "*container.Split":
		obj := &container.Split{}
		info := m["Struct"].(map[string]interface{})
		if info["Horizontal"].(bool) {
			obj.Horizontal = true
		}
		if off, ok := info["Offset"]; ok {
			obj.Offset = off.(float64)
		}
		if info["Leading"] != nil {
			child, _ := DecodeMap(info["Leading"].(map[string]interface{}), meta)
			obj.Leading = child
		}
		if info["Trailing"] != nil {
			child, _ := DecodeMap(info["Trailing"].(map[string]interface{}), meta)
			obj.Trailing = child
		}

		props := map[string]string{}
		if name, ok := m["Name"]; ok {
			props["name"] = name.(string)
		}

		meta[obj] = props
		return obj, nil
	case "*canvas.Rectangle":
		obj := &canvas.Rectangle{}
		e := reflect.ValueOf(obj).Elem()

		err := decodeFields(e, m["Struct"].(map[string]interface{}))
		return obj, err
	}

	obj := decodeWidget(m)
	obj.Refresh()
	props := map[string]string{}
	if name, ok := m["Name"]; ok {
		props["name"] = name.(string)
	}

	if set, ok := m["Actions"]; ok {
		if actions, ok := set.(map[string]any); ok {
			for k, v := range actions {
				props[k] = v.(string)
			}
		}
	}

	meta[obj] = props
	return obj, nil
}

// EncodeObject writes a JSON stream for the tree of `CanvasObject` elements provided.
// If an error occurs it will be returned, otherwise nil.
func EncodeObject(obj fyne.CanvasObject, meta map[fyne.CanvasObject]map[string]string, w io.Writer) error {
	guidefs.InitOnce()

	if meta == nil {
		meta = make(map[fyne.CanvasObject]map[string]string)
	}
	tree, _ := EncodeMap(obj, meta)

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(tree)
}

// EncodeMap returns a JSON map for the tree of `CanvasObject` elements provided, using additional metadata if required.
// If an error occurs it will be returned, otherwise nil.
func EncodeMap(obj fyne.CanvasObject, meta map[fyne.CanvasObject]map[string]string) (interface{}, error) {
	guidefs.InitOnce()

	props := meta[obj]
	name := ""
	actions := map[string]string{}
	if props == nil {
		props = make(map[string]string)
		meta[obj] = props
	} else {
		name = props["name"]

		for k, v := range props {
			if len(k) > 2 && k[0:2] == "On" {
				actions[k] = v
			}
		}
	}

	switch c := obj.(type) {
	case *widget.Button:
		if c.Icon == nil {
			return encodeWidget(c, name, actions), nil
		}

		ic := c.Icon
		c.Icon = guidefs.WrapResource(c.Icon)
		wid := encodeWidget(c, name, actions)
		go func() { // TODO find a better way to reset this after encoding
			time.Sleep(time.Millisecond * 100)
			c.Icon = ic
		}()
		return wid, nil
	case *widget.Icon:
		if c.Resource == nil {
			return encodeWidget(c, name, actions), nil
		}

		ic := c.Resource
		c.Resource = guidefs.WrapResource(c.Resource)
		wid := encodeWidget(c, name, actions)
		go func() { // TODO find a better way to reset this after encoding
			time.Sleep(time.Millisecond * 100)
			c.Resource = ic
		}()
		return wid, nil
	case *widget.Toolbar:
		for id, i := range c.Items {
			switch t := i.(type) {
			case *widget.ToolbarAction:
				ic := t.Icon
				t.Icon = guidefs.WrapResource(t.Icon)
				go func() { // TODO find a better way to reset this after encoding
					time.Sleep(time.Millisecond * 100)
					t.Icon = ic
				}()
			case *widget.ToolbarSeparator:
				c.Items[id] = toolbarItem{Type: "Separator"}
			case *widget.ToolbarSpacer:
				c.Items[id] = toolbarItem{Type: "Spacer"}
			}
		}

		return encodeWidget(c, name, actions), nil
	case *container.Scroll:
		node := &cntObj{Struct: make(map[string]interface{})}
		node.Type = "*container.Scroll"
		node.Struct["Direction"] = c.Direction
		node.Name = name

		node.Struct["Content"], _ = EncodeMap(c.Content, meta)

		return &node, nil
	case *container.Split:
		node := &cntObj{Struct: make(map[string]interface{})}
		node.Type = "*container.Split"
		node.Struct["Horizontal"] = c.Horizontal
		node.Struct["Offset"] = c.Offset
		node.Name = name

		node.Struct["Leading"], _ = EncodeMap(c.Leading, meta)
		node.Struct["Trailing"], _ = EncodeMap(c.Trailing, meta)

		return &node, nil
	case fyne.Widget:
		if form, ok := c.(*widget.Form); ok {
			return encodeForm(form, name), nil
		}
		return encodeWidget(c, name, actions), nil
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
			enc, _ := EncodeMap(o, meta)
			node.Objects = append(node.Objects, enc)
		}
		node.Properties = meta[c]
		return &node, nil
	}

	return &canvObj{Type: reflect.TypeOf(obj).String(), Name: name, Struct: obj}, nil
}

func encodeForm(obj *widget.Form, name string) interface{} {
	var items []*formItem
	for _, o := range obj.Items {
		items = append(items,
			&formItem{
				HintText: o.HintText,
				Text:     o.Text,
				Widget:   encodeWidget(o.Widget, "", nil),
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

func encodeWidget(obj fyne.CanvasObject, name string, actions map[string]string) *canvObj {
	w := &canvObj{Type: reflect.TypeOf(obj).String(), Name: name, Struct: obj}

	if len(actions) > 0 {
		w.Actions = actions
	}

	return w
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
		switch val.Type().Kind() {
		case reflect.Ptr:
			continue
		case reflect.Uint8:
			val.SetUint(uint64(reflect.ValueOf(v).Float()))
		default:
			val.Set(reflect.ValueOf(v))
		}
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

func decodePosition(m map[string]interface{}) fyne.Position {
	x := m["X"].(float64)
	y := m["Y"].(float64)

	return fyne.NewPos(float32(x), float32(y))
}

func decodeToolbarItem(m map[string]interface{}) widget.ToolbarItem {
	if v, ok := m["Type"]; ok {
		switch v {
		case "Separator":
			return widget.NewToolbarSeparator()
		default:
			return widget.NewToolbarSpacer()
		}
	}

	return widget.NewToolbarAction(guidefs.Icons[m["Icon"].(string)], nil)
}

func decodeRichTextStyle(m map[string]interface{}) (s widget.RichTextStyle) {
	for k, v := range m {
		switch k {
		case "TextStyle":
			s.TextStyle = decodeTextStyle(v.(map[string]interface{}))
			// TODO more!
		}
	}

	return
}

func decodeFields(e reflect.Value, in map[string]interface{}) error {
	for k, v := range in {
		f := e.FieldByName(k)

		typeName := f.Type().String()
		switch typeName {
		case "fyne.TextAlign", "fyne.TextTruncation", "fyne.TextWrap", "widget.ButtonAlign", "widget.ButtonImportance",
			"widget.ButtonIconPlacement", "widget.Importance", "widget.Orientation", "widget.ScrollDirection":
			f.SetInt(int64(reflect.ValueOf(v).Float()))
		case "fyne.TextStyle":
			f.Set(reflect.ValueOf(decodeTextStyle(reflect.ValueOf(v).Interface().(map[string]interface{}))))
		case "widget.RichTextStyle":
			f.Set(reflect.ValueOf(decodeRichTextStyle(reflect.ValueOf(v).Interface().(map[string]interface{}))))
		case "fyne.Position":
			f.Set(reflect.ValueOf(decodePosition(reflect.ValueOf(v).Interface().(map[string]interface{}))))
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
		case "[]widget.ToolbarItem":
			var items []widget.ToolbarItem
			for _, item := range reflect.ValueOf(v).Interface().([]interface{}) {
				items = append(items, decodeToolbarItem(item.(map[string]interface{})))
			}
			f.Set(reflect.ValueOf(items))
		case "[]widget.RichTextSegment":
			var items []widget.RichTextSegment
			for _, item := range reflect.ValueOf(v).Interface().([]interface{}) {
				obj := &widget.TextSegment{}
				_ = decodeFields(reflect.ValueOf(obj).Elem(), item.(map[string]interface{}))
				items = append(items, obj)
			}
			f.Set(reflect.ValueOf(items))
		case "fyne.CanvasObject":
			return errors.New("unsupported object type")
		case "*url.URL":
			u := &url.URL{}
			decodeFromMap(reflect.ValueOf(v).Interface().(map[string]interface{}), u)
			f.Set(reflect.ValueOf(u))
		case "[]string":
			anySlice := reflect.ValueOf(v).Interface().([]interface{})
			strings := make([]string, len(anySlice))
			for i, a := range anySlice {
				strings[i] = a.(string)
			}
			f.Set(reflect.ValueOf(strings))
		case "time.Time":
			s := reflect.ValueOf(v).String()

			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				fyne.LogError("Failed parse time "+s, err)
			} else {
				f.Set(reflect.ValueOf(t))
			}
		case "*time.Time":
			s := reflect.ValueOf(v).String()

			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				fyne.LogError("Failed parse time "+s, err)
			} else {
				f.Set(reflect.ValueOf(&t))
			}
		case "color.Color":
			c := &color.NRGBA{}
			decodeFromMap(reflect.ValueOf(v).Interface().(map[string]interface{}), c)
			f.Set(reflect.ValueOf(c))
		default:
			if strings.Index(typeName, "int") == 0 {
				f.SetInt(int64(reflect.ValueOf(v).Float()))
			} else if typeName == "float32" {
				f.SetFloat(reflect.ValueOf(v).Float())
			} else if v != nil {
				f.Set(reflect.ValueOf(v))
			}
		}
	}

	return nil
}

func decodeWidget(m map[string]interface{}) fyne.Widget {
	class := m["Type"].(string)
	obj := guidefs.Lookup(class).Create().(fyne.Widget)
	e := reflect.ValueOf(obj).Elem()

	err := decodeFields(e, m["Struct"].(map[string]interface{}))
	if err != nil {
		fyne.LogError("Failed to handle type "+class, err)
	}
	return obj
}

type toolbarItem struct {
	Type string
}

func (toolbarItem) ToolbarObject() fyne.CanvasObject { return nil }
