## GodoOS应用商店开发教程

### 快速开始
1. 下载[mysql5.7的zip包](https://downloads.mysql.com/archives/get/p/23/file/mysql-5.7.44-winx64.zip)，并解压到用户目录.godoos/run/windows/目录下，文件夹命名为mysql5.7
2. 复制本程序docs/demo/mysql5.7到用户目录.godoos/run/windows/mysql5.7目录下
3. 打开应用商店，添加应用，选择开发模式，本地路径输入mysql5.7，点击确定

### 开发条件

1. 会简单的html开发
2. 熟悉可执行文件的启动流程，然后根据如下的流程配置json文件

### 如何添加一个应用

1. 在用户目录.godoos/run/windows/下建立一个应用文件夹。
2. 一个应用需要包含两个配置文件，在应用根目录下创建`install.json`和`store.json`两个文件。 配置文件格式如下：
- install.json [样本](./demo/mysql5.7/install.json)
```json
{
    "name": "",             // string, 应用程序名称。
    "url": "",              // string, 应用程序下载地址或适配包的下载地址。
    "pkg": "",              // string, 应用程序的官方下载地址。可为空。
    "webUrl":"",            // string, 如何设置，应用程序将显示到桌面。
    "isDev": true,          // boolean, 是否为开发环境。如果设置为true将不会下载数据。
    "version": "1.0.0",     // string, 应用程序版本。
    "icon": "",             // string, 应用程序图标，为可访问的网络地址。
    "hasStart": true,       // boolean, 表示是否显示启动和停止。
    "hasRestart": true,     // boolean, 是否需要重启。
    "setting": true         // boolean, 是否需要设置。只有应用停止才会显示。
}

```
- 备注：如果不需要后端进程的web应用，store.json可以不用配置。

install.json的结构体为：
```json
type InstallInfo struct {
	Name         string           `json:"name"`         // 应用程序名称。
	URL          string           `json:"url"`          // 应用程序下载地址或适配包的下载地址。
	Pkg          string           `json:"pkg"`          // 应用程序的官方下载地址。
	WebUrl       string           `json:"webUrl"`       // 应用程序的网页地址。
	IsDev        bool             `json:"isDev"`        // 标志位，表示是否为开发者版本。
	Version      string           `json:"version"`      // 应用程序的版本号。
	Desc         string           `json:"desc"`         // 应用程序的描述信息。
	Icon         string           `json:"icon"`         // 应用程序的图标路径。
	HasStart     bool             `json:"hasStart"`     // 标志位，表示是否显示启动和停止。
	HasRestart   bool             `json:"hasRestart"`   // 标志位，表示安装后是否需要重启。
	Setting      bool             `json:"setting"`      // 标志位，表示是否需要配置。
	Dependencies []Item           `json:"dependencies"` // 依赖项。
	Categrory    string           `json:"category"`     // 应用程序的分类。
	History      []InstallHastory `json:"history"`      // 应用程序的历史版本。
}
type InstallHastory struct {
	Version string `json:"version"`
	URL     string `json:"url"`
	Pkg     string `json:"pkg"` // 应用程序的官方下载地址。
}
```

- store.json [样本](./demo/mysql5.7/store.json)

```json
{
    "name": "mysql5.7",     // string, 应用程序名称。可不设置，会继承自install.json
    "setting": {
        "binPath": "{exePath}/bin/mysqld.exe", // string, 重要，必须设置。为启动程序路径。
        "confPath": "{exePath}/my.ini",// string, 可为空。为配置文件路径。
        "progressName": "mysqld.exe",// string, 进程名称。如果为单线程可不设置。
        "isOn": true // boolean, 是否启动守护进程
    },
    "config": { // object, 配置文件。里面的配置可以任意填写，和commands配合使用。可以通过http设置里面的参数。
    },
    "commands": {},// object, 命令列表集。可供install里的installCmds调用，也可以通过外部http请求调用。
    "install": { // object, 安装配置。
        "installEnvs": [],// object[], 环境变量。
        "installCmds": []// object[], 启动命令。可调用命令列表集commands里面的命令。
    },
    "start": {
        "startEnvs": [],
        "beforeCmds": [],// 启动前需要执行的命令列表。可调用命令列表集commands里面的命令。
        "startCmds": [// object[], 纯参数命令集。将启动`setting.binPath`,不可调用命令列表集`commands`里面的命令。
            "--defaults-file={exePath}/my.ini"
        ],
        "AfterCmds": []// 启动后需要执行的命令列表。可调用命令列表集commands里面的命令。
    }
}
```

- 备注：核心的替换参数为`{exePath}`即程序的执行目录。其他`{参数}`对应`store.json`里的`config`。

store.json的结构体为：

```json
type StoreInfo struct {
	Name     string           `json:"name"`     // 应用程序商店的名称。
	Setting  Setting          `json:"setting"`  // 应用程序商店的配置信息。
	Config   map[string]any   `json:"config"`   // 应用程序的配置信息映射。
	Commands map[string][]Cmd `json:"commands"` // 应用程序的命令集合。
	Install  InstallStruct    `json:"install"`  // 安装应用程序的信息。
	Start    StartStruct      `json:"start"`    // 启动应用程序的信息。
}
```
Setting的结构体为：
```json
// 包含应用程序的二进制文件路径、配置文件路径等关键设置信息。
type Setting struct {
	BinPath      string `json:"binPath"`      // 应用程序二进制文件的路径。
	ConfPath     string `json:"confPath"`     // 应用程序配置文件的路径。
	ProgressName string `json:"progressName"` // 进程的名称。
	IsOn         bool   `json:"isOn"`         //是否守护进程运行。
}
```
Cmd的结构体为：
```json
type Cmd struct {
	Name     string   `json:"name"`               // 命令的名称。
	FilePath string   `json:"filePath,omitempty"` // 命令文件的路径。
	Content  string   `json:"content,omitempty"`  // 命令的内容。
	BinPath  string   `json:"binPath,omitempty"`  // 执行命令的二进制文件路径。
	TplPath  string   `json:"tplPath,omitempty"`  // 命令的模板路径。
	Cmds     []string `json:"cmds,omitempty"`     // 要执行的子命令列表。
	Waiting  int      `json:"waiting"`            // 等待的时间。
	Kill     bool     `json:"kill"`               // 标志位，表示是否需要终止之前的命令。如果设置了`content`中的进程名，则优先以进程名字杀死进程
	Envs     []Item   `json:"envs"`               // 命令执行时的环境变量。
}
```
Install的结构体为：
```json
// InstallStruct 描述了安装过程中的环境变量和命令列表。
type InstallStruct struct {
	InstallEnvs []Item   `json:"installEnvs"` // 安装过程中需要的环境变量。
	InstallCmds []string `json:"installCmds"` // 安装过程中需要执行的命令列表。
}
type StartStruct struct {
	StartEnvs  []Item   `json:"startEnvs"`  // 启动过程中需要的环境变量。
	BeforeCmds []string `json:"beforeCmds"` // 启动前需要执行的命令列表。可调用命令列表集commands里面的命令。
	StartCmds  []string `json:"startCmds"`  // 启动过程中需要执行的命令列表。
	AfterCmds  []string `json:"afterCmds"`  // 启动后需要执行的命令列表。可调用命令列表集commands里面的命令。
}
// Item 是一个通用的键值对结构体，用于表示配置项或环境变量等。
type Item struct {
	Name  string `json:"name"`  // 配置项的名称。
	Value any    `json:"value"` // 配置项的值。
}
```

3. 在应用商店添加应用，选择本地添加，输入应用的目录名字（不需要填写整个目录）。

### 配置文件`store.json`说明

1. `install`可以调用`commands`里的命令。
2. `commands`里的命令也可以调用自身的命令。
3. 所有的命令都可以串联使用。

### 如何设置配置
1. 在应用目录下建立`static`目录，并创建`index.html`文件。`install.json`中`setting`设置为`true`。
前端配置样列 [样本](./demo/mysql5.7/static/index.html)
```js
const postData = {
dataDir: dataDir,// 此处对应store.json中的config配置项
logDir: logDir,// 此处对应store.json中的config配置项
port: port, // 此处对应store.json中的config配置项
name: "mysql5.7",// 应用名称
cmdKey: "setting"// 命令键，cmds的name
 }
const comp = await fetch('http://localhost:56780/store/setting', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify(postData)
});
```
- 关键要设置好`name`和`cmdKey`，`name`为应用名称，`cmdKey`对应`store.json`中的`commands`对象中的键，一个对象中又可以配置一系列的命令，可以参考Cmd的结构体，样例：
```json
"commands": {
        "initData": [
            {
                "name": "exec",
                "binPath": "{exePath}/bin/mysqld.exe",
                "cmds": [
                    "--defaults-file={exePath}/my.ini",
                    "--initialize"
                ],
                "waiting": 1 //为等待的秒数
            },
            {
                "name": "exec",
                "binPath": "{exePath}/bin/mysqld.exe",
                "cmds": [
                    "--defaults-file={exePath}/my.ini",
                    "--init-file={exePath}/password.txt"
                ],
                "waiting": 3,
                "content": "mysqld.exe",
                "kill": true
            },
            {
                "name": "start"
            }
        ],
		 "setting": [
            {
                "name": "changeFile",
                "tplPath": "{exePath}/my.ini.tpl",
                "filePath": "{exePath}/my.ini"
            },
            {
                "name": "initData"
            }
        ],
}
```
- 上述样例中setting又调用了initData命令。
- post的固定地址为`http://localhost:56780/store/setting`
- 原理是通过`http`请求`/store/setting`接口，将配置信息发送给`store`服务，然后`store`服务会根据配置信息，自动更换配置信息，并启动应用

### static目录说明
1. `index.html`是应用的首页
2. `static`目录下的在执行install时文件会自动复制到`.godoos/static/`应用目录下
3. `store.json`如果设置了`icon`，并且`static`目录下存在该文件，则应用图标为该文件。否则为`install.json`中的icon	

### `commands`内置函数说明
- 系统封装了一些用于处理进程控制和文件操作的功能函数，以下是各函数的详细描述：
1. `start` 启动应用
2. `stop` 停止应用
3. `restart` 重启应用
4. `exec` 执行命令 必须设置binPath和cmds
5. `writeFile` 写入文件 必须设置filePath和content，以config为依据，对content中的`{参数}`进行替换
6. `changeFile` 修改文件 必须设置filePath和tplPath，以config为依据，在模板文件中进行替换
7. `deleteFile` 删除文件
8. `unzip` 解压文件 必须设置filePath和content，content为解压目录 
9. `zip` 压缩文件 必须设置filePath和content，filePath为将要压缩的文件夹，content为压缩后的文件名 
10. `mkdir` 创建文件夹 必须设置FilePath，FilePath为创建的文件夹路径  
11. `startApp` 启动其他应用，content为应用名
12. `stopApp` 停止其他应用，content为应用名

- 保留命令`uninstall`，如果设置了`uninstall`，系统在卸载的时候将会执行它

### 进阶操作
1. 下载[mysql8.0](https://dev.mysql.com/get/Downloads/MySQL-8.0/mysql-8.0.39-winx64.zip)
2. 参考demo下的mysql8.0目录，尝试自己制作安装包

