# 静态文件服务器

#### 功能 
    实现文件的上传下载功能, 以及js, css的静态文件

### 环境 go 1.10

### 编译运行
     go build
     ./filesvr    fancy.filesvr >> log 2>&1 &

### 在浏览器输入ip:8001 即可看到上传下载页面

### 实现服务发现机制, 目前使用consul监控filesvr的状态(所有域名均需要自己配置hosts, consul服务需要提前部署好)
- 本地服务ID "svr.file.1"
- 本地服务名 "filesvr"
- 本地服务域名 "fancygo.cn"
- 本地服务端口 8001
- consul服务域名 "linuxfj.com"
- consul服务端口 8500
