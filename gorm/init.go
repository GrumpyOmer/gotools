package Gorm

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type (
	connect struct {
		master *gorm.DB
		slave  []*gorm.DB //支持多从库
	}
	Config struct {
		Master struct {
			User   string
			Pass   string
			Ip     string
			Port   string
			DBName string
		}
		Slave []struct {
			User   string
			Pass   string
			Ip     string
			Port   string
			DBName string
		}
	}
)

var (
	// 实例对象
	db = connect{}
	// 配置对象
	cf = Config{}
)

// 主库对象
func (c *connect) Master() (*gorm.DB, error) {
	var (
		err error
	)
	if c.master != nil {
		return c.master, nil
	} else {
		// get config init connection
	}
	return nil, err
}

// 从库对象
func (c *connect) Slave() (*gorm.DB, error) {
	var (
		err error
	)
	if c.slave != nil && len(c.slave) != 0 {
		// 随机选择一个从库
		// seed函数是用来创建随机数的种子,如果不执行该步骤创建的随机数是一样的，因为默认Go会使用一个固定常量值来作为随机种子。
		rand.Seed(time.Now().UnixNano())
		return c.slave[rand.Intn(len(c.slave))], nil
	} else {
		// get config init connection
	}
	return nil, err
}

// 数据库配置信息初始化
func ConfigInit(c Config) error {
	// 可支持动态获取db配置 暂时需要外部引入
	// init
	cf = c
	return nil
}

// 获取数据库连接实例
func Connection() *connect {
	return &db
}
