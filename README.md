# 云文档

## 介绍

一个简单的 Markdown、OpenAPI 文档服务器

## 软件架构

- go1.21
- sqlite
- vue3

## 编译步骤

- 安装 go 环境
- 安装 nodejs 环境
- 进入目录，添加权限并执行 **build.bat** 或 **build.sh**
- 编译后的可执行文件在 **./md** 目录下

## 命令行参数

- `-p`：监听端口号。默认值：**9900**
- `-log`：日志目录，存放近 30 天的日志。默认值：**./logs**
- `-data`：数据目录，存放数据库文件和图片。默认值：**./data**
- `-reg`：是否允许注册（即使禁止注册，在没有任何用户的情况时仍可注册）。默认值：**true**
- `-pg_host`：postgres 主机地址
- `-pg_port`：postgres 端口
- `-pg_user`：postgres 用户
- `-pg_password`：postgres 密码
- `-pg_db`：postgres 数据库名

### 数据库选择

当 postgres 相关的 5 个命令行参数全部填写时，将使用 postgres 数据库，否则使用默认的 sqlite 数据库

## docker 镜像

- [https://hub.docker.com/r/streamerzero/md](https://hub.docker.com/r/streamerzero/md)

## 隐含的功能按钮

- 编辑模式：“云文档”标题
- 收缩边栏：“云文档”标题左侧图标
- 发布链接：文档列表蓝色倒三角
- 公开文档列表：登录页“云文档”标题
