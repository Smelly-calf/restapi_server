# restapi_server

基于 Restful api 和 PostgreSQL 的用户关系服务。

#### Restful api 风格规范：
1. 基于统一资源标识操作
2. 使用标准方法对资源进行增删改查。
3. 无状态：每个RESTful API请求都包含了所有足够完成本次操作的信息，服务端无需保持 Session。

#### 技术选型
对象关系型数据库 PostgreSQL 和 Gin 框架。

#### 模块设计
开发、测试、构建。

REST 错误响应格式规范。

TODO：并发请求测试。

