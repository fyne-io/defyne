package gui

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/fyne-io/defyne/internal/guidefs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const labelJSON = `{
  "Type": "*widget.Label",
  "Name": "%s",
  "Struct": {
    "Hidden": false,
    "Text": "Hi",
    "Alignment": 1,
    "Wrapping": 0,
    "TextStyle": {
      "Bold": true,
      "Italic": false,
      "Monospace": false,
      "Symbol": false,
      "TabWidth": 0
    },
    "Truncation": 0,
    "Importance": 0
  }
}`

func labelJSONWith(indent string) string {
	in := fmt.Sprintf(labelJSON, "")
	out := ""

	rows := strings.Split(in, "\n")
	for i, row := range rows {
		in := indent
		if i == 0 {
			in = ""
		}
		end := "\n"
		if i >= len(rows)-1 {
			end = ""
		}
		out = out + in + row + end
	}

	return out
}

var splitJSON = `{
  "Type": "*container.Split",
  "Name": "mySplit",
  "Struct": {
    "Horizontal": true,
    "Leading": ` + labelJSONWith("    ") + `,
    "Offset": 0.75,
    "Trailing": ` + labelJSONWith("    ") + `
  }
}
`

func TestDecodeObject(t *testing.T) {
	guidefs.InitOnce()

	buf := bytes.NewReader([]byte(fmt.Sprintf(labelJSON, "myLabel")))
	obj, meta, err := DecodeObject(buf)
	assert.Nil(t, err)

	l, ok := obj.(*widget.Label)
	require.True(t, ok)
	assert.Equal(t, "Hi", l.Text)
	assert.Equal(t, "myLabel", meta[l]["name"])
	assert.Equal(t, fyne.TextAlignCenter, l.Alignment)
	assert.Equal(t, fyne.TextStyle{Bold: true}, l.TextStyle)
}

func TestDecodeSplit(t *testing.T) {
	buf := bytes.NewReader([]byte(splitJSON))
	obj, meta, err := DecodeObject(buf)
	assert.Nil(t, err)

	s, ok := obj.(*container.Split)
	require.True(t, ok)
	assert.True(t, s.Horizontal)
	assert.Equal(t, 0.75, s.Offset)
	assert.NotNil(t, s.Leading)
	assert.NotNil(t, s.Trailing)
	assert.Equal(t, "mySplit", meta[s]["name"])

	o1, ok := s.Leading.(*widget.Label)
	assert.True(t, ok)
	o2, ok := s.Trailing.(*widget.Label)
	assert.True(t, ok)
	assert.Equal(t, "Hi", o1.Text)
	assert.Equal(t, "Hi", o2.Text)
	assert.Equal(t, fyne.TextStyle{Bold: true}, o1.TextStyle)
	assert.Equal(t, fyne.TextStyle{Bold: true}, o2.TextStyle)
}

func TestEncodeObject(t *testing.T) {
	l := widget.NewLabelWithStyle("Hi", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	props := map[string]string{"name": "myLabel"}
	meta := map[fyne.CanvasObject]map[string]string{l: props}

	var buf bytes.Buffer
	err := EncodeObject(l, meta, &buf)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf(labelJSON, "myLabel")+"\n", buf.String())
}

func TestEncodeSplit(t *testing.T) {
	l1 := widget.NewLabelWithStyle("Hi", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	l2 := widget.NewLabelWithStyle("Hi", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	s := container.NewVSplit(l1, l2)
	s.Horizontal = true
	s.Offset = 0.75

	props := map[string]string{"name": "mySplit"}
	meta := map[fyne.CanvasObject]map[string]string{s: props}

	var buf bytes.Buffer
	err := EncodeObject(s, meta, &buf)
	assert.Nil(t, err)
	assert.Equal(t, splitJSON, buf.String())
}
