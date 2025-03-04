package services

import (
	"path"
	"strings"
)

// Return basename and extension
func GetBasenameAndExtension(file string) (string, string) {
	ext := path.Ext(file)
	return strings.TrimSuffix(file, ext), ext
}
