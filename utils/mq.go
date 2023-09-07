package utils

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var Rdb *redis.Client

func InitMq() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:       "newk8s.ferdinandaedth.top:6379",
		DB:         0,
		Password:   "123456",
		PoolSize:   10, // 设置连接池大小
		MaxRetries: 3,  // 设置最大重试次数

	})

}
func Publish(ctx context.Context, channel string, payload string) error {
	var err error
	logrus.Debugf("[Redis] publish [%s]: %s", channel, payload)
	err = Rdb.Publish(ctx, channel, payload).Err()
	if err != nil {
		logrus.Errorf("[Redis] pulish error: %s", err.Error())
		return err
	}
	return err
}

// Subscribe 订阅redis消息
// channel是订阅的目标信道
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Rdb.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		logrus.Errorf("[Redis] subscribe [%s]", channel)
		return "", err
	}
	logrus.Debugf("[Redis] subscribe [%s]: %s", channel, msg.String())
	return msg.Payload, err
}
