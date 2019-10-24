package redis_storage

import (
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"

	"github.com/go-redis/cache/v7"
	"instatasks/config"
	"log"
	"os"
)

var RedisCache *cache.Codec

func InitCache() *cache.Codec {
	var codec = RedisCache
	var options *redis.Options

	config := config.InitConfig()
	appEnv := config.AppEnv

	if appEnv == "staging" {
		redisURL := os.Getenv("REDIS_URL")
		log.Println("REDIS_URL : redisURL")
		options, _ = redis.ParseURL(redisURL)
	} else {
		options = &redis.Options{
			Addr:         config.Redis.Addr,
			Password:     config.Redis.Password,
			DB:           0,
			PoolSize:     config.Redis.PoolSize,
			MinIdleConns: config.Redis.MinIdleConns,
		}
	}

	client := redis.NewClient(options)
	// ring := redis.NewRing(&redis.RingOptions{
	// 	Addrs: map[string]string{
	// 		"redis": ":6379",
	// 		"redis2": ":6380",
	// 	},
	// })

	pong, err := client.Ping().Result()
	log.Println(pong)

	if err != nil {
		panic(err)
		panic("failed to connect database")
	}

	codec = &cache.Codec{
		// Redis: ring,
		Redis: client,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

	RedisCache = codec

	return RedisCache
}

func GetRedisCache() *cache.Codec {
	return RedisCache
}
