package dingtalk

const (
	Api    = "https://oapi.dingtalk.com" // 钉钉ap域名
	NewApi = "https://api.dingtalk.com"  // 全新一代的钉钉API域名

	GetTokenKey       = "/gettoken"         // 获取access_token
	GetSSOTokenKey    = "/sso/gettoken"     // 获取微应用后台免登的access_token
	GetJsApiTicketKey = "/get_jsapi_ticket" // 获取jsapi_ticket

	MicroAppListKey              = "/microapp/list"           // 获取应用列表
	MicroAppVisibleScopesKey     = "/microapp/visible_scopes" // 获取应用可见范围
	UserMicroAppVisibleScopesKey = "/microapp/list_by_userid" // 获取员工可见的应用列表

	SuiteAccessToken        = "/service/get_suite_token"    // 获取第三方企业应用的
	CorpAccessToken         = "/service/get_corp_token"     // 服务商获取第三方应用授权企业的access_token
	GetAuthInfo             = "/service/get_auth_info"      // 获取企业授权信息
	ActivateSuiteKey        = "/service/activate_suite"     // 激活应用
	GetAgentKey             = "/service/get_agent"          // 获取授权应用的基本信息
	GetUnactiveCorpKey      = "/service/get_unactive_corp"  // 获取应用未激活的企业列表
	ReauthCorpKey           = "/service/reauth_corp"        // 重新授权未激活应用的企业
	GetCorpPermanentCodeKey = "/service/get_permanent_code" // 获取授权企业的永久授权码

	CreateUserKey             = "/topapi/v2/user/create"       // 创建用户
	DeleteUserKey             = "/topapi/v2/user/delete"       // 删除用户
	UpdateUserKey             = "/topapi/v2/user/update"       // 更新用户详情
	GetUserDetailKey          = "/topapi/v2/user/get"          // 获取用户详情
	GetDeptUserIdKey          = "/topapi/user/listid"          // 获取部门用户userid列表
	GetUserIdByUnionIdKey     = "/topapi/user/getbyunionid"    // 根据unionid获取userid
	GetUserIdByMobileKey      = "/topapi/v2/user/getbymobile"  // 根据手机号获取userid
	GetDeptSimpleUserKey      = "/topapi/user/listsimple"      // 获取部门用户基础信息
	GetDeptDetailUserKey      = "/topapi/v2/user/list"         // 获取部门用户详情
	GetOrgAdminUserKey        = "/topapi/user/listadmin"       // 获取管理员列表
	GetUserCanAccessAppletKey = "/user/can_access_microapp"    // 获取管理员的应用管理权限
	GetUserCountKey           = "/topapi/user/count"           // 获取员工人数
	GetInactiveUserKey        = "/topapi/inactive/user/v2/get" // 获取未登录钉钉的员工列表
	GetOrgAdminScopeKey       = "/topapi/user/get_admin_scope" // 获取管理员通讯录权限范围
	GetUserInfoByCodeKey      = "/topapi/v2/user/getuserinfo"  // 通过免登码获取用户信息
	GetSSOUserInfoKey         = "/sso/getuserinfo"             // 获取应用管理员的身份信息
	GetSNSUserInfoKey         = "/sns/getuserinfo_bycode"      // 根据sns临时授权码获取用户信息

	CreateDeptKey            = "/topapi/v2/department/create"           // 创建部门
	DeleteDeptKey            = "/topapi/v2/department/delete"           // 删除部门
	UpdateDeptKey            = "/topapi/v2/department/update"           // 更新部门
	GetDeptDetailKey         = "/topapi/v2/department/get"              // 获取部门详情
	GetDeptListKey           = "/topapi/v2/department/listsub"          // 获取部门详情
	GetSubDeptListKey        = "/topapi/v2/department/listsubid"        // 获取子部门列表
	GetParentDeptsByUserKey  = "/topapi/v2/department/listparentbyuser" // 查询指定用户的所有上级父部门路径
	GetParentDeptsByDeptKey  = "/topapi/v2/department/listparentbydept" // 查询部门的所有上级父部门路径
	GetIndustryDeptDetailKey = "/topapi/industry/department/get"        // 获取外部联系人详情
	GetIndustryDeptKey       = "/topapi/industry/department/list"       // 行业获取部门列表
	FetchDeptListKey         = "/department/list"                       // 支持递归获取部门列表

	CreateExtContactKey    = "/topapi/extcontact/create"          // 添加外部联系人
	DeleteExtContactKey    = "/topapi/extcontact/delete"          // 删除外部联系人
	UpdateExtContactKey    = "/topapi/extcontact/update"          // 更新外部联系人
	GetExtContactKey       = "/topapi/extcontact/list"            // 获取外部联系人列表
	GetExtContactLabelKey  = "/topapi/extcontact/listlabelgroups" // 获取外部联系人标签列表
	GetExtContactDetailKey = "/topapi/extcontact/get"             // 获取外部联系人详情

	GetRoleListKey               = "/topapi/role/list"               // 获取角色列表
	GetRoleUserListKey           = "/topapi/role/simplelist"         // 获取指定角色的员工列表
	GetRoleGroupKey              = "/topapi/role/getrolegroup"       // 获取角色组
	GetRoleDetailKey             = "/topapi/role/getrole"            // 获取角色详情
	CreateRoleGroupKey           = "/role/add_role_group"            // 创建角色组
	CreateRoleKey                = "/role/add_role"                  // 创建角色
	UpdateRoleKey                = "/role/update_role"               // 更新角色
	DeleteRoleKey                = "/topapi/role/deleterole"         // 删除角色
	RoleBatchAddUserKey          = "/topapi/role/addrolesforemps"    // 批量增加员工角色
	RoleBatchRemoveUserKey       = "/topapi/role/removerolesforemps" // 批量删除员工角色
	RoleUpdateUserManageScopeKey = "/topapi/role/scope/update"       // 设定角色成员管理范围

	CreateChatKey       = "/chat/create"                            // 创建群
	GetChatInfoKey      = "/chat/get"                               // 获取群会话
	UpdateChatKey       = "/chat/update"                            // 修改群会话
	SendChatMessageKey  = "/chat/send"                              // 发送消息到企业群
	GetChatReadUserKey  = "/chat/getReadList"                       // 查询群消息已读人员列表
	ChatFriendSwitchKey = "/topapi/chat/member/friendswitch/update" // 设置禁止群成员私聊
	GetChatQRCodeKey    = "/topapi/chat/qrcode/get"                 // 获取入群二维码链接
	ChatSetSubAdminKey  = "/topapi/chat/subadmin/update"            // 设置群管理员
	ChatSetUserNickKey  = "/topapi/chat/updategroupnick"            // 设置群成员昵称

	SendToConversationKey          = "/message/send_to_conversation"                      // 发送普通消息
	SendTemplateMessageKey         = "/topapi/message/corpconversation/sendbytemplate"    // 使用模板发送工作通知消息
	MessageProgressKey             = "/topapi/message/corpconversation/getsendprogress"   // 获取工作通知消息的发送进度
	GetMessageSendResultKey        = "/topapi/message/corpconversation/getsendresult"     // 获取工作通知消息的发送结果
	SendCorpConversationMessageKey = "/topapi/message/corpconversation/asyncsend_v2"      // 发送工作通知
	UpdateCorpConvMessageStatusKey = "/topapi/message/corpconversation/status_bar/update" // 更新工作通知状态栏
	RecallCorpConvMessageKey       = "/topapi/message/corpconversation/recall"            // 撤回工作通知消息

	SendRobotMsgKey = "/robot/send" // 发送机器人消息

	TranslateKey      = "/topapi/ai/mt/translate"          // 文本翻译
	OrcRecognizeKey   = "/topapi/ocr/structured/recognize" // OCR文字识别
	VoiceTranslateKey = "/topapi/asr/voice/translate"      // ASR 一句话语音识别

	corpConversation          = "/topapi/message/corpconversation"    // 工作通知
	SendCorpConversationKey   = corpConversation + "/asyncsend_v2"    // 发送工作通知
	GetSendProgressKey        = corpConversation + "/getsendprogress" // 获取工作通知消息的发送进度
	GetSendResultKey          = corpConversation + "/getsendresult"   // 获取工作通知消息的发送结果
	RecallCorpConversationKey = corpConversation + "/recall"          // 撤回工作通知

	MediaUploadKey = "/media/upload" // 上传媒体文件

	callback                   = "/call_back"                              // 事件回调模板
	RegisterCallBackKey        = callback + "/register_call_back"          // 注册事件回调接口
	GetCallBackKey             = callback + "/get_call_back"               // 查询事件回调接口
	UpdateCallBackKey          = callback + "/update_call_back"            // 更新事件回调接口
	GetCallBackFailedDetailKey = callback + "/get_call_back_failed_result" // 获取回调失败的结果
	DeleteCallBackKey          = callback + "/delete_call_back"            // 删除事件回调接口

	GetUserAccessToken     = "/v1.0/oauth2/userAccessToken" // 通过oauth2临时授权码获取用户Token
	GetContactUser         = "/v1.0/contact/users/%s"       // 获取用户通讯录个人信息
	CreateCooperateCorpKey = "/v1.0/contact/cooperateCorps" // 创建合作空间

	CreateDriveSpacesKey  = "/v1.0/drive/spaces"    // 新建空间
	DeleteDriveSpacesKey  = "/v1.0/drive/spaces/%s" // 删除空间
	GetDriveSpacesKey     = "/v1.0/drive/spaces"    // 获取空间列表
	GetDriveSpacesInfoKey = "/v1.0/drive/spaces/%s" // 获取空间信息

	GetDriveSpacesFilesKey    = "/v1.0/drive/spaces/%s/files"           // 查询文件（夹）列表
	GetDriveSpacesFileInfoKey = "/v1.0/drive/spaces/%s/files/%s"        // 查询文件（夹）信息
	CreateDriveSpacesFileKey  = "/v1.0/drive/spaces/%s/files"           // 添加文件（夹）
	DeleteDriveSpacesFilesKey = "/v1.0/drive/spaces/%s/files/%s"        // 删除文件（夹）
	MoveDriveSpacesFilesKey   = "/v1.0/drive/spaces/%s/files/%s/move"   // 移动文件（夹）
	RenameDriveSpacesFilesKey = "/v1.0/drive/spaces/%s/files/%s/rename" // 修改文件（夹）名

	GetStorageSpacesFilesKey            = "/v1.0/storage/spaces/%s/dentries"                        // 获取存储空间列表
	GetStorageSpacesFileInfoKey         = "/v1.0/storage/spaces/%s/dentries/%s/query"               // 获取文件信息
	GetStorageSpacesFileDownloadInfoKey = "/v1.0/storage/spaces/%s/dentries/%s/downloadInfos/query" // 获取文件下载信息

	GetDriveSpacesFilesDownloadInfoKey = "/v1.0/drive/spaces/%s/files/%s/downloadInfos" // 获取文件下载信息
	GetDriveSpacesFilesUploadInfoKey   = "/v1.0/drive/spaces/%s/files/%s/uploadInfos"   // 获取文件上传信息

	GetDriveSpacesFilesPermissionsKey = "/v1.0/drive/spaces/%s/files/%s/permissions" // 获取权限列表
	AddDriveSpacesFilesPermissionsKey = "/v1.0/drive/spaces/%s/files/%s/permissions" // 添加权限

	GetDingIndexKey = "/v1.0/contact/dingIndexs" // 获取企业最新钉钉指数信息

	CreateTodoKey            = "/v1.0/todo/users/%s/tasks"                   // 新增钉钉待办任务
	GetTodoDetailKey         = "/v1.0/todo/users/%s/tasks/%s"                // 获取钉钉待办任务详情
	DeleteTodoKey            = "/v1.0/todo/users/%s/tasks/%s"                // 删除钉钉待办任务
	UpdateTodoKey            = "/v1.0/todo/users/%s/tasks/%s"                // 更新钉钉待办任务
	UpdateTodoDoneKey        = "/v1.0/todo/users/%s/tasks/%s/executorStatus" // 更新钉钉待办执行者状态
	GetTodoListBySourceIdKey = "/v1.0/todo/users/%s/tasks/sources/%s"        // 根据sourceId获取钉钉待办任务详情
	GetTodoListKey           = "/v1.0/todo/users/%s/org/tasks/query"         // 查询企业下用户待办列表

	GetHrmEmployeeKey           = "/topapi/smartwork/hrm/employee/queryonjob"     // 获取在职员工列表
	GetHrmToBeHiredEmployeeKey  = "/topapi/smartwork/hrm/employee/querypreentry"  // 获取待入职员工列表
	GetHrmResignEmployeeKey     = "/topapi/smartwork/hrm/employee/querydimission" // 获取离职员工列表
	GetHrmResignEmployeeInfoKey = "/topapi/smartwork/hrm/employee/listdimission"  // 获取员工离职信息
	HrmCreateEmployeeKey        = "/topapi/smartwork/hrm/employee/addpreentry"    // 添加企业待入职员工
	GetHrmEmpLeaveRecordsKey    = "/v1.0/contact/empLeaveRecords"                 // 获取时间范围内离职员工记录

	GetHrmFieldKey            = "/topapi/smartwork/hrm/employee/field/grouplist" // 获取花名册字段组详情
	GetHrmEmployeeFieldKey    = "/topapi/smartwork/hrm/employee/v2/list"         // 获取员工花名册字段信息
	UpdateHrmEmployeeFieldKey = "/topapi/smartwork/hrm/employee/v2/update"       // 更新员工花名册信息
	GetHrmMetaKey             = "/topapi/smartwork/hrm/roster/meta/get"          // 获取花名册元数据

	GetAttendanceGroupsKey          = "/topapi/attendance/getsimplegroups"       // 批量获取考勤组详情
	GetAttendanceUserGroupKey       = "/topapi/attendance/getusergroup"          // 获取用户考勤组
	GetAttendanceGroupMinimalismKey = "/topapi/attendance/group/minimalism/list" // 批量获取考勤组摘要
	GetAttendanceGroupDetailKey     = "/topapi/attendance/group/query"           // 获取考勤组详情
	SearchAttendanceGroupKey        = "/topapi/attendance/group/search"          // 搜索考勤组摘要
	CreateAttendanceGroupKey        = "/topapi/attendance/group/add"             // 创建考勤组
)

const DingDingAccessTokenCachekey = "ding_ding_access_token"
