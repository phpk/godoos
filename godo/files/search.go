package files

import (
	"encoding/json"
	"godo/libs"
	"godo/office"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func HandleSarch(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	query := r.URL.Query().Get("query")
	if err := validateFilePath(path); err != nil {
		libs.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	if query == "" {
		libs.HTTPError(w, http.StatusBadRequest, "query is empty")
		return
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	searchPath := filepath.Join(basePath, path)
	var osFileInfos []*OsFileInfo

	// Initialize tokenizer
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}

	// Tokenize the query
	queryTokens := t.Tokenize(query)
	queryWords := make([]string, 0, len(queryTokens))
	for _, token := range queryTokens {
		queryWords = append(queryWords, token.Surface)
	}

	err = filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		osFileInfo, err := GetFileInfo(path, "", "")
		if err != nil {
			return err
		}
		if !info.IsDir() {
			doc, err := office.GetDocument(path)
			if err == nil {
				osFileInfo.Content = doc.Content

				// Tokenize the file content
				contentTokens := t.Tokenize(osFileInfo.Content)
				contentWords := make([]string, 0, len(contentTokens))
				for _, token := range contentTokens {
					contentWords = append(contentWords, token.Surface)
				}

				// Check if any query word is in the content words
				for _, queryWord := range queryWords {
					if containsIgnoreCase(contentWords, queryWord) {
						//osFileInfo.Path = strings.TrimPrefix(path, basePath)
						osFileInfos = append(osFileInfos, osFileInfo)
						break
					}
				}
			}
		} else {
			// Check if the path or filename contains the query
			if strings.Contains(strings.ToLower(path), strings.ToLower(query)) {
				//osFileInfo.Path = strings.TrimPrefix(path, basePath)
				osFileInfos = append(osFileInfos, osFileInfo)
			}
		}

		return nil
	})
	if err != nil {
		libs.HTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for i, osFileInfo := range osFileInfos {
		osFileInfos[i].Path = strings.TrimPrefix(osFileInfo.Path, basePath)
		osFileInfos[i].OldPath = strings.TrimPrefix(osFileInfo.OldPath, basePath)
	}
	res := libs.APIResponse{
		Message: "Directory read successfully.",
		Data:    osFileInfos,
	}
	json.NewEncoder(w).Encode(res)
}

// Helper function to check if a slice contains a string (case insensitive)
func containsIgnoreCase(slice []string, item string) bool {
	itemLower := strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == itemLower {
			return true
		}
	}
	return false
}
