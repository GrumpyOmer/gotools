# gotools
#### 注: 一些本人在golang开发期间项目里自己封装的提升效率的小工具
## 目录简介
#### <- mysql ->
### *基于gorm封装的一套mysql组件*
##### 首先在使用之前,需要通过ConfigInit([]byte)方法把db的配置以基于json字符串转换后的[]byte格式传递进来,转换前的json格式如下 (支持一主多从, 配置自定义根据自身需求主库从库都能留空):
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
##### 实例的初始化会在主从库各自第一次调用GetMaster()和GetSlave()时进行,初始化成功后连接池相关配置也会建立,后续再调用以上两个方法就会直接从连接池种获取对应实例的长连接以供应用层使用了
#### <- goroutine ->
### *自己实现的一套基于任务分发的协程池*
##### 项目初始化时可以顺便引入包初始化,这将会初始化协程任务的消费者,就像这样: 
    _ "github.com/GrumpyOmer/gotools/goroutine"
##### 首先在使用之前我们需要根据实际情况去设置协程的最大可控数量,通过调用SetGoroutineNumber(uint64)方法
##### 之后使用的话需要再向上抽象一层,调用MakeTask(func(*sync.WaitGroup) error, *sync.WaitGroup)方法, 我们需要将在应用层封装的协程的上下文func(*sync.WaitGroup) error函数,以及管控协程的等待组两个参数传递进来, 该方法会响应一个error来告知应用层任务是否成功丢进协程池(任务执行是异步的), 外部可以自己用等待组的Wait()方法去等待协程执行结束,执行完成后,可以通过获取上面传入的协程上下文func(*sync.WaitGroup) error函数的返回值error来做到业务层面的判断,来达到业务层面预期的效果
#### <- redis ->
### *基于redigo封装的一套redis组件*
