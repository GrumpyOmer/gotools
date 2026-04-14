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
		mu     sync.Mutex
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

// GetMaster 主库对象
func (c *client) GetMaster() (*gorm.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Master != nil {
		return c.Master, nil
	}
	db, err := initDB(cf.Master)
	if err != nil {
		return nil, err
	}
	c.Master = db
	return c.Master, nil
}

// GetSlave 从库对象
func (c *client) GetSlave() (*gorm.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.Slave) != 0 {
		// 随机选择一个从库 (Go 1.20+ rand.Seed已废弃，直接使用rand.Intn即可)
		return c.Slave[rand.Intn(len(c.Slave))], nil
	}
	if len(cf.Slave) == 0 {
		return nil, errors.New("无可用从库!!")
	}
	var lastErr error
	for _, v := range cf.Slave {
		db, err := initDB(v)
		if err != nil {
			lastErr = err
			continue
		}
		c.Slave = append(c.Slave, db)
	}
	if len(c.Slave) != 0 {
		return c.Slave[rand.Intn(len(c.Slave))], nil
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("无可用从库!!")
}

// ConfigInit 数据库配置信息初始化
func ConfigInit(c []byte) error {
	// 外部传入json字符串配置
	var next config
	err := json.Unmarshal(c, &next)
	if err != nil {
		// 初始化失败
		return err
	}
	cf = next
	dbClient.reset()
	return nil
}

// Client 获取数据库连接实例
func Client() *client {
	return &dbClient
}

func (c *client) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Master != nil {
		if sqldb, err := c.Master.DB(); err == nil {
			_ = sqldb.Close()
		}
		c.Master = nil
	}
	for _, slave := range c.Slave {
		if slave == nil {
			continue
		}
		if sqldb, err := slave.DB(); err == nil {
			_ = sqldb.Close()
		}
	}
	c.Slave = nil
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
