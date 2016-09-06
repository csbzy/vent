package redisapi

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

var (
	redisSubChanChecker RedisSubChanChecker
)

func createPubSubConnPool() *redis.Pool {
	return CreateRedisPool(":6379", 10, 3, true)
}

func setUpServerMsgChecker() {
	redisSubChanChecker = NewRedisSubChanChecker(createPubSubConnPool, "notify_queue", checkingFunc)
	redisSubChanChecker.Start()
}
func stopServerMsgChecker() {
	go redisSubChanChecker.Stop()
}

type MyData struct {
	//
	Name string
}

func checkingFunc(dataFromRedis interface{}) {
	var value MyData // anyting
	err := json.Unmarshal(dataFromRedis.([]byte), &value)
	if err != nil {
		return
	}
	if value.Name == "foo" {
		// do sth
	}
}
