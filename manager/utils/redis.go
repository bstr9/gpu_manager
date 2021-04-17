package redis 

import (
	goredis "github.com/go-redis/redis"
)


type RedisClient struct {
	client 	*goredis.Client
	name 	string
	addr 	string
	db 		int
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error){
	cmd := c.client.Get(key)
	return cmd.Result()
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, timeout time.Duration) (string, error) {
	cmd := c.client.Set(key, value, timeout)

	return cmd.Result()
}

func (c *RedisClient) HMGet(ctx context.Context, key string, values ...string) ([]interface{}, error) {
	cmd := c.client.HMGet(key, values...)
	return cmd.Result()
}
