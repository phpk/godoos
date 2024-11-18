package model

import (
	"fmt"
	"godo/libs"
	"net/http"
	"os"
)

func ConvertOllama(w http.ResponseWriter, r *http.Request, req ReqBody) {
	modelFile := "FROM " + req.Paths[0] + "\n"
	modelFile += `TEMPLATE """` + req.Info["template"].(string) + `"""`
	if parameters, ok := req.Info["parameters"].([]interface{}); ok {
		for _, param := range parameters {
			if strParam, ok := param.(string); ok {
				modelFile += "\nPARAMETER " + strParam
			} else {
				// 处理非字符串的情况，根据需要可以选择忽略或报告错误
				fmt.Fprintf(os.Stderr, "Unexpected parameter type: %T\n", param)
			}
		}
	}

	url := GetOllamaUrl() + "/api/create"
	postParams := map[string]string{
		"name":      req.Model,
		"modelfile": modelFile,
	}
	ForwardHandler(w, r, postParams, url, "POST")
	modelDir, err := GetModelDir(req.Model)
	if err != nil {
		libs.ErrorMsg(w, "GetModelDir")
		return
	}

	// modelFilePath := filepath.Join(modelDir, "Modelfile")
	// if err := os.WriteFile(modelFilePath, []byte(modelFile), 0644); err != nil {
	// 	ErrMsg("WriteFile", err, w)
	// 	return
	// }
	err = os.RemoveAll(modelDir)
	if err != nil {
		libs.ErrorMsg(w, "Error removing directory")
		return
	}
}
