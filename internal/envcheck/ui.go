package envcheck

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ShowSummaryDialog shows a new dialog in the specified window and runs the environment checks.
// After displaying the process will start and the UI will update as it progresses.
func ShowSummaryDialog(w fyne.Window) {
	content := makeEnvcheckForm()
	d := dialog.NewCustom("Checking Development Environment", "Cancel", content, w)
	d.Show()
	go content.runChecks(func() {
		d.SetDismissText("Done")
	})
}

// ShowSummaryWindow shows a new window for the specified app and runs the environment checks.
// After displaying the process will start and the UI will update as it progresses.
func ShowSummaryWindow(a fyne.App) {
	w := a.NewWindow("Environment Check")
	content := makeEnvcheckForm()
	d := dialog.NewCustom("Checking Development Environment", "Cancel", content, w)
	d.SetOnClosed(w.Close)
	d.Show()
	w.Resize(d.MinSize().AddWidthHeight(theme.Padding()*4, theme.Padding()*4))
	w.Show()
	go content.runChecks(func() {
		d.SetDismissText("Done")
	})
}

type checkForm struct {
	widget.Form
}

func makeEnvcheckForm() *checkForm {
	items := make([]*widget.FormItem, len(tasks))
	for i, t := range tasks {
		check := newCheckLine(t)
		row := widget.NewFormItem(t.name, check)
		row.HintText = t.hint
		items[i] = row
	}

	f := &checkForm{}
	f.Items = items
	f.ExtendBaseWidget(f)

	return f
}

func (c *checkForm) runChecks(done func()) {
	for _, f := range c.Items {
		line := f.Widget.(*checkLine)
		line.icon.SetResource(theme.ViewRefreshIcon())

		time.Sleep(time.Second / 10)
		out, err := line.task.test()
		if err != nil {
			line.err = err
			line.icon.SetResource(theme.NewErrorThemedResource(theme.CancelIcon()))
			line.label.SetText("Failed")
			if line.errCB != nil {
				line.errCB(err)
			}
		} else {
			th := theme.NewThemedResource(theme.ConfirmIcon())
			th.ColorName = theme.ColorNameSuccess
			line.icon.SetResource(th)
			line.label.SetText(out)
		}
	}

	done()
}

type checkLine struct {
	widget.BaseWidget
	err   error
	errCB func(error)
	task  *task

	icon  *widget.Icon
	label *widget.Label
}

func newCheckLine(t *task) *checkLine {
	l := &checkLine{icon: widget.NewIcon(nil),
		label: widget.NewLabel("Waiting..."), task: t}
	l.ExtendBaseWidget(l)
	return l
}

func (l *checkLine) Validate() error {
	return l.err
}

func (l *checkLine) SetOnValidationChanged(cb func(error)) {
	l.errCB = cb
}

func (l *checkLine) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, nil, l.icon,
		l.label))
}
