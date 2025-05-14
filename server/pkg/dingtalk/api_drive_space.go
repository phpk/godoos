package dingtalk

import (
	"fmt"
	"godocms/common"
	"godocms/pkg/dingtalk/payload"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GetDriveSpaces 获取空间列表
func (ding *DingTalk) GetDriveSpaces(unionId string, spaceType payload.SpaceType, token string,
	size int) (rsp payload.GetDriveSpaces, err error,
) {
	query := url.Values{}
	query.Set("spaceType", string(spaceType))
	query.Set("unionId", unionId)
	query.Set("nextToken", token)
	query.Set("maxResults", strconv.Itoa(size))

	return rsp, ding.Request(http.MethodGet, GetDriveSpacesKey, query, nil, &rsp)
}

// GetStorageSapceFiles 获取空间文件列表
func (ding *DingTalk) GetStorageSapceFiles(req *payload.GetStorageSpacesFilesReq) (rsp payload.FileDentriesResponse, err error) {
	query := url.Values{}
	query.Set("parentId", req.ParentId)
	query.Set("nextToken", req.NextToken)
	query.Set("maxResults", strconv.Itoa(req.Size))
	query.Set("orderBy", req.OrderBy)
	query.Set("order", string(req.Order))
	query.Set("withThumbnail", strconv.FormatBool(req.WithThumbnail))
	query.Set("unionId", req.UnionId)

	return rsp, ding.Request(http.MethodGet, fmt.Sprintf(GetStorageSpacesFilesKey, req.SpaceId), query, nil, &rsp)
}

// GetStorageFileInfo 获取文件信息
func (ding *DingTalk) GetStorageFileInfo(spaceId, dentryId, unionId string) (rsp payload.GetStorageFileInfoResponse, err error) {
	query := url.Values{}
	query.Set("unionId", unionId)

	return rsp, ding.Request(http.MethodPost, fmt.Sprintf(GetStorageSpacesFileInfoKey, spaceId, dentryId), query, nil, &rsp)
}

// GetStorageSpacesFileDownloadInfo 获取文件下载信息
func (ding *DingTalk) GetStorageSpacesFileDownloadInfo(spaceId, fileId, unionId string) (rsp payload.GetStorageSpacesFileDownloadInfo, err error) {
	query := url.Values{}
	query.Set("unionId", unionId)
	return rsp, ding.Request(http.MethodPost, fmt.Sprintf(GetStorageSpacesFileDownloadInfoKey, spaceId, fileId),
		query, nil, &rsp)
}

type FileNode struct {
	ID       string      `json:"id"`
	SpaceID  string      `json:"space_id"`
	ParentID string      `json:"parent_id"`
	Type     string      `json:"type"`
	Name     string      `json:"name"`
	Size     int64       `json:"size"`
	Path     string      `json:"path"`
	Status   string      `json:"status"`             // NORMAL DELETED EXPIRED
	Children []*FileNode `json:"children,omitempty"` // 子节点
}

// 获取用户钉盘文件列表
func (ding *DingTalk) GetUserDingSpace(unionid string) ([]*FileNode, error) {
	var nextToken string
	var fileTree []*FileNode

	for {
		rsp, err := ding.GetDriveSpaces(unionid, payload.Org, nextToken, 10)
		if err != nil {
			return nil, err
		}
		for _, space := range rsp.Spaces {
			fmt.Printf("space: %+v\n", space)
			rootNode := &FileNode{
				ID:       "0", // 假设根目录的 ParentId 为 "0"
				SpaceID:  space.SpaceId,
				ParentID: "0",
				Type:     payload.DentryFolder,
				Name:     space.Name, // 根目录的名称
			}

			if err := ding.buildFileTree(unionid, space.SpaceId, "0", rootNode); err != nil {
				return nil, err
			}

			fileTree = append(fileTree, rootNode)
		}

		nextToken = rsp.Token
		if nextToken == "" {
			break
		}
	}

	return fileTree, nil
}

func (ding *DingTalk) buildFileTree(unionid, spaceID, parentID string, parentNode *FileNode) error {
	var nextToken string
	for {
		files, err := ding.GetStorageSapceFiles(payload.NewGetStorageSpacesFilesReq(spaceID, unionid, 10,
			payload.WithStorageFilesParentId(parentID), payload.WithStorageFilesNextToken(nextToken)))
		if err != nil {
			return err
		}

		for _, v := range files.Dentries {
			// 创建新的文件或文件夹节点
			node := &FileNode{
				ID:       v.ID,
				SpaceID:  v.SpaceId,
				ParentID: v.ParentId,
				Type:     v.Type,
				Name:     v.Name,
			}

			// 如果是文件夹，递归构建文件树
			if v.Type == payload.DentryFolder {
				if err := ding.buildFileTree(unionid, spaceID, v.ID, node); err != nil {
					return err
				}
			}

			// 将当前节点添加到父节点的子节点列表
			parentNode.Children = append(parentNode.Children, node)
		}

		nextToken = files.Token
		if nextToken == "" {
			break
		}
	}

	return nil
}

type SpaceChannel struct {
	node    *FileNode
	unionid string
	path    []string
}

func (ding *DingTalk) JoinDingSpaceFileDwonloadTask(unionid string, files []*FileNode) error {
	signal := make(chan *SpaceChannel)
	fileLog := make(chan *WriteFileLog)
	done := make(chan struct{}) // 用于通知文件处理完毕
	fileCache := make([]*WriteFileLog, 0)

	go ding.listenToChannel(signal, fileLog)

	var wg sync.WaitGroup
	wg.Add(1)

	// 处理文件下载
	go func() {
		defer wg.Done()
		dfs(unionid, files, signal, []string{})
		signal <- nil // 发送 EOF 信号
		close(signal)
		close(fileLog)
	}()

	// 处理下载日志文件
	go func() {
		defer close(done) // 完成时通知
		for msg := range fileLog {
			if msg == nil {
				fmt.Println("-------------------------日志读取结束--------------------------")
				return
			}
			slog.Info("Write file log", "file", msg.FileName, "path", msg.Path, "remote", msg.Remote, "time", msg.Time)
			fileCache = append(fileCache, msg)
		}
	}()

	// 等待所有工作完成
	wg.Wait()
	<-done // 等待 fileLog goroutine 完成

	// 写入缓存
	common.Cache.Delete(fmt.Sprintf("ding_space_%s", unionid))
	common.Cache.Set(fmt.Sprintf("ding_space_%s", unionid), fileCache, 24*60)
	fmt.Println("-------------------------日志文件写入完成--------------------------")

	return nil
}

func dfs(unionid string, nodes []*FileNode, ch chan *SpaceChannel, path []string) {
	for _, node := range nodes {
		if node == nil {
			continue
		}
		// 将当前节点的 ID 添加到路径中
		newPath := append([]string(nil), path...) // 创建一个新切片来避免修改原路径
		newPath = append(newPath, node.Name)
		fmt.Println("parent path: ", newPath, "this node: ", node.ID, " name: ", node.Name)
		// 发送当前节点到 channel
		ch <- &SpaceChannel{node: node, unionid: unionid, path: newPath}

		// 如果有子节点，递归遍历子节点
		if len(node.Children) > 0 {
			dfs(unionid, node.Children, ch, newPath)
		}
	}
}

// 监听 channel，处理节点
func (ding *DingTalk) listenToChannel(ch chan *SpaceChannel, fileLog chan *WriteFileLog) {
	for signal := range ch {
		if signal == nil {
			// 处理 EOF 信号
			slog.Info("EOF signal received")
			break
		}
		if signal.node.Type == payload.DentryFolder {
			folderPath := filepath.Join(signal.path...)
			// 新建文件夹
			realFolderPath := filepath.Join("data", "upload", "dingspace", folderPath)
			err := os.MkdirAll(realFolderPath, 0755)
			if err != nil {
				fmt.Println("Error creating directory:", err)
			}
			slog.Info("New folder", slog.Any("folderPath", realFolderPath))
		}
		if signal.node.Type == payload.DentryFile {
			// 获取文件信息,下载文件
			filePath := filepath.Join(signal.path...)
			realFilePath := filepath.Join("data", "upload", "dingspace", filePath)
			resp, err := ding.GetStorageSpacesFileDownloadInfo(signal.node.SpaceID, signal.node.ID, signal.unionid)
			if err != nil {
				slog.Error("GetStorageFileInfo:", slog.Any("err", err))
			}
			headers := resp.HeaderSignatureInfo.Headers

			for _, url := range resp.HeaderSignatureInfo.ResourceUrls {
				wfl := &WriteFileLog{
					Action:   "download file info",
					FileName: signal.node.Name,
					Header:   headers,
					Path:     realFilePath,
					Remote:   url,
					Time:     time.Now().Local().Format("2006-01-02 15:04:05"),
					Unionid:  signal.unionid,
					SpaceID:  signal.node.SpaceID,
					FileID:   signal.node.ID,
				}
				fileLog <- wfl
			}
		}
	}
}

// 打印文件树结构
func printFileTree(fileTree []*FileNode, level int) {
	for _, node := range fileTree {
		// 打印当前节点信息
		fmt.Printf("%s[%s] (ID: %s, SpaceID: %s, Type: %s)\n", getIndentation(level), node.Name, node.ID, node.SpaceID, node.Type)
		// 递归打印子节点
		printFileTree(node.Children, level+1)
	}
}

// 获取缩进字符串，方便打印树形结构
func getIndentation(level int) string {
	return strings.Repeat(" ", level*2) // 每级缩进 2 个空格
}
