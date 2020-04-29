## 部署
**环境准备**
> redis

> nginx

[下载DNSLOG前端和后端程序](https://github.com/joke0jie/DNSLOG/releases)

#### 1.安装redis

#### 2.安装nginx
**nginx 配置**
```
server {
    listen       80;
    access_log  /var/log/nginx/host.access.log  main;
    location / {
        root   /root/dnslog/dnslog;   		#更改为你自己前端所在的目录 
        index  index.html index.htm;
    }
	    location  /api/ {
	        proxy_pass http://127.0.0.1:443/;  #配置GDNSLOG的转发端口一般默认就行了
	    	}
	}
```


#### 3.运行GDNSLOG

可以使用定时任务轮循

```crontab -e```

```30 * * * * /bin/bash /root/damon.sh```

```#!/bin/bash

COUNT=$(ps -ef |grep GDNslog_linux |grep -v "grep" |wc -l)
echo $COUNT
if [ $COUNT -eq 0 ]; then
        cd /root/dnslog
        ./GDNslog_linux
        cd /usr/local/bin
        ./redis-server
else
        echo not run
fi```


