package files

import (
	"encoding/base64"
	"fmt"
	"godo/libs"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HandleReadFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		libs.ErrorMsg(w, "文件路径不能为空")
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	fullFilePath := filepath.Join(basePath, path)
	file, err := os.Open(fullFilePath)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	fileData, err := io.ReadAll(file)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	text := string(fileData)
	//log.Printf("isPwd: %v", strData)
	if strings.HasPrefix(text, "link::") {
		libs.SuccessMsg(w, text, "文件读取成功")
		return
	}
	// 判断是否为加密文件
	isPwd := libs.IsEncryptedFile(text)

	if !isPwd {
		// 未加密文件，直接返回
		text = base64.StdEncoding.EncodeToString(fileData)
		// if len(text)%8 == 0 {
		// 	text += " "
		// }
		//log.Printf("fileData: %s", content)
		libs.SuccessMsg(w, text, "文件读取成功")
		return
	}
	filePwd := r.Header.Get("pwd")
	// 获取文件密钥
	fileSecret := libs.GetConfigString("filePwd")
	if filePwd == "" && fileSecret == "" {
		libs.Error(w, "加密文件，需要密码", "needPwd")
		return
	}
	decodeFile := ""
	needPwd := false
	if fileSecret != "" {
		decodeFile, err = libs.DecodeFile(fileSecret, text)
		if err != nil {
			libs.Error(w, "请输入文件密码", "needPwd")
			return
		}
		needPwd = true
	}
	if filePwd != "" && !needPwd {
		decodeFile, err = libs.DecodeFile(filePwd, text)
		if err != nil {
			libs.Error(w, "文件密码输入错误", "needPwd")
			return
		}
	}
	// if len(decodeFile)%8 == 0 {
	// 	decodeFile = decodeFile + " "
	// }
	libs.SuccessMsg(w, decodeFile, "加密文件读取成功")
}
func HandleWriteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if err := validateFilePath(path); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fullFilePath := filepath.Join(basePath, path)
	//  获取文件内容
	fileContent, _, err := r.FormFile("content")
	if err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer fileContent.Close()
	content, err := io.ReadAll(fileContent)
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 打开文件，如果不存在则创建
	file, err := os.OpenFile(fullFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		libs.ErrorMsg(w, "Failed to create file.")
		return
	}
	defer file.Close()

	oldContent, err := io.ReadAll(file)
	text := string(oldContent)
	//log.Printf("text: %s", text)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	isPwdFile := libs.IsEncryptedFile(text)

	filePwd := r.Header.Get("pwd")
	// 获取文件密钥
	fileSecret := libs.GetConfigString("filePwd")
	haslink := strings.HasPrefix(string(content), "link::")
	//log.Printf("haslink:%v,fileSecret: %s,isPwdFile:%v,filePwd:%s", haslink, fileSecret, isPwdFile, filePwd)
	needPwd := false
	if !haslink {
		needPwd = true
	}
	if fileSecret != "" || filePwd != "" {
		needPwd = true
	}
	if isPwdFile {
		needPwd = true
	}

	// 即不是加密用户又不是加密文件
	if !needPwd {
		// 直接写入新内容
		file.Truncate(0)
		file.Seek(0, 0)
		_, err = file.Write(content)
		if err != nil {
			libs.ErrorMsg(w, "Failed to write file content.")
			return
		}
		CheckAddDesktop(path)
		libs.SuccessMsg(w, "", "文件写入成功")
		return
	}

	//log.Printf("fileSecret:%s,filePwd:%s,ispwdfile:%v", fileSecret, filePwd, isPwdFile)

	// 是加密文件,写入需继续加密
	pwd := ""
	if isPwdFile {
		//先尝试系统解密
		if fileSecret != "" {
			_, err := libs.DecodeFile(fileSecret, text)
			if err == nil {
				pwd = fileSecret
			}
		}
		//先用户输入解密
		if pwd == "" && filePwd != "" {
			_, err := libs.DecodeFile(filePwd, text)
			if err == nil {
				pwd = filePwd
			}
		}

	} else {
		// 不是加密文件，先判断是否有用户输入
		if filePwd != "" {
			pwd = filePwd
		} else {
			// 没有用户输入，则使用系统默认密码
			if fileSecret != "" {
				pwd = fileSecret
			}
		}
	}
	if pwd == "" {
		libs.Error(w, "文件密码错误", "needPwd")
		//log.Printf("pwd is empty, returning error")
		return
	}
	//log.Printf("pwd: %s", pwd)

	encryData, err := libs.EncodeFile(pwd, string(content))
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}

	// 清空文件内容
	file.Truncate(0)
	file.Seek(0, 0)

	_, err = file.Write([]byte(encryData))
	if err != nil {
		libs.ErrorMsg(w, fmt.Sprintf("文件内容写入失败: %s", err.Error()))
		return
	}
	CheckAddDesktop(path)
	libs.SuccessMsg(w, "", "文件写入成功")
	/*
		//log.Printf("fileData: %s", string(fileData))
		configPwd, ishas := libs.GetConfig("filePwd")
		// 如果不是加密文件或者exe文件
		if !ishas || strings.HasPrefix(string(fileData), "link::") {
			// 没开启加密，直接明文写入
			_, err := newFile.Write(fileData)
			if err != nil {
				libs.HTTPError(w, http.StatusInternalServerError, "数据写入失败")
				return
			}
			err = CheckAddDesktop(path)
			if err != nil {
				log.Printf("Error adding file to desktop: %s", err.Error())
			}
			libs.SuccessMsg(w, "", "文件写入成功")
			return
		} else {
			// 开启加密后，写入加密数据
			configPwdStr, ok := configPwd.(string)
			if !ok {
				libs.HTTPError(w, http.StatusInternalServerError, "配置文件密码格式错误")
				return
			}
			// 拼接密码和加密后的数据
			passwordPrefix := fmt.Sprintf("@%s@", configPwdStr)
			// _, err = newFile.WriteString(fmt.Sprintf("@%s@", configPwdStr))
			// if err != nil {
			// 	libs.HTTPError(w, http.StatusInternalServerError, "密码写入失败")
			// 	return
			// }
			entryData, err := libs.EncryptData(fileData, []byte(configPwdStr))
			if err != nil {
				libs.HTTPError(w, http.StatusInternalServerError, "文件加密失败")
				return
			}
			// 将密码前缀和加密数据拼接成一个完整的字节切片
			completeData := []byte(passwordPrefix + string(entryData))
			// 一次性写入文件
			_, err = newFile.Write(completeData)
			if err != nil {
				libs.HTTPError(w, http.StatusInternalServerError, "文件写入失败")
				return
			}
			err = CheckAddDesktop(path)
			if err != nil {
				log.Printf("Error adding file to desktop: %s", err.Error())
			}
			libs.SuccessMsg(w, "", "文件写入成功")
		}*/
}

/*
	func resEncode(w http.ResponseWriter, fileData []byte, filePwd string) {
		decryptData, err := libs.DecryptData(fileData[34:], []byte(filePwd))
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}

		//content := base64.StdEncoding.EncodeToString(decryptData)
		if len(decryptData)%8 == 0 {
			decryptData = append(decryptData, ' ')
		}
		// res := libs.APIResponse{Message: "加密文件读取成功", Data: decryptData}
		// json.NewEncoder(w).Encode(res)
		libs.SuccessMsg(w, decryptData, "加密文件读取成功")
	}
*/
func HandleSetFilePwd(w http.ResponseWriter, r *http.Request) {
	Pwd := r.Header.Get("Pwd")
	if Pwd == "" {
		err := libs.DeleteConfig("filePwd")
		if err != nil {
			libs.ErrorMsg(w, fmt.Sprintf("取消加密失败:%s", err.Error()))
			return
		}
		libs.SuccessMsg(w, nil, "取消加密成功")
		return
	}
	err := libs.SetConfigByName("filePwd", Pwd)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, nil, "设置密码成功")
}
