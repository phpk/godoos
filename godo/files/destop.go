/*
 * GodoOS - A lightweight cloud desktop
 * Copyright (C) 2024 https://godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package files

import (
	"strings"
	"time"
)

// RootObject represents the root object in the JSON data.
type RootObject struct {
	UpdateTime time.Time    `json:"updatetime"`
	Desktop    []OsFileInfo `json:"apps"`
	Menulist   []OsFileInfo `json:"menulist"`
}

var RootAppList = []map[string]string{
	{"name": "computer", "icon": "diannao", "position": "Desktop,Menulist"},
	{"name": "appstore", "icon": "store", "position": "Desktop,Menulist"},
	{"name": "localchat", "icon": "chat", "position": "Desktop"},
	{"name": "document", "icon": "word", "position": "Desktop"},
	{"name": "excel", "icon": "excel", "position": "Desktop"},
	{"name": "markdown", "icon": "markdown", "position": "Desktop"},
	{"name": "mindmap", "icon": "mindexe", "position": "Desktop"},
	{"name": "ppt", "icon": "pptexe", "position": "Desktop"},
	{"name": "fileEditor", "icon": "editorbt", "position": "Desktop"},
	{"name": "board", "icon": "kanban", "position": "Desktop"},
	{"name": "whiteBoard", "icon": "baiban", "position": "Desktop"},
	{"name": "piceditor", "icon": "picedit", "position": "Desktop"},
	{"name": "gantt", "icon": "gant", "position": "Desktop"},
	{"name": "browser", "icon": "brower", "position": "Desktop,Menulist"},
	{"name": "setting", "icon": "setting", "position": "Menulist"},
	{"name": "system.version", "icon": "info", "position": "Menulist"},
	{"name": "process.title", "icon": "progress", "position": "Menulist"},
	{"name": "calculator", "icon": "calculator", "position": "Menulist"},
	{"name": "calendar", "icon": "calendar", "position": "Menulist"},
	{"name": "musicStore", "icon": "music", "position": "Menulist"},
	{"name": "gallery", "icon": "gallery", "position": "Menulist"},
	{"name": "aiHelper", "icon": "aiassistant", "position": "Desktop"},
	{"name": "aiModule", "icon": "aidown", "position": "Desktop"},
	{"name": "aiSetting", "icon": "aisetting", "position": "Menulist"},
}

// GetInitRootList constructs the initial root list.
func GetInitRootList() RootObject {
	var desktopApps []OsFileInfo
	var menulistApps []OsFileInfo
	nowtime := time.Now()
	var id = 1
	for _, app := range RootAppList {
		positions := strings.Split(app["position"], ",")
		content := "link::Desktop::" + app["name"] + "::" + app["icon"]
		for _, pos := range positions {
			switch pos {
			case "Desktop":
				newApp := OsFileInfo{
					Name:       app["name"] + ".exe",
					Path:       "/C/Users/Desktop/" + app["name"] + ".exe",
					OldPath:    "/C/Users/Desktop/" + app["name"] + ".exe",
					ParentPath: "/C/Users/Desktop",
					Content:    content,
					Ext:        "exe",
					Title:      app["name"],
					ID:         id,
					IsFile:     true,
					IsDir:      false,
					IsSymlink:  false,
					Size:       int64(len(content)), // Size can be set to any value
					ModTime:    nowtime,
					AccessTime: nowtime,
					CreateTime: nowtime,
					Mode:       511,
				}
				desktopApps = append(desktopApps, newApp)
			case "Menulist":
				newApp := OsFileInfo{
					Name:       app["name"] + ".exe",
					Path:       "/C/Users/Menulist/" + app["name"] + ".exe",
					OldPath:    "/C/Users/Menulist/" + app["name"] + ".exe",
					ParentPath: "/C/Users/Menulist",
					Content:    content,
					Ext:        "exe",
					Title:      app["name"],
					ID:         id,
					IsFile:     true,
					IsDir:      false,
					IsSymlink:  false,
					Size:       int64(len(content)), // Size can be set to any value
					ModTime:    nowtime,
					AccessTime: nowtime,
					CreateTime: nowtime,
					Mode:       511,
				}
				menulistApps = append(menulistApps, newApp)
			}
		}
		id++
	}

	return RootObject{
		UpdateTime: nowtime,
		Desktop:    desktopApps,
		Menulist:   menulistApps,
	}
}
