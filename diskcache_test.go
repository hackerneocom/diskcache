package diskcache

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

var storageDir string

var blobs = map[string][]byte{
	"null":       {},
	"empty.text": []byte(``),
	"hello":      []byte(`Hello World`),
}

// clearStorage empties the temporary storage directory
func clearStorage() {
	err := os.RemoveAll(storageDir)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(storageDir, 0777)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	// Create a temporary storage directory for tests
	name, err := ioutil.TempDir("", "cache")
	if err != nil {
		log.Fatal(err)
	}
	storageDir = name
	defer os.RemoveAll(name)

	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	for i, c := range []struct {
		dir     string
		maxSize int64
		err     error
	}{
		{
			dir:     "",
			maxSize: 1024,
			err:     ErrBadDir,
		},
		{
			dir:     storageDir,
			maxSize: 0,
			err:     ErrBadSize,
		},
	} {
		clearStorage()

		_, err := New(c.dir, c.maxSize)
		if err != c.err {
			t.Fatalf("#%d: Expected err == %q, got %q", i+1, c.err, err)
		}
	}
}

func TestPut(t *testing.T) {
	clearStorage()

	cache, err := New(storageDir, 1024)
	catch(err)

	for k, b := range blobs {
		err := cache.Put(k, b)
		catch(err)
	}

	for k, b := range blobs {
		path := completeFileName(storageDir, k)
		data, err := ioutil.ReadFile(path)
		catch(err)

		if !bytes.Equal(data, b) {
			t.Fatalf("Expected data == %q, got %q", b, data)
		}
	}
}

func TestPutOverSize(t *testing.T) {
	clearStorage()

	cache, err := New(storageDir, 4)
	catch(err)

	err = cache.Put("a", []byte("12345"))
	if err != ErrTooLarge {
		t.Fatalf("Expected err == %q, got %q", ErrTooLarge, err)
	}
}

func TestGetNotExistKey(t *testing.T) {
	clearStorage()

	cache, err := New(storageDir, 4)
	catch(err)

	reader, err := cache.Get("not_exist_key")
	if err != ErrNotFound {
		t.Fatalf("Expected err == %q, got %q", ErrNotFound, err)
	}
	if reader != nil {
		t.Fatalf("Expected reader is nil, but not is")
	}
}

func TestGetAndRePutSameKey(t *testing.T) {
	clearStorage()

	cache, err := New(storageDir, 4)
	catch(err)

	cachedData := []byte("123")
	err = cache.Put("a", cachedData)
	catch(err)
	reader, err := cache.Get("a")
	catch(err)
	data, err := ioutil.ReadAll(reader)
	catch(err)
	if !reflect.DeepEqual(data, cachedData) {
		t.Fatalf("Expected data == %q, got %q", cachedData, data)
	}

	reCachedData := []byte("1234")
	err = cache.Put("a", reCachedData)
	catch(err)
	reader, err = cache.Get("a")
	catch(err)
	data, err = ioutil.ReadAll(reader)
	catch(err)
	if !reflect.DeepEqual(data, reCachedData) {
		t.Fatalf("Expected data == %q, got %q", reCachedData, data)
	}
}

func TestEvictLast(t *testing.T) {
	clearStorage()

	cache, err := New(storageDir, 8)
	catch(err)

	err = cache.Put("a", []byte("123456"))
	catch(err)
	err = cache.Put("b", []byte("7"))
	catch(err)
	assertKeys(t, cache.Keys(), []string{"b", "a"})

	err = cache.Put("c", []byte("8"))
	catch(err)
	assertKeys(t, cache.Keys(), []string{"c", "b", "a"})

	err = cache.Put("d", []byte("9"))
	catch(err)
	assertKeys(t, cache.Keys(), []string{"d", "c", "b"})

	err = cache.Put("a", []byte("12345678"))
	catch(err)
	assertKeys(t, cache.Keys(), []string{"a"})
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func assertKeys(t *testing.T, keys []string, expected []string) {
	if len(keys) != len(expected) {
		t.Fatalf("Expected %d key(s), got %d", len(expected), len(keys))
	}
	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("Expected keys == %q, got %q", expected, keys)
	}
}
