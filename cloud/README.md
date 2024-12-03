## 安装帮助

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

- 首先安装mingw-w64，进入命令行界面

```
./build.sh
```

- 打包成功后每个系统的版本在dist目录下

