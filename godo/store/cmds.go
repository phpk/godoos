package store

import (
	"fmt"
	"godo/files"
	"godo/libs"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func runStop(storeInfo StoreInfo) error {
	return StopCmd(storeInfo.Name)
}
func runStart(storeInfo StoreInfo) error {
	err := SetEnvs(storeInfo.Start.Envs)
	if err != nil {
		return fmt.Errorf("failed to set start environment variable %s: %w", storeInfo.Name, err)
	}
	if !libs.PathExists(storeInfo.Setting.BinPath) {
		return fmt.Errorf("script file %s does not exist", storeInfo.Setting.BinPath)
	}
	var cmd *exec.Cmd
	if len(storeInfo.Start.Cmds) > 0 {
		cmd = exec.Command(storeInfo.Setting.BinPath, storeInfo.Start.Cmds...)
	} else {
		cmd = exec.Command(storeInfo.Setting.BinPath)
	}
	if runtime.GOOS == "windows" {
		// 在Windows上，通过设置CreationFlags来隐藏窗口
		cmd = SetHideConsoleCursor(cmd)
	}
	if err := cmd.Start(); err != nil {
		//log.Println("Error starting script:", err)
		return fmt.Errorf("error starting script: %v", err)
	}

	// 启动脚本命令并返回可能的错误
	go func(cmd *exec.Cmd) {
		RegisterProcess(storeInfo.Name, storeInfo.Setting.ProgressName, storeInfo.Setting.IsOn, cmd)
		if err := cmd.Wait(); err != nil {
			log.Println("Error starting script:", err)
			return
		}

	}(cmd)

	return nil
}
func runRestart(storeInfo StoreInfo) error {
	err := runStop(storeInfo)
	if err != nil {
		return fmt.Errorf("failed to stop process %s: %w", storeInfo.Name, err)
	}
	return runStart(storeInfo)
}
func runExec(storeInfo StoreInfo, cmdParam Cmd) error {
	err := SetEnvs(cmdParam.Envs)
	if err != nil {
		return fmt.Errorf("failed to set start environment variable %s: %w", storeInfo.Name, err)
	}
	log.Printf("bin path:%v", cmdParam.BinPath)
	log.Printf("cmds:%v", cmdParam.Cmds)
	cmd := exec.Command(cmdParam.BinPath, cmdParam.Cmds...)
	if runtime.GOOS == "windows" {
		// 在Windows上，通过设置CreationFlags来隐藏窗口
		cmd = SetHideConsoleCursor(cmd)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run exec process %s: %w", storeInfo.Name, err)
	}
	go func() {
		if cmdParam.Waiting > 0 {
			if err = cmd.Wait(); err != nil {
				log.Printf("failed to wait exec process %s: %s", storeInfo.Name, err.Error())
				return
			}
			time.Sleep(time.Second * time.Duration(cmdParam.Waiting))
		}
		if cmdParam.Kill {
			if storeInfo.Setting.ProgressName != "" {
				if err := KillProcessByName(storeInfo.Setting.ProgressName); err != nil {
					log.Printf("failed to kill process %s: %s", storeInfo.Setting.ProgressName, err.Error())
				}
			} else {
				if err := cmd.Process.Kill(); err != nil {
					log.Printf("failed to kill process %s: %s", storeInfo.Name, err.Error())
					return
				}
			}

		}
		//log.Printf("run exec process %s, name is %s", storeInfo.Name, cmdParam.Name)
	}()

	return nil
}
func WriteFile(cmd Cmd) error {
	if cmd.FilePath != "" {
		content := cmd.Content
		if content != "" {
			err := os.WriteFile(cmd.FilePath, []byte(content), 0644)
			if err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}
	}
	return nil
}
func ChangeFile(storeInfo StoreInfo, cmd Cmd) error {
	if cmd.FilePath != "" && cmd.TplPath != "" && libs.PathExists(cmd.TplPath) && libs.PathExists(cmd.FilePath) {
		content, err := os.ReadFile(cmd.TplPath)
		if err != nil {
			return fmt.Errorf("failed to read tpl file: %w", err)
		}
		exePath := GetExePath(storeInfo.Name)
		contentstr := strings.ReplaceAll(string(content), "{exePath}", exePath)
		contentstr = ChangeConfig(storeInfo, contentstr)
		if err := os.WriteFile(cmd.FilePath, []byte(contentstr), 0644); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
	return nil
}
func DeleteFile(cmd Cmd) error {
	if cmd.FilePath != "" {
		err := os.Remove(cmd.FilePath)
		if err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}
	return nil
}
func MkDir(cmd Cmd) error {
	err := os.MkdirAll(cmd.FilePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to make dir: %w", err)
	}
	return nil
}
func Unzip(cmd Cmd) error {
	return files.Decompress(cmd.FilePath, cmd.Content)
}
func Zip(cmd Cmd) error {
	return files.Encompress(cmd.FilePath, cmd.Content)
}
