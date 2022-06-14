package configController

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type (
	clientStruct struct {
		xmlConfig    map[string]string
		jsonConfig   map[string]string
		sync.RWMutex // 保证读取修改配置内存安全性
	}

	configStruct struct {
		name    string
		modTime int64
	}
)

var (
	Client                   clientStruct
	pubPath                  = "./config" //默认配置文件目录
	xmlConfig                = configStruct{name: "config.xml"}
	jsonConfig               = configStruct{name: "config.json"}
	l                        sync.Mutex
	UpdatePubPathChan        = make(chan struct{}, 1)
	UpdateXmlConfigNameChan  = make(chan struct{})
	UpdateJsonConfigNameChan = make(chan struct{})
)

// 初始化配置信息
func init() {
	// 初始化
	UpdatePubPathChan <- struct{}{}
	// 文件监听事件
	go func() {
		for {
			select {
			// 公共配置目录修改事件
			case <-UpdatePubPathChan:
				l.Lock()
				defer l.Unlock()
				if !Exists(pubPath) {
					continue
				}
				if Exists(pubPath + "/" + xmlConfig.name) {
					Client.initXmlConfig(pubPath + "/" + xmlConfig.name)
				}
				if Exists(pubPath + "/" + jsonConfig.name) {
					Client.initJsonConfig(pubPath + "/" + jsonConfig.name)
				}
				// json配置文件修改事件
			case <-UpdateJsonConfigNameChan:
				l.Lock()
				defer l.Unlock()
				if Exists(pubPath + "/" + jsonConfig.name) {
					Client.initJsonConfig(pubPath + "/" + jsonConfig.name)
				}
				// xml配置文件修改事件
			case <-UpdateXmlConfigNameChan:
				l.Lock()
				defer l.Unlock()
				if Exists(pubPath + "/" + xmlConfig.name) {
					Client.initXmlConfig(pubPath + "/" + xmlConfig.name)
				}
			}
		}
	}()

	// 内容监听定时事件
	go func() {
		ticker := time.NewTicker(30 * time.Second) //定时检测配置文件是否改动
		defer ticker.Stop()
		for {
			<-ticker.C
			l.Lock()
			defer l.Unlock()
			// 检查文件是否更新过 （第一次初始化必须重新加载一次）
			if xmlInfo, err := os.Stat(pubPath + "/" + xmlConfig.name); err == nil {
				if xmlInfo.ModTime().Unix() != xmlConfig.modTime {
					xmlConfig.modTime = xmlInfo.ModTime().Unix()
					// 重新根据配置文件生成配置信息
					Client.initXmlConfig(pubPath + "/" + xmlConfig.name)
				}
			}

			if jsonInfo, err := os.Stat(pubPath + "/" + jsonConfig.name); err == nil {
				if jsonInfo.ModTime().Unix() != jsonConfig.modTime {
					jsonConfig.modTime = jsonInfo.ModTime().Unix()
					// 重新根据配置文件生成配置信息
					Client.initJsonConfig(pubPath + "/" + jsonConfig.name)
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
		UpdatePubPathChan <- struct{}{}
	}()
	pubPath = path
}

// 自定义xml文件名
func SetXmlConfigName(configName string) {
	l.Lock()
	defer func() {
		l.Unlock()
		// 触发更新配置信息事件
		UpdateXmlConfigNameChan <- struct{}{}
	}()
	xmlConfig.name = configName
}

// 自定义json文件名
func SetJsonConfigName(configName string) {
	l.Lock()
	defer func() {
		l.Unlock()
		// 触发更新配置信息事件
		UpdateJsonConfigNameChan <- struct{}{}
	}()
	jsonConfig.name = configName
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (c *clientStruct) initXmlConfig(path string) error {
	c.Lock()
	defer c.Unlock()
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("load xml conf failed: " + err.Error())
	}
	if err := xml.Unmarshal(buf, &c.xmlConfig); err != nil {
		return err
	}
	return nil
}

func GetXmlField(field string) string {
	Client.RLock()
	defer Client.RUnlock()
	return Client.xmlConfig[field]
}

func (c *clientStruct) initJsonConfig(path string) error {
	c.Lock()
	defer c.Unlock()
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("load json conf failed: " + err.Error())
	}
	if err := json.Unmarshal(buf, &c.jsonConfig); err != nil {
		return err
	}
	return nil
}

func GetJsonField(field string) string {
	Client.RLock()
	defer Client.RUnlock()
	return Client.jsonConfig[field]
}
