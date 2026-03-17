package logCenter

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type (
	stackInfo struct {
		// 日志内容
		Content string
		// 日志打印位置(文件/行)
		Location string
		// 日志打印函数名
		FunctionName string
		// 日志打印时间
		LogTime string
	}
)

var (
	// 默认10W记录缓存
	channel = make(chan stackInfo, 100000)
	// 是否将日志强制写入磁盘，默认不强制
	fSync        bool
	fSyncMutex   sync.RWMutex
	// 文件操作相关的并发安全控制
	fileMutex    sync.Mutex   // 保护文件操作的互斥锁
	currentFile  *os.File     // 当前日志文件句柄
	currentDate  string       // 当前日志文件对应的日期
	logDir       = "./log"   // 日志目录
)

func init() {
	fmt.Println("log center init")
	// 启动时确保日志目录存在
	if err := os.MkdirAll(logDir, 0777); err != nil {
		fmt.Println("创建日志目录失败：", err.Error())
	}
	// 单消费者模式：日志写入是IO密集型，单消费者避免文件竞争，更高效
	go func() {
		for {
			data := <-channel
			writeLog(data)
		}
	}()
}

// writeLog 写入日志（线程安全）
func writeLog(data stackInfo) {
	str, _ := json.Marshal(data)
	today := time.Now().Format("20060102")

	fileMutex.Lock()
	defer fileMutex.Unlock()

	// 检查是否需要切换文件（日期变化或文件未打开）
	if currentDate != today || currentFile == nil {
		// 关闭旧文件
		if currentFile != nil {
			currentFile.Close()
		}
		// 打开新文件
		filename := logDir + "/" + today + ".text"
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("打开日志文件失败：", err.Error())
			return
		}
		currentFile = f
		currentDate = today
	}

	// 写入日志
	_, err := currentFile.WriteString(string(str) + "\n")
	if err != nil {
		fmt.Println("写入日志失败：", err.Error())
		return
	}

	// 检查是否需要刷盘
	fSyncMutex.RLock()
	shouldSync := fSync
	fSyncMutex.RUnlock()
	if shouldSync {
		if err := currentFile.Sync(); err != nil {
			fmt.Println("刷入磁盘失败：", err.Error())
		}
	}
}

func Add(content string) {
	pc, codePath, codeLine, _ := runtime.Caller(1)
	info := stackInfo{
		Content: content,
		// 拼接文件名与所在行
		Location: fmt.Sprintf("%s:%d", codePath, codeLine),
		// 根据PC获取函数名
		FunctionName: runtime.FuncForPC(pc).Name(),
		LogTime:      time.Now().String(),
	}
	channel <- info
}

func SaveFSync(o bool) {
	fSyncMutex.Lock()
	fSync = o
	fSyncMutex.Unlock()
}
