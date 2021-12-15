package toolib

import (
	"github.com/go-redis/redis"
)

func NewRedisClient(addr, password string, dbNum int) (*redis.Client, error) {
	red := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbNum,
	})
	if err := red.Ping().Err(); err != nil {
		return nil, err
	}
	return red, nil
}
