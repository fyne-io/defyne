package guibuilder

import (
	"fmt"
	"image/color"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var current fyne.CanvasObject

type jsonResource struct {
	fyne.Resource `json:"-"`
}

func (r *jsonResource) MarshalJSON() ([]byte, error) {
	icon := "\"" + iconReverse[fmt.Sprintf("%p", r.Resource)] + "\""
	return []byte(icon), nil
}

type overlayContainer struct {
	widget.BaseWidget
	name      string
	c, parent *fyne.Container
}

func (o *overlayContainer) CreateRenderer() fyne.WidgetRenderer {
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 4
	return &overRender{p: o, c: o.c, r: border}
}

func (o *overlayContainer) GoString() string {
	code := widgets["*fyne.Container"].gostring(o.c)
	if o.name != "" {
		defs[o.name] = code
		return "g." + o.name
	}
	return code
}

func (o *overlayContainer) MinSize() fyne.Size {
	min := o.c.MinSize()
	if min.IsZero() {
		return fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	}

	return min
}

func (o *overlayContainer) Move(p fyne.Position) {
	o.c.Move(p)
	o.BaseWidget.Move(p)
}

func (o *overlayContainer) Refresh() {
	o.BaseWidget.Refresh()
	o.c.Refresh()
}

func (o *overlayContainer) Resize(s fyne.Size) {
	o.c.Resize(s)
	o.BaseWidget.Resize(s)
}

func (o *overlayContainer) Tapped(e *fyne.PointEvent) {
	setCurrent(o)
	choose(o)
}

func (o *overlayContainer) Object() fyne.CanvasObject {
	return o.c
}

type overlayWidget struct {
	widget.BaseWidget
	name string

	child  fyne.Widget
	parent *fyne.Container
}

func (w *overlayWidget) CreateRenderer() fyne.WidgetRenderer {
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 4

	return &overRender{p: w, r: border}
}

func (w *overlayWidget) GoString() string {
	var code string
	name := reflect.TypeOf(w.child).String()
	if widgets[name].gostring != nil {
		code = widgets[name].gostring(w.child)
	} else {
		code = fmt.Sprintf("%#v", w.child)
	}

	if w.name != "" {
		defs[w.name] = code
		return "g." + w.name
	}
	return code
}

func (w *overlayWidget) Object() fyne.CanvasObject {
	return w.child
}

func (w *overlayWidget) Packages() []string {
	name := reflect.TypeOf(w.child).String()
	if widgets[name].packages != nil {
		return widgets[name].packages(w.child)
	}

	return []string{"widget"}
}

func (w *overlayWidget) Refresh() {
	w.BaseWidget.Refresh()
	w.child.Refresh()
}

func (w *overlayWidget) Tapped(e *fyne.PointEvent) {
	setCurrent(w)
	choose(w)
}

type overRender struct {
	p fyne.CanvasObject
	c *fyne.Container
	r *canvas.Rectangle
}

func (o overRender) BackgroundColor() color.Color {
	return color.Transparent
}

func (o overRender) Destroy() {
}

func (o overRender) Layout(s fyne.Size) {
	o.r.Resize(s)
}

func (o overRender) MinSize() fyne.Size {
	return fyne.Size{}
}

func (o overRender) Objects() []fyne.CanvasObject {
	if o.c == nil {
		return []fyne.CanvasObject{o.r}
	}

	return append([]fyne.CanvasObject{o.r}, o.c.Objects...)
}

func (o overRender) Refresh() {
	if o.p == current {
		o.r.StrokeColor = theme.PrimaryColor()
	} else {
		o.r.StrokeColor = color.Transparent
	}
	o.r.Refresh()
}

func setCurrent(o fyne.CanvasObject) {
	old := current
	current = o
	if old != nil {
		old.Refresh()
	}
	current.Refresh()
}

func wrapContent(o fyne.CanvasObject, parent *fyne.Container) fyne.CanvasObject {
	switch obj := o.(type) {
	case *fyne.Container:
		var c *fyne.Container
		if obj.Layout == nil {
			c = container.NewWithoutLayout()
		} else {
			c = container.New(obj.Layout)
		}
		items := make([]fyne.CanvasObject, len(obj.Objects))
		for i, child := range obj.Objects {
			items[i] = wrapContent(child, c)
		}
		c.Objects = items

		o := &overlayContainer{c: c, parent: parent}
		layoutProps[o.c] = map[string]string{"layout": "VBox"}
		o.ExtendBaseWidget(o)
		return o
	case fyne.Widget:
		return wrapWidget(obj, parent)
	}

	return nil //?
}

func wrapWidget(w fyne.Widget, parent *fyne.Container) fyne.CanvasObject {
	switch t := w.(type) {
	case *widget.Icon:
		t.Resource = wrapResource(t.Resource)
	case *widget.Button:
		if t.Icon != nil {
			t.Icon = wrapResource(t.Icon)
		}
	}
	o := &overlayWidget{child: w, parent: parent}
	o.ExtendBaseWidget(o)
	return container.NewMax(w, o)
}

func wrapResource(r fyne.Resource) fyne.Resource {
	return &jsonResource{r}
}

// unwrap gets the wrapping container or widget for a passed CanvasObject.
func unwrap(o fyne.CanvasObject) (*overlayContainer, *overlayWidget) {
	if c, ok := o.(*fyne.Container); ok { // the content of an overlayWidget container
		return nil, c.Objects[1].(*overlayWidget)
	} else if w, ok := o.(*overlayWidget); ok {
		return nil, w
	} else if c, ok := o.(*overlayContainer); ok {
		return c, nil
	}

	return nil, nil
}
