package diskcache

import (
	"io"
	"os"
	"path/filepath"
)

var (
	defaultFilePerm os.FileMode = 0666
	defaultPathPerm os.FileMode = 0777
)

func writeFile(dir, key string, r io.Reader) (string, int64, error) {
	err := os.MkdirAll(dir, defaultPathPerm)
	if err != nil {
		return "", 0, nil
	}

	fileName := completeFileName(dir, key)
	mode := os.O_WRONLY | os.O_CREATE | os.O_TRUNC // overwrite if exists
	f, err := os.OpenFile(fileName, mode, defaultFilePerm)
	defer f.Close()
	if err != nil {
		return "", 0, err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return "", 0, err
	}
	return fileName, n, nil
}

func completeFileName(dir, name string) string {
	return filepath.Join(dir, name)
}

func validDir(dir string) bool {
	return dir != ""
}
