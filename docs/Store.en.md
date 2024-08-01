## GodoOS App Store Development Tutorial

### Quick Start
1. Download the [mysql5.7 zip package](https://downloads.mysql.com/archives/get/p/23/file/mysql-5.7.44-winx64.zip), and extract it to the `.godoos/run/windows/` directory under the user directory, naming the folder `mysql5.7`.
2. Copy the `mysql5.7` from this program's `docs/demo/mysql5.7` to the `.godoos/run/windows/mysql5.7` directory under the user directory.
3. Open the app store, add an app, select development mode, input `mysql5.7` as the local path, and click OK.

### Development Requirements

1. Basic HTML development skills.
2. Familiarity with the executable file startup process, then configure the JSON file according to the following flow.

### How to Add an Application

1. Create an application folder under the `.godoos/run/windows/` directory in the user directory.
2. An application requires two configuration files; create the `install.json` and `store.json` files at the root of the application directory. The configuration file formats are as follows:
- `install.json` [Sample](./demo/mysql5.7/install.json)
```json
{
    "name": "",             // string, Name of the application.
    "url": "",              // string, Download URL of the application or adapter package.
    "pkg": "",              // string, Official download URL of the application. Can be empty.
    "webUrl":"",            // string, If set, the application will be displayed on the desktop.
    "isDev": true,          // boolean, Whether it is a development environment. If set to true, data will not be downloaded.
    "version": "1.0.0",     // string, Version of the application.
    "icon": "",             // string, Icon of the application, accessible network address.
    "hasStart": true,       // boolean, Indicates whether startup and shutdown are shown.
    "hasRestart": true,     // boolean, Whether a restart is needed.
    "setting": true         // boolean, Whether settings are needed. Only shown when the application is stopped.
}
```

- Note: If a web application does not require a backend process, store.json does not need to be configured.
The structure of install.json is:

```json
type InstallInfo struct {
	Name          string `json:"name"`          // Name of the application. Important, must match the directory name of the application.
	URL           string `json:"url"`           // Download URL of the application or adapter package.
	Pkg           string `json:"pkg"`           // Official download URL of the application.
	WebUrl        string `json:"webUrl"`        // Web address of the application.
	IsDev         bool   `json:"isDev"`         // Flag indicating whether it is a developer version.
	Version       string `json:"version"`       // Version number of the application.
	Desc          string `json:"desc"`          // Description information of the application.
	Icon          string `json:"icon"`          // Path to the application icon.
	HasStart      bool   `json:"hasStart"`      // Flag indicating whether startup and shutdown are shown.
	HasRestart    bool   `json:"hasRestart"`    // Flag indicating whether a restart is needed after installation.
	Setting       bool   `json:"setting"`       // Flag indicating whether configuration is needed.
	Dependencies  []Item `json:"dependencies"`  // Dependencies.
    History      []InstallHastory `json:"history"`// History
}
type InstallHastory struct {
	Version string `json:"version"`
	URL     string `json:"url"`
	Pkg     string `json:"pkg"` // Official download URL of the application.
}
```
- store.json [Sample](./demo/mysql5.7/store.json)

```json
{
    "setting": {
        "binPath": "{exePath}/bin/mysqld.exe", // string, Important, must be set. Path to the startup program.
        "confPath": "{exePath}/my.ini",        // string, Can be empty. Path to the configuration file.
        "progressName": "mysqld.exe",          // string, Process name. Not required if single-threaded.
        "isOn": true                           // boolean, Whether to start the daemon process.
    },
    "config": {                                 // object, Configuration file. Any configuration inside can be filled out, used in conjunction with commands. Can be set via HTTP.
    },
    "commands": {},                             // object, List of commands. Available for invocation by `installCmds` inside `install`, also callable through external HTTP requests.
    "install": {                                // object, Installation configuration.
        "installEnvs": [],                      // object[], Environment variables.
        "installCmds": []                       // object[], Startup commands. Can invoke commands from the command list set `commands`.
    },
    "start": {
        "startEnvs": [],
        "beforeCmds": [],                       // List of commands to execute before startup. Can invoke commands from the command list set `commands`.
        "startCmds": [                          // object[], Pure parameter command set. Will start `setting.binPath`, cannot invoke commands from the `commands` list.
            "--defaults-file={exePath}/my.ini"
        ],
        "AfterCmds": []                        // List of commands to execute after startup. Can invoke commands from the command list set `commands`.
    }
}
```
- Note: The core replacement parameter is `{exePath}`, which is the execution directory of the program. Other `{parameters}` correspond to the config in `store.json`.


The structure for `store.json` is as follows:

```json
type StoreInfo struct {
	Setting  Setting          `json:"setting"`  // Configuration information for the application store.
	Config   map[string]any   `json:"config"`   // Mapping of application configuration information.
	Commands map[string][]Cmd `json:"commands"` // Collection of application commands.
	Install  InstallStruct    `json:"install"`  // Information for installing the application.
	Start    StartStruct      `json:"start"`    // Information for starting the application.
}
```
The structure for Setting is:
```json
// Contains critical setting information such as the binary file path and configuration file path of the application.
type Setting struct {
	BinPath      string `json:"binPath"`      // Path to the application's binary file.
	ConfPath     string `json:"confPath"`     // Path to the application's configuration file.
	ProgressName string `json:"progressName"` // Name of the process.
	IsOn         bool   `json:"isOn"`         // Indicates if the daemon process is running.
}
```
The structure for Cmd is:
```json
type Cmd struct {
	Name     string   `json:"name"`               // Name of the command.
	FilePath string   `json:"filePath,omitempty"` // Path to the command file.
	Content  string   `json:"content,omitempty"`  // Content of the command.
	BinPath  string   `json:"binPath,omitempty"`  // Path to the binary file for executing the command.
	TplPath  string   `json:"tplPath,omitempty"`  // Template path for the command.
	Cmds     []string `json:"cmds,omitempty"`     // List of subcommands to be executed.
	Waiting  int      `json:"waiting"`            // Waiting time.
	Kill     bool     `json:"kill"`               // Flag indicating whether to terminate previous commands. If the process name is set in `content`, priority is given to killing the process by name.
	Envs     []Item   `json:"envs"`               // Environment variables during command execution.
}
```
The structure for Install is:
```json
// `InstallStruct` describes environment variables and command lists during the installation process.
type InstallStruct struct {
	InstallEnvs []Item   `json:"installEnvs"` // Environment variables required during installation.
	InstallCmds []string `json:"installCmds"` // List of commands to execute during installation.
}
type StartStruct struct {
	StartEnvs  []Item   `json:"startEnvs"`  // Environment variables required during startup.
	BeforeCmds []string `json:"beforeCmds"` // List of commands to execute before startup. Commands can be invoked from the `commands` list.
	StartCmds  []string `json:"startCmds"`  // List of commands to execute during startup.
	AfterCmds  []string `json:"afterCmds"`  // List of commands to execute after startup. Commands can be invoked from the `commands` list.
}
// `Item` is a generic key-value pair structure used to represent configuration items or environment variables, etc.
type Item struct {
	Name  string `json:"name"`  // Name of the configuration item.
	Value any    `json:"value"` // Value of the configuration item.
}
```

3. To add an application in the app store, select local addition and input the name of the application directory (no need to fill in the entire directory).

### Explanation of the Configuration File `store.json`

1. install can invoke commands from the commands list.
2. Commands in the commands list can also invoke their own commands.
3. All commands can be chained together.

### How to Set Configuration
1. Create a static directory under the application directory and create an index.html file. Set setting to true in install.json. Frontend configuration example [Sample](./demo/mysql5.7/static/index.html)
```js
const postData = {
	dataDir: dataDir, // Corresponds to the config configuration item in store.json
	logDir: logDir,   // Corresponds to the config configuration item in store.json
	port: port,       // Corresponds to the config configuration item in store.json
	name: "mysql5.7", // Application name
	cmdKey: "setting" // Command key, name of cmds
};
const comp = await fetch('http://localhost:56780/store/setting', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify(postData)
});
```
- It is crucial to set name and cmdKey correctly; name is the application name, and cmdKey corresponds to the key in the commands object in store.json. An object can configure a series of commands, which can refer to the structure of Cmd, sample:
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
                "waiting": 1 //Waiting seconds
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
- In the above sample, setting invokes the initData command again.
- The fixed address for POST is http://localhost:56780/store/setting
- The principle works by sending configuration information to the store service through an http request to the /store/setting interface. Then, the store service will automatically update the configuration information and start the application based on the configuration information.

### Explanation of the static Directory
1. `index.html` is the homepage of the application.
2. When install is executed, files under the static directory are automatically copied to the .godoos/static/ application directory.
3. If store.json sets icon and the file exists in the static directory, the application icon is that file. Otherwise, it is the icon in install.json.

### Description of Built-in Applications
- The system encapsulates some functions for handling process control and file operations. Below are detailed descriptions of each function:
1. `start` Start the application.
2. `stop`  Stop the application.
3. `restart` Restart the application.
4. `exec` Execute a command. Must set binPath and cmds.
5. `writeFile` Write to a file. Must set filePath and content. Based on config, replaces {parameters} in content.
6. `changeFile` Modify a file. Must set filePath and tplPath. Based on config, performs replacements in the template file.
7. `deleteFile` Delete a file.
8. `unzip`  Unzip a file. Must set filePath and content. content is the extraction directory.
9. `zip` Compress a file. Must set filePath and content. filePath is the folder to be compressed, and content is the name of the compressed file. 
10. `mkdir`  Create a directory. Must set FilePath, which is the path of the directory to be created.
11. `startApp`  Start another application. content is the name of the application.
12. `stopApp`  Stop another application. content is the name of the application.

### Advanced Operations
1. Download [mysql8.0](https://dev.mysql.com/get/Downloads/MySQL-8.0/mysql-8.0.39-winx64.zip)
2. Refer to the mysql8.0 directory under demo and attempt to create your own installation package.

