## godocloud客户端安装帮助

- cloud目录为企业版打包客户端而设定，可以自由定制客户端，根据以下步骤打包

### 第一步：安装nodejs

```
cd ../frontend
npm i
npm run build
```
- 打包成功后复制/godo/deps/dist目录下所有文件到cloud/deps/dist目录

### 第二步：安装golang环境打包

#### linux/mac环境下打包

```
sudo chmod +x build.sh
./build.sh
```
#### windows环境下打包

- 首先安装mingw-w32，进入命令行界面

```
./build.sh
```

- 打包成功后每个系统的版本在dist目录下

### 第三步：部署

- 启动文件
```
sudo chmod +x ./godocloud_linux_amd64
./godocloud_linux_amd64
```
- 访问地址：http://localhost:56781

- 备注：如需更改端口地址，修改cloud/cmd/main.go中，改完之后要重新打包部署

```
const serverAddress = ":56781"
```

