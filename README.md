### 7天用Go从零实现Web框架 - Gee

Gee 是一个模仿 [gin](https://github.com/gin-gonic/gin) 实现的 Web 框架。

- 第一天：前置知识(http.Handler接口) | [Code](gee-web/day1-http-base)
- 第二天：上下文设计(Context) | [Code](gee-web/day2-context)
- 第三天：Trie树路由(Router) | [Code](gee-web/day3-router)
- 第四天：分组控制(Group) | [Code](gee-web/day4-group)
- 第五天：中间件(Middleware) | [Code](gee-web/day5-middleware)
- 第六天：HTML模板(Template) | [Code](gee-web/day6-template)
- 第七天：错误恢复(Panic Recover) | [Code](gee-web/day7-panic-recover)

 ### 7天用Go从零实现ORM框架 GeeORM

GeeORM 是一个模仿 [gorm](https://github.com/jinzhu/gorm) 和 [xorm](https://github.com/go-xorm/xorm) 的 ORM 框架
geeorm 接口设计上主要参考了 xorm，一些细节实现上参考了 gorm。  
gorm 目前支持的特性有：  
表的创建、删除、迁移。  
记录的增删查改，查询条件的链式操作。  
单一主键的设置(primary key)。  
钩子(在创建/更新/删除/查找之前或之后)  
事务(transaction)。  

- 第一天：database/sql 基础 | [Code](gee-orm/day1-database-sql)
- 第二天：对象表结构映射 | [Code](gee-orm/day2-对象表结构映射)
- 第三天：记录新增和查询 | [Code](gee-orm/day3-记录新增和查询)
- 第四天：链式操作与更新删除 | [Code](gee-orm/day4-链式操作与更新删除)
- 第五天：实现钩子(Hooks) | [Code](gee-orm/day5-实现钩子)
- 第六天：支持事务(Transaction) | [Code](gee-orm/day6-支持事务)
- 第七天：数据库迁移(Migrate) | [Code](gee-orm/day7-数据库迁移)

### 7天用Go从零实现分布式缓存 GeeCache
 
GeeCache 是一个模仿 [groupcache](https://github.com/golang/groupcache) 实现的分布式缓存系统  
支持特性有：  

单机缓存和基于 HTTP 的分布式缓存  
最近最少访问(Least Recently Used, LRU) 缓存策略  
使用 Go 锁机制防止缓存击穿  
使用一致性哈希选择节点，实现负载均衡  
使用 protobuf 优化节点间二进制通信  

- 第一天：LRU 缓存淘汰策略 | [Code](gee-cache/day1-lru)
- 第二天：单机并发缓存 | [Code](gee-cache/day2-单机并发缓存)
- 第三天：HTTP 服务端 | [Code](gee-cache/day3-HTTP服务端)
- 第四天：一致性哈希(Hash) | [Code](gee-cache/day4-一致性哈希)
- 第五天：分布式节点 | [Code](gee-cache/day5-分布式节点)
- 第六天：防止缓存击穿 | [Code](gee-cache/day6-防止缓存击穿)
- 第七天：使用 Protobuf 通信 | [Code](gee-cache/day7-使用Protobuf通信)




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
