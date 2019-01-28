package redis

import (
	"fmt"
	"go-admin-starter/utils/config"
	"gopkg.in/redis.v5"
)

var conf = config.New()

func factory(name string) *redis.Client {
	host := conf.GetString("redis." + name + ".host")
	port := conf.GetString("redis." + name + ".port")
	password := conf.GetString("redis." + name + ".password")
	poolSize := conf.GetInt("redis." + name + ".maxactive")
	fmt.Printf("conf-redis: %s:%s - %s\r\n", host, port, password)

	address := fmt.Sprintf("%s:%s", host, port)
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
		PoolSize: poolSize,
	})
}

/**
 * 获取连接
 */
func getRedis(name string) *redis.Client {
	return factory(name)
}

/**
 * 获取master连接
 */
func Master() *redis.Client {
	return getRedis("master")
}

/**
 * 获取slave连接
 */
func Slave() *redis.Client {
	return getRedis("slave")
}
