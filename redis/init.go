package redis

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

type (
	client struct {
		Master *redis.Pool
		Slave  []*redis.Pool //支持多从库
		mu     sync.Mutex
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
		Auth    string `json:"auth"`
		User    string `json:"user"` // redis 6.0支持用户名登录 兼容一下
		Pass    string `json:"pass"`
		Db      int    `json:"db"`
		Network string `json:"network"`
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

// ConfigInit redis配置信息初始化
func ConfigInit(c []byte) error {
	// 外部传入json字符串配置
	var next config
	err := json.Unmarshal(c, &next)
	if err != nil {
		// 初始化失败
		return err
	}
	cf = next
	redisClient.reset()
	return nil
}

// Client 获取数据库连接实例
func Client() *client {
	return &redisClient
}

// GetMaster 获取redis连接 / master
func (c *client) GetMaster() (redis.Conn, error) {
	c.mu.Lock()
	if c.Master == nil {
		c.Master = initPool(cf.Master)
	}
	pool := c.Master
	c.mu.Unlock()

	conn := pool.Get()
	if conn.Err() != nil {
		// 连接不可用
		return nil, conn.Err()
	}

	return conn, nil
}

// GetSlave 获取redis连接 / slave
func (c *client) GetSlave() (redis.Conn, error) {
	c.mu.Lock()
	if len(c.Slave) == 0 && len(cf.Slave) != 0 {
		for _, v := range cf.Slave {
			c.Slave = append(c.Slave, initPool(v))
		}
	}
	slaves := c.Slave
	c.mu.Unlock()

	if len(slaves) != 0 {
		// 随机选择一个从库的连接 (Go 1.20+ rand.Seed已废弃，直接使用rand.Intn即可)
		conn := slaves[rand.Intn(len(slaves))].Get()
		if conn.Err() != nil {
			return nil, conn.Err()
		}
		return conn, nil
	}

	return nil, errors.New("无可用从库!!")
}

func (c *client) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Master != nil {
		_ = c.Master.Close()
		c.Master = nil
	}
	for _, slave := range c.Slave {
		if slave != nil {
			_ = slave.Close()
		}
	}
	c.Slave = nil
}

// 初始化连接池
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
			if cf.Auth != "" {
				_, err = pool.Do("Auth", cf.Auth)
				if err != nil {
					return pool, err
				}
			}
			return pool, err
		},
	}
}
