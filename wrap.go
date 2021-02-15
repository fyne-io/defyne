package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var current fyne.CanvasObject

func setCurrent(o fyne.CanvasObject) {
	old := current
	current = o
	if old != nil {
		old.Refresh()
	}
	current.Refresh()
}

type overlayContainer struct {
	widget.BaseWidget
	c *fyne.Container
}

func (o *overlayContainer) CreateRenderer() fyne.WidgetRenderer {
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 4
	return &overRender{p: o, c: o.c, r: border}
}

func (o *overlayContainer) GoString() string {
	objs := ""
	for _, obj := range o.c.Objects {
		objs += fmt.Sprintf("\t%#v,\n", obj.(*fyne.Container).Objects[0])
	}
	l := layoutProps[o.c]["layout"]
	return fmt.Sprintf("&fyne.Container{Layout: layout.New%s(), Objects: []fyne.CanvasObject{\n%s}}", l, objs)
}

func (o *overlayContainer) Move(p fyne.Position) {
	o.c.Move(p)
	o.BaseWidget.Move(p)
}

func (o *overlayContainer) Resize(s fyne.Size) {
	o.c.Resize(s)
	o.BaseWidget.Resize(s)
}

func (o *overlayContainer) Tapped(e *fyne.PointEvent) {
	setCurrent(o)
	choose(o.c)
}

func (o *overlayContainer) Object() fyne.CanvasObject {
	return o.c
}

type overlayWidget struct {
	widget.BaseWidget
	child fyne.Widget
}

func (w *overlayWidget) CreateRenderer() fyne.WidgetRenderer {
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 4

	return &overRender{p: w, r: border}
}

func (w *overlayWidget) GoString() string {
	return fmt.Sprintf("%#v", w.child)
}

func (w *overlayWidget) Object() fyne.CanvasObject {
	return w.child
}

func (w *overlayWidget) Tapped(e *fyne.PointEvent) {
	setCurrent(w)
	choose(w.child)
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

func wrapContent(o fyne.CanvasObject) fyne.CanvasObject {
	switch obj := o.(type) {
	case *fyne.Container:
		items := make([]fyne.CanvasObject, len(obj.Objects))
		for i, child := range obj.Objects {
			items[i] = wrapContent(child)
		}

		o := &overlayContainer{c: container.New(obj.Layout, items...)}
		layoutProps[o.c] = map[string]string{"layout": "VBox"}
		o.ExtendBaseWidget(o)
		return o
	case fyne.Widget:
		return wrapWidget(obj)
	}

	return nil //?
}

func wrapWidget(w fyne.Widget) fyne.CanvasObject {
	o := &overlayWidget{child: w}
	o.ExtendBaseWidget(o)
	return container.NewMax(w, o)
}
