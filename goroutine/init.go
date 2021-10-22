package Goroutine

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type (
	queueStruct struct {
		function  func(*sync.WaitGroup) error
		waitGroup *sync.WaitGroup
	}
	goroutineManager struct {
		*sync.Mutex
		max     uint64             // 最大协程数
		current uint64             // 当前正在运行的协程
		queue   chan (queueStruct) // 装载 当前/等待 执行的任务
	}
)

var (
	manager goroutineManager
)

// 初始化
func init() {
	// 初始化通道 默认缓冲1000
	manager.queue = make(chan queueStruct, 1000)
	for v := range manager.queue {
		go func(function func(*sync.WaitGroup) error, w *sync.WaitGroup) {
			defer func() {
				// 执行结束修改当前协程数信息 原子操作保证一致性
				temp := int64(-1)
				// 多一步绕过编译器....
				dec := uint64(temp)
				atomic.AddUint64(&manager.current, dec)
				if err := recover(); err != nil {
					// 记录任务错误 防止进程重启
					fmt.Printf("recover(): %v\n", recover())
				}
			}()
			// running task
			function(w)
		}(v.function, v.waitGroup)
	}
}

// 设置最大协程数
func SetGoroutineNumber(max uint64) {
	// 防止动态修改造成竞态问题 改为原子操作
	atomic.StoreUint64(&manager.max, max)
}

// 生成协程任务
func MakeTask(task func(*sync.WaitGroup) error, w *sync.WaitGroup) error {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	if manager.current == manager.max {
		return errors.New("当前协程数量已达限制")
	}
	// 更新当前协程数信息 原子操作保证一致性
	atomic.AddUint64(&manager.current, 1)
	// 任务写入通道
	manager.queue <- queueStruct{
		function:  task,
		waitGroup: w,
	}
	return nil
}

//查看当前堆积任务
func SearchCurTask() int {
	return len(manager.queue)
}
