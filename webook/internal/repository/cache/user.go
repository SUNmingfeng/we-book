package cache

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrorKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, userid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
}
type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

// NewUserCache 在外面实例化好cmd再传进来，因为可能需要许多不同的参数
func NewUserCache(cmdable redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmdable,
		expiration: time.Minute * 15,
	}
}

func (c *RedisUserCache) Get(ctx context.Context, userid int64) (domain.User, error) {
	key := c.Key(userid)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (c *RedisUserCache) Key(userid int64) string {
	return fmt.Sprintf("user:info:%d", userid)
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := c.Key(du.Id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}
