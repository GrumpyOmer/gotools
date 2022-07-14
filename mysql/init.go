package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	client struct {
		Master *gorm.DB
		Slave  []*gorm.DB //支持多从库
		m      sync.Once  //初始化master配置
		s      sync.Once  //初始化slave配置
	}
	// 连接池相关配置
	connPool struct {
		// 设置空闲连接池中连接的最大数量
		MaxIgleConns int `json:"max_igle_conn"`
		// 设置打开数据库连接的最大数量
		MaxOpenConns int `json:"max_open_conn"`
		// 设置了连接可复用的最大时间
		ConnMaxLifetime int `json:"conn_max_life_time"`
	}
	dbConfig struct {
		User   string `json:"user"`
		Pass   string `json:"pass"`
		Ip     string `json:"ip"`
		Port   string `json:"port"`
		DBName string `json:"db_name"`
		connPool
	}
	config struct {
		Master dbConfig   `json:"master"`
		Slave  []dbConfig `json:"slave"`
	}
)

var (
	// 实例对象
	dbClient = client{}
	// 配置对象
	cf = config{}
)

// 主库对象
func (c *client) GetMaster() (*gorm.DB, error) {
	var (
		err error
	)

	// init once
	c.m.Do(func() {
		// get config init connection
		if db, error := initDB(cf.Master); error != nil {
			err = error
		} else {
			c.Master = db
		}
	})

	if err != nil {
		return nil, err
	}

	if c.Master != nil {
		return c.Master, nil
	}

	return nil, errors.New("无可用主库!!")
}

// 从库对象
func (c *client) GetSlave() (*gorm.DB, error) {
	var (
		err error
	)

	// init once
	c.s.Do(func() {
		// get config init connection
		if len(cf.Slave) != 0 {
			for _, v := range cf.Slave {
				if db, error := initDB(v); error != nil {
					err = error
				} else {
					c.Slave = append(c.Slave, db)
				}
			}
		}
	})

	if err != nil {
		return nil, err
	}

	if len(c.Slave) != 0 {
		// 随机选择一个从库
		// seed函数是用来创建随机数的种子,如果不执行该步骤创建的随机数是一样的，因为默认Go会使用一个固定常量值来作为随机种子。
		rand.Seed(time.Now().UnixNano())
		return c.Slave[rand.Intn(len(c.Slave))], nil
	}

	return nil, errors.New("无可用从库!!")
}

// 数据库配置信息初始化
func ConfigInit(c []byte) error {
	// 外部传入json字符串配置
	err := json.Unmarshal(c, &cf)
	if err != nil {
		// 初始化失败
		return err
	}
	return nil
}

// 获取数据库连接实例
func Client() *client {
	return &dbClient
}

// init db
func initDB(config dbConfig) (*gorm.DB, error) {
	var (
		err   error
		db    *gorm.DB
		sqldb *sql.DB
	)
	dsn := config.User +
		":" +
		config.Pass +
		"@tcp(" +
		config.Ip +
		":" +
		config.Port +
		")/" +
		config.DBName +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		// 初始化失败
		return nil, err
	}
	if sqldb, err = db.DB(); err != nil {
		// 初始化失败
		return nil, err
	}
	// 连接池相关配置
	if config.MaxIgleConns != 0 {
		sqldb.SetMaxIdleConns(config.MaxIgleConns)
	}
	if config.MaxOpenConns != 0 {
		sqldb.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.ConnMaxLifetime != 0 {
		sqldb.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	}
	return db, nil
}
