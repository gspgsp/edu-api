package services

import (
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v3"
	"time"
	jwt2 "github.com/dgrijalva/jwt-go"
	"edu_api/models"
	"edu_api/utils"
	"strconv"
)

/**
redis缓存类
 */
//redis配置对象
type RedisConf struct {
	Redis redisYaml
}

//redis配置
type redisYaml struct {
	CacheHost     string `yaml:"cache_host"`
	CachePort     string `yaml:"cache_port"`
	CacheDatabase string `yaml:"cache_database"`
	CacheUsername string `yaml:"cache_username"`
	CachePassword string `yaml:"cache_password"`
}

//获取redis配置对象
func getRedisConf() (redisConf RedisConf, err error) {
	conf := RedisConf{}
	cacheFile, err := ioutil.ReadFile("./src/edu_api/config/redis.yaml")

	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(cacheFile, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

//redis 线程池
func redisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		MaxActive:   1000,
		Dial: func() (redis.Conn, error) {
			conf, err := getRedisConf()
			if err != nil {
				//记录日志
				log.Printf("get conf error:%v", err)
			}

			c, err := redis.Dial("tcp", conf.Redis.CacheHost+":"+conf.Redis.CachePort)
			if err != nil {
				c.Close()
				log.Printf("dial error:%v", err)
			}

			if _, err := c.Do("auth", conf.Redis.CachePassword); err != nil {
				c.Close()
				log.Printf("auth error:%v", err)
			}

			if _, err := c.Do("select", conf.Redis.CacheDatabase); err != nil {
				c.Close()
				log.Printf("select db error:%v", err)
			}

			return c, nil
		},
	}
}

//获取redis实例
func GetRedisConnection() redis.Conn {
	pool := redisPool()
	return pool.Get()
}

//获取redis缓存下的信息:
func GetRedisCache(authHeaderToken string, act, val string) (info string) {
	var j models.JwtClaim
	token := j.ParseToken(authHeaderToken)

	switch value := token.Claims.(jwt2.MapClaims)["id"].(type) {
	case float64:
		conn := GetRedisConnection()
		defer conn.Close()

		key := utils.ContactHashKey([]string{"user:", strconv.FormatFloat(value, 'f', -1, 64)}...)
		info, _ := redis.String(conn.Do(act, key, val))

		return info
	}

	return
}

/**
mongoDb缓存类
 */


