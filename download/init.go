package download

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

type writeCounter struct {
	total    uint64 //文件总大小
	current  uint64 //当前已成功下载大小
	progress int32  //当前下载进度（百分比）
	err      error  //错误原因
	success  bool   //是否上传成功
	inUse    bool   //是否被使用
	sync.Mutex
}

func NewWc() *writeCounter {
	return &writeCounter{}
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	atomic.AddUint64(&wc.current, uint64(n))
	atomic.CompareAndSwapInt32(&wc.progress, wc.progress, int32(float64(atomic.LoadUint64(&wc.current))/float64(wc.total)*100))
	return n, nil
}

func (wc *writeCounter) FetchProgress() int32 {
	return atomic.LoadInt32(&wc.progress)
}

func (wc *writeCounter) FetchCurrent() uint64 {
	return atomic.LoadUint64(&wc.current)
}

func (wc *writeCounter) FetchTotal() uint64 {
	return atomic.LoadUint64(&wc.total)
}

func (wc *writeCounter) DownloadRes() (bool, error) {
	return wc.success, wc.err
}

func (wc *writeCounter) initProperty() {
	wc.inUse = true
	wc.success = false
	wc.total = 0
	wc.current = 0
	wc.progress = 0
	wc.err = nil
}

func (wc *writeCounter) DownloadFile(filepath string, url string) error {
	// 异步下载文件 为了保障安全性 一个wc对象同一时刻只能下载一个文件
	wc.Lock()
	defer wc.Unlock()
	if wc.inUse {
		return errors.New("当前wc对象正在使用")
	}
	// 开始占用wc
	wc.initProperty()
	go func(string, string) {
		var (
			err         error
			DownloadMap sync.Map
		)
		defer func() {
			wc.err = err
			wc.inUse = false
		}()
		DownloadMap.Store(filepath, wc)

		out, err := os.Create(filepath + ".tmp")
		if err != nil {
			return
		}
		resp, err := http.Get(url)
		if err != nil {
			out.Close()
			return
		}
		defer resp.Body.Close()
		tmpTotal, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
		// 获取文件大小
		wc.total = uint64(tmpTotal)
		if _, err = io.Copy(out, io.TeeReader(resp.Body, wc)); err != nil {
			out.Close()
			return
		}
		out.Close()
		if err = os.Rename(filepath+".tmp", filepath); err == nil {
			wc.success = true
		}
		return
	}(filepath, url)
	return nil
}
