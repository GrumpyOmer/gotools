package configController

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

type(
	configStruct struct {
		xmlConfig map[string]string
		jsonConfig map[string]string
		sync.RWMutex	// 保证读取修改配置内存安全性
	}
)
var(
	Client configStruct
	//默认配置文件目录
	pubPath  = "./config"
	xmlConfigName = "config.xml"
	jsonConfigName = "config.json"
	l sync.Mutex
	UpdatePubPath  = make(chan struct{},1)
	UpdateXmlConfigName  = make(chan struct{})
	UpdateJsonConfigName  = make(chan struct{})

)

// 初始化配置信息
func init() {
	// 初始化
	UpdatePubPath<- struct{}{}
	go func() {
		for {
			select {
					// 公共配置目录修改事件
				case <-UpdatePubPath:
					l.Lock()
					defer l.Unlock()
					if !Exists(pubPath) {
						continue
					}
					if Exists(pubPath+"/"+xmlConfigName) {
						Client.initXmlConfig(pubPath+"/"+xmlConfigName)
					}
					if Exists(pubPath+"/"+jsonConfigName) {
						Client.initJsonConfig(pubPath+"/"+jsonConfigName)
					}
					// json配置文件修改事件
				case <-UpdateJsonConfigName:
					l.Lock()
					defer l.Unlock()
					if Exists(pubPath+"/"+jsonConfigName) {
						Client.initJsonConfig(pubPath+"/"+jsonConfigName)
					}
					// xml配置文件修改事件
				case <-UpdateXmlConfigName:
					l.Lock()
					defer l.Unlock()
					if Exists(pubPath+"/"+jsonConfigName) {
						Client.initJsonConfig(pubPath+"/"+jsonConfigName)
					}
			}
		}
	}()
}

// 自定义配置目录
func SetPubDir(path string) {
	l.Lock()
	defer func() {
		l.Unlock()
		// 触发更新配置信息事件
		UpdatePubPath<- struct{}{}
	}()
	pubPath = path
}

// 自定义xml文件名
func SetXmlConfigName(configName string) {
	l.Lock()
	defer func() {
		l.Unlock()
		// 触发更新配置信息事件
		UpdateXmlConfigName<- struct{}{}
	}()
	xmlConfigName = configName
}

// 自定义json文件名
func SetJsonConfigName(configName string) {
	l.Lock()
	defer func() {
		l.Unlock()
		// 触发更新配置信息事件
		UpdateJsonConfigName<- struct{}{}
	}()
	jsonConfigName = configName
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func(c *configStruct) initXmlConfig(path string) error {
	c.Lock()
	defer c.Unlock()
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("load xml conf failed: "+err.Error())
	}
	if err:= xml.Unmarshal(buf, &c.xmlConfig); err != nil {
		return err
	}
	return nil
}

func GetXmlField(field string) string {
	Client.RLock()
	defer Client.RUnlock()
	return Client.xmlConfig[field]
}

func(c *configStruct) initJsonConfig(path string) error {
	c.Lock()
	defer c.Unlock()
	buf, err := ioutil.ReadFile(path)

	if err != nil {
		return errors.New("load json conf failed: "+err.Error())
	}
	if err:= json.Unmarshal(buf, &c.jsonConfig); err != nil {
		return err
	}
	return nil
}

func GetJsonField(field string) string {
	Client.RLock()
	defer Client.RUnlock()
	return Client.jsonConfig[field]
}

