package lib

import (
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectRedis() {
	opt, err := redis.ParseURL(EnvConfig.RedisUrl)

	if err != nil {
		log.Fatalln("Unable to connect to redis:", err)
	}

	Redis = redis.NewClient(opt)
}
