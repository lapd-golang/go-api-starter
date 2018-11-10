package redis

import (
	"admin-server/pkg/config"
	"fmt"
	"gopkg.in/redis.v5"
)

func factory(name string) *redis.Client {
	host := config.Conf.GetString("redis." + name + ".host")
	port := config.Conf.GetString("redis." + name + ".port")
	password := config.Conf.GetString("redis." + name + ".password")
	poolSize := config.Conf.GetInt("redis." + name + ".maxactive")
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
