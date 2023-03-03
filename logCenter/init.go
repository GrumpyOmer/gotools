package logCenter

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
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
)

func init() {
	fmt.Println("log center init")
	// 默认10个消费者
	for i := 0; i < 10; i++ {
		go func() {
			for {
				data := <-channel
				func() {
					str, _ := json.Marshal(data)
					fileName := time.Now().Format("20060102")
					var (
						filename = "./log/" + fileName + ".text"
						f        *os.File
						err1     error
					)
					if _, err := os.Stat(filename); os.IsNotExist(err) {
						if err1 = os.MkdirAll("./log", 0777); err1 != nil {
							fmt.Println(err1.Error())
							return
						}
						f, err1 = os.Create(filename) //创建文件
					} else {
						f, err1 = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
					}
					if err1 != nil {
						fmt.Println(err1.Error())
						return
					}
					defer func() {
						if err1 = f.Close(); err1 != nil {
							fmt.Println("关闭文件失败，err：" + err1.Error())
						}
					}()
					_, err1 = f.WriteString(string(str) + "\n") //写入文件(字符串)
					if err1 != nil {
						fmt.Println(err1.Error())
						return
					}
					// 刷入磁盘
					if err1 = f.Sync(); err1 != nil {
						fmt.Println("刷入磁盘失败，err：" + err1.Error())
					}
				}()
			}
		}()
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
