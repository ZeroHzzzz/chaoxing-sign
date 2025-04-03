package utils

import (
	"strings"
)

func ParseAnalysis(data string) string {
	codeStart := strings.Index(data, "code='+''")
	if codeStart == -1 {
		return ""
	}

	codeStart += 8
	code := data[codeStart:]

	codeEnd := strings.Index(code, "'")
	if codeEnd == -1 {
		return ""
	}

	code = code[:codeEnd]
	return code
}
