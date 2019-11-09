package redis_storage

import (
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"

	"github.com/go-redis/cache/v7"
	"instatasks/config"
	"log"
	"os"
)

type (
	CacheCodec = cache.Codec
	CacheItem  = cache.Item
)


var CacheCodecs = make(map[string]*cache.Codec)

func GetCacheCodec(key string) *cache.Codec {
	return CacheCodecs[key]
}

func GetOptions(db int) *redis.Options {

	var opt *redis.Options

	conf := config.GetConfig()

	if conf.AppEnv == "staging" {
		redisURL := os.Getenv("REDIS_URL")
		opt, _ = redis.ParseURL(redisURL)
		opt.DB = db
	} else {
		opt = &redis.Options{
			Addr:         conf.Redis.Addr,
			Password:     conf.Redis.Password,
			DB:           db,
			PoolSize:     conf.Redis.PoolSize,
			MinIdleConns: conf.Redis.MinIdleConns,
		}
	}

	return opt
}

func GetCodec(client *redis.Client) *cache.Codec {
	return &cache.Codec{

		Redis: client,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func RegisterCacheCodec(key string) *cache.Codec {
	db_count := len(CacheCodecs)

	options := GetOptions(db_count)


	client := redis.NewClient(options)

	pong, err := client.Ping().Result()
	log.Println(pong)

	if err != nil {
		panic(err)
		panic("failed to connect database")
	}

	cdc := GetCodec(client)

	CacheCodecs[key] = cdc

	return cdc
}
