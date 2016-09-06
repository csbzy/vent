package redisapi

import (
	"fmt"
	"time"
)

type RedisQueueChecker struct {
	Conn         Redis
	QueueName    string
	fun          func(interface{})
	keepChecking bool
}

func NewRedisQueueChecker(redisConn Redis, QueueName string, fun func(interface{})) (checker RedisQueueChecker) {
	checker = RedisQueueChecker{
		Conn:         redisConn,
		QueueName:    QueueName,
		fun:          fun,
		keepChecking: true,
	}
	return
}

func (checker *RedisQueueChecker) Stop() {
	checker.keepChecking = false
}
func (checker *RedisQueueChecker) Start() {
	go checker.run()
}
func (checker *RedisQueueChecker) run() {
	for checker.keepChecking {
		ret, err := checker.Conn.Brpop(checker.QueueName, 0)
		if err != nil {
			fmt.Println("read %s from redis error  ", checker.QueueName)
			time.Sleep(time.Second * 5)
		} else {
			protectFunc(func() { checker.fun(ret) })
		}
	}

}

func protectFunc(fun func()) {
	defer func() {
		recover()
	}()
	fun()
}
