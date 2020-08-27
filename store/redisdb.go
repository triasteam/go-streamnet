package store

import(
	"log"
	streamnet_conf "github.com/triasteam/go-streamnet/config"
	"github.com/go-redis/redis"
)

func NewRedisStore() *redisStore {
	return &redisStore{}
}

type redisStore struct {}

var redisdb *redis.Client

func init(){
	redisdb = redis.NewClient(&redis.Options{
        Addr:     streamnet_conf.EnvConfig.Redis.Url, 				// use default Addr
        Password: streamnet_conf.EnvConfig.Redis.Password,          // no password set
        DB:       streamnet_conf.EnvConfig.Redis.DB,               // use default DB
	})
	pong, err := redisdb.Ping().Result();
	log.Println(pong, err);
}

func (redis_store *redisStore) Set(key string, val string) string{
	result,err := redisdb.Set(key,val,0).Result()

	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (redis_store *redisStore) Get(key string) string{
	result,err := redisdb.Get(key).Result()

	if err != nil {
		log.Fatal(err)
	}
	return result
}