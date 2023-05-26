package configController

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sbabiv/xml2map"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type (
	clientStruct struct {
		xmlConfig    map[string]string
		jsonConfig   map[string]string
		envConfig    map[string]string
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
	envConfig                = configStruct{name: ".env"}
	l                        sync.Mutex
	UpdatePubPathChan        = make(chan struct{}, 1)
	UpdateXmlConfigNameChan  = make(chan struct{})
	UpdateJsonConfigNameChan = make(chan struct{})
	UpdateEnvConfigNameChan  = make(chan struct{})
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
				func() {
					l.Lock()
					defer l.Unlock()
					if Exists(pubPath + "/" + xmlConfig.name) {
						if err := Client.initXmlConfig(pubPath + "/" + xmlConfig.name); err != nil {
							fmt.Println("init xml err: ", err.Error())
						}
					}
					if Exists(pubPath + "/" + jsonConfig.name) {
						if err := Client.initJsonConfig(pubPath + "/" + jsonConfig.name); err != nil {
							fmt.Println("init json err: ", err.Error())
						}
					}
					if Exists(pubPath + "/" + envConfig.name) {
						if err := Client.initEnvConfig(pubPath + "/" + envConfig.name); err != nil {
							fmt.Println("init .env err: ", err.Error())
						}
					}
				}()
				// json配置文件修改事件
			case <-UpdateJsonConfigNameChan:
				func() {
					l.Lock()
					defer l.Unlock()
					if Exists(pubPath + "/" + jsonConfig.name) {
						if err := Client.initJsonConfig(pubPath + "/" + jsonConfig.name); err != nil {
							fmt.Println("init json err: ", err.Error())
						}
					}
				}()
				// xml配置文件修改事件
			case <-UpdateXmlConfigNameChan:
				func() {
					l.Lock()
					defer l.Unlock()
					if Exists(pubPath + "/" + xmlConfig.name) {
						if err := Client.initXmlConfig(pubPath + "/" + xmlConfig.name); err != nil {
							fmt.Println("init xml err: ", err.Error())
						}
					}
				}()
				// .env配置文件修改事件
			case <-UpdateEnvConfigNameChan:
				func() {
					l.Lock()
					defer l.Unlock()
					if Exists(pubPath + "/" + envConfig.name) {
						if err := Client.initEnvConfig(pubPath + "/" + envConfig.name); err != nil {
							fmt.Println("init .env err: ", err.Error())
						}
					}
				}()
			}
		}
	}()

	// 内容监听定时事件
	go func() {
		ticker := time.NewTicker(30 * time.Second) //定时检测配置文件是否改动
		defer ticker.Stop()
		for {
			func() {
				// fmt.Println("定期配置更新检查....")
				<-ticker.C
				l.Lock()
				defer l.Unlock()
				// 检查文件是否更新过 （第一次初始化必须重新加载一次）
				if xmlInfo, err := os.Stat(pubPath + "/" + xmlConfig.name); err == nil {
					if xmlInfo.ModTime().Unix() != xmlConfig.modTime {
						xmlConfig.modTime = xmlInfo.ModTime().Unix()
						// 重新根据配置文件生成配置信息
						if err := Client.initXmlConfig(pubPath + "/" + xmlConfig.name); err != nil {
							fmt.Println("init xml err: ", err.Error())
						}
					}
				}

				if jsonInfo, err := os.Stat(pubPath + "/" + jsonConfig.name); err == nil {
					if jsonInfo.ModTime().Unix() != jsonConfig.modTime {
						jsonConfig.modTime = jsonInfo.ModTime().Unix()
						// 重新根据配置文件生成配置信息
						if err := Client.initJsonConfig(pubPath + "/" + jsonConfig.name); err != nil {
							fmt.Println("init json err: ", err.Error())
						}
					}
				}

				if envInfo, err := os.Stat(pubPath + "/" + envConfig.name); err == nil {
					if envInfo.ModTime().Unix() != envConfig.modTime {
						envConfig.modTime = envInfo.ModTime().Unix()
						// 重新根据配置文件生成配置信息
						if err := Client.initEnvConfig(pubPath + "/" + envConfig.name); err != nil {
							fmt.Println("init env err: ", err.Error())
						}
					}
				}
			}()
		}
	}()
}

// SetPubDir 自定义配置目录
func SetPubDir(path string) {
	l.Lock()
	defer func() {
		l.Unlock()
		// 触发更新配置信息事件
		UpdatePubPathChan <- struct{}{}
	}()
	pubPath = path
}

// SetXmlConfigName 自定义xml文件名
func SetXmlConfigName(configName string) {
	l.Lock()
	defer func() {
		l.Unlock()
		// 触发更新配置信息事件
		UpdateXmlConfigNameChan <- struct{}{}
	}()
	xmlConfig.name = configName
}

// SetJsonConfigName 自定义json文件名
func SetJsonConfigName(configName string) {
	l.Lock()
	defer func() {
		l.Unlock()
		// 触发更新配置信息事件
		UpdateJsonConfigNameChan <- struct{}{}
	}()
	jsonConfig.name = configName
}

// SetEnvConfigName 自定义.env文件名
func SetEnvConfigName(configName string) {
	l.Lock()
	defer func() {
		l.Unlock()
		// 触发更新配置信息事件
		UpdateEnvConfigNameChan <- struct{}{}
	}()
	envConfig.name = configName
}

// Exists 判断所给路径文件/文件夹是否存在
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
	decoder := xml2map.NewDecoder(strings.NewReader(string(buf)))
	res, err := decoder.Decode()
	if err != nil {
		return err
	}
	tmp := res["root"].(map[string]interface{})
	//init map
	c.xmlConfig = make(map[string]string)
	for k, v := range tmp {
		c.xmlConfig[k] = v.(string)
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

func (c *clientStruct) initEnvConfig(path string) error {
	c.Lock()
	defer c.Unlock()
	tmpMap, err := godotenv.Read(path)
	if err != nil {
		return errors.New("load .env conf failed: " + err.Error())
	}
	c.envConfig = tmpMap
	return nil
}

func GetEnvField(field string) string {
	Client.RLock()
	defer Client.RUnlock()
	return Client.envConfig[field]
}
