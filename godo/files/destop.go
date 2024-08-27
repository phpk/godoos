package files

import (
	"strings"
	"time"
)

type AppInfo struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	OldPath     string    `json:"oldPath"`
	ParentPath  string    `json:"parentPath"`
	Content     string    `json:"content"`
	Ext         string    `json:"ext"`
	Title       string    `json:"title"`
	IsSys       int       `json:"isSys"`
	ID          int       `json:"id"`
	IsFile      bool      `json:"isFile"`
	IsDirectory bool      `json:"isDirectory"`
	IsSymlink   bool      `json:"isSymlink"`
	Size        int       `json:"size"`
	Mtime       time.Time `json:"mtime"`     // Modified time
	Atime       time.Time `json:"atime"`     // Access time
	Birthtime   time.Time `json:"birthtime"` // Creation time
	Mode        int       `json:"mode"`
	Rdev        int       `json:"rdev"`
}

// RootObject represents the root object in the JSON data.
type RootObject struct {
	UpdateTime time.Time `json:"updatetime"`
	Desktop    []AppInfo `json:"apps"`
	Menulist   []AppInfo `json:"menulist"`
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
	{"name": "process.title", "icon": "progress", "position": "Menulist"},
}

// GetInitRootList constructs the initial root list.
func GetInitRootList() RootObject {
	var desktopApps []AppInfo
	var menulistApps []AppInfo
	nowtime := time.Now()
	var id = 1
	for _, app := range RootAppList {
		positions := strings.Split(app["position"], ",")
		content := "link::Desktop::" + app["name"] + "::" + app["icon"]
		for _, pos := range positions {
			switch pos {
			case "Desktop":
				newApp := AppInfo{
					Name:        app["name"] + ".exe",
					Path:        "/C/Users/Desktop/" + app["name"] + ".exe",
					OldPath:     "/C/Users/Desktop/" + app["name"] + ".exe",
					ParentPath:  "/C/Users/Desktop",
					Content:     content,
					Ext:         "exe",
					Title:       app["name"],
					IsSys:       1,
					ID:          id,
					IsFile:      true,
					IsDirectory: false,
					IsSymlink:   false,
					Size:        len(content), // Size can be set to any value
					Mtime:       nowtime,
					Atime:       nowtime,
					Birthtime:   nowtime,
					Mode:        511,
					Rdev:        0,
				}
				desktopApps = append(desktopApps, newApp)
			case "Menulist":
				newApp := AppInfo{
					Name:        app["name"] + ".exe",
					Path:        "/C/Users/Menulist/" + app["name"] + ".exe",
					OldPath:     "/C/Users/Menulist/" + app["name"] + ".exe",
					ParentPath:  "/C/Users/Menulist",
					Content:     content,
					Ext:         "exe",
					Title:       app["name"],
					IsSys:       1,
					ID:          id,
					IsFile:      true,
					IsDirectory: false,
					IsSymlink:   false,
					Size:        len(content),
					Mtime:       nowtime,
					Atime:       nowtime,
					Birthtime:   nowtime,
					Mode:        511,
					Rdev:        0,
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
