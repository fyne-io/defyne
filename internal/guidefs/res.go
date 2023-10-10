package guidefs

import (
	"fmt"

	"fyne.io/fyne/v2"
)

type jsonResource struct {
	fyne.Resource `json:"-"`
}

func (r *jsonResource) MarshalJSON() ([]byte, error) {
	icon := "\"" + IconReverse[fmt.Sprintf("%p", r.Resource)] + "\""
	return []byte(icon), nil
}

func WrapResource(r fyne.Resource) fyne.Resource {
	return &jsonResource{r}
}
