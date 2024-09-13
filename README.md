# gin-enhance
## gin 增强（包括swagger文档生成、自定义注解等）

### 项目说明：
- 通过注解增强，可实现类似java的控制器开发方式；
- 减少代码重复编写，提高开发效率；
- 一键生成swagger文档；
- 支持自定义注解 可通过自定义注解的方式，实现日志记录、权限控制等功能；

### 范例：
- [github.com/celt237/gin-enhance.demo](https://github.com/celt237/gin-enhance.demo)

### 支持的注解标签：
- @zService
    - 用于生成该service对应的handler和router代码，可标注在service结构体上，以及对应的swagger文档注释
    - 格式：// @zService
- @zResult
    - 可标注在方法及service结构体上，标注接口要返回的外层结构体类型，方法上标注优先级高于service结构体上标注
    - 格式：// @zResult 返回值类型
    - 例：// @zResult model.Result
- @zSummary
    - 标注接口的简要描述，可标注在方法上 同swagger的@Summary
    - 格式：// @zSummary 接口简要描述
    - 例：// @zSummary 获取用户信息
- @zDescription
    - 标注接口的详细描述，可标注在方法上 同swagger的@Description
    - 格式：// @zDescription 接口详细描述
    - 例：// @zDescription 根据用户id获取用户信息
- @zTags
    - 标注接口的标签，可标注在方法上 同swagger的@Tags
    - 格式：// @zTags 标签1,标签2
    - 例：// @zTags 用户,信息
- @zParam
    - 标注接口的参数，可标注在方法上 同swagger的@Param
    - 格式：// @zParam 参数名 参数请求方式 参数类型 是否必须 参数描述
    - 例：// @zParam id path int true "用户id"
- @zResultData
    - 标注接口的返回数据类型，可标注在方法上 默认使用方法返回的类型，如果方法类型与返回类型不一致则需要标注
    - 格式：// @zResultData 数据类型
    - 例：// @zResultData model.User "用户信息
- @zAccept
    - 标注接口的请求类型，可标注在方法上，默认值为：application/json 同swagger的@Accept
    - 格式：// @zAccept 请求类型
    - 例：// @zAccept application/json
- @zProduce
    - 标注接口的响应类型，可标注在方法上，默认值为：application/json 同swagger的@Produce
    - 格式：// @zProduce 响应类型
    - 例：// @zProduce application/json
- @zRouter
    - 标注接口的路由及请求类型，可标注在方法上 同swagger的@Router
    - 格式：// @zRouter 路由 [请求类型]
    - 例：// @zRouter /api/v1/user/{id} [get]

### 命令使用方式：
- 1、安装命令
```shell
go get github.com/celt237/gin-enhance/cmd/gin-enhance@latest
```

- 2、在项目中添加 github.com/celt237/gin-enhance 依赖


- 3、运行
```shell
gin-enhance --servicePath=xxx --handlePath=xxx --result=xxx --errorCode=xxx
```
- 参数说明：
    - servicePath：service文件目录路径 例如：./service
    - handlePath：handle文件目录路径 例如：./handle
    - result：返回值类型（zResult同时也支持通过结构体或方法注释中的@zResult标签进行设定）
    - errorCode：错误码，默认500
- 如果运行上面命令找不到 `gin-enhance` , 将 GOPATH 下的 bin 加入到环境变量中, 修改完后重启 idea
```shell
export PATH=$PATH:$GOPATH/bin
```



### MIT LICENSE
[LICENSE](./LICENSE)

### Links
- knife4j：https://github.com/xiaomin/knife4j
- go-annotation: https://github.com/celt237/go-annotation