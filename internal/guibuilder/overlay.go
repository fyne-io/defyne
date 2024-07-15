package guibuilder

import (
	"image/color"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fyne-io/defyne/internal/guidefs"
)

type overlay struct {
	widget.BaseWidget

	b         *Builder
	indicator *canvas.Rectangle
}

func newOverlay(b *Builder) *overlay {
	o := &overlay{b: b}
	o.ExtendBaseWidget(o)
	return o
}

func (o *overlay) CreateRenderer() fyne.WidgetRenderer {
	r := canvas.NewRectangle(color.Transparent)
	r.StrokeColor = color.Transparent
	r.StrokeWidth = 4
	r.Resize(fyne.NewSize(20, 10))
	r.Move(fyne.NewPos(10, 10))

	o.indicator = r

	return widget.NewSimpleRenderer(container.NewWithoutLayout(r))
}

func (o *overlay) Tapped(pe *fyne.PointEvent) {
	rootPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(o.b.root)
	pos := pe.AbsolutePosition.Subtract(rootPos)
	obj := findObject(o.b.root, pos)

	// TODO update when an item is removed, inserted, or if the UI resizes
	o.indicator.StrokeColor = theme.PrimaryColor()
	objAbsPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(obj)
	objPos := objAbsPos.Subtract(rootPos)
	o.indicator.Move(objPos)
	o.indicator.Resize(obj.Size())

	o.b.choose(obj)
}

func findObject(o fyne.CanvasObject, p fyne.Position) fyne.CanvasObject {
	switch w := o.(type) {
	case *fyne.Container:
		for _, child := range w.Objects {
			if !insideObject(child, p) {
				continue
			}

			match := findObject(child, p.Subtract(child.Position()))
			if match != nil && isContainerOrWidget(match) {
				return match
			}
			if isContainerOrWidget(child) {
				return child
			}
		}
		return w
	case fyne.Widget:
		class := reflect.TypeOf(o).String()
		info, ok := guidefs.Widgets[class]
		if !ok || !info.IsContainer() {
			return nil
		}

		for _, child := range info.Children(o) {
			if !insideObject(child, p) {
				continue
			}

			match := findObject(child, p.Subtract(child.Position()))
			if match != nil && isContainerOrWidget(match) {
				return match
			}
			if isContainerOrWidget(child) {
				return child
			}
		}
	}

	return nil
}

func insideObject(o fyne.CanvasObject, p fyne.Position) bool {
	pos := o.Position()
	if p.X < pos.X || p.Y < pos.Y {
		return false
	}

	size := o.Size()
	return p.X < pos.X+size.Width && p.Y < pos.Y+size.Height
}

func isContainerOrWidget(o fyne.CanvasObject) bool {
	switch o.(type) {
	case fyne.Widget, *fyne.Container:
		return true
	default:
		return false
	}
}
