package redisapi

import (
	"errors"
	"reflect"
	"strings"
	"time"
	"github.com/garyburd/redigo/redis"
)

type RedisClient struct {
	pool *redis.Pool
	addr string
}


//REDIS GEO
type GeoOption func(* geoOption)
type geoOption struct {
	distanceUnit string // m | km | mi | ft
	sort string   // ASC | DESC
	withOpt string // WITHDIST | WITHCOORD | WITHHASH
}
type Coordinate struct {
	Longitude interface{}
	Latitude interface{}
}

type GeoRadiusEle struct {
	value interface{}
	distance interface{}
	coordinae Coordinate
}
const(
	WITHDIST="WITHDIST"
	WITHCOORD="WITHCOORD"
	WITHBOth = " WITHDIST WITHCOORD"
	//WITHHASH = "WITHHASH"
	ASC = "ASC"
	DESC = "DESC"
	DistanceUnitM = "m"
	DistanceUnitKM = "km"
	DistanceUnitFT = "ft"
	DistanceUnitMI = "mi"
)

func SetDistancUnit(distanceUnit string ) func (*geoOption){
	var distanceUnitTmp string
	if (distanceUnit == DistanceUnitM || distanceUnit == DistanceUnitKM || distanceUnit == DistanceUnitMI || distanceUnit == DistanceUnitFT){
		distanceUnitTmp = distanceUnit
	}else {
		distanceUnitTmp = DistanceUnitM
	}
	return func(o *geoOption){
		o.distanceUnit = distanceUnitTmp
	}
}

func SetSort(sort string) func (* geoOption){
	if (sort == ASC || sort == DESC ){
		return func(o *geoOption){
			o.sort = sort
		}
	}
	return func(o *geoOption){
	}
}

func SetWith(with string) func (* geoOption){
	if(with == WITHDIST || with == WITHCOORD || with == WITHBOth){
		return func(o *geoOption){
			o.withOpt = with
		}
	}
	return func(o *geoOption){

	}
}

func (rc RedisClient) GeoAdd(key string  ,c Coordinate ,value string) error{
	conn := rc.connectInit()
	defer conn.Close()
	_, err := conn.Do("GEOADD",key ,c.Longitude,c.Latitude,value)

	return err
}

func (rc RedisClient) GeoPos(key string,value string) ( *Coordinate,error){
	conn := rc.connectInit()
	defer conn.Close()
	v,err := redis.Values(conn.Do("GEOPOS",key,value))
	if (err != nil || len(v) != 2){
		return nil,err
	}

	return &Coordinate{
		Latitude:v[0],
		Longitude:v[1],
	},nil
}

func (rc RedisClient) GeoDist(key string,value1 ,value2 ,distancUnit string) (float64,error){
	conn := rc.connectInit()
	defer conn.Close()
	v,err := conn.Do("GEODIST",key,value1,value2,distancUnit)
	if err != nil{
		return 0,err
	}
	return redis.Float64(v,nil)
}

func (rc RedisClient) GeoRadius(key string, c Coordinate,maxDis float64,opt...GeoOption) (*[]GeoRadiusEle,error){
	conn := rc.connectInit()
	defer conn.Close()
	geoOpt := &geoOption{}
	for _,optFun := range opt{
		optFun(geoOpt)
	}
	var v []interface{}
	var err error
	if geoOpt.withOpt == WITHBOth{
		v,err = redis.Values(conn.Do("GEORADIUS",key,c.Longitude,c.Latitude,maxDis,geoOpt.distanceUnit,WITHCOORD,WITHDIST,geoOpt.sort))
	}else{
		v,err = redis.Values(conn.Do("GEORADIUS",key,c.Longitude,c.Latitude,maxDis,geoOpt.distanceUnit,geoOpt.withOpt,geoOpt.sort))
	}
	if err != nil {
		return nil,err
	}
	geoRadius := makeGeoRadiusEle(v,geoOpt)
	return geoRadius,nil
}

func (rc RedisClient)GeoRadiusByMember(key string,member string, maxDis float64,opt ...GeoOption) (*[]GeoRadiusEle,error){
	conn := rc.connectInit()
	defer conn.Close()
	geoOpt := &geoOption{}
	for _,optFun := range opt{
		optFun(geoOpt)
	}
	var v []interface{}
	var err error
	if geoOpt.withOpt == WITHBOth{
		v,err = redis.Values(conn.Do("GEORADIUSBYMEMBER",key,member,maxDis,geoOpt.distanceUnit,WITHCOORD,WITHDIST,geoOpt.sort))
	}else{
		v,err = redis.Values(conn.Do("GEORADIUSBYMEMBER",key,member,maxDis,geoOpt.distanceUnit,geoOpt.withOpt,geoOpt.sort))
	}
	if err != nil {
		return nil,err
	}
	geoRadius := makeGeoRadiusEle(v,geoOpt)
	return geoRadius,nil
}

func makeGeoRadiusEle(v []interface{},geoOpt *geoOption) *[]GeoRadiusEle{
	length := len(v)
	geoRadius := make([]GeoRadiusEle,length)
	switch geoOpt.withOpt {
	case WITHDIST:
		for i:=0;i < length;i++{
			vv,err := redis.Values(v[i],nil)
			if err  != nil {
				break
			}
			geoRadius[i].value = vv[0]
			geoRadius[i].distance = vv[1]
		}
	case WITHCOORD:
		for i:=0;i<length;i++{
			vv,err := redis.Values(v[i],nil)
			if err  != nil {
				break
			}
			geoRadius[i].value = vv[0]
			vvv,err := redis.Values(vv[1],nil)
			if err  != nil {
				break
			}
			coordinate := Coordinate{Longitude:vvv[0] ,Latitude:vvv[1]}
			geoRadius[i].coordinae = coordinate
		}
	case WITHBOth:
		for i:=0;i<length;i++{
			vv,err := redis.Values(v[i],nil)
			if err  != nil {
				break
			}
			geoRadius[i].value = vv[0]
			geoRadius[i].value = vv[0]
			geoRadius[i].distance = vv[1]
			vvv,err := redis.Values(vv[2],nil)
			if err  != nil {
				break
			}
			coordinate := Coordinate{Longitude:vvv[0] ,Latitude:vvv[1]}
			geoRadius[i].coordinae = coordinate
		}
	}
	return &geoRadius
}

func (rc RedisClient) Exists(key string) bool {
	conn := rc.connectInit()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (rc RedisClient) Lpush(key string, value interface{}) error {
	conn := rc.connectInit()
	defer conn.Close()

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		s := reflect.ValueOf(value)
		values := make([]interface{}, s.Len()+1)
		values[0] = key
		for i := 1; i <= s.Len(); i++ {
			values[i] = s.Index(i - 1).Interface()
		}
		_, err := conn.Do("LPUSH", values...)
		return err
	} else {
		_, err := conn.Do("LPUSH", key, value)
		return err
	}
}

func (rc RedisClient) Rpush(key string, value interface{}) error {
	conn := rc.connectInit()
	defer conn.Close()

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		s := reflect.ValueOf(value)
		values := make([]interface{}, s.Len()+1)
		values[0] = key
		for i := 1; i <= s.Len(); i++ {
			values[i] = s.Index(i - 1).Interface()
		}
		_, err := conn.Do("RPUSH", values...)
		return err
	} else {
		_, err := conn.Do("RPUSH", key, value)
		return err
	}
}

func (rc RedisClient) Lrange(key string, start, end int) ([]interface{}, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("LRANGE", key, start, end)
	return v.([]interface{}), err
}
func (rc RedisClient) Rpop(key string) (interface{}, error) {
	conn := rc.connectInit()
	defer conn.Close()

	value, err := conn.Do("RPOP", key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	return value, err
}

func (rc RedisClient) Llen(key string) (int, error) {
	conn := rc.connectInit()
	defer conn.Close()

	length, err := redis.Int(conn.Do("LLEN", key))
	if err != nil {
		return 0, err
	}
	return length, nil
}

func (this RedisClient) Lset(key string, index int, value interface{}) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("LSET", key, index, value)
	return err
}

func (this RedisClient) Ltrim(key string, start, end int) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("LTRIM", key, start, end)
	return err
}

func (rc RedisClient) Brpop(key string, timeoutSecs int) (interface{}, error) {
	conn := rc.connectInit()
	defer conn.Close()

	var val interface{}
	var err error
	if timeoutSecs < 0 {
		val, err = conn.Do("BRPOP", key, 0)
	} else {
		val, err = conn.Do("BRPOP", key, timeoutSecs)
	}
	values, err := redis.Values(val, err)
	if err != nil {
		return nil, err
	}
	return string(values[1].([]byte)), err
}

func (rc RedisClient) Lrem(key string, value interface{}, remType int) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("LREM", key, remType, value)
	return err
}

func (rc RedisClient) Set(key string, value []byte) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

func (rc RedisClient) Get(key string) ([]byte, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("GET", key)
	if err != nil || v == nil {
		return nil, err
	}
	return v.([]byte), err
}

func (rc RedisClient) Delete(key string) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (rc RedisClient) Incr(key string, step uint64) (int64, error) {
	conn := rc.connectInit()
	defer conn.Close()

	value, err := conn.Do("INCRBY", key, step)
	if err != nil {
		return 0, nil
	}
	return value.(int64), err
}

func (rc RedisClient) Decr(key string, step uint64) (int64, error) {
	conn := rc.connectInit()
	defer conn.Close()

	value, err := conn.Do("DECRBY", key, step)
	if err != nil {
		return 0, nil
	}
	return value.(int64), err
}

func (rc RedisClient) Expire(key string, sec int64) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, sec)
	return err
}

func (rc RedisClient) Ttl(key string) (int32, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := redis.Int(conn.Do("TTL", key))
	return int32(v), err
}

func (rc RedisClient) Persist(key string) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("PERSIST", key)
	return err
}

func (rc RedisClient) Keys(key string) ([]string, error) {
	conn := rc.connectInit()
	defer conn.Close()
	
	vList, err := redis.Strings(conn.Do("KEYS", key))
	return vList, err
}

func (rc RedisClient) MultiGet(keys []interface{}) ([]interface{}, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("MGET", keys...)
	return v.([]interface{}), err
}

func (rc RedisClient) MultiSet(kvMap map[string][]byte) error {
	conn := rc.connectInit()
	defer conn.Close()

	var values []interface{}
	for key, value := range kvMap {
		values = append(values, key)
		values = append(values, value)
	}

	_, err := conn.Do("MSET", values...)
	return err
}

func (rc RedisClient) ClearAll() error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("FLUSHALL")
	return err
}

func (this RedisClient) Sadd(key string, value ...interface{}) error {
	conn := this.connectInit()
	defer conn.Close()

	params := []interface{}{key}
	params = append(params, value...)
	_, err := conn.Do("SADD", params...)
	return err
}

func (this RedisClient) Scard(key string) (int, error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Int(conn.Do("SCARD", key))
	return int(v), err
}

func (rc RedisClient) SisMember(key string, value interface{}) bool {
	conn := rc.connectInit()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("SIsMember", key, value))
	if err != nil {
		return false
	}
	return v
}

func (rc RedisClient) Smembers(key string) (members []interface{}, err error) {
	conn := rc.connectInit()
	defer conn.Close()

	members, err = redis.Values(conn.Do("SMembers", key))
	return
}
func (rc RedisClient) SmembersAsString(key string) (members []string, err error) {
	ms, err := rc.Smembers(key)
	if err != nil {
		return
	}
	members = make([]string, len(ms), len(ms))
	for idx, m := range ms {
		members[idx] = string(m.([]uint8))
	}
	return
}

func (this RedisClient) Spop(key string) (value interface{}, err error) {
	conn := this.connectInit()
	defer conn.Close()

	return conn.Do("SPop", key)

}
func (this RedisClient) SpopAsString(key string) (value string, err error) {
	return redis.String(this.Spop(key))

}
func (this RedisClient) SrandMember(key string) (value interface{}, err error) {
	conn := this.connectInit()
	defer conn.Close()

	return conn.Do("SRandMember", key)

}
func (this RedisClient) SrandMemberAsString(key string) (value string, err error) {
	conn := this.connectInit()
	defer conn.Close()

	return redis.String(conn.Do("SRandMember", key))

}
func (this RedisClient) Srem(key string, value ...interface{}) error {
	conn := this.connectInit()
	defer conn.Close()

	params := []interface{}{key}
	params = append(params, value...)

	_, err := conn.Do("SRem", params...)
	return err
}

// order set begin
func (this RedisClient) Zadd(key string, score interface{}, value interface{}) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("ZADD", key, score, value)
	return err
}

func (this RedisClient) ZaddBatch(key string, list []ScoreStruct) error {
	conn := this.connectInit()
	defer conn.Close()
	var cmdArgs []interface{} = make([]interface{}, 2*len(list)+1)
	cmdArgs[0] = key
	var idx int = 1
	for _, scoreStruct := range list {
		cmdArgs[idx] = scoreStruct.GetScore()
		idx++
		cmdArgs[idx] = scoreStruct.GetMember()
		idx++

	}
	_, err := conn.Do("ZADD", cmdArgs...)
	return err
}

func (this RedisClient) ZScore(key string, value interface{}) (int, error) {
	conn := this.connectInit()
	defer conn.Close()

	return redis.Int(conn.Do("ZSCORE", key, value))

}

func (this RedisClient) ZScoreAsFloat64(key string, value interface{}) (float64, error) {
	conn := this.connectInit()
	defer conn.Close()

	return redis.Float64(conn.Do("ZSCORE", key, value))
}

func (this RedisClient) Zrem(key string, value interface{}) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("ZREM", key, value)
	return err
}

func (this RedisClient) ZRemRangeByRank(key string, start, end int) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("ZREMRANGEBYRANK", key, start, end)
	return err
}
func (this RedisClient) ZRemRangeByScore(key string, min, max interface{}) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("ZREMRANGEBYSCORE", key, min, max)
	return err
}

func (this RedisClient) Zcard(key string) (int, error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Int(conn.Do("ZCARD", key))
	return int(v), err
}

func (this RedisClient) ZRrank(key string, value interface{}) (int, error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Int(conn.Do("ZRANK", key, value))
	return int(v), err
}

func (this RedisClient) ZRrange(key string, begin int, end int) (scoreStructList []ScoreStruct, err error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Values(conn.Do("ZRANGE", key, begin, end, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	length := len(v)
	scoreStructList = make([]ScoreStruct, length/2)
	for i := 0; i < length/2; i++ {
		scoreStructList[i].Member = v[i*2]
		scoreStructList[i].Score = v[i*2+1]
	}
	return
}

func (this RedisClient) ZRevRrank(key string, value interface{}) (int, error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Int(conn.Do("ZREVRANK", key, value))
	return int(v), err
}

func (this RedisClient) ZRevRrange(key string, begin int, end int) (scoreStructList []ScoreStruct, err error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Values(conn.Do("ZREVRANGE", key, begin, end, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	length := len(v)
	scoreStructList = make([]ScoreStruct, length/2)
	for i := 0; i < length/2; i++ {
		scoreStructList[i].Member = v[i*2]
		scoreStructList[i].Score = v[i*2+1]
	}
	return
}

// order set end

// hash set begin

func (rc RedisClient) Hexist(table, key string) bool {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := redis.Bool(conn.Do("HEXISTS", table, key))
	if err != nil {
		return false
	}
	return v
}

func (rc RedisClient) Hgetall(table string) ([]string, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := redis.Strings(conn.Do("HGETALL", table))
	return v, err
}

func (rc RedisClient) Hkeys(table string) ([]string, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := redis.Strings(conn.Do("HKEYS", table))
	return v, err
}

func (rc RedisClient) Hlen(table string) (int, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := redis.Int(conn.Do("HLEN", table))
	return v, err
}

func (rc RedisClient) Hdel(table, key string) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("HDEL", table, key)
	return err
}

func (rc RedisClient) Hset(table, key string, value interface{}) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("HSET", table, key, value)
	return err
}

func (rc RedisClient) HMset(table string, scoreList []ScoreStruct) error {
	conn := rc.connectInit()
	defer conn.Close()

	var interList []interface{}
	interList = append(interList, table)
	for _, score := range scoreList {
		interList = append(interList, score.GetMember())
		interList = append(interList, score.GetScore())
	}
	_, err := conn.Do("HMSET", interList...)
	return err
}

func (rc RedisClient) Hget(table, key string) (interface{}, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("HGET", table, key)
	if v == nil {
		return nil, err
	} else {
		return v.(interface{}), err
	}
}

func (rc RedisClient) HMget(table string, keys ...string) ([]ScoreStruct, error) {
	conn := rc.connectInit()
	defer conn.Close()

	argList := make([]interface{}, len(keys)+1)
	argList[0] = table
	for i := 1; i <= len(keys); i++ {
		argList[i] = keys[i-1]
	}
	var scoreList []ScoreStruct
	values, err := redis.Strings(conn.Do("HMGET", argList...))
	if err != nil {
		return scoreList, err
	}

	if len(keys) != len(values) {
		return scoreList, errors.New("count error for hmget result.")
	}

	length := len(keys)
	scoreList = make([]ScoreStruct, length)
	for i := 0; i < length; i++ {
		scoreList[i].Member = keys[i]
		scoreList[i].Score = values[i]
	}
	return scoreList, nil
}

// hash set end

func (rc RedisClient) Pub(key string, value interface{}) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("PUBLISH", key, value)
	return err
}

func (rc RedisClient) Sub(keys ...string) ([]string, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("SUBSCRIBE", keys)
	value := string(v.([]interface{})[1].([]byte))
	values := strings.Split(value[1:len(value)-1], " ")
	return values, err
}

func (rc RedisClient) UnSub(keys ...string) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("UNSUBSCRIBE", keys)
	return err
}

func (rc RedisClient) connectInit() redis.Conn {
	conn := rc.pool.Get()
	return conn
}
func (rc RedisClient) Ping() bool {
	conn := rc.connectInit()
	defer conn.Close()
	pong, err := conn.Do("PING")
	if err != nil {
		// fmt.Println(err)
		return false
	}
	return pong == "PONG"
}
func (rc RedisClient) DisConnect() {
	rc.pool.Close()
}

func CreateRedisPool(addr string,DBNum, maxActive, maxIdle int, wait bool) (pool *redis.Pool) {
	msgRedisConfig := addr
	pool = &redis.Pool{
		MaxActive:   maxActive,
		MaxIdle:     maxIdle,
		Wait:        wait,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", msgRedisConfig)
			if err != nil {
				return nil, err
			}
			_,err =c.Do("SELECT",DBNum)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return
}

func InitRedisClient(addr string, DB,MaxActive, MaxIdle int, Wait bool) (*RedisClient, error) {
	pool := CreateRedisPool(addr,DB, MaxActive, MaxIdle, Wait)
	return &RedisClient{pool, addr}, nil
}

func InitDefaultClient(addr string) (Redis, error) {
	return InitRedisClient(addr, 0,0, 3, true)
}



