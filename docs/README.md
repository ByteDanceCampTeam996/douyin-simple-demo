# 极简版抖音服务端项目开发文档

### 一.功能实现:

目前极简版抖音项目要求的所有服务端接口功能都已实现完成,主要功能包括：

1. 用户注册与登陆，登陆用户可以查看自己的基本信息。
2. 所有用户都可以首页刷抖音视频，按投稿时间倒序推出。
3. 登录用户可以自己拍视频投稿。发布后可在个人主页查看自己的投稿列表。
4. 登录用户可以对视频点赞，在个人主页能够查看所有点赞的视频列表。
5. 登陆用户可以在视频下进行评论。
6. 登录用户可以关注其他用户，在个人信息页查看自己关注数和粉丝数。点击可以打开关注列表和粉丝列表。

**技术栈：** Gin  + GORM + MySQL + Redis

### 二. MySQL数据库表设计：

**1. db_user_infos 用户信息表**

| 序号 | 字段名        | 数据类型   | 索引类型 | 说明       |
| ---- | ------------- | ---------- | -------- | ---------- |
| 1    | id            | bigint(20) | 主键     | 用户ID     |
| 2    | user_name     | longtext   | 普通索引 | 用户名     |
| 3    | password_hash | longtext   |          | 加密后密码 |
| 4    | token         | longtext   |          | 校验token  |

**2. db_video_infos 视频信息表**

| 序号 | 字段名       | 数据类型    | 索引类型 | 说明              |
| ---- | ------------ | ----------- | -------- | ----------------- |
| 1    | video_id     | bigint(20)  | 自增主键 | 视频ID            |
| 2    | user_id      | bigint(20)  |          | 视频作者ID        |
| 3    | play_url     | longtext    |          | 视频播放URL地址   |
| 4    | cover_url    | longtext    |          | 视频封面图URL地址 |
| 5    | created_time | datetime(3) |          | 视频创建时间      |

**3. db_comments 评论表**

| 序号 | 字段名      | 数据类型   | 索引类型 | 说明       |
| ---- | ----------- | ---------- | -------- | ---------- |
| 1    | id          | bigint(20) | 自增主键 | 评论ID     |
| 2    | vid         | bigint(20) |          | 视频ID     |
| 3    | uid         | bigint(20) |          | 评论用户ID |
| 4    | content     | longtext   |          | 评论内容   |
| 5    | create_date | longtext   |          | 评论时间   |

**4. db_favorites 视频点赞表**

| 序号 | 字段名 | 数据类型   | 索引类型 | 说明                     |
| ---- | ------ | ---------- | -------- | ------------------------ |
| 1    | uid    | bigint(20) | 复合主键 | 用户ID                   |
| 2    | vid    | bigint(20) | 复合主键 | 视频ID                   |
| 3    | status | bigint(20) |          | 是否点赞，0未点赞，1点赞 |

**5. follows 关注表**

| 序号 | 字段名     | 数据类型    | 索引类型 | 说明                          |
| ---- | ---------- | ----------- | -------- | ----------------------------- |
| 1    | user_id    | bigint(20)  | 复合主键 | 用户ID                        |
| 2    | follow_id  | bigint(20)  | 复合主键 | 关注对象ID                    |
| 3    | status     | bigint(20)  |          | 关注状态，0取关，1关注，2互关 |
| 4    | created_at | datetime(3) |          | 创建时间                      |
| 5    | updated_at | datetime(3) |          | 更新时间                      |

**6. user_follow_infos 用户关注和粉丝列表**

| 序号 | 字段名         | 数据类型    | 索引类型 | 说明     |
| ---- | -------------- | ----------- | -------- | -------- |
| 1    | user_id        | bigint(20)  | 主键     | 用户ID   |
| 2    | name           | longtext    |          | 用户名   |
| 3    | follow_count   | bigint(20)  |          | 关注数   |
| 4    | follower_count | bigint(20)  |          | 粉丝数   |
| 5    | created_at     | datatime(3) |          | 创建时间 |
| 6    | updated_at     | datatime(3) |          | 更新时间 |

### 三.代码质量：

1. 代码结构上分为dao,service和controller三层。

- dao层封装针对视频信息表、用户信息表、评论表、点赞表、关注表的单表增删改查操作。
- service层对一个或多个dao进行再次封装为服务，实现具体的业务逻辑。
- controller层负责请求参数的校验，将参数传给Service处理，再统一接口结果返回。

2. 代码规范上使用VsCode自带的Go fmt模块保持代码风格的一致，统一使用中文注释。

### 四.单元测试：

单元测试分为三个部分，分别是对controller，dao和service模块的单元测试。
* 对dao模块的测试使用sqltest模拟mysql数据库，测试数据库操作逻辑,参见user_db_test.go，运行测试使用
```shell
go test -v github.com/ByteDanceCampTeam996/douyin-simple-demo/dao
```
* 对service的测试使用Go monkey进行打桩，以测试业务逻辑是否工作正常，运行测试使用
```shell
go test -v github.com/ByteDanceCampTeam996/douyin-simple-demo/service -gcflags=all=-l
```
* 对controller的测试使用httptest模拟http请求，使用测试用数据库进行测试，参见main_test.go
```shell
go test -v github.com/ByteDanceCampTeam996/douyin-simple-demo
```

### 五.服务性能优化：

* 针对视频列表涉及的多表查询，通过Go协程并行提高查询效率。
* 通过FFmpeg完成视频封面图的高效截取和图片压缩。
* 针对视频评论添加Redis缓存，提高视频评论接口响应速度。

### 六.安全：

* 用户的账号密码信息存储时，使用SHA-256对密码进行哈希，从用户名中提取信息用作盐值，将密码和盐值拼接后哈希，可以有效抵抗服务器端的数据泄露以及简单碰撞，彩虹表等针对用户密码的攻击行为。
* 项目使用Gorm访问数据库，GORM使用database/sql的参数占位符来构造 SQL 语句，这可以自动转义参数，避免SQL注入。但对于用户输入的主键检索等记录由于功能上的需要不支持自动转义，项目中实现了简单的基于正则表达式的sql过滤，检测危险输入中是否包含sql命令以及特殊符号，如果存在，则判定为危险输入，不予执行。
* 对于需要登陆才能访问的接口，均做了token校验。token校验失败的会抛返回错误状态码提示。
* 视频发布接口上传视频格式的校验，文件格式正确的再存储到服务器防止文件攻击。

