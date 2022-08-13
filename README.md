# gotools
##### Language: 🇨🇳 | [🇺🇸](./README-EN.md)
[![OSCS Status](https://www.oscs1024.com/platform/badge/GrumpyOmer/gotools.svg?size=small)](https://www.oscs1024.com/project/GrumpyOmer/gotools?ref=badge_small)
#### 注：一些在使用golang期间封装的提升效率的小工具
## 简介
## Installation
```
    go get github.com/GrumpyOmer/gotools@latest
```
# <- goroutine ->
### *基于任务分发的协程池*
##### 引入时先进行包初始化，这将会初始化协程任务的消费线程
```
    _ "github.com/GrumpyOmer/gotools/goroutine"
```
##### 使用之前需要根据实际情况去设置协程的最大数量，通过调用  SetGoroutineNumber(uint64)  方法
```
    SetGoroutineNumber(CPUCoreNum)
```
##### 使用时通过调用  MakeTask(func(), *sync.WaitGroup) error  方法。需要将业务逻辑封装进 func() 函数内，以及生成的管控业务逻辑的WaitGroup两个作为参数传递进去，该方法会响应一个error来告知应用层任务是否成功丢进协程池(任务执行是异步的)，外部如果有业务需要必须自己用等待组的Wait()方法去等待业务执行结束。
##### 执行完成后，需要通过自行获取传入的 func() 函数内部的error来做到业务层面的判断，来达到业务层面预期的效果，使用例子如下：
```
    var (
        a int
        b int
        err1 error
        err2 error
        w = sync.WaitGroup{}
        tempFunc1 = func() {
            b = a+1
            //自行传入错误和捕捉错误
            err1 = errors.New("this is a demo")
        }
        
        tempFunc2 = func() {
            b = a+1
            //自行传入错误和捕捉错误
            err2 = errors.New("this is a demo")
        }
    )
    w.Add(2)
    if test1:= goroutine.MakeTask(tempFunc1, &w); test1 != nil {
        // 下发失败
        w.Done()
        fmt.Println(test1)
    }
    
    if test2:= goroutine.MakeTask(tempFunc2, &w); test2 != nil {
        // 下发失败
        w.Done()
        fmt.Println(test2)
    }
    
    w.Wait()
    fmt.Println(a)
    fmt.Println(b)
    //自行捕获连带func()一起传入协程池的err
    fmt.Println(err1)
    fmt.Println(err2)
```

# <- mysql ->
### *基于gorm封装的一套mysql组件*
##### 首先在使用之前，需要通过ConfigInit([]byte)方法把db的配置以基于json字符串转换后的[]byte格式传递进来，转换前的json格式如下 (支持一主多从，配置自定义根据自身需求主库从库都能留空)：
```json
 {
    "master": {
        "user": "x",
        "pass": "x",
        "ip": 123,
        "port": 456,
        "db_name": "x",
        "max_igle_conn": 1, // 设置空闲连接池中连接的最大数量
        "max_open_conn": 1, // 设置打开数据库连接的最大数量
        "conn_max_life_time": 1 // 设置了连接可复用的最大时间/秒
        },
    "slave": [
        {
           "user": "x",
           "pass": "x",
           "ip": 123,
           "port": 456,
           "db_name": "x",
           "max_igle_conn": 1, // 设置空闲连接池中连接的最大数量
           "max_open_conn": 1, // 设置打开数据库连接的最大数量
           "conn_max_life_time": 1 // 设置了连接可复用的最大时间/秒
        }
    ]
 }
```
##### 实例的初始化会在主从库各自第一次调用GetMaster()和GetSlave()时进行，初始化成功后连接池相关配置也会建立，后续再调用以上两个方法就会直接从连接池种获取对应实例的长连接以供应用层使用了


# <- redis ->
### *基于redigo封装的一套redis组件*
##### 使用逻辑方式与上面的mysql大致相同(初始化，获取主从实例方法都一样)，也是基于redigo实现的一套连接池，只不过json配置不同，以下给出对应的json格式：
```json
 {
    "master": {
        "host": "",
        "port": "",
        "user": "", // 兼容redis 6.0支持用户名登录
        "pass": "",
        "auth": "",
        "db": 0,
        "network": "tcp",
        "max_idle": 1, // 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
        "max_active": 1, // 最大的激活连接数，表示同时最多有N个连接
        "idle_timeout": 1 // 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭 / 秒
        },
    "slave": [
        {
           "host": "x",
           "port": "x",
           "user": "x", // 兼容redis 6.0支持用户名登录
           "pass": "x",
           "auth": "",
           "db": 0,
           "network": "tcp",
           "max_idle": 1, // 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
           "max_active": 1, // 最大的激活连接数，表示同时最多有N个连接
           "idle_timeout": 1 // 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭 / 秒
        },
    ]
 }
```
##### 注意: 获取到的连接使用结束后需要应用层手动调用Close()方法放回连接池，不然会导致连接池内存泄漏
##### 获取到连接后的cmd操作参考redigo文档：https://pkg.go.dev/github.com/gomodule/redigo/redis#hdr-Connections

# <- es ->
### *基于 olivere/elastic 封装的一套es组件*
##### 使用逻辑方式与上面的mysql/redis大致相同(初始化，es不需要区分主从服务，只需要配置里把集群机器的ip:port都配置好)，以下给出对应的json格式：
```json
 {
    "address": [
      "http://127.0.0.1:9200",
      "http://127.0.0.2:9200",
      "http://127.0.0.3:9200"
    ]
 }
```
##### 获取到连接后的cmd操作参考文档：https://pkg.go.dev/github.com/olivere/elastic#section-readme


# <- configController ->
### *实现本地配置文件热更新*
##### 支持本地配置文件的热更新（修改配置文件目录或配置文件内容）
##### 项目初始化时必须引入包初始化，这将会初始化配置任务的监听者（实现热更新），就像这样：
```
    _ "github.com/GrumpyOmer/gotools/configController"
```
目前只支持json和xml两种配置文件格式,初始化后可通过对应的方法获取
```
    GetJsonField(field string)
    GetXmlField(field string)
```
配置的热更新除了支持手动修改配置文件内容，还支持在程序动态运行过程中调用以下方法修改配置文件名称或目录来触发更新
```
    // 自定义配置目录
    SetPubDir(path string)
    // 自定义xml文件名
    SetXmlConfigName(configName string)
    // 自定义json文件名
    SetJsonConfigName(configName string)
```
配置文件内容格式如下: 
```
    [json] :
    `
    {
      "test1": "1",
      "test2": "2"
    }
    `
    
    [xml] :
    `
    <?xml version="1.0" encoding="UTF-8" ?>
    <root>              //xml配置文件必须以root作为根元素，否则无法解析
        <test1>1</test1>
        <test2>2</test2>
    </root>
    `
    
    以上成功序列化后为以下map
    map[string]string{
        "test1": "1",
        "test2": "2",
    }
```
# <- download ->
### *文件下载组件*
##### 文件下载组件，支持实时下载进度查询
```
    实时下载进度查询,内部用协程异步处理
    wc:= NewWc()
    wc.DownloadFile(filepath string, url string) error
    filepath: 保存路径+文件名
    url: 远程图片url
    error: 只需要关注下err是否返回当前wc已在使用中 （一个wc对象同时只能支持一个文件下载，批量同时下载可以创建多个wc对象）
```
图片文件下载过程中可以通过以下方法实时观察进度
```
    获取文件下载百分比 （%）
    wc.FetchProgress() int32 
    获取文件下载实际大小 （Bytes）
    wc.FetchCurrent() uint64 
    获取文件总大小 （Bytes）
    wc.FetchTotal() uint64
    获取文件是否下载成功 （bool, error message）
    wc.DownloadRes() (bool, error)
```
# 更新记录
# v1.0.7 
##### feat: 协程池新增批量生成任务方法 (BatchMakeTask)
##### 这种方法有一种好处就是，在当前协程数量已经饱和的情况下，批量生成任务中溢出来的任务不会直接向调用方抛出错误，而是通过一个新增的通道保存下来，然后与其他下发的任务一起竞争协程资源

# v1.0.10
##### feat: 新增热更新配置文件组件
##### 目前版本支持两种配置文件格式（json, xml），默认指定配置文件位于程序启动目录下的config目录下的config.xml和config.json，可通过SetPubDir()，SetXmlConfigName()，SetJsonConfigName()方法自定义它们
##### （比如配置文件位于./app/config/jsonConfig.json，可通过分别调用SetPubDir("./app/config")以及SetJsonConfigName("jsonConfig.json")，则可以使插件自己找新配置文件信息并后续维护它们，xml配置方式雷同，使用SetXmlConfigName()方法）

# v1.0.11
##### feat: 文件下载（实时获取下载信息）
##### pref: 协程池优化

# v1.0.12
##### fix: 文件下载重复使用wc属性未重置bug修复

# v1.0.13
##### feat: 新增es组件

# v1.0.14
##### fix: 修复es组件初始化失败bug

# v1.0.15
##### fix: 修复一些小问题

# v1.0.16
##### fix: 修复一些小问题

# v1.0.17
##### optimize: 优化download组件使用方法

# v1.0.18
##### optimize: 引入sync.once优化（redis, mysql, es）组件初始化