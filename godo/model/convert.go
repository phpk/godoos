package model

import (
	"godo/libs"
	"net/http"
	"os"
	"strings"
)

func ConvertOllama(w http.ResponseWriter, r *http.Request, req ReqBody) {
	modelFile := "FROM " + req.Info.Path[0] + "\n"
	modelFile += `TEMPLATE """` + req.Info.Template + `"""`
	if req.Info.Parameters != "" {
		parameters := strings.Split(req.Info.Parameters, "\n")
		for _, param := range parameters {
			modelFile += "\nPARAMETER " + param
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
