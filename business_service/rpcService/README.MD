
1. rpcService 从consul获取自己需要创建的grpc service信息,内容json如下:
```json
	[{"serviceName":"registerServer","port":8002,"redisConfig":{"host":"127.0.0.1:6379","dbNum":"0"}}]
```
通过serviceName 指定自己需要调用的grpc service 有哪些