package guidefs

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	// GraphicsNames is an array with the list of names of all the graphical primitives
	GraphicsNames []string

	// Graphics provides the info about the type of canvas object primitives
	Graphics map[string]WidgetInfo
)

func initGraphics() {
	Graphics = map[string]WidgetInfo{
		"*canvas.Rectangle": {
			Name: "Rectangle",
			Create: func() fyne.CanvasObject {
				return canvas.NewRectangle(color.Black)
			},
			Edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
				r := obj.(*canvas.Rectangle)
				return []*widget.FormItem{
					widget.NewFormItem("Fill", newColorButton(r.FillColor, func(c color.Color) {
						r.FillColor = c
						r.Refresh()
					})),
					widget.NewFormItem("Corner", newSliderButton(r.CornerRadius, func(f float32) {
						r.CornerRadius = f
						r.Refresh()
					})),
					widget.NewFormItem("Stroke", newSliderButton(r.StrokeWidth, func(f float32) {
						r.StrokeWidth = f
						r.Refresh()
					})),
					widget.NewFormItem("Color", newColorButton(r.StrokeColor, func(c color.Color) {
						r.StrokeColor = c
						r.Refresh()
					})),
				}
			},
		},
	}

	GraphicsNames = extractNames(Graphics)
}

// TODO tidy the API and move to a widget package

func newColorButton(c color.Color, fn func(color.Color)) fyne.CanvasObject {
	// TODO get the window passed in somehow
	w := fyne.CurrentApp().Driver().AllWindows()[0]

	input := widget.NewEntry()
	input.SetText(formatColor(c))
	preview := newColorTapper(c, func(col color.Color) {
		raw := formatColor(col)
		input.SetText(raw)
		fn(col)
	}, w)

	input.OnChanged = func(raw string) {
		c := parseColor(raw)
		preview.setColor(c)
		fn(c)
	}
	return container.NewBorder(nil, nil, preview, nil, input)
}

type colorTapper struct {
	widget.BaseWidget

	r   *canvas.Rectangle
	fn  func(color.Color)
	win fyne.Window
}

func newColorTapper(c color.Color, fn func(color.Color), win fyne.Window) *colorTapper {
	preview := canvas.NewRectangle(c)
	preview.SetMinSize(fyne.NewSquareSize(32))

	t := &colorTapper{r: preview, fn: fn, win: win}
	t.ExtendBaseWidget(t)
	return t
}

func (c *colorTapper) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.r)
}

func (c *colorTapper) Tapped(_ *fyne.PointEvent) {
	d := dialog.NewColorPicker("Choose Color", "Pick a color", c.fn, c.win)
	d.Advanced = true
	d.Show()
}

func (c *colorTapper) setColor(col color.Color) {
	c.r.FillColor = col
	c.r.Refresh()
}

func newSliderButton(f float32, fn func(float32)) fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetText(strconv.Itoa(int(f)))
	slide := widget.NewSlider(0, 32)
	slide.SetValue(float64(f))

	input.OnChanged = func(s string) {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return
		}

		slide.SetValue(f)
		fn(float32(f))
	}
	slide.OnChanged = func(f float64) {
		input.SetText(fmt.Sprintf("%0.0f", f)) // format like an int
		fn(float32(f))
	}
	return container.NewBorder(nil, nil, input, nil, slide)
}

func parseColor(s string) color.Color {
	if s == "" {
		return color.Black
	}

	var rgb int
	_, err := fmt.Sscanf(s, "#%x", &rgb)
	if err != nil {
		return color.Transparent
	}

	hasAlpha := len(s) > 7
	a := 0xff
	offset := 0
	if hasAlpha {
		a = rgb & 0xff
		offset = 8
	}

	b := rgb >> offset & 0xff
	gg := rgb >> (offset + 8) & 0xff
	r := rgb >> (offset + 16) & 0xff
	return color.NRGBA{R: uint8(r), G: uint8(gg), B: uint8(b), A: uint8(a)}
}

func formatColor(c color.Color) string {
	if c == nil {
		return "#000000"
	}
	ch := color.RGBAModel.Convert(c).(color.RGBA)
	if ch.A == 0xff {
		return fmt.Sprintf("#%.2x%.2x%.2x", ch.R, ch.G, ch.B)
	}

	return fmt.Sprintf("#%.2x%.2x%.2x%.2x", ch.R, ch.G, ch.B, ch.A)
}
