## 应用管理

### 如何添加一个应用

1. 在用户目录.godoos/run下建立一个应用文件夹。
2. 一个应用需要包含两个配置文件，在应用根目录下创建`install.json`和`store.json`两个文件。 配置文件格式如下：
- install.json
```json
{
    "name": "",             // string, 应用程序名称。
    "url": "",              // string, 应用程序下载地址。
    "webUrl":"",            // string, 如何设置，应用程序将显示到桌面。
    "isDev": true,          // boolean, 是否为开发环境。
    "needDownload": true,   // boolean, 是否需要下载。
    "needInstall": true,    // boolean, 是否需要安装。
    "version": "1.0.0",    // string, 应用程序版本。
    "icon": "",             // string, 应用程序图标，为可访问的网络地址。
    "checkProgress": true,  // boolean, 表示是否显示启动和停止。
    "hasRestart": true,     // boolean, 是否需要重启。
    "setting": true         // boolean, 是否需要设置。
}

```
如果不需要后端进程的web应用，store.json可以不用配置。

install.json的struct为：
```json
type InstallInfo struct {
	Name          string `json:"name"`          // 应用程序名称。重要，必须和应用的目录名称一致。
	URL           string `json:"url"`           // 应用程序下载地址。
	WebUrl        string `json:"webUrl"`        // 应用程序的网页地址。
	IsDev         bool   `json:"isDev"`         // 标志位，表示是否为开发者版本。
	NeedDownload  bool   `json:"needDownload"`  // 标志位，表示是否需要下载。
	NeedInstall   bool   `json:"needInstall"`   // 标志位，表示是否需要安装。
	Version       string `json:"version"`       // 应用程序的版本号。
	Desc          string `json:"desc"`          // 应用程序的描述信息。
	Icon          string `json:"icon"`          // 应用程序的图标路径。
	CheckProgress bool   `json:"checkProgress"` // 标志位，表示是否显示启动和停止。
	HasRestart    bool   `json:"hasRestart"`    // 标志位，表示安装后是否需要重启。
	Setting       bool   `json:"setting"`       // 标志位，表示是否需要配置。
	Dependencies  []Item `json:"dependencies"`  // 依赖项。
}
```

- store.json

```json
{
    "name": "mysql5.7",     // string, 应用程序名称。可不设置，会继承自install.json
    "icon": "mysql.png",    // string, 应用程序图标，一般放在应用程序static目录下。
    "setting": {
        "binPath": "{exePath}/bin/mysqld.exe", // string, 重要，必须设置。为启动程序路径。
        "confPath": "{exePath}/my.ini",// string, 重要，必须设置。为配置文件路径。
        "progressName": "mysqld.exe",// string, 进程名称。如果为单线程可不设置。
        "isOn": true // boolean, 是否启动守护进程
    },
    "config": { // object, 配置文件。里面的配置可以任意填写，和cmds配合使用。
    },
    "cmds": {},// object, 命令集。
    "install": { // object, 安装配置。
        "envs": [],// object[], 环境变量。
        "cmds": []// object[], 启动命令。
    },
    "start": {// object, 启动配置。
        "envs": [],// object[], 环境变量。
        "cmds": []// object[], 启动命令。
    }
}
```
store.json的struct为
```json
type StoreInfo struct {
	Name    string           `json:"name"`    // 应用程序的名称。
	Icon    string           `json:"icon"`    // 应用程序的图标路径。
	Setting Setting          `json:"setting"` // 应用程序的配置信息。
	Config  map[string]any   `json:"config"`  // 应用程序的配置信息映射。
	Cmds    map[string][]Cmd `json:"cmds"`    // 应用程序的命令集合。
	Install Install          `json:"install"` // 安装应用程序的信息。
	Start   Install          `json:"start"`   // 启动应用程序的信息。
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
	Kill     bool     `json:"kill"`               // 标志位，表示是否需要终止之前的命令。如果setting中设置了progressName会优先杀死整个进程
	Envs     []Item   `json:"envs"`               // 命令执行时的环境变量。
}
```
Install的struct为：
```json
// Install 描述了安装过程中的环境变量和命令列表。
type Install struct {
	Envs []Item   `json:"envs"` // 安装过程中需要的环境变量。
	Cmds []string `json:"cmds"` // 安装过程中需要执行的命令列表。
}
// Item 是一个通用的键值对结构体，用于表示配置项或环境变量等。
type Item struct {
	Name  string `json:"name"`  // 配置项的名称。
	Value any    `json:"value"` // 配置项的值。
}
```

3. 在应用商店添加应用，选择本地添加，输入应用的目录名字（不需要填写整个目录）。

### 配置文件`store.json`说明

1. `install`和`start`配置可以调用`cmds`里的命令。
2. `cmds`里的命令也可以调用自身的命令。
3. 所有的命令都可以串联使用。

### 如何设置配置
1. 在应用目录下建立`static`目录，并创建`index.html`文件。`install.json`中`setting`设置为`true`。
前端配置样列
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
关键要设置好`name`和`cmdKey`，`cmdKey`对应`store.json`中的`cmds`数组配置项中的`name`，一个数组中又可以配置一系列的命令，可以参考Cmd的结构体

### 内置应用说明
- 系列封装了一些用于处理进程控制和文件操作的功能函数，以下是各函数的详细描述：
1. `start` 启动应用
2. `stop` 停止应用
3. `restart` 重启应用
4. `exec` 执行命令 必须设置binPath和cmds
5. `writeFile` 写入文件 必须设置filePath和content，以config为依据，对content进行替换
6. `changeFile` 修改文件 必须设置filePath和tplPath，以config为依据，在模板文件中进行替换
7. `deleteFile` 删除文件
8. `unzip` 解压文件 必须设置filePath和content，content为解压目录 
9. `zip` 压缩文件 必须设置filePath和content，filePath为将要压缩的文件夹，content为压缩后的文件名 
10. `mkdir` 创建文件夹 必须设置FilePath，FilePath为创建的文件夹路径  
