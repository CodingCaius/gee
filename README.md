### 7天用Go从零实现Web框架 - Gee

Gee 是一个模仿 [gin](https://github.com/gin-gonic/gin) 实现的 Web 框架。

- 第一天：前置知识(http.Handler接口) | [Code](gee-web/day1-http-base)
- 第二天：上下文设计(Context) | [Code](gee-web/day2-context)
- 第三天：Trie树路由(Router) | [Code](gee-web/day3-router)
- 第四天：分组控制(Group) | [Code](gee-web/day4-group)
- 第五天：中间件(Middleware) | [Code](gee-web/day5-middleware)
- 第六天：HTML模板(Template) | [Code](gee-web/day6-template)
- 第七天：错误恢复(Panic Recover) | [Code](gee-web/day7-panic-recover)



### 7天用Go从零实现RPC框架 GeeRPC

GeeRPC 是一个基于 [net/rpc](https://github.com/golang/go/tree/master/src/net/rpc) 开发的 RPC 框架 GeeRPC 是基于 Go 语言标准库 `net/rpc` 实现的，添加了协议交换、服务注册与发现、负载均衡等功能。  
GeeRPC 选择从零实现 Go 语言官方的标准库 net/rpc，并在此基础上，新增了协议交换(protocol exchange)、注册中心(registry)、服务发现(service discovery)、负载均衡(load balance)、超时处理(timeout processing)等特性。分七天完成，最终代码约 1000 行

- 第一天 : 服务端与消息编码 | [Code](https://github.com/geektutu/7days-golang/blob/master/gee-rpc/day1-codec)
- 第二天 : 支持并发与异步的客户端 | [Code](https://github.com/geektutu/7days-golang/blob/master/gee-rpc/day2-client)
- 第三天 : 服务注册 | [Code](https://github.com/geektutu/7days-golang/blob/master/gee-rpc/day3-service)
- 第四天 : 超时处理(timeout) | [Code](https://github.com/geektutu/7days-golang/blob/master/gee-rpc/day4-timeout)
- 第五天 : 支持HTTP协议 | [Code](https://github.com/geektutu/7days-golang/blob/master/gee-rpc/day5-http-debug)
- 第六天 : 负载均衡(load balance) | [Code](https://github.com/geektutu/7days-golang/blob/master/gee-rpc/day6-load-balance)
- 第七天 : 服务发现与注册中心 | [Code](https://github.com/geektutu/7days-golang/blob/master/gee-rpc/day7-registry)
