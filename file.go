package diskcache

import (
	"io"
	"os"
	"path/filepath"
)

func writeFile(dir, key string, r io.Reader) (string, int64, error) {
	path := filePath(dir, key)

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return "", 0, err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return "", 0, err
	}
	return path, n, nil
}

func filePath(dir, name string) string {
	return filepath.Join(dir, name)
}

func validDir(dir string) bool {
	return dir != ""
}
