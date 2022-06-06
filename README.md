# 仿抖音项目服务端

### 一.项目文档:

1. **项目文档：**[题目一：抖音项目【青训营】 - 飞书文档 (feishu.cn)](https://bytedance.feishu.cn/docx/doxcnbgkMy2J0Y3E6ihqrvtHXPg)
2. **接口文档：**[视频流接口 - 抖音极简版 (apifox.cn)](https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18345145)
3. **客户端下载：** [极简抖音App使用说明 - 飞书文档 (feishu.cn)](https://bytedance.feishu.cn/docx/doxcnZd1RWr6Wpd1WVfntGabCFg)

### 二.项目运行说明：

### 项目启动：

**环境准备：** Go 1.7 + MySQL 8.0 + Redis

上传视频提取封面需要额外安装FFmepg：[Download FFmpeg](http://ffmpeg.org/download.html)

启动需先配置数据库连接，在controller/db.go文件内修改数据库连接账号和密码，默认账号为root，密码为123456

```shell
1.进入项目主目录
3.go mod tidy
2.go build
3.启动生成douyin-simple-demo.exe文件
```

### 模拟器调试：

1.下载手机模拟器如夜神模拟器。[夜神安卓模拟器-安卓模拟器电脑版下载_安卓手游模拟器_手机模拟器_官网 (yeshen.com)](https://www.yeshen.com/)

2.下载官方提供的Apk文件 [Docs (feishu.cn)](https://bytedance.feishu.cn/docx/doxcnZd1RWr6Wpd1WVfntGabCFg)（更新比较频繁，下载最新的），打开夜神模拟器后将其拉拽到桌面后自动会下载安装

3.打开安装好的"抖声"APP，双击两下右下角“我的”打开高级配置，在本机输入 ipconfig 查看本地机的Ipv4地址后更改保存重启即可！注意不能用127.0.0.1！

### 功能说明

* 用户数据使用mysql和redis存储，启动时可自行恢复
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可


### 三.项目结构说明：

```bash
.
├── docs                 // 项目文档存放
├── public               // 公用模块，存放视频和mock数据
├── controller           // API逻辑处理
│   └── user.go          // 实现用户登录部分的HTTP解析与响应
│   └── ......           // 各个模块的HTTP解析与响应实现
├── dao                  // 实现数据库操作逻辑
│   └── db.go            // 数据库连接文件，包含数据库配置和表的创建等初始化
│   └── user_db.go       // 实现用户登录信息部分与数据库交互的逻辑
│   └── ......           // 各个模块的数据库交互实现
├── service              // 业务逻辑
│   └── user.go          // 实现用户登录部分业务逻辑实现
│   └── ......           // 各个模块的业务逻辑实现
├── model                // 存放数据结构
├── router.go            // 定义http接口路由
└── main.go              // 项目主文件，入口
```

分层：**controller**处理接口处理HTTP请求，解析后将其转发给**service**层做业务逻辑的处理，**service**层在处理过程中调用**dao**层实现与数据库的交互，处理完成后返回给**controller**层，**controller**生成HTTP响应，**model**存放业务和与数据库交互的数据结构，不同模块的具体业务实现在各层中用不同文件分离。

### 测试

* dao：使用sqltest模拟mysql数据库，测试数据库操作逻辑,参见user_db_test.go
* service：使用Go monkey对service层调用的dao层逻辑进行打桩，以测试业务逻辑是否工作正常，参见user_test.go
* controller：使用httptest模拟http请求，使用测试用数据库进行测试，参见main_test.go

### 四.技术栈：

1. HTTP框架：Gin    相关文档：https://learnku.com/docs/gin-gonic/1.7
2. 持久层：GORM    相关文档：[GORM 指南 | GORM - The fantastic ORM library for Golang, aims to be developer friendly.](https://gorm.cn/zh_CN/docs/index.html)
3. 数据库： MySQL   相关文档：[MySQL 教程 | 菜鸟教程 (runoob.com)](https://www.runoob.com/mysql/mysql-tutorial.html)
4. 协作管理：Git 相关文档：[Git教程 - 廖雪峰的官方网站 (liaoxuefeng.com)](https://www.liaoxuefeng.com/wiki/896043488029600)
