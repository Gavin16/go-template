# go-template 项目
定义GO Web 开发项目包模板,按项目一般组织结构封装配置读取,接口定义,中间件使用等


## 包模块
作为模版项目主要包含如下模块
* docs路径: 该部分保存 gin-swagger 生成文档数据
* pck路径: 存放项目主要源代码
  * 缓存访问: 访问缓存接口封装
  * RPC客户端封装: 微服务调用客户端封装
  * API模块: API接口暴露
  * Service模块: 实际业务逻辑代码
  * DAL模块: 数据库访问CRUD 接口封装
  * 配置模块: 读取并使用配置初始化项目
  * 模型定义模块: 数据实体模型定义
  * 工具模块: utils提供编解码,加解密等功能封装
  * 中间件模块: 提供对请求类似切面处理
* cmd路径: 存放指令或者项目启动入口(main.go)
* Makefile文件: 清理、构建项目(拉取依赖,生成文档,启动项目等)


## 依赖组件
| 组件          | 描述       | 版本      |
|:------------|:---------|---------|
| gin         | 轻量级Web框架 | v1.10.0 |
| GORM        | ORM框架    |         |
| go-Redis    | Redis客户端 |         | 
| viper       | 读取配文件工具  | v1.19.0 | 
| gin-swagger | 提供在线接口文档 | v1.6.0  | 
| gods        | 常见集合工具封装 |         |


## 部署使用
### 本地启动
在根路径下直接使用 `make` 指令(环境需要支撑)，即可完成模块依赖下载, API文档生成 和 服务启动。  
若只希望下载依赖和生成API文档，可以执行 `make init`. 之后可以直接启动项目。

也可以手动执行如下指令启动项目
```bash
  go mod tidy
  go get -u github.com/swaggo/swag/cmd/swag
  go run run github.com/swaggo/swag/cmd/swag init
```
然后便可在本地启动项目。


### 服务器部署
宿主机部署
```bash
  cd go-template
  make build
  nohup ./go-template > go-template.log 2>&1 &
```

容器化部署
```bash
  cd go-template
  docker build -t yourRepo/go-template .
  docker run -d --name go-template -p 8799:8799 yourRepo/go-template
```

## 接口文档
服务部署好之后,本地可通过 http://localhost:8000/swagger/index.html 访问