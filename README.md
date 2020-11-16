# restapi_server

基于 Restful api 和 PostgreSQL 的用户关系服务。

#### 环境搭建
使用 go module

go get github.com/xxx 失败，配置代理：
```
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=gosum.io+ce6e7565+AY5qEHUk/qmHc5btzW45JVoENfazw8LielDsaI+lEbq6
```
Goland -> Go Modules： Environment 也要配置。

#### Restful api 风格规范：
1. 基于统一资源标识操作
2. 使用标准方法对资源进行增删改查。
3. 无状态：每个REST 请求都包含了所有信息，服务端无需保持 Session。

#### 技术选型
go 操作 pg：github.com/lib/pq；

REST api 框架：github.com/gin-gonic/gin；

单元测试工具：httpexpect


#### 模块设计
```

.
├── Makefile    ------ 编译/构建命令
├── README.md   
├── bin         ------ 编译后可执行文件
│   └── users
├── common      ------ 基础服务
│   ├── runner.go
│   └── server.go
├── db.sql      ------ postgres DDL 语句
├── go.mod      
├── go.sum
├── main        ------ 程序唯一入口
│   └── main.go
├── pkg         ------ 源码包
│   ├── config      
│   │   ├── config.go
│   │   └── postgres.go  ----- DB 相关配置   
│   ├── model            ----- 数据访问对象层       
│   │   ├── base.go      ----- DB 对象池初始化
│   │   ├── relationships.go    
│   │   └── user.go
│   └── service          ----- 路由和处理函数
│       ├── errors.go
│       ├── handler.go
│       └── router.go
└─ test                 ----- 单元测试模块
    ├── model_test.go
    └── service_test.go

```