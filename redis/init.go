package redis

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
)

type (
	client struct {
		Master *redis.Pool
		Slave  []*redis.Pool //支持多从库
	}
	// 连接池配置
	poolConfig struct {
		MaxIdle     int `json:"max_idle"`     // 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxActive   int `json:"max_active"`   // 最大的激活连接数，表示同时最多有N个连接
		IdleTimeout int `json:"idle_timeout"` // 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭 / 秒
	}
	// 连接配置
	redisConfig struct {
		Host    string `json:"host"`
		Port    string `json:"port"`
		User    string `json:"user"` // redis 6.0支持用户名登录 兼容一下
		Pass    string `json:"pass"`
		Db      int    `json:"db"`
		Network string `json:"network`
		poolConfig
	}
	config struct {
		Master redisConfig   `json:"master"`
		Slave  []redisConfig `json:"slave"`
	}
)

var (
	cf          = config{}
	redisClient = client{}
)

// redis配置信息初始化
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
	return &redisClient
}

// 获取redis连接 / master
func (c *client) GetMaster() (*redis.Conn, error) {
	var (
		conn redis.Conn
	)
	// 连接池未始化
	if c.Master == nil {
		// get config init master connPool
		c.Master = initPool(cf.Master)
	}
	conn = c.Master.Get()
	if conn.Err() != nil {
		// 连接不可用
		return nil, c.Master.Get().Err()
	}
	return &conn, nil
}

// 获取redis连接 / slave
func (c *client) GetSlave() (*redis.Conn, error) {
	var (
		conn redis.Conn
	)
	if len(c.Slave) != 0 {
		// 随机选择一个从库的连接
		// seed函数是用来创建随机数的种子,如果不执行该步骤创建的随机数是一样的，因为默认Go会使用一个固定常量值来作为随机种子。
		rand.Seed(time.Now().UnixNano())
		conn = c.Slave[rand.Intn(len(c.Slave))].Get()
		if conn.Err() != nil {
			return nil, conn.Err()
		}
		return &conn, nil
	} else {
		// get config init connection
		if len(cf.Slave) != 0 {
			for _, v := range cf.Slave {
				c.Slave = append(c.Slave, initPool(v))
			}
			conn = c.Slave[0].Get()
			if conn.Err() != nil {
				return nil, conn.Err()
			}
			return &conn, nil
		}
	}
	return &conn, errors.New("无可用从库!!")
}

//初始化连接池
func initPool(cf redisConfig) *redis.Pool {
	var (
		DialOptionSlice []redis.DialOption
	)
	if cf.Pass != "" {
		DialOptionSlice = append(DialOptionSlice, redis.DialPassword(cf.Pass))
	}
	if cf.User != "" {
		DialOptionSlice = append(DialOptionSlice, redis.DialUsername(cf.User))
	}
	return &redis.Pool{
		MaxIdle:     cf.MaxIdle,
		MaxActive:   cf.MaxActive,
		IdleTimeout: time.Duration(cf.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			pool, err := redis.Dial(cf.Network, cf.Host+":"+cf.Port, DialOptionSlice...)
			if err != nil {
				return pool, err
			}
			return pool, err
		},
	}
}
