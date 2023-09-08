package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "redis:6379"
	RedisKeyPrefix       = "default:"
)

type Cache struct {
	Client     *redis.Client
	JSONHelper *rejson.Handler
	Context    context.Context
	KeyPrefix  string
}

func NewCache(prefix string) (*Cache, error) {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}

	return newCacheInstance(redisUrl, prefix)
}

func newCacheInstance(url string, prefix string) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: url,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	return &Cache{
		Client:     client,
		JSONHelper: jsonHelper,
		Context:    ctx,
		KeyPrefix:  prefix,
	}, nil
}

func IsRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

func (c *Cache) KeyFromId(id uint) string {
	return fmt.Sprintf("%s:%d", c.KeyPrefix, id)
}

func (c *Cache) GetItemFromRedis(key string, item interface{}) error {
	itemObject, err := c.JSONHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return err
	}

	return nil
}
