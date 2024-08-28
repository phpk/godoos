## 是否支持浏览器访问？
- 支持。
```
cd frontend
pnpm build
```
然后复制打包后的dist目录到运行程序的根目录，然后重启程序。访问地址为http://localhost:56780/

## 为什么找不到本地文件？
程序默认为浏览器存储。进入系统设置，修改存储方式为本地存储。

## 是否支持切换存储目录？
-支持。
进入系统设置页面，修改存储目录即可。修改后程序会重启一次。

## docker部署失败？
daemon.json文件配置:
```json
{
  "builder": {
    "gc": {
      "defaultKeepStorage": "20GB",
      "enabled": true
    }
  },
  "experimental": false,
  "registry-mirrors": [
    "https://docker.m.daocloud.io"
  ]
}

```

### linux版本如何做守护进程
- 下载web版
- 假设运行目录为/root/godoos
```
chmod a+x /root/godoos
vim /etc/systemd/system/godoos.service
```
- 编辑服务文件
```
[Unit]
Description=Godoos Service
After=network.target

[Service]
ExecStart=/root/godoos
Restart=on-failure
User=myuser
Environment=ENV_VAR=value

[Install]
WantedBy=multi-user.target
```
- 加载服务文件
```
systemctl daemon-reload
systemctl enable godoos
```
- 启动服务
```
systemctl start godoos
```
- 查看服务状态
```
systemctl status godoos
```
- 停止服务
```
systemctl stop godoos
```


