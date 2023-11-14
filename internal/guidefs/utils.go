package guidefs

import "strings"

func escapeLabel(inStr string) (outStr string) {
	outStr = strings.ReplaceAll(inStr, "\"", "\\\"")
	outStr = strings.ReplaceAll(outStr, "\n", "\\n")
	return
}
