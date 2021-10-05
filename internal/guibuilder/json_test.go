package guibuilder

import (
	"bytes"
	"testing"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const labelJSON = `{
  "Type": "*widget.Label",
  "Struct": {
    "Hidden": false,
    "Text": "Hi",
    "Alignment": 1,
    "Wrapping": 0,
    "TextStyle": {
      "Bold": true,
      "Italic": false,
      "Monospace": false,
      "TabWidth": 0
    }
  }
}
`

func TestDecodeJSON(t *testing.T) {
	initIcons()
	initWidgets()

	buf := bytes.NewReader([]byte(labelJSON))
	obj := DecodeJSON(buf)

	l, ok := obj.(*widget.Label)
	require.True(t, ok)
	assert.Equal(t, "Hi", l.Text)
	assert.Equal(t, fyne.TextAlignCenter, l.Alignment)
	assert.Equal(t, fyne.TextStyle{Bold: true}, l.TextStyle)
}

func TestEncodeJSON(t *testing.T) {
	l := widget.NewLabelWithStyle("Hi", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	var buf bytes.Buffer
	EncodeJSON(l, &buf)
	assert.Equal(t, labelJSON, buf.String())
}
