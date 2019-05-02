package cache

import (
	"github.com/aloxc/gappuser/config"
	"github.com/garyburd/redigo/redis"
	"github.com/smallnest/rpcx/log"
	"os"
	"strconv"
	"time"
)

var (
	// 定义常量
	RedisClient *redis.Pool
)

func init() {
	host := os.Getenv(config.REDIS_HOST)
	mi := os.Getenv(config.REDIS_MAXIDLE)
	ma := os.Getenv(config.REDIS_MAXACTIVE)
	if host == "" {
		host = config.REDIS_HOST_DEFAULT
	}
	maxIdle := config.REDIS_MAXIDLE_DEFAULT
	maxActive := config.REDIS_MAXACTIVE_DEFAULT
	var err error
	if mi != "" {
		maxIdle, err = strconv.Atoi(mi)
		if err != nil {
			log.Info("REDIS_MAXIDLE只能是数字")
			os.Exit(1)
		}
	}
	if ma != "" {
		maxActive, err = strconv.Atoi(ma)
		if err != nil {
			log.Info("REDIS_MAXACTIVE只能是数字")
			os.Exit(1)
		}
	}
	// 建立连接池
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			// 选择db
			//c.Do("SELECT", config.REDIS_DB)
			return c, nil
		},
	}
}
