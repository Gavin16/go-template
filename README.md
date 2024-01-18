# go-template 项目
定义GO Web CRUD开发项目包模板,包括实用模块定义,实用工具封装,数据访问,远程及微服务调用等


## 包模块
作为模版项目主要包含如下模块
* 数据库访问: 数据库访问SQL封装
* 缓存访问: 访问缓存接口封装
* RPC客户端封装: 微服务调用客户端封装
* API模块: API接口暴露
* Service模块: 实际业务逻辑代码
* DAL模块: 数据库访问CRUD 接口封装
* 配置模块: 读取并使用配置初始化项目
* 模型定义模块: 数据实体模型定义
* 工具模块: utils提供编解码,加解密等功能封装
* docs模块: 该部分保存 gin-swagger 生成文档数据
* 中间件模块: 提供对请求类似切面处理

## 依赖组件
* Gin: 轻量级Web框架
* GORM: ORM框架, 用于数据库连接
* Go-Redis: Go Redis客户端
* godotenv: 环境变量工具,方便获取环境变量
* gin-swagger: 提供在线接口文档(需要执行swag安装: go get -u )
* gods(Go Data Structures) : 提供常见集合工具List, Set, Stack,Tree,Queue的封装


## 接口文档
服务部署好之后,本地可通过 http://localhost:8000/swagger/index.html 访问
