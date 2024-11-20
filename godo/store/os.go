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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// DetectOperatingSystem returns the operating system name and version.
func DetectOperatingSystem() (name string, version string, err error) {
	switch runtime.GOOS {
	case "darwin":
		// macOS specific detection
		cmd := exec.Command("sw_vers")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", "", err
		}
		output := string(out)
		lines := strings.Split(output, "\n")
		if len(lines) >= 2 {
			name = strings.TrimSpace(strings.Split(lines[0], ":")[1])
			version = strings.TrimSpace(strings.Split(lines[1], ":")[1])
		}
	case "linux":
		// Linux specific detection
		file, err := os.Open("/etc/os-release")
		if err != nil {
			return "", "", err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "ID=") {
				name = strings.TrimPrefix(line, "ID=")
				name = strings.Trim(name, "\"")
			} else if strings.HasPrefix(line, "VERSION_ID=") {
				version = strings.TrimPrefix(line, "VERSION_ID=")
				version = strings.Trim(version, "\"")
			}
		}

		if err := scanner.Err(); err != nil {
			return "", "", err
		}

		if name == "" || version == "" {
			return "", "", fmt.Errorf("could not determine Linux distribution from /etc/os-release")
		}
	case "windows":
		// Windows specific detection
		cmd := exec.Command("cmd", "/c", "ver")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", "", err
		}
		version = strings.TrimSpace(string(out))
		name = "Windows"
	default:
		// For other operating systems, just return the GOOS value
		name = runtime.GOOS
		version = ""
	}

	return name, version, nil
}

func DetectPackageManager() (pm string, err error) {
	switch runtime.GOOS {
	case "darwin":
		// macOS specific detection
		cmd := exec.Command("which", "brew")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("homebrew not found: %w", err)
		}
		pm = strings.TrimSpace(string(out))
		if pm != "" {
			pm = "brew"
		}
	case "linux":
		// Linux specific detection
		file, err := os.Open("/etc/os-release")
		if err != nil {
			return "", err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "ID=") {
				id := strings.TrimPrefix(line, "ID=")
				id = strings.Trim(id, "\"")
				switch id {
				case "debian", "ubuntu", "linuxmint", "elementary":
					pm = "apt"
				case "centos", "fedora", "rhel":
					cmd := exec.Command("which", "dnf")
					out, err := cmd.CombinedOutput()
					if err != nil {
						cmd = exec.Command("which", "yum")
						out, err = cmd.CombinedOutput()
						if err != nil {
							return "", fmt.Errorf("neither dnf nor yum found: %w", err)
						}
					}
					pm = strings.TrimSpace(string(out))
					if pm != "" {
						pm = "dnf"
					} else {
						pm = "yum"
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return "", err
		}
	case "windows":
		// Windows specific detection - Not applicable for this context
		return "", fmt.Errorf("windows does not use apt, yum, or brew")
	default:
		// For other operating systems, just return an error
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return pm, nil
}

func CheckPathInstalled(execName string, execCommands string) bool {
	cmd := exec.Command(execName, execCommands)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// 如果命令执行失败，可能是因为brew没有安装
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 判断是否是因为找不到命令导致的错误
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.ExitStatus() == 127 {
				return false // brew没有安装
			}
		}
	}
	fmt.Printf("Homebrew version: %s\n", out)
	return true
}
