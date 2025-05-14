package payload

type SpaceType string

const (
	// Personal 私人空间
	Personal = SpaceType("personal")

	// Org 企业空间
	Org = SpaceType("org")
)

type GetDriveSpaces struct {
	Response

	Spaces []struct {
		// 空间ID
		SpaceId string `json:"spaceId"`

		// 空间名称
		Name string `json:"spaceName"`

		// 空间类型
		Type string `json:"spaceType"`

		// 空间总额度
		Quota int `json:"quota"`

		// 空间已使用额度
		UsedQuota int `json:"usedQuota"`

		//授权模式，取值：
		//
		//acl：acl授权
		//custom：自定义授权
		PermissionMode string `json:"permissionMode"`

		// 创建时间
		CreateTime string `json:"createTime"`

		// 修改时间
		ModifyTime string `json:"modifyTime"`
	} `json:"spaces"`

	Token string `json:"nextToken"`
}

type GetStorageSpacesFilesReq struct {
	// 钉盘空间ID
	SpaceId string `json:"spaceId" validate:"required"`

	// 用户unionId
	UnionId string `json:"unionId" validate:"required"`

	// 父目录ID
	ParentId string `json:"parentId,omitempty"`

	// 分页游标
	NextToken string `json:"nextToken,omitempty"`

	// 分页大小
	Size int `json:"maxResults"`

	// MODIFIED_TIME：最后修改时间，默认值
	// CREATE_TIME：创建时间
	// NAME：名称
	// SIZE：大小
	OrderBy string `json:"orderBy,omitempty"`

	// 排序 ASC：升序, DESC：降序，默认值
	Order string `json:"order,omitempty"`

	// 是否获取文件缩略图临时链接。按需获取，会影响接口耗时。
	WithThumbnail bool `json:"withThumbnail,omitempty"`
}

type Option func(*GetStorageSpacesFilesReq)

func NewGetStorageSpacesFilesReq(spaceId, unionId string, size int, opts ...Option) *GetStorageSpacesFilesReq {
	req := &GetStorageSpacesFilesReq{
		SpaceId: spaceId,
		UnionId: unionId,
		Size:    size,
		OrderBy: "MODIFIED_TIME", // Default value
		Order:   "ASC",           // Default value
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// Option to set ParentId
func WithStorageFilesParentId(parentId string) Option {
	return func(req *GetStorageSpacesFilesReq) {
		req.ParentId = parentId
	}
}

// Option to set NextToken
func WithStorageFilesNextToken(nextToken string) Option {
	return func(req *GetStorageSpacesFilesReq) {
		req.NextToken = nextToken
	}
}

// Option to set OrderBy
func WithStorageFilesOrderBy(orderBy string) Option {
	return func(req *GetStorageSpacesFilesReq) {
		req.OrderBy = orderBy
	}
}

// Option to set Order
func WithStorageFilesOrder(order string) Option {
	return func(req *GetStorageSpacesFilesReq) {
		req.Order = order
	}
}

// Option to set WithThumbnail
func WithStorageFilesThumbnail(withThumbnail bool) Option {
	return func(req *GetStorageSpacesFilesReq) {
		req.WithThumbnail = withThumbnail
	}
}

type FileDentriesResponse struct {
	Response

	// 下一页的游标，为空字符串则表示分页结束
	Token string `json:"nextToken"`

	// 文件列表
	Dentries []Dentry `json:"dentries"`
}

const (
	DentryFile   = "FILE"
	DentryFolder = "FOLDER"
)

const (
	// 文件状态
	NormalFile  = "NORMAL"
	DeletedFile = "DELETED"
	ExpiredFile = "EXPIRED"
)

// Dentry represents each entry (file/folder) in the response
type Dentry struct {
	ID            string          `json:"id"`
	SpaceId       string          `json:"spaceId"`
	ParentId      string          `json:"parentId"`
	Type          string          `json:"type"`
	Name          string          `json:"name"`
	Size          int             `json:"size"`
	Path          string          `json:"path"`
	Version       int             `json:"version"`
	Status        string          `json:"status"`
	Extension     string          `json:"extension"`
	CreatorId     string          `json:"creatorId"`
	ModifierId    string          `json:"modifierId"`
	CreateTime    string          `json:"createTime"`
	ModifiedTime  string          `json:"modifiedTime"`
	Properties    map[string]bool `json:"properties"`    // readOnly: true/false
	AppProperties map[string]any  `json:"appProperties"` // name value visibility
	UUID          string          `json:"uuid"`
	PartitionType string          `json:"partitionType"`
	StorageDriver string          `json:"storageDriver"`
	Thumbnail     *Thumbnail      `json:"thumbnail,omitempty"`
}

// Thumbnail represents the image thumbnail details for a dentry
type Thumbnail struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

// GetStorageResponse represents the structure of the API response
type GetStorageResponse struct {
	Dentries  []Dentry `json:"dentries"`
	NextToken string   `json:"nextToken"`
}

type GetStorageFileInfoResponse struct {
	Response
	Dentry Dentry `json:"dentry"`
}

type GetStorageSpacesFileDownloadInfo struct {
	Response
	Protocol            string              `json:"protocol"`
	HeaderSignatureInfo HeaderSignatureInfo `json:"headerSignatureInfo"`
}

// HeaderSignatureInfo represents the details of the header signature information
type HeaderSignatureInfo struct {
	ResourceUrls         []string          `json:"resourceUrls"`
	Headers              map[string]string `json:"headers"`
	ExpirationSeconds    int               `json:"expirationSeconds"`
	Region               string            `json:"region"`
	InternalResourceUrls []string          `json:"internalResourceUrls"`
}
