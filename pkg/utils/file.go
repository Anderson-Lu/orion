package utils

import (
	"path/filepath"
	"strings"
)

type FileUtil struct{}

func (f *FileUtil) SeperatePath(p string) []string {
	return strings.Split(p, string(filepath.Separator))
}
