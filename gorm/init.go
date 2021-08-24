package Gorm

import (
	"gorm.io/gorm"
)

type (
	connect struct {
		master *gorm.DB
		slave  *gorm.DB
	}
	Config struct {
		Master struct {
			User   string
			Pass   string
			Ip     string
			Port   string
			DBName string
		}
		Slave struct {
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

func (c *connect) Master() (*gorm.DB, error) {
	var (
		err error
	)
	if c.master != nil {
		return c.master,nil
	} else {
		 // get config init connection
	}
	return nil,err
}

func (c *connect) Slave() (*gorm.DB, error) {
	var (
		err error
	)
	if c.slave != nil {
		return c.slave,nil
	} else {
		// get config init connection
	}
	return nil,err

}

func ConfigInit(c Config) error {
	// 可支持动态获取db配置 暂时需要外部引入
	// init
	cf = c
	return nil
}

func Connection() *connect {
	return &db
}
