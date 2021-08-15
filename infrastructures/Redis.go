package infrastructures

import (
	"os"

	redis "github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// IRedis interface
type IRedis interface {
	Client() *redis.Client
}

var redisClient *redis.Client

// Redis struct
type Redis struct{}

// Client return new redis client from given address and password
func (r *Redis) Client() *redis.Client {

	if redisClient == nil {
		client := redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis.address"),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		_, err := client.Ping().Result()
		if err != nil {
			log.WithFields(log.Fields{
				"code":  5500,
				"error": err,
			}).Error(err)
			os.Exit(1)
		}
		redisClient = client
	}

	return redisClient
}
