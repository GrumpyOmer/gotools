package download

import (
"io"
"net/http"
"os"
"strconv"
"sync"
"sync/atomic"
)

var (
	DownloadMap sync.Map
)

type WriteCounter struct {
	Total uint64	//文件总大小
	Current uint64	//当前已成功下载大小
	Progress int32	//当前下载进度（百分比）
	Err	error	//错误原因
	Success	bool	//是否上传成功
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
	return wc.Total
}

func DownloadFile(filepath string, url string) {
	var(
		err error
	)
	counter := &WriteCounter{}
	DownloadMap.Store(filepath, counter)
	defer func() {
		counter.Err = err
	}()

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
	tmpTotal,_:= strconv.Atoi(resp.Header.Get("Content-Length"))
	// 获取文件大小
	counter.Total = uint64(tmpTotal)
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return
	}
	out.Close()
	if err = os.Rename(filepath+".tmp", filepath); err == nil {
		counter.Success = true
	}
	return
}
