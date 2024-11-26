package files

import (
	"fmt"
	"godo/libs"
	"io"
	"os"
	"testing"
)

func Test_WriteFile(t *testing.T) {
	filePath := "test.txt"
	content := "hello world"
	configPwd := "b854f92a0bb8462ad75239152081c12b"
	file, err := os.Create(filePath)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("@%s@", configPwd))
	if err != nil {
		t.Error(err)
	}
	encryData, err := libs.EncryptData([]byte(content), []byte(configPwd))
	if err != nil {
		t.Error(err)
	}
	_, err = file.Write(encryData)
	if err != nil {
		t.Error(err)
	}
}

func Test_ReadFile(t *testing.T) {
	filePath := "test.txt"
	pwd := "b854f92a0bb8462ad75239152081c12b"
	file, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}
	if pwd == "" {
		fmt.Println(string(data))
		return
	}
	realPwd := data[1:33]
	if pwd != string(realPwd) {
		t.Error("密码错误")
		return
	}
	content := data[34:]
	decryData, err := libs.DecryptData(content, []byte(pwd))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(decryData))
}
