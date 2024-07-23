package app

import (
	"context"
	"errors"
	"godoos/cmd"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx   context.Context
	exDir string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	go cmd.OsStart()
}
func (a *App) Shutdown(ctx context.Context) {
	a.ctx = ctx
	cmd.OsStop()
}
func (a *App) OpenDirDialog() string {
	path, err := wruntime.OpenDirectoryDialog(a.ctx, wruntime.OpenDialogOptions{
		Title: "Select Folder",
	})
	if err != nil {
		wruntime.LogErrorf(a.ctx, "Error: %+v\n", err)
	}
	return path
}

func (a *App) GetAbsPath(path string) (string, error) {
	var absPath string
	var err error
	if filepath.IsAbs(path) {
		absPath = filepath.Clean(path)
	} else {
		absPath, err = filepath.Abs(filepath.Join(a.exDir, path))
		if err != nil {
			return "", err
		}
	}
	absPath = strings.ReplaceAll(absPath, "/", string(os.PathSeparator))
	//println("GetAbsPath:", absPath)
	return absPath, nil
}

func (a *App) RestartApp() error {
	if runtime.GOOS == "windows" {
		name, err := os.Executable()
		if err != nil {
			return err
		}
		exec.Command(name, os.Args[1:]...).Start()
		wruntime.Quit(a.ctx)
		return nil
	}
	return errors.New("unsupported OS")
}

func (a *App) GetPlatform() string {
	return runtime.GOOS
}
