package connection

import (
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func RedisConnection(config RedisConfig) *redis.Client {

	redisOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	}

	if config.TLSEnable {
		redisOptions.TLSConfig = &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		}
	}

	rdb := redis.NewClient(redisOptions)

	return rdb
}
