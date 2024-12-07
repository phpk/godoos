<p align="center">
    <img src="./build/appicon.png" width="120" height="120">
</p>

<h1 align="center">GodoOS</h1>
一款高效的内网办公操作系统，内含word/excel/ppt/pdf/内网聊天/白板/思维导图等多个办公系统工具，支持原生文件存储。平台界面精仿windows风格，操作简便，同时保持低资源消耗和高性能运行。无需注册即可自动连接内网用户，实现即时通讯和文件共享。灵活高配置的应用商店，可无限扩展。
<div align="center">


[English](README.en.md) | 简体中文

[使用文档](https://docs.godoos.com/zh/godoos.html) | [FAQ](./docs/Faq.md) | [应用开发](./docs/Store.md)

</div>

## 🎉 V1.0.3更新日志

- 新增ai模型管理，可下载管理ollama模型(需要先安装ollama)
- 新增ai助手，可控制整个系统的prompt
- word新增ai优化/续写/纠错/翻译/总结，生成大纲，根据大纲一键创建文章
- markdown更换为更实用的cherry-markdown，支持draw.io绘图，支持导出为思维导图/pdf/长图/docx/md/html格式
- 修复截图/截屏路径
- 修复更换存储将路径后系统不重置的问题
- 新增文件密码箱（系统设置里），可根据不同文件进行加密存储
- 美化日程提醒弹窗
- 修复word格式问题以及导出名字不对
- markdown新增ai优化/续写/纠错/翻译/总结，生成大纲，根据大纲一键创建文章
- 更改文档存储方式，支持选择文件夹
- 内网聊天新增ai对话，可保存对话历史，可更换模型和prompt
- 新增可定义端口和访问路径，支持web端系统重启
- 新增每个文件可独立设置密码，支持不可逆加密文件（加密文件后不可更改密码）

## 🏭 第三阶段目标（十二月底发布）
1. **文档处理与Markdown智能升级**：（已完成）
	- **AI续写**：借助先进的自然语言处理技术，让您的文档创作灵感不断，续写流畅无阻。
	- **智能总结**：一键提取文档精髓，快速生成精炼总结，助力高效阅读与信息提炼。
	- **纠错优化**：智能识别并纠正文档中的语法、拼写错误，确保内容准确无误。
	- **智能提纲生成**：自动梳理文章结构，生成逻辑清晰的提纲，助您轻松驾驭复杂文档。

2. **本地文件级知识库管理**：
	- 引入全新的知识库管理系统，实现对本地文件的智能分类、标签化管理与高效检索，让您的知识积累更加有序、便捷。

3. **图形处理一键生图**：
	- 创新功能上线，只需简单操作，即可根据文字描述或数据自动生成高质量图表与图像，为报告、演示增添视觉亮点。	

4. **Markdown扩展功能**：
	- **思维导图生成**：支持Markdown内容直接转换为思维导图，可视化呈现信息架构，提升思维整理效率。（已完成）
	- **PPT一键制作**：无缝衔接Markdown文档，轻松导出专业级PPT，让汇报与分享更加生动、专业。

5. **文字转声音功能**：
	- 新增文字朗读服务，支持多种语音风格与语速调节，无论是阅读文档、学习资料还是辅助视力障碍者，都能享受前所未有的便捷与舒适。

### 📥 下载安装(v1.0.2)

1. 💻 **Windows 用户**:
   
- Windows (AMD64) [**Web版**](https://godoos.com/upload/godoos/1.0.2/web/godoos_web_windows_amd64.exe) [**桌面版**](https://godoos.com/upload/godoos/1.0.2/desktop/godoos-amd64-installer.exe)
- Windows (ARM64) [**Web版**](https://godoos.com/upload/godoos/1.0.2/web/godoos_web_windows_arm64.exe) [**桌面版**](https://godoos.com/upload/godoos/1.0.2/desktop/godoos-arm64-installer.exe)


2. 💼 **MacOS 用户**:

- MacOS (AMD64) [**Web版**](https://godoos.com/upload/godoos/1.0.2/web/godoos_web_darwin_amd64)
- MacOS (ARM64) [**Web版**](https://godoos.com/upload/godoos/1.0.2/web/godoos_web_darwin_arm64)

提示：下载后以godoos_web_darwin_amd64为例，命令行：
```
sudo chmod +x godoos_web_darwin_amd64
sudo ./godoos_web_darwin_amd64
```

3. 💽 **Linux 用户**:

- Linux (AMD64) [**Web版**](https://godoos.com/upload/godoos/1.0.2/web/godoos_web_linux_amd64)
- Linux (ARM64) [**Web版**](https://godoos.com/upload/godoos/1.0.2/web/godoos_web_linux_arm64)
提示：下载后以godoos_web_darwin_amd64为例，root账号登录，命令行：
```
chmod +x godoos_web_darwin_amd64
./godoos_web_darwin_amd64
```

- 备注：web版下载后启动服务端。访问地址为：http://localhost:56780/。

### 🚢 **Docker安装**

#### 构建并启动服务

```
cd frontend
pnpm i
pnpm build
cd ..
docker-compose up --build
```

或者直接拉取(v1.0.1)
```
docker run -d -p 56780:56780 --name godoos godoos/godoos:latest
```

- 如果设置本地存储，存储地址为 /root/.godoos/os，设置成功后保存



## 💝 亮点
- ***无需联网使用，全开源***
- ***零配置，无需注册，下载即用***
- ***零污染，无插件依赖***
- ***精小，打包后仅70M，却包含了所有的办公套件***
- ***可无限扩展，支持自定义应用***
- ***golang开发后端，低资源消耗和高性能***
- ***支持多平台，Windows、Linux、MacOS***
- ***完善的应用商店体系，简单学习一下[应用商店配置](./docs/Store.md)即可开发出复杂的应用***

## 💖 开源地址
- [Gitee](https://gitee.com/ruitao_admin/godoos)
- [Github](https://github.com/phpk/godoos)

## 🚀 演示视频
- [全程操作](https://www.bilibili.com/video/BV1NdvaeEEz3/?vd_source=739e0e59aeefdb2e9f760e5037d00245)

## 🚧 开发进程
- 2024年11月15日，发布v1.0.2版本，企业版跟随发布。
- 2024年8月1日，发布v1.0.0版本，发布后，项目进入第二阶段。

###  🎉 v1.0.2更新日志
- 新增本地文件加密存储
- 新增企业端接口（聊天/工作流/文件分享/文件加密）
- 重构本地聊天，修改发现机制（基于ip扫描和arp过滤）
- 本地聊天可批量发送图片/文件夹，修改发送机制，消息基于udp发送，文件基于tcp发送
- 修复word导入格式丢失问题
- 修复文件重命名错误的bug
- 修复拖拽上传中断的bug
- 新增可手动关闭广告
- 优化初始化系统，初始化系统时只请求读写一次
- 去除浏览器存储
- 开源核心底层源码
- 优化思维导图和文件读取

###  🎉 v1.0.1更新日志

- 优化初始化系统，初始化系统时只请求读写一次，确保1秒内打开
- 去除浏览器存储
- 内网聊天增加手工添加ip，跨网段通信在ping通的前提下如果发现不了对方可手工添加对方ip
- 修复思维导图保存的文件每次打开主题又会变成默认主题
- 新增webdav客户端
- 新增远程存储
- 修改选择文件夹会删除文件夹内的文件



## ⚡ 功能说明和预览

### 一、系统桌面
- 精仿windows风格
- 桌面文件管理
- 支持原生文件拖拽上传

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/home.png" width="600" />

### 二、文件管理
- 文件拖拽上传
- 文件搜索
- 原生文件存储
- 直接压缩/解压文件夹（本地存储支持zip/tar/gz/bz2）

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/file.png" width="600" />

### 三、内网聊天
- 无需注册流程，只需在同一内网，即可自动发现并列出所有可用的聊天对象，支持基于局域网的即时消息传输、文件传输等功能。

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/localchat.png" width="600" />

### 四、文档
- 简便的word编辑器，原生存储，支持二维码、手写签名，导入导出

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/doc.png" width="600" />

### 五、表格
- 原生存储，Excel编辑器，支持导入、导出，支持图片、公式

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/excel.png" width="600" />

### 六、markdown
- 原生存储，采用vditor，支持导入、导出，支持大纲、实时预览

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/markdown.png" width="600" />

### 七、思维导图
- 内置多种主题；支持快捷键；节点内容支持图片、图标、超链接、备注、标签。

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/mind.png" width="600" />

### 八、演示文稿
- 原生存储，采用pptist，支持文字、图片、形状、线条、图表、表格、视频、公式等。

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/ppt.png" width="600" />

### 九、文件编辑器
- 原生存储，支持打开text/html/css/js/svg/xml/md等，可以当作一个简单的在线editplus。

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/fileeditor.png" width="600" />

### 十、白板
- 集自由布局、画笔、便签多种创意表达能力于一体，激发团队创造力，随时随地，围绕一块白板沟通。

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/baiban.png" width="600" />

### 十一、图片编辑
- 一个小型的photoshop，原生存储，支持搜索图片，支持图片裁剪、旋转、缩放、滤镜等功能

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/pic.png" width="600" />

### 十二、甘特图
- 项目管理必备工具，支持自定义项目人员和角色，支持拖拽/管理分配（资源、角色、工作）等。

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/gant.png" width="600" />

### 十三、浏览器
- 一款简单的内置浏览器

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/ie.png" width="600" />


### 十四、系统设置
- 可在这里切换存储方式，可切换系统背景。

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/setting-store.png" width="600" />

### 十五、应用商店
- 应用商店管理，丰富的外部接口，可导入/添加/下载外部应用。支持依赖库安装/卸载。

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/store.png" width="600" />

### 十六、截图
- 一个简单的截图工具。截图后文件存到本地。

### 十七、录屏
- 一个简单的录屏工具。录屏后录后文件存到本地。

### 十八、计算器
- 一个仿windows10的计算器，支持历史记录。

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/cal.png" width="600" />

### 十九、音乐库
- 一个简单的声音存储库，支持播放音乐。

### 二十、图片库
- 一个简单的图片存储库，支持查看图片

### 二十一、看板
- 支持标准看板，项目管理必备工具，可快速创建看板并放置到不同的文件夹

<img src="https://gitee.com/ruitao_admin/godoos-image/raw/master/img/kanban.png" width="600" />

### 二十二、进程管理
- 支持进程管理，可以查看进程列表，杀死进程

## 🏆 开发
### 1.进入godo/deps/找到对应系统的文件夹，直接手工打zip压缩包
### 2.构建
- 前端构建（必须）
```bash
cd frontend
pnpm i
pnpm build
```
- 桌面端构建
```bash
# go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails build
# wails build -nsis -upx //you need install nsis and upx
```
- web端构建
```bash
cd godo
chmod +x quick_build.sh //linux or mac必须有执行权限，windows不需要
./quick_build.sh
```

## 📊 帮助

1. 是否支持切换存储目录？
- 支持。进入系统设置页面，修改存储目录即可。修改后程序会重启一次。

2. 如何上传文件？
- 支持拖拽上传。

## 📆 使用场景：
1. 对办公安全要求严苛的企业，比如不许连外网。
2. 对办公存储有特殊需求的企业，比如要求员工的数据必须存储到对应的地方。
3. 对办公office有极客思维的企业，office太过庞大，而godoos仅60多M。

## GodoOS企业版介绍

[使用文档](https://docs.godoos.com/zh/godoos/enterprise/)

### 一、分客户端和服务端

- 客户端为开源版，用户所有数据存储到服务端
- 服务端支持windows/linux和docker安装，安装端支持web端
- 客户端支持全平台（windows/macos/linux），安装端支持桌面端和web端

### 二、完善的用户权限管理

- 可指定用户组存储目录
- 可设定用户组存储空间大小
- 可设定用户组角色权限，细分到每一个接口
- 可以根据部门设定工作流
- 用户文件可同部门分享

### 三、可配置的工作流引擎
- 可视化的自定义工作流引擎
- 先定义表单数据，再定义工作流，表单和工作流业务逻辑隔离
- 支持任务流和审计流
- 任务流支持考试/签到/工作日志/周报/季报等
- 任务流支持自动提醒/迟到提交/自定义周期
- 任务流根据表单字段支持分数/答案自定义，阅卷支持手工打分和自动打分
- 审计流支持审批/驳回等
- 审计支持部门和用户自定义权限
- 支持会签（通过需全员否决只需一人）/或签（一人通过或否决）/民主签（少数服从多数，平票自动打回）
- 支持站内消息和邮件通知
- 支持条件分支判断和自定义抄送人
- 支持手工确认/阅读即审批/手写签名
- 支持驳回到上一个节点或发起人
- 条件判断支持和表单联动，支持发起人过滤，分数判断（任务流）

后续开发：
- 支持表单数据导入/导出，支持自定义导出字段
- 支持工作流数据统计
- 支持自定义横向/纵向数据分析
- 支持表单数据和审计数据同步到其他接口
- 审计流支持打分和评论

### 四、完善的表单管理体系
- 支持表单设计器，可自定义表单字段，支持表单联动，支持表单校验
- 支持表单数据归档，支持自定义时间范围
- 支持数据查看/编辑，搜索定义，自动类型定义（数字和字符串）

### 五、完善的企业聊天沟通工具

- 和本地聊天完全隔离
- 支持群聊/单聊
- 支持文件发送和图片发送
- 在线聊天消息不存储服务端（离线消息加密存储）
- 支持消息提醒

后续开发：
- 支持视频聊天/视频会议
- 支持远程协助
- 支持远程文档协作

### 六、支持对接钉钉和企业微信
- 支持钉钉和企微H5应用对接
- 支持钉钉和企微自动登录
- 支持钉钉和企微用户数据同步（第一次登录同步）

### 六、强大的消息通知
- 支持邮件发送
- 支持站内消息
- 支持弹窗提醒

### 七、强大的数据统计功能
- 支持统计在线人数/应用数
- 支持统计存储空间/运行时长/内存使用情况

### 八、强大的加密功能
- 所有系统配置均加密存储
- 所有用户数据一旦后台加密，仅用户可查看，每个文件可设置二级密码

### 九、强大的本地化AI支持（后续版本）
后续开发：
- 支持文档创作/翻译/润色/总结
- 支持语音识别/合成
- 支持图片识别/生成/训练

### 十、支持本地更新客户端
- 可设置版本号，上传版本，自动提示更新
- 可管理不同操作系统版本的更新策略

### 十一、支持本地化应用商店
- 支持应用商店，用户可下载应用，后台可上传应用
- 支持不同操作系统对应不同的应用版本
- 支持应用版本依赖

## 💻 **企业版下载试用**:
   
- Windows (AMD64) [**Web版**](https://godoos.com/upload/godoos-server/1.0.0/close-pro/osadmin_windows_amd64.exe)
- Windows (ARM64) [**Web版**](https://godoos.com/upload/godoos-server/1.0.0/close-pro/osadmin_windows_arm64.exe)
- MacOS (AMD64) [**Web版**](https://godoos.com/upload/godoos-server/1.0.0/close-pro/osadmin_darwin_amd64)
- MacOS (ARM64) [**Web版**](https://godoos.com/upload/godoos-server/1.0.0/close-pro/osadmin_darwin_arm64)
- Linux (AMD64) [**Web版**](https://godoos.com/upload/godoos-server/1.0.0/close-pro/osadmin_linux_amd64)
- Linux (ARM64) [**Web版**](https://godoos.com/upload/godoos-server/1.0.0/close-pro/osadmin_linux_arm64)

提示：下载后以osadmin_linux_amd64为例，root账号登录，命令行：
```
chmod +x osadmin_linux_amd64
./osadmin_linux_amd64
```

- 备注：企业版为server端，需要配合开源版（客户端）一起使用。需先安装mysql，测试版试用期为一个月。访问地址为：http://localhost:8816/。

## ❤️ 感谢
- [element-plus](http://element-plus.org/)
- [vue3](https://v3.cn.vuejs.org/)
- [wails](https://wails.io/)
- [pptist](https://github.com/pipipi-pikachu/PPTist)
- [cherry-markdown](https://github.com/Tencent/cherry-markdown)
- [mind-map](https://github.com/wanglin2/mind-map)
- [canvas-editor](https://github.com/Hufe921/canvas-editor)
- [Luckysheet](https://gitee.com/mengshukeji/Luckysheet/)

## 💕 关联项目
- [godoai](https://gitee.com/ruitao_admin/godoai)
- [godooa](https://gitee.com/ruitao_admin/gdoa)
- [gododb](https://gitee.com/ruitao_admin/gododb)

## 微信群
<img src="https://gitee.com/ruitao_admin/gdoa/raw/master/docs/wx.png" width="150" />

## 开源

- 承诺永久开源
- 允许企业/个人单独使用，但需保留版权信息
- 如用于商业活动或二次开发后发售，请购买相关版权
- 不提供私下维护工作，如有bug请 [issures](https://gitee.com/ruitao_admin/godoos/issues) 提交
- 请尊重作者的劳动成果

## 💌 支持作者

如果觉得不错，或者已经在使用了，希望你可以去 
<a target="_blank" href="https://gitee.com/ruitao_admin/godoos">Gitee</a> 帮我点个 ⭐ Star，这将是对我极大的鼓励与支持。