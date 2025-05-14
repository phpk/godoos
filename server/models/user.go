package model

const TableNameUser = "user"

// User 租户用户表
type User struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键" json:"id"` // 主键
	GroupID   int32  `gorm:"column:group_id;default:1;comment:租户ID" json:"group_id"`       // 租户ID
	Username  string `gorm:"column:username;not null;comment:用户姓名--例如张三" json:"username"`  // 用户姓名--例如张三
	Nickname  string `gorm:"column:nickname;not null;comment:用户真实姓名" json:"nickname"`      // 用户真实姓名
	Sex       int32  `gorm:"column:sex;comment:性别，0：男，1：女" json:"sex"`                     // 性别，0：男，1：女
	Password  string `gorm:"column:password;comment:登陆密码" json:"password"`                 // 登陆密码
	Salt      string `gorm:"column:salt;not null;comment:salt校验" json:"salt"`              // salt校验
	Email     string `gorm:"column:email;comment:电子邮箱" json:"email"`                       // 电子邮箱
	Phone     string `gorm:"column:phone;comment:手机号码" json:"phone"`                       // 手机号码
	Status    int32  `gorm:"column:status;comment:状态，0：正常，1：删除，2封禁" json:"status"`         // 状态，0：正常，1：删除，2封禁
	Desc      string `gorm:"column:desc;comment:用户描述信息" json:"desc"`                       // 用户描述信息
	Remark    string `gorm:"column:remark;comment:备注" json:"remark"`                       // 备注
	AddTime   int32  `gorm:"column:add_time;default:0;comment:添加时间" json:"add_time"`       // 添加时间
	UpTime    int32  `gorm:"column:up_time;default:0;comment:更新时间" json:"up_time"`         // 更新时间
	Avatar    string `gorm:"column:avatar;comment:头像地址" json:"avatar"`                     // 头像地址
	JobNumber string `gorm:"column:job_number;comment:员工工号" json:"job_number"`             // 员工工号
	WorkPlace string `gorm:"column:work_place;comment:办公地点" json:"work_place"`             // 办公地点
	LoginNum  int32  `gorm:"column:login_num;default:1" json:"login_num"`
	HiredDate int64  `gorm:"column:hired_date;default:0;comment:入职时间" json:"hired_date"` // 入职时间
	LoginIP   string `gorm:"column:login_ip;comment:登录ip" json:"login_ip"`
	DeptId    int32  `gorm:"column:dept_id;comment:部门id" json:"dept_id"`
	RoleId    int32  `gorm:"column:role_id;comment:角色id" json:"role_id"`
	UseSpace  int32  `gorm:"column:use_space;comment:已用空间" json:"use_space"`
	DeptAuth  string `gorm:"column:dept_auth;comment:部门权限-1领导0自己多个数字加，为可查看人" json:"dept_auth"`

	ThirdUserID string `gorm:"-" json:"third_user_id"` // 额外字段不入库，注册使用
	UnionID     string `gorm:"-" json:"union_id"`      // 额外字段不入库，注册使用
	Platform    string `gorm:"-" json:"patform"`       // 额外字段不入库，注册使用
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
