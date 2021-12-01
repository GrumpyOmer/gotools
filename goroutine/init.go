package Goroutine

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type (
	queueStruct struct {
		function  func(*sync.WaitGroup)
		waitGroup *sync.WaitGroup
	}
	goroutineManager struct {
		sync.Mutex
		max     uint64             // 最大协程数
		current uint64             // 当前正在运行的协程
		queue   chan (queueStruct) // 执行的任务
		waitQueue chan (queueStruct) // 等待执行的任务
	}
)

var (
	manager goroutineManager
)

// 初始化
func init() {
	// 初始化通道 默认缓冲1000
	manager.queue = make(chan queueStruct, 1000)
	// 初始化待消费通道 默认缓冲3000
	manager.waitQueue = make(chan queueStruct, 3000)
	// 初始化任务消费者
	go func() {
		for v := range manager.queue {
			go func(function func(*sync.WaitGroup), w *sync.WaitGroup) {
				defer func() {
					// 执行结束修改当前协程数信息 原子操作保证一致性
					temp := int64(-1)
					// 多一步绕过编译器....
					dec := uint64(temp)
					atomic.AddUint64(&manager.current, dec)
					// 等待组处理
					w.Done()
					if err := recover(); err != nil {
						// 记录任务错误 防止进程重启
						fmt.Printf("recover(): %v\n", recover())
					}
				}()
				// running task
				function(w)
			}(v.function, v.waitGroup)
		}
	}()

	// 初始化待执行任务消费者
	go func() {
		for {
			current := atomic.LoadUint64(&manager.current)
			// 利用cas保证协程数量在控制范围内
			if len(manager.waitQueue) > 0 &&
				current < atomic.LoadUint64(&manager.max) &&
				atomic.CompareAndSwapUint64(&manager.current, current, current+1) {
				currentTask:=<-manager.waitQueue
				go func(function func(*sync.WaitGroup), w *sync.WaitGroup) {
					defer func() {
						// 执行结束修改当前协程数信息 原子操作保证一致性
						temp := int64(-1)
						// 多一步绕过编译器....
						dec := uint64(temp)
						atomic.AddUint64(&manager.current, dec)
						// 等待组处理
						w.Done()
						if err := recover(); err != nil {
							// 记录任务错误 防止进程重启
							fmt.Printf("recover(): %v\n", recover())
						}
					}()
					// running task
					function(w)
				}(currentTask.function, currentTask.waitGroup)
			}
		}
	}()
}

// 设置最大协程数
func SetGoroutineNumber(max uint64) error {
	// 防止动态修改造成竞态问题 改为原子操作
	// 防止要修改的最大协程数量小于当前已有协程数量
	if max < atomic.LoadUint64(&manager.current) {
		return errors.New("非法值：此刻已有协程数大于当前设置协程数")
	}
	atomic.StoreUint64(&manager.max, max)
	return nil
}

// 生成协程任务
func MakeTask(task func(*sync.WaitGroup), w *sync.WaitGroup) error {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()

	if atomic.LoadUint64(&manager.current) == atomic.LoadUint64(&manager.max) {
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

// 批量生成协程任务
func BatchMakeTask(tasks []func(*sync.WaitGroup), w *sync.WaitGroup) error {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	for _,v:= range tasks {
		if atomic.LoadUint64(&manager.current) == atomic.LoadUint64(&manager.max) {
			// 当前协程已跑满 任务保存进执行待执行任务通道
			manager.waitQueue <- queueStruct{
				function:  v,
				waitGroup: w,
			}
			continue
		}
		// 更新当前协程数信息 原子操作保证一致性
		atomic.AddUint64(&manager.current, 1)
		// 任务写入通道
		manager.queue <- queueStruct{
			function:  v,
			waitGroup: w,
		}
	}

	return nil
}

//查看当前堆积任务
func SearchCurTask() int {
	return len(manager.queue)
}
