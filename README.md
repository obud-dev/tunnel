[![Docker Image CI](https://github.com/obud-dev/tunnel/actions/workflows/docker-image.yml/badge.svg)](https://github.com/obud-dev/tunnel/actions/workflows/docker-image.yml)

### 概念职责
- 公网服务器。提供api接口管理tunnel及routes，通过routes规则转发服务到对应tunnel  
- 内网代理服务。连接公网，接受公网请求转发到内网目标服务器，接受目标服务器返回转发到公网  
- 内网目标服务。  

#### 1
tunnel单通道传输信息  
