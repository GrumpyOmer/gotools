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

type WriteCounter struct {
	Total    uint64 //文件总大小
	Current  uint64 //当前已成功下载大小
	Progress int32  //当前下载进度（百分比）
	Err      error  //错误原因
	Success  bool   //是否上传成功
	InUse    bool   //是否被使用
	sync.Mutex
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	atomic.AddUint64(&wc.Current, uint64(n))
	atomic.CompareAndSwapInt32(&wc.Progress, wc.Progress, int32(float64(atomic.LoadUint64(&wc.Current))/float64(wc.Total)*100))
	return n, nil
}

func (wc *WriteCounter) FetchProgress() int32 {
	return atomic.LoadInt32(&wc.Progress)
}

func (wc *WriteCounter) FetchCurrent() uint64 {
	return atomic.LoadUint64(&wc.Current)
}

func (wc *WriteCounter) FetchTotal() uint64 {
	return atomic.LoadUint64(&wc.Total)
}

func (wc *WriteCounter) DownloadRes() (bool, error) {
	return wc.Success, wc.Err
}

func (wc *WriteCounter) initProperty() {
	wc.InUse = true
	wc.Success = false
	wc.Total = 0
	wc.Current = 0
	wc.Progress = 0
	wc.Err = nil
}
func (wc *WriteCounter) DownloadFile(filepath string, url string) error {
	// 异步下载文件 为了保障安全性 一个wc对象同一时刻只能下载一个文件
	wc.Lock()
	defer wc.Unlock()
	if wc.InUse {
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
			wc.Err = err
			wc.InUse = false
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
		wc.Total = uint64(tmpTotal)
		if _, err = io.Copy(out, io.TeeReader(resp.Body, wc)); err != nil {
			out.Close()
			return
		}
		out.Close()
		if err = os.Rename(filepath+".tmp", filepath); err == nil {
			wc.Success = true
		}
		return
	}(filepath, url)
	return nil
}
