package cache

import (
	"errors"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/thread"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var CacheMaxFileSize int64 = 1024 * 1024 //
var CacheFileTimeout int64 = 60 * 3      //秒

type CacheFileItem struct {
	Info         os.FileInfo
	LastReadTime time.Time
	Byte         []byte
}
type CacheFileByte struct {
	m map[string]*CacheFileItem
	sync.RWMutex
}

func (c *CacheFileByte) GetAllCache() map[string]*CacheFileItem {
	return c.m
}
func (c *CacheFileByte) readFile(path string) (*CacheFileItem, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, errors.New("目标路径是一个文件夹：" + path)
	}

	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() > CacheMaxFileSize {

		return &CacheFileItem{Info: fileInfo, LastReadTime: time.Now(), Byte: fileByte}, nil

	}

	item := &CacheFileItem{}
	item.Info = fileInfo
	item.LastReadTime = time.Now()

	item.Byte = fileByte

	c.Lock()
	defer c.Unlock()
	if c.m == nil {
		c.m = make(map[string]*CacheFileItem)
	}
	c.m[path] = item
	return item, nil
}
func (c *CacheFileByte) RemoveTimeout(path string) {

	cache := c.get(path)
	if cache != nil {

		if time.Now().Unix()-cache.LastReadTime.Unix() > CacheFileTimeout {

			c.Lock()
			defer c.Unlock()
			delete(c.m, path)

		}

	}

}
func (c *CacheFileByte) Read(path string) (*CacheFileItem, error) {

	cache := c.get(path)
	if cache == nil {
		var err error
		cache, err = c.readFile(path)
		if err != nil {
			return nil, err
		}
	} else {

		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if cache.Info.Size() != fileInfo.Size() || fileInfo.ModTime().Unix() != cache.Info.ModTime().Unix() {
			_cache, _err := c.readFile(path)
			if _err != nil {
				return nil, _err
			} else {
				cache = _cache
			}
		}

	}

	return cache, nil
}

func (c *CacheFileByte) get(path string) *CacheFileItem {

	c.RLock()
	defer c.RUnlock()
	return c.m[path]

}

var cache = &CacheFileByte{}

func Read(path string) (*CacheFileItem, error) {
	return cache.Read(path)
}
func init() {

	thread.NewCoroutine(func(option *thread.Option) {
		for {
			allCache := cache.GetAllCache()
			for path := range allCache {

				cache.RemoveTimeout(path)

			}

			time.Sleep(time.Second * 3)
		}

	}, func(option *thread.Option) {

		glog.Trace("检测缓存文件协程，意外重启")

	}, &thread.Option{ReRun: true})

}
