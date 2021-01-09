# n9e-probe

## 功能
ping 和 http get 请求探测
适配 [nightingale](https://github.com/didi/nightingale)

### 指标
#### ping
|metric|说明|
|--|--|
|ping.latency|ping 请求的延迟，单位是毫秒。-1 表示 ping 不通|

|tag|说明|
|--|--|
|ip|探测的目标 ip|
|region|如果配置了，则插入 region tag|

#### url
|metric|说明|
|--|--|
|url.latency|http 请求的延迟，单位是毫秒。-1 表示无法访问|
|url.cert|证书探测。1正常，-1不正常。http 站点则是0|
|url.status_code|返回的状态码|

|tag|说明|
|--|--|
|host|目标 host|
|scheme|目标 scheme|
|path|目标的 path|
|region|如果配置了，则插入 region tag|

### 配置
#### probe.yml
```
logger:
  dir: logs/
  level: INFO
  keepHours: 24

probe:
  # 如果需要区分来自不同区域的探针，可以通过在配置 region 来插入 tag
  #region: default
  timeout: 5 # 探测的超时时间，单位是秒
  limit: 10 # 并发限制
  interval: 30 # 请求的间隔
  headers: # 插入到 http 请求中的 headers，可以多条
    user-agent: Mozilla/5.0 (Linux; Android 6.0.1; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Mobile Safari/537.36 Edg/87.0.664.66

ping:
  107: # n9e 节点上的 nid 号
    - 114.114.114.114 # 要探测的 ip 地址列表
    - 114.114.115.115

url:
  107: # n9e 节点上的 nid 号
    - https://www.baidu.com # 要探测的 ip 地址列表
    - https://www.sjtu.edu.cn/
    - https://bbs.ngacn.cc
    - https://www.163.com
```

## 编译
```
# cd /home
# git clone https://github.com/shanghai-edu/n9e-probe.git
# cd n9e-probe
# ./control build
```
也可以直接在 release 中下载打包好的二进制
## 运行
### 支持 `systemctl` 的操作系统，如 `CentOS7`
执行 `install.sh` 脚本即可，`systemctl` 将托管运行

```
# ./install.sh 
Created symlink from /etc/systemd/system/multi-user.target.wants/agent.service to /usr/lib/systemd/system/agent.service.
```
后续可通过 `systemctl start/stop/restart probe` 来进行服务管理

注意如果没有安装在 `/home` 路径上，则需要修改 `service/agent.service` 中的相关路径，否则 `systemctl` 注册时会找不到

### 不支持 systemctl 的操作系统
执行 `./control start` 启动即可
```
# ./control start
probe started
```
后续可通过 `./control start/stop/restart` 来进行服务管理