package redisapi

import (
	"fmt"
	"strconv"
)

type SetRedis interface {
	Sadd(key string, value ...interface{}) error
	Scard(key string) (int, error)
	SisMember(key string, value interface{}) bool
	Smembers(key string) (members []interface{}, err error)
	SmembersAsString(key string) (members []string, err error)
	Spop(key string) (interface{}, error)
	SpopAsString(key string) (string, error)
	Srem(key string, value ...interface{}) error
	SrandMember(key string) (interface{}, error)
	SrandMemberAsString(key string) (string, error)
}

type OrderSetRedis interface {
	Zadd(key string, score interface{}, value interface{}) error
	ZScore(key string, value interface{}) (int, error)
	ZScoreAsFloat64(key string, value interface{}) (float64, error)
	ZaddBatch(key string, list []ScoreStruct) error
	Zcard(key string) (int, error)
	ZRrange(key string, begin int, end int) ([]ScoreStruct, error)
	ZRevRrange(key string, begin int, end int) ([]ScoreStruct, error)
	ZRrank(key string, value interface{}) (int, error)
	ZRevRrank(key string, value interface{}) (int, error)
	Zrem(key string, value interface{}) error
	ZRemRangeByRank(key string, begin int, end int) error
	ZRemRangeByScore(key string, min, max interface{}) error
}

type HashRedis interface {
	Hexist(table, key string) bool
	Hgetall(table string) ([]string, error)
	Hdel(table string, key string) error
	Hset(table, key string, value interface{}) error
	HMset(table string, scoreList []ScoreStruct) error
	Hget(table, key string) (interface{}, error)
	HMget(table string, keys ...string) ([]ScoreStruct, error)
	Hkeys(table string) ([]string, error)
	Hlen(table string) (int, error)
}

type QueueRedis interface {
	Lpush(key string, value interface{}) error
	Rpush(key string, value interface{}) error
	Lrange(key string, start, end int) ([]interface{}, error)
	Rpop(key string) (interface{}, error)
	Llen(key string) (int, error)
	Lset(key string, index int, value interface{}) error
	Ltrim(key string, start, end int) error
	Brpop(key string, timeoutSecs int) (interface{}, error)
	Lrem(key string, value interface{}, remType int) error
}

type Redis interface {
	QueueRedis
	SetRedis
	OrderSetRedis
	HashRedis
	Ping() bool
	Exists(key string) bool

	Set(key string, value []byte) error

	Get(key string) ([]byte, error)

	Delete(key string) error

	Incr(key string, step uint64) (int64, error)

	Decr(key string, step uint64) (int64, error)

	Expire(key string, sec int64) error

	Ttl(key string) (int32, error)

	Persist(key string) error
	
	Keys(key string) ([]string, error)

	MultiGet(keys []interface{}) ([]interface{}, error)

	MultiSet(kvMap map[string][]byte) error

	ClearAll() error

	Pub(key string, value interface{}) error
	Sub(keys ...string) ([]string, error)
	UnSub(keys ...string) error
}
type ScoreStruct struct {
	Member interface{}
	Score  interface{}
}

func (ss ScoreStruct) GetMember() interface{} {
	return ss.Member
}
func (ss ScoreStruct) GetScore() interface{} {
	return ss.Score
}

func (ss ScoreStruct) GetMemberAsString() string {
	return string(ss.Member.([]uint8))
}
func (ss ScoreStruct) GetMemberAsUint64() (member uint64) {
	var i int
	i, _ = strconv.Atoi(string(ss.Member.([]uint8)))
	return uint64(i)
}

func (ss ScoreStruct) GetScoreAsInt() (score int) {
	var e error
	score, e = strconv.Atoi(string(ss.Score.([]uint8)))
	if e != nil {
		f, e := strconv.ParseFloat((string(ss.Score.([]uint8))), 10)
		if e != nil {
			fmt.Println("redisapi.GetScoreAsInt error", e)
		}
		score = int(f)
	}
	return
}
