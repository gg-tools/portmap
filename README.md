# PortMap

轻量级超强反向代理软件，用go语言开发，通过公共可知的服务器端口进行代理，把局域网内任何主机的本地服务发布给局域网外的用户，可用来跨越各种防火墙。

## Usage

**例如如下场景：**

局域网内的主机A(192.168.1.2)上开启HTTP服务，外部网络的主机B希望访问A的服务。由于A被防火墙保护，局域网外的主机完全无法访问A，使用PortMap可实现这一要求：

1. 找一台A和B都能访问的内网或公网服务器P(61.1.1.1)，在其上启动portmap服务程序 

```
./server -p 8009 -up 8010 -pwd mypassword

# 或使用Docker
# docker run -d --name portmaps --net=host liamylian/portmaps -p 8009 -up 8010 -pwd mypassword
```

2. 在A上启动PortMap客户程序 

```
./client -l 127.0.0.1:80 -r 61.1.1.1:8009 -pwd mypassword

# 或使用Docker
# docker run -d --name portmapc --net=host liamylian/portmapc -l 127.0.0.1:80 -r 51.1.1.1:8009 -pwd mypassword
```

3. 外部的任何主机直接通过http://61.1.1.1:8010即可访问A的HTTP服务