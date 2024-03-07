package guidefs

import (
	"fyne.io/fyne/v2"
)

type jsonResource struct {
	fyne.Resource `json:"-"`
}

func (r *jsonResource) MarshalJSON() ([]byte, error) {
	icon := "\"" + IconName(r.Resource) + "\""

	return []byte(icon), nil
}

func WrapResource(r fyne.Resource) fyne.Resource {
	return &jsonResource{r}
}

func IconName(res fyne.Resource) string {
	name := res.Name()
	// strip prefix numbers to unwrap
	for name[0] >= '0' && name[0] <= '9' {
		name = name[1:]
	}

	ret, ok := IconReverse[name]
	if !ok {
		return "BrokenImageIcon"
	}

	return ret
}
