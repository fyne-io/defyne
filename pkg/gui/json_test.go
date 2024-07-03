package gui

import (
	"bytes"
	"testing"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/fyne-io/defyne/internal/guidefs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const labelJSON = `{
  "Type": "*widget.Label",
  "Name": "myLabel",
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
      "TabWidth": 0,
      "Underline": false
    },
    "Truncation": 0,
    "Importance": 0
  }
}
`

func TestDecodeObject(t *testing.T) {
	guidefs.InitOnce()

	buf := bytes.NewReader([]byte(labelJSON))
	obj, meta, err := DecodeObject(buf)
	assert.Nil(t, err)

	l, ok := obj.(*widget.Label)
	require.True(t, ok)
	assert.Equal(t, "Hi", l.Text)
	assert.Equal(t, "myLabel", meta[l]["name"])
	assert.Equal(t, fyne.TextAlignCenter, l.Alignment)
	assert.Equal(t, fyne.TextStyle{Bold: true}, l.TextStyle)
}

func TestEncodeObject(t *testing.T) {
	l := widget.NewLabelWithStyle("Hi", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	props := map[string]string{"name": "myLabel"}
	meta := map[fyne.CanvasObject]map[string]string{l: props}

	var buf bytes.Buffer
	err := EncodeObject(l, meta, &buf)
	assert.Nil(t, err)
	assert.Equal(t, labelJSON, buf.String())
}
