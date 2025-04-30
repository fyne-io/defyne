package guidefs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type layoutListItem struct {
	widget.BaseWidget
	icon     *widget.Icon
	label    *widget.Label
	OnTapped func()
}

func newLayoutListItem(r fyne.Resource, s string, cb func()) *layoutListItem {
	item := &layoutListItem{
		icon:     widget.NewIcon(r),
		label:    widget.NewLabel(s),
		OnTapped: cb,
	}
	item.ExtendBaseWidget(item)

	return item
}

func (item *layoutListItem) SetIcon(r fyne.Resource) {
	item.icon.SetResource(r)
}

func (item *layoutListItem) SetText(s string) {
	item.label.SetText(s)
}

func (item *layoutListItem) Tapped(ev *fyne.PointEvent) {
	if item.OnTapped == nil {
		return
	}
	item.OnTapped()
}

func (item *layoutListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewHBox(item.icon, item.label))
}
