
1. rpcService 从consul获取自己需要创建的grpc service信息,内容json如下:
```json
	[
    {"serviceName":"UserService","port":8002,"redisConfig":{"host":"127.0.0.1:6379","dbNum":0}},
    {"serviceName":"relationService","port":8102,"redisConfig":{"host":"127.0.0.1:6379","dbNum":1}},
    {"serviceName":"authService","port":8202,"redisConfig":{"host":"127.0.0.1:6379","dbNum":3}},
    {"serviceName":"captchaService","port":8302,"redisConfig":{"host":"127.0.0.1:6379","dbNum":2}}
    ]
```
通过`serviceName` 指定自己需要调用的`grpc service `有哪些