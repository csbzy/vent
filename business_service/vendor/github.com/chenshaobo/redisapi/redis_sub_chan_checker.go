package redisapi

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

// sub pub 消息checker
type RedisSubChanChecker struct {
	SubChanName     string
	fun             func(interface{})
	keepChecking    bool
	CreateRedisPool func() *redis.Pool
	Conn            redis.PubSubConn
}

func NewRedisSubChanChecker(createPubSubConnPoolFunc func() *redis.Pool, SubChanName string, fun func(interface{})) (checker RedisSubChanChecker) {
	checker = RedisSubChanChecker{
		CreateRedisPool: createPubSubConnPoolFunc,
		Conn:            redis.PubSubConn{Conn: createPubSubConnPoolFunc().Get()},
		SubChanName:     SubChanName,
		fun:             fun,
		keepChecking:    true,
	}
	return
}

func (checker *RedisSubChanChecker) Stop() {
	checker.keepChecking = false
	// checker.Conn.Unsubscribe(checker.SubChanName)
	go checker.Conn.Close()
}
func (checker *RedisSubChanChecker) Start() {
	go checker.run()
	checker.Conn.Subscribe(checker.SubChanName)
}
func (checker *RedisSubChanChecker) run() {
	for checker.keepChecking {
		switch n := checker.Conn.Receive().(type) {
		case redis.Message:
			//  care this
			protectFunc(func() { checker.fun(n.Data) })
		case redis.PMessage:
			// donot care this
			fmt.Printf("PMessage: %s %s %s\n", n.Pattern, n.Channel, n.Data)
		case redis.Subscription:
			// // donot care this
			// fmt.Printf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
			// if n.Count == 0 {
			// 	return
			// }
		case error:
			fmt.Println("read %s from redis error  ", checker.SubChanName)
			time.Sleep(time.Second * 30)
			go checker.Conn.Close()
			checker.Conn = redis.PubSubConn{Conn: checker.CreateRedisPool().Get()}
			checker.Conn.Subscribe(checker.SubChanName)
		}
	}

}
