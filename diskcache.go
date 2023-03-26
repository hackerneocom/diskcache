package diskcache

import (
	"bytes"
	"container/list"
	"io"
	"os"
	"strings"
	"sync"
)

type Item struct {
	Key  string
	Size int64
	Path string
}

type Cache struct {
	dir     string // Path to storage directory
	maxSize int64  // Total size of files allowed, unit byte

	sizeUsed int64 // Total size of files added

	list    *list.List
	fileMap map[string]*list.Element

	mutex sync.RWMutex
}

func New(dir string, sz int64) (*Cache, error) {
	if !validDir(dir) {
		return nil, ErrBadDir
	}
	if sz <= 0 {
		return nil, ErrBadSize
	}
	dir = strings.TrimRight(dir, string(os.PathSeparator))

	return &Cache{
		dir:     dir,
		maxSize: sz,
		list:    list.New(),
		fileMap: make(map[string]*list.Element),
	}, nil
}

func (c *Cache) Get(key string) (io.ReadCloser, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if element, ok := c.fileMap[key]; ok {
		c.list.MoveToFront(element)
		item := element.Value.(*Item)
		return os.Open(item.Path)
	}
	return nil, ErrNotFound
}

func (c *Cache) Keys() []string {
	keys := make([]string, len(c.fileMap))
	i := 0
	for element := c.list.Front(); element != nil; element = element.Next() {
		keys[i] = element.Value.(*Item).Key
		i++
	}
	return keys
}

func (c *Cache) Put(key string, data []byte) error {
	size := int64(len(data))
	return c.putReader(key, size, bytes.NewBuffer(data))
}

func (c *Cache) putReader(key string, size int64, r io.Reader) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.validate(size); err != nil {
		return err
	}

	path, size, err := writeFile(c.dir, key, r)
	if err != nil {
		return err
	}
	c.addItem(key, path, size)
	return nil
}

func (c *Cache) validate(size int64) error {
	if size > c.maxSize {
		return ErrTooLarge
	}

	for c.sizeUsed+size > c.maxSize {
		err := c.evictLast()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) evictLast() error {
	if last := c.list.Back(); last != nil {
		item := last.Value.(*Item)
		if err := os.Remove(item.Path); err != nil {
			return err
		}
		c.list.Remove(last)
		delete(c.fileMap, item.Key)
		c.sizeUsed -= item.Size
		return nil
	}
	return nil
}

func (c *Cache) addItem(key, path string, size int64) {
	c.sizeUsed += size
	if elemnt, ok := c.fileMap[key]; ok {
		c.list.Remove(elemnt)
	}

	item := &Item{
		Key:  key,
		Path: path,
		Size: size,
	}
	element := c.list.PushFront(item)
	c.fileMap[key] = element
}
