# 仿抖音项目服务端

### 一.项目文档:

1. **项目文档：**[题目一：抖音项目【青训营】 - 飞书文档 (feishu.cn)](https://bytedance.feishu.cn/docx/doxcnbgkMy2J0Y3E6ihqrvtHXPg)
2. **接口文档：**[视频流接口 - 抖音极简版 (apifox.cn)](https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18345145)
3. **客户端下载：** [极简抖音App使用说明 - 飞书文档 (feishu.cn)](https://bytedance.feishu.cn/docx/doxcnZd1RWr6Wpd1WVfntGabCFg)

### 二.项目运行说明：

基于提供的demo项目进行二次开发

启动需先配置数据库连接，controller/db.go文件修改数据库连接账号和密码，默认账号为root，密码为123456

上传视频提取封面需要额外安装FFmepg：[Download FFmpeg](http://ffmpeg.org/download.html)

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

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

### 测试数据

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试

### 三.项目结构说明：

```bash
.
├── config               // 项目配置
│   └── settings.yml     // 项目配置文件，如数据库连接配置
├── docs                 // 项目文档存放
├── public               // 公用模块，存放视频和mock数据
├── controller           // API逻辑处理
│   └── db.go            // 数据库连接文件，包含数据库配置和表的创建
│   └── model.go         // 存放一些公用的模型
├── sql                  // 存放数据库文件
├── temp                 // 临时文件存放
│   └── logs             // 存放生成日志文件
├── router.go            // 定义http接口路由
└── main.go              // 项目主文件，入口
```

分层：**controller**处理接口接受、验证和返回结果，**models**定义与数据库交互实体。当服务复杂可以添加service分离具体业务实现。

### 四.技术栈：

1. HTTP框架：Gin    相关文档：https://learnku.com/docs/gin-gonic/1.7
2. 持久层：GORM    相关文档：[GORM 指南 | GORM - The fantastic ORM library for Golang, aims to be developer friendly.](https://gorm.cn/zh_CN/docs/index.html)
3. 数据库： MySQL   相关文档：[MySQL 教程 | 菜鸟教程 (runoob.com)](https://www.runoob.com/mysql/mysql-tutorial.html)
4. 协作管理：Git 相关文档：[Git教程 - 廖雪峰的官方网站 (liaoxuefeng.com)](https://www.liaoxuefeng.com/wiki/896043488029600)

### 五.项目更新记录：
