# gotools
#### 注: 一些本人在golang开发期间项目里自己封装的提升效率的小工具
## 目录简介
#### <- mysql ->
### 基于gorm封装的一套mysql组件
首先在使用之前,需要通过ConfigInit([]byte)方法把db的配置以基于json字符串转换后的[]byte格式传递进来,转换前的json格式如下 (支持一主多从, 配置自定义根据自身需求主库从库都能留空):
###### {
###### &emsp;"master": {
###### &emsp;&emsp;&emsp;"user": "x",
###### &emsp;&emsp;&emsp;"pass": "x",
###### &emsp;&emsp;&emsp;"ip": 123,
###### &emsp;&emsp;&emsp;"port": 456,
###### &emsp;&emsp;&emsp;"db_name": "x",
###### &emsp;&emsp;&emsp;"max_igle_conn": 1, // 设置空闲连接池中连接的最大数量
###### &emsp;&emsp;&emsp;"max_open_conn": 1, // 设置打开数据库连接的最大数量
###### &emsp;&emsp;&emsp;"conn_max_life_time": 1 // 设置了连接可复用的最大时间/秒
######	},
###### &emsp;"slave": [
###### &emsp;&emsp;&emsp;{
###### &emsp;&emsp;&emsp;"user": "x",
###### &emsp;&emsp;&emsp;"pass": "x",
###### &emsp;&emsp;&emsp;"ip": 123,
###### &emsp;&emsp;&emsp;"port": 456,
###### &emsp;&emsp;&emsp;"db_name": "x",
###### &emsp;&emsp;&emsp;"max_igle_conn": 1, // 设置空闲连接池中连接的最大数量
###### &emsp;&emsp;&emsp;"max_open_conn": 1, // 设置打开数据库连接的最大数量
###### &emsp;&emsp;&emsp;"conn_max_life_time": 1 // 设置了连接可复用的最大时间/秒
###### &emsp;&emsp;&emsp;}
###### &emsp;]
###### }

实例的初始化会在主从库各自第一次调用GetMaster()和GetSlave()时进行,初始化成功后连接池相关配置也会建立,后续再调用以上两个方法就会直接从连接池种获取对应实例的长连接以供应用层使用了
#### <- goroutine ->
### 自己实现的一套基于任务分发的协程池

#### <- redis ->
### 基于redigo封装的一套redis组件