package diskcache

import (
	"io"
	"os"
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
	return dir + string(os.PathListSeparator) + name
}

func validDir(dir string) bool {
	return dir != ""
}
