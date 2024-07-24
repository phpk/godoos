---
title: 文件系统接口
icon: circle-info
---

### 读取目录

#### HTTP 方法
`GET`

#### 路径
`/read`

#### 请求参数
- **Query 参数**: `path` (目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 目录内容列表

---

### 获取文件或目录状态

#### HTTP 方法
`GET`

#### 路径
`/stat`

#### 请求参数
- **Query 参数**: `path` (文件或目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 文件或目录的状态信息

---

### 更改文件权限

#### HTTP 方法
`POST`

#### 路径
`/chmod`

#### 请求体
- **Content-Type**: `application/json`
- **Body**: `{ "path": "string", "mode": "string" }`

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

---

### 检查文件或目录是否存在

#### HTTP 方法
`GET`

#### 路径
`/exists`

#### 请求参数
- **Query 参数**: `path` (文件或目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 存在性检查结果

---

### 读取文件内容

#### HTTP 方法
`GET`

#### 路径
`/readfile`

#### 请求参数
- **Query 参数**: `path` (文件路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 文件内容

---

### 删除文件

#### HTTP 方法
`GET`

#### 路径
`/unlink`

#### 请求参数
- **Query 参数**: `path` (文件路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

---

### 清空文件系统

#### HTTP 方法
`GET`

#### 路径
`/clear`

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

---

### 重命名文件或目录

#### HTTP 方法
`GET`

#### 路径
`/rename`

#### 请求参数
- **Query 参数**: `oldPath` (原文件或目录路径), `newPath` (新文件或目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

---

### 创建目录

#### HTTP 方法
`POST`

#### 路径
`/mkdir`

#### 请求参数
- **Query 参数**: `dirPath` (目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

---

### 删除目录

#### HTTP 方法
`GET`

#### 路径
`/rmdir`

#### 请求参数
- **Query 参数**: `dirPath` (目录路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

---

### 复制文件

#### HTTP 方法
`GET`

#### 路径
`/copyfile`

#### 请求参数
- **Query 参数**: `srcPath` (源文件路径), `dstPath` (目标文件路径)

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

---

### 写入文件

#### HTTP 方法
`POST`

#### 路径
`/writefile`

#### 请求参数
- **Query 参数**: `filePath` (文件路径)

#### 请求体
- **Content-Type**: `multipart/form-data`
- **Body**: 包含 `content` 的表单数据

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

---

### 追加文件内容

#### HTTP 方法
`POST`

#### 路径
`/appendfile`

#### 请求参数
- **Query 参数**: `filePath` (文件路径)

#### 请求体
- **Content-Type**: `multipart/form-data`
- **Body**: 包含 `content` 的表单数据

#### 响应
- **Content-Type**: `application/json`
- **响应体**: 操作结果

---

### 文件系统事件监听

#### 功能
监听文件系统变化事件

#### 参数
- **path** (监听的文件或目录路径)
- **callback** (事件回调函数)
- **errback** (错误回调函数)
