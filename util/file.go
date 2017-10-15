package util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileNameNoExt(fpath string) string {
	base := filepath.Base(fpath)
	return strings.TrimSuffix(base, filepath.Ext(fpath))
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
