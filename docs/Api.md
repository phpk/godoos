## API
- 系统接口地址: http://localhost:56780

### 读取目录

#### HTTP 方法
`GET`

#### 路径
`/files/read`

#### 请求参数
- **Query 参数**: `path` (目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 目录内容列表

#### 示例
```bash
PATH_PARAM="/path/to/directory"
curl -X GET "http://localhost:56780/files/read?path=$PATH_PARAM"
```
---

### 获取文件或目录状态

#### HTTP 方法
`GET`

#### 路径
`/files/stat`

#### 请求参数
- **Query 参数**: `path` (文件或目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 文件或目录的状态信息

#### 示例
```bash
# 获取文件或目录状态
PATH_PARAM="/path/to/file_or_directory"
curl -X GET "http://localhost:56780/files/stat?path=$PATH_PARAM"
```
---

### 更改文件权限

#### HTTP 方法
`POST`

#### 路径
`/files/chmod`

#### 请求体
- **Content-Type**: `application/json`
- **Body**: `{ "path": "string", "mode": "string" }`

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

#### 示例
```bash
# 更改文件权限
JSON='{"path":"/path/to/file","mode":"0644"}'
curl -X POST -H "Content-Type: application/json" -d "$JSON" "http://localhost:56780/files/chmod"
```
---

### 检查文件或目录是否存在

#### HTTP 方法
`GET`

#### 路径
`/files/exists`

#### 请求参数
- **Query 参数**: `path` (文件或目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 存在性检查结果
#### 示例
```bash
# 检查文件或目录是否存在
PATH_PARAM="/path/to/file_or_directory"
curl -X GET "http://localhost:56780/files/exists?path=$PATH_PARAM"
```
---

### 读取文件内容

#### HTTP 方法
`GET`

#### 路径
`/files/readfile`

#### 请求参数
- **Query 参数**: `path` (文件路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 文件内容
#### 示例
```bash
# 读取文件内容
PATH_PARAM="/path/to/file"
curl -X GET "http://localhost:56780/files/readfile?path=$PATH_PARAM"
```
---

### 删除文件

#### HTTP 方法
`GET`

#### 路径
`/files/unlink`

#### 请求参数
- **Query 参数**: `path` (文件路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果
#### 示例
```bash
# 删除文件
PATH_PARAM="/path/to/file"
curl -X GET "http://localhost:56780/files/unlink?path=$PATH_PARAM"
```
---

### 清空文件系统

#### HTTP 方法
`GET`

#### 路径
`/files/clear`

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果
#### 示例
```bash
# 清空文件系统
curl -X GET "http://localhost:56780/files/clear"
```
---

### 重命名文件或目录

#### HTTP 方法
`GET`

#### 路径
`/files/rename`

#### 请求参数
- **Query 参数**: `oldPath` (原文件或目录路径), `newPath` (新文件或目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果
#### 示例
```bash
# 重命名文件或目录
OLD_PATH_PARAM="/path/to/old_file_or_directory"
NEW_PATH_PARAM="/path/to/new_file_or_directory"
curl -X GET "http://localhost:56780/files/rename?oldPath=$OLD_PATH_PARAM&newPath=$NEW_PATH_PARAM"
```
---

### 创建目录

#### HTTP 方法
`POST`

#### 路径
`/files/mkdir`

#### 请求参数
- **Query 参数**: `dirPath` (目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果
#### 示例
```bash
# 创建目录
DIR_PATH_PARAM="/path/to/new_directory"
curl -X POST "http://localhost:56780/files/mkdir?dirPath=$DIR_PATH_PARAM"
```
---

### 删除目录

#### HTTP 方法
`GET`

#### 路径
`/files/rmdir`

#### 请求参数
- **Query 参数**: `dirPath` (目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果
#### 示例
```bash
# 删除目录
DIR_PATH_PARAM="/path/to/directory"
curl -X GET "http://localhost:56780/files/rmdir?dirPath=$DIR_PATH_PARAM"
```
---

### 复制文件

#### HTTP 方法
`GET`

#### 路径
`/files/copyfile`

#### 请求参数
- **Query 参数**: `srcPath` (源文件路径), `dstPath` (目标文件路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果
#### 示例
```bash
# 复制文件
SRC_PATH_PARAM="/path/to/source_file"
DST_PATH_PARAM="/path/to/destination_file"
curl -X GET "http://localhost:56780/files/copyfile?srcPath=$SRC_PATH_PARAM&dstPath=$DST_PATH_PARAM"
```
---

### 写入文件

#### HTTP 方法
`POST`

#### 路径
`/files/writefile`

#### 请求参数
- **Query 参数**: `filePath` (文件路径)

#### 请求体
- **Content-Type**: `multipart/form-data`
- **Body**: 包含 `content` 的表单数据

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果
#### 示例
```bash
# 写入文件
FILE_PATH_PARAM="/path/to/file"
CONTENT="This is the content to be written."
curl -X POST -F "content=@-$CONTENT" "http://localhost:56780/files/writefile?filePath=$FILE_PATH_PARAM"
```
---

### 追加文件内容

#### HTTP 方法
`POST`

#### 路径
`/files/appendfile`

#### 请求参数
- **Query 参数**: `filePath` (文件路径)

#### 请求体
- **Content-Type**: `multipart/form-data`
- **Body**: 包含 `content` 的表单数据

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果
#### 示例
```bash
# 追加文件内容
FILE_PATH_PARAM="/path/to/file"
CONTENT="This is the content to be appended."
curl -X POST -F "content=@-$CONTENT" "http://localhost:56780/files/appendfile?filePath=$FILE_PATH_PARAM"
```
---

### 文件系统事件监听

#### 功能
监听文件系统变化事件

#### 参数
- **path** (监听的文件或目录路径)
- **callback** (事件回调函数)
- **errback** (错误回调函数)

