# gotools
#### 注：一些本人在golang开发期间项目里自己封装的提升效率的小工具
## 目录简介
#### <- goroutine ->
### *自己实现的一套基于任务分发的协程池*
##### 项目初始化时可以顺便引入包初始化，这将会初始化协程任务的消费者，就像这样：
    _ "github.com/GrumpyOmer/gotools/goroutine"
##### 首先在使用之前我们需要根据实际情况去设置协程的最大可控数量，通过调用SetGoroutineNumber(uint64)方法
##### 之后使用的话需要再向上抽象一层，调用 MakeTask(func()， *sync.WaitGroup)方法，我们需要将在应用层封装的协程的上下文 func() 函数，以及管控协程的等待组两个参数传递进来，该方法会响应一个error来告知应用层任务是否成功丢进协程池(任务执行是异步的)，外部必须自己用等待组的Wait()方法去等待协程执行结束。
##### 执行完成后，需要通过自行获取上面传入的协程上下文 func() 函数内部的error来做到业务层面的判断，来达到业务层面预期的效果，使用例子如下：
```
    var (
        a int
        b int
        err error
        w = sync.WaitGroup{}
        tempFunc = func() {
            b = a+1
            //自行传入错误和捕捉错误
            err = errors.New("this is a demo")
        }
    )
    w.Add(1)
    goroutine.MakeTask(tempFunc, &w)
    w.Wait()
    fmt.Println(a)
    fmt.Println(b)
    //自行传入错误和捕捉错误
    fmt.Println(err)
```

#### <- mysql ->
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

#### <- redis ->
### *基于redigo封装的一套redis组件*
##### 使用逻辑方式与上面的mysql大致相同(初始化，获取主从实例方法都一样)，也是基于redigo实现的一套连接池，只不过json配置不同，以下给出对应的json格式：
```json
 {
    "master": {
        "host": "x",
        "port": "x",
        "user": "x",
        "pass": "x",
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
           "user": "x",
           "pass": "x",
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

#### <- configController ->
### *实现本地配置文件热更新的组件*
##### 支持本地配置文件的热更新（修改指向配置文件或配置文件内容）


### 更新记录
#### v1.0.7 
##### 协程池新增批量生成任务方法 (BatchMakeTask)
##### 这种方法有一种好处就是，在当前协程数量已经饱和的情况下，批量生成任务中溢出来的任务不会直接向调用方抛出错误，而是通过一个新增的通道保存下来，然后与其他下发的任务一起竞争协程资源

#### v1.0.10
##### 新增热更新配置文件组件
##### 目前版本支持两种配置文件格式（json, xml），默认指定配置文件位于程序启动目录下的config目录下的config.xml和config.json，可通过SetPubDir()，SetXmlConfigName()，SetJsonConfigName()方法自定义它们
##### （比如配置文件位于./app/config/jsonConfig.json，可通过分别调用SetPubDir("./app/config")以及SetJsonConfigName("jsonConfig.json")，则可以使插件自己找新配置文件信息并后续维护它们，xml配置方式雷同，使用SetXmlConfigName()方法）
