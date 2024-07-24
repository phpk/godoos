<p align="center">
    <img src="./build/appicon.png" width="120" height="120">
</p>

<h1 align="center">GodoOS</h1>
An efficient intranet office platform that includes various tools such as Word, Excel, PPT, PDF, intranet chat, whiteboard, and mind mapping, and supports native file storage. The platform interface closely resembles the Windows style, featuring easy operation while maintaining low resource consumption and high performance. Automatically connects to intranet users without registration, enabling instant messaging and file sharing.

<div align="center">

[![license][license-image]][license-url] 

English | [简体中文](README.md)

### Install

[![Windows][Windows-image]][Windows-url]
[![MacOS][MacOS-image]][MacOS-url]
[![Linux][Linux-image]][Linux-url]


[license-image]: ./docs/img/license_%20MIT.svg

[license-url]: https://spdx.org/licenses/MIT.html


[Windows-image]: ./docs/img/Windows.svg

[Windows-url]: https://gitee.com/ruitao_admin/godoos/releases/download/v1.0.0/godoos-windows.exe

[MacOS-image]: ./docs/img/MacOS.svg

[MacOS-url]: https://gitee.com/ruitao_admin/godoos/releases/download/v1.0.0/godoos-macos.dmg

[Linux-image]: ./docs/img/Linux.svg

[Linux-url]: https://gitee.com/ruitao_admin/godoos/releases/download/v1.0.0/godoos-linux

</div>

## Highlights
- No need for internet connection, fully open source
- Zero configuration, no registration required, download and use immediately
- Zero pollution, no plugin dependency
- Small in size, packaged for only 61M, it does include all the office suites
- Unlimited scalability, supports custom applications
- Golang develops backend with low resource consumption and high performance

## Function Description

### 1、 System Desktop
- Exquisite imitation of Windows style
- Desktop file management
- Support native file drag and drop upload
<img src="./docs/img/home.png" width="600" />

### 2、 File management
- Drag and drop file upload
- File Search
- Native file storage
<img src="./docs/img/file.png" width="600" />

### 3、 Internal chat
- No need for complicated registration process, simply discover and list all available chat partners within the same local area network, and start instant messaging immediately. Support LAN based instant messaging, file transfer, and other functions to facilitate seamless communication and collaboration within the team.
<img src="./docs/img/localchat.png" width="600" />

### 4、 Documents
- Simple Word editor, native storage, supports QR codes, handwritten signatures, import and export
<img src="./docs/img/doc.png" width="600" />

### 5、 Table
- Native storage, Excel editor, supports import and export, supports images and formulas
<img src="./docs/img/excel.png" width="600" />

### 6、 Markdown
- Native storage, using VDitors, supporting import and export, outline and real-time preview
<img src="./docs/img/markdown.png" width="600" />

### 7、 Mind map
- Built in multiple themes; Support shortcut keys; Node content supports images, icons, hyperlinks, notes, and tags.
<img src="./docs/img/mind.png" width="600" />

### 8、 Presentation Presentation
- Native storage, using pptist, supporting text, images, shapes, lines, charts, tables, videos, formulas, etc.
<img src="./docs/img/ppt.png" width="600" />

### 9、 File Editor
- Native storage, supports opening text/html/css/js/svg/xml/md, etc., can be used as a simple online editplus.
<img src="./docs/img/fileeditor.png" width="600" />

### 10、 Whiteboard
- Integrating various creative expression abilities such as free layout, paintbrushes, and notes, it inspires team creativity and enables communication around a whiteboard anytime, anywhere.
<img src="./docs/img/baiban.png" width="600" />

### 11、 Image editing
- A small Photoshop with native storage, supporting image search, cropping, rotation, scaling, filtering, and other functions
<img src="./docs/img/pic.png" width="600" />

### 12、 Gantt Chart
- A must-have tool for project management, supporting custom project personnel and roles, and supporting drag and drop/management allocation (resources, roles, work), etc.
<img src="./docs/img/gant.png" width="600" />

### 13、 Browser
- A simple built-in browser
<img src="./docs/img/ie.png" width="600" />


### 14、 System settings
- You can switch storage methods and system backgrounds here.
<img src="./docs/img/setting-store.png" width="600" />

### 15、 App Store
- App store management allows for the addition of external applications.

### 16、 Screenshot
- A simple screenshot tool. Save the screenshot file locally.

### 17、 Screen recording
- A simple screen recording tool. After recording the screen, save the recorded file locally.

### 18、 Calculator
- A calculator that mimics Windows 10 and supports historical records.
<img src="./docs/img/cal.png" width="600" />

### 19、 Music Library
- A simple sound storage library that supports playing music.

### 20、 Picture Library
- A simple image repository that supports viewing images

### 21、 Kanban board
- Support standard kanban, a must-have tool for project management, which can quickly create kanban and place it in different folders
<img src="./docs/img/kanban.png" width="600" />

## Development
### Build
- Front end construction
```bash
cd frontend
pnpm i
pnpm build
```
- Backend construction
```bash
# go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails build
# wails build -nsis -upx //you need install nsis and upx
```
## thank
- [element-plus](http://element-plus.org/)
- [vue3](https://v3.cn.vuejs.org/)
- [wails](https://wails.io/)
- [pptist](https://github.com/pipipi-pikachu/PPTist)
- [vditor](https://github.com/Vanessa219/vditor)
- [mind-map](https://github.com/wanglin2/mind-map)
- [canvas-editor](https://github.com/Hufe921/canvas-editor)