# gotools
##### Language: [🇨🇳](./README.md) | 🇺🇸
[![OSCS Status](https://www.oscs1024.com/platform/badge/GrumpyOmer/gotools.svg?size=small)](https://www.oscs1024.com/project/GrumpyOmer/gotools?ref=badge_small)
#### Note: Some productivity boosting gadgets are packaged during Golang use
#### （golang 1.16 and above is recommended, partly depending on the version limitation）
## introduction
## Installation
```
    go get github.com/GrumpyOmer/gotools@latest
```
## Module List
「  
&ensp;&ensp;&ensp;[goroutine](#goroutine)  
&ensp;&ensp;&ensp;[mysql](#mysql)  
&ensp;&ensp;&ensp;[redis](#redis)  
&ensp;&ensp;&ensp;[es](#es)  
&ensp;&ensp;&ensp;[configController](#configController)  
&ensp;&ensp;&ensp;[download](#download)  
&ensp;&ensp;&ensp;[logCenter](#logCenter)  
」
# <span id="goroutine"><- goroutine -></span>
### *Coroutine pooling based on task distribution*
##### When introduced, the first line is initialized, which will initialize the consumer thread of the coroutine task
```
    _ "github.com/GrumpyOmer/gotools/goroutine"
```
##### You need to set the maximum number of coroutines as required by calling the SetGoroutineNumber(uint64) method

```
    SetGoroutineNumber(CPUCoreNum)
```
##### This is done by calling the MakeTask(func(), * sync.waitgroup) error method. The business logic needs to be encapsulated in the func() function, and the generated WaitGroup which controls the business logic is passed in as parameters. This method will respond with an error to tell the application layer whether the task is successfully thrown into the coroutine pool (the task execution is asynchronous). If there is a business need, the external must use the Wait() method of the Wait group to Wait for the end of the business execution.
##### After the execution is complete, you need to obtain the error inside the func() function to make the judgment at the service level and achieve the expected effect at the service level. The following is an example:

```
    var (
        a int
        b int
        err1 error
        err2 error
        w = sync.WaitGroup{}
        tempFunc1 = func() {
            b = a+1
            //Self-pass errors and catch errors
            err1 = errors.New("this is a demo")
        }
        
        tempFunc2 = func() {
            b = a+1
            //Self-pass errors and catch errors
            err2 = errors.New("this is a demo")
        }
    )
    w.Add(2)
    if test1:= goroutine.MakeTask(tempFunc1, &w); test1 != nil {
        // Issued by the failure
        w.Done()
        fmt.Println(test1)
    }
    
    if test2:= goroutine.MakeTask(tempFunc2, &w); test2 != nil {
        // Issued by the failure
        w.Done()
        fmt.Println(test2)
    }
    
    w.Wait()
    fmt.Println(a)
    fmt.Println(b)
    //The ERRs that are passed into the coroutine pool along with func() are captured by themselves
    fmt.Println(err1)
    fmt.Println(err2)
```

# <span id="mysql"><- mysql -></span>
### *A set of mysql components packaged based on GORM*
##### Before using the DB configuration, you need to use the ConfigInit([]byte) method to transfer the DB configuration to the converted []byte format based on the JSON string. The converted JSON format is as follows (one master and multiple slaves are supported, and the master and slave libraries can be left blank according to the configuration requirements) :
```json
 {
    "master": {
        "user": "x",
        "pass": "x",
        "ip": "127.0.0.1",
        "port": "456",
        "db_name": "x",
        "max_igle_conn": 1, // Set the maximum number of connections in the free connection pool
        "max_open_conn": 1, // Set the maximum number of open database connections
        "conn_max_life_time": 1 // The maximum time/second that a connection can be reused is set
        },
    "slave": [
        {
           "user": "x",
           "pass": "x",
           "ip": "127.0.0.1",
           "port": "456",
           "db_name": "x",
           "max_igle_conn": 1, // Set the maximum number of connections in the free connection pool
           "max_open_conn": 1, // Set the maximum number of open database connections
           "conn_max_life_time": 1 // The maximum time/second that a connection can be reused is set
        }
    ]
 }
```
##### The instance is initialized when the master and slave libraries call GetMaster() and GetSlave() for the first time. After the initialization is successful, the connection pool configuration is established. Subsequent calls of the two methods will directly obtain the long connection of the corresponding instance from the connection pool for use by the application layer


# <span id="redis"><- redis -></span>
### *A set of Redis components packaged based on Redigo*
##### Using the same logic as mysql above (initialization, master/slave instance method is the same), also based on Redigo implementation of a set of connection pooling, but the JSON configuration is different, the corresponding JSON format is given below:
```json
 {
    "master": {
        "host": "",
        "port": "",
        "user": "", // Compatible with Redis 6.0 supports user name login
        "pass": "",
        "auth": "",
        "db": 0,
        "network": "tcp",
        "max_idle": 1, // The maximum number of idle connections indicates that N idle connections can be maintained even if there is no REDIS connection.
        "max_active": 1, // Maximum number of active connections: indicates that there are at most N connections at the same time
        "idle_timeout": 1 // The maximum idle connection wait time beyond which the idle connection will be closed/SEC
        },
    "slave": [
        {
           "host": "x",
           "port": "x",
           "user": "x", // Compatible with Redis 6.0 supports user name login
           "pass": "x",
           "auth": "",
           "db": 0,
           "network": "tcp",
           "max_idle": 1, // The maximum number of idle connections indicates that N idle connections can be maintained even if there is no REDIS connection.
           "max_active": 1, // Maximum number of active connections: indicates that there are at most N connections at the same time
           "idle_timeout": 1 // The maximum idle connection wait time beyond which the idle connection will be closed/SEC
        },
    ]
 }
```
##### Note: After the connection is used, the application layer needs to manually call the Close() method to put it back into the connection pool. Otherwise, the connection pool memory will leak
##### Access to the connection of CMD operation reference redigo document: https://pkg.go.dev/github.com/gomodule/redigo/redis#hdr-Connections

# <span id="es"><- es -></span>
### *A suite of ES components based on Olivere/Elastic packaging*
##### The logical method is roughly the same as mysql/redis above (initialization, ES does not need to distinguish between master and slave services, only need to configure the cluster machine IP :port). The corresponding JSON format is given as follows:
```json
 {
    "address": [
      "http://127.0.0.1:9200",
      "http://127.0.0.2:9200",
      "http://127.0.0.3:9200"
    ]
 }
```
##### Access to the connection of CMD operating reference documentation: https://pkg.go.dev/github.com/olivere/elastic#section-readme for the connection of CMD operation reference documentation


# <span id="configController"><- configController -></span>
### *Implement hot update of the local configuration file*
##### Support hot update of local configuration files (modify configuration file directory or configuration file contents)
##### The package initialization must be introduced when the project is initialized, which will initialize the listener of the configuration task (implementing the hot update), like this:
```
    _ "github.com/GrumpyOmer/gotools/configController"
```
Currently, only JSON, XML and .env configuration files are supported. After initialization, you can obtain the configuration file in the corresponding method
```
    GetJsonField(field string)
    GetXmlField(field string)
    GetEnvField(field string)
```
In addition to manually modifying the content of the configuration file, the hot update of the configuration file can also trigger the update by calling the following methods to change the configuration file name or directory during the dynamic running of the program
```
    // Customize the configuration directory
    SetPubDir(path string)
    // Customize the xml file name
    SetXmlConfigName(configName string)
    // Customize the json file name
    SetJsonConfigName(configName string)
    // Customize the .env file name
    SetEnvConfigName(configName string)
```
The format of the configuration file is as follows:
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
    <root>              //The XML configuration file must have root as the root element or it cannot be parsed
        <test1>1</test1>
        <test2>2</test2>
    </root>
    `
    
    [.env] :
    `
    test1="1"
    test2="2"
    `
    The above was successfully serialized to the following map
    map[string]string{
        "test1": "1",
        "test2": "2",
    }
```
# <span id="download"><- download -></span>
### *File Download Component*
##### File download component, supports real-time download progress query
```
    Real-time download progress query, internal asynchronous processing with coroutine
    wc:= NewWc()
    wc.DownloadFile(filepath string, url string) error
    filepath: Save path + file name
    url: Remote image URL
    error: Only need to pay attention to whether the ERR returns that the current WC is in use (one WC object can only support one file download at the same time, batch download at the same time can create multiple WC objects).
```
The following methods can be used to observe the progress of the image file in real time
```
    Get the file download percentage （%）
    wc.FetchProgress() int32 
    Get the actual file download size （Bytes）
    wc.FetchCurrent() uint64 
    Gets the total file size （Bytes）
    wc.FetchTotal() uint64
    Check whether the file is downloaded successfully （bool, error message）
    wc.DownloadRes() (bool, error)
```
# <span id="logCenter"><- logCenter -></span>
### *Log component*

##### Records the log information in code in a detailed format

##### For the convenience of the current version, there is only one way to save logs as written files in the log directory of the startup file directory (files are divided according to the date). Components will be added in later (Es, Mysql...).

```

{

"Content":"omer test1",                                                     // Indicates the log Content

"Location" : "D: / server/gotools/logCenter/config_init_test go: 9",        // log print position (file/line)

"FunctionName" : "github.com/GrumpyOmer/gotools/logCenter.TestConfigInit",  // log print function name

"LogTime":"2022-09-19 18:24:05.5941501 +0800 CST m=+0.006999901"            // Indicates the log printing time

}

```

It is very simple to use, just need to call Add(string) function to pass the log content, the coroutine will be asynchronous to record these logs, do not worry about the impact of writing logs on service performance

```

Add("omer test1")

```
# tag update record
# v1.0.1
##### feat: feature release
# v1.0.2
##### feat: add logCenter
# v1.1.1
##### feat: support .env file by configController  
# v1.1.3
##### fix: fix logCenter memory leak of fd
# v1.1.4
##### refactor: shield configuration check prompt
