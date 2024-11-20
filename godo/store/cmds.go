/*
 * godoos - A lightweight cloud desktop
 * Copyright (C) 2024 godoos.com
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
	err := SetEnvs(storeInfo.Start.StartEnvs)
	if err != nil {
		return fmt.Errorf("failed to set start environment variable %s: %w", storeInfo.Name, err)
	}
	if len(storeInfo.Start.BeforeCmds) > 0 {
		for _, cmdKey := range storeInfo.Start.BeforeCmds {
			if _, ok := storeInfo.Commands[cmdKey]; ok {
				// 如果命令存在，你可以进一步处理 cmds
				RunCmds(storeInfo, cmdKey)
			}
		}

	}
	cmd, err := GetRunCmd(storeInfo.Setting.BinPath, storeInfo.Start.StartCmds)
	if err != nil {
		return fmt.Errorf("failed to get run cmd: %w", err)
	}
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting script: %v\nDetailed error: %v", err, cmd.Stderr)
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
	if len(storeInfo.Start.AfterCmds) > 0 {
		for _, cmdKey := range storeInfo.Start.AfterCmds {
			if _, ok := storeInfo.Commands[cmdKey]; ok {
				// 如果命令存在，你可以进一步处理 cmds
				RunCmds(storeInfo, cmdKey)
			}
		}

	}
	return nil
}
func GetRunCmd(binPath string, cmds []string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	log.Printf("run script: %s", strings.Join(cmds, " "))
	if binPath != "" {
		if !libs.PathExists(binPath) {
			return cmd, fmt.Errorf("script file %s does not exist", binPath)
		}
		// 如果是非Windows环境，设置可执行权限
		if runtime.GOOS != "windows" {
			// 设置文件权限，0755是一个常见的可执行文件权限掩码
			if err := os.Chmod(binPath, 0755); err != nil {
				return cmd, fmt.Errorf("failed to set executable permissions on %s: %w", binPath, err)
			}
		}

		if len(cmds) > 0 {
			cmd = exec.Command(binPath, cmds...)
		} else {
			cmd = exec.Command(binPath)
		}
	} else {
		if len(cmds) == 0 {
			return cmd, fmt.Errorf("no commands provided")
		}
		log.Printf("run script: %s", strings.Join(cmds, " "))
		if len(cmds) > 1 {
			cmd = exec.Command(cmds[0], cmds[1:]...)
		} else {
			cmd = exec.Command(cmds[0])
		}
	}

	if runtime.GOOS == "windows" {
		// 在Windows上，通过设置CreationFlags来隐藏窗口
		cmd = SetHideConsoleCursor(cmd)
	}
	return cmd, nil
}
func RunStartApp(appName string) error {
	return ExecuteScript(appName)
}

func RunStopApp(appName string) error {
	return StopCmd(appName)
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

	cmd, err := GetRunCmd(cmdParam.BinPath, cmdParam.Cmds)
	if err != nil {
		return fmt.Errorf("failed to get run cmd: %w", err)
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
			if cmdParam.Content != "" {
				if err := KillProcessByName(cmdParam.Content); err != nil {
					log.Printf("failed to kill process %s: %s", cmdParam.Content, err.Error())
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
	_, err := files.Decompress(cmd.FilePath, cmd.Content)
	return err
}
func Zip(cmd Cmd) error {
	return files.Encompress(cmd.FilePath, cmd.Content)
}
