package model

import (
	"fmt"
	"godocms/pkg/dbfactory"
)

const TableNameUserDept = "user_dept"

// UserDept 租户部门表
type UserDept struct {
	ID         int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	GroupID    int32  `gorm:"column:group_id;default:1;comment:租户ID" json:"group_id"` // 租户ID
	Remark     string `gorm:"column:remark;comment:描述" json:"remark"`                 // 描述
	Name       string `gorm:"column:name;comment:部门名称" json:"name"`                   // 部门名称
	Pid        int32  `gorm:"column:pid;default:0;comment:上级部门ID" json:"pid"`         // 上级部门ID
	OrderNum   int32  `gorm:"column:order_num;comment:排序" json:"order_num"`           // 排序
	FromType   int32  `gorm:"column:from_type;comment:来源1钉钉2企业微信" json:"from_type"`   // 来源1钉钉2企业微信
	AddTime    int32  `gorm:"column:add_time;default:0;comment:添加时间" json:"add_time"` // 添加时间
	UpTime     int32  `gorm:"column:up_time;default:0;comment:更新时间" json:"up_time"`   // 更新时间
	Duty       string `gorm:"column:duty;comment:职责" json:"duty"`                     // 职责
	DingDeptID string `gorm:"column:ding_dept_id" json:"ding_dept_id"`                //  钉钉部门ID
}

// TableName UserDept's table name
func (*UserDept) TableName() string {
	return TableNameUserDept
}

type DeptNode struct {
	ID         int32       `json:"id"`
	Name       string      `json:"name"`
	PID        int32       `json:"pid"`
	DingDeptID string      `json:"ding_dept_id"`
	Child      []*DeptNode `json:"child"`
}

func (*UserDept) BuildTree() ([]*DeptNode, error) {
	var (
		deptList  []UserDept
		deptMap   = make(map[int32]*DeptNode) // 用于存储 ID 到 DeptNode 的映射
		rootNodes []*DeptNode                 // 根节点列表
	)

	// 查询部门列表
	result := dbfactory.Db.GetList("user_dept", nil, &deptList)
	if result != nil {
		return nil, result
	}

	// 创建部门节点并存入 deptMap
	for _, item := range deptList {
		deptNode := &DeptNode{
			ID:         item.ID,
			Name:       item.Name,
			PID:        item.Pid,
			DingDeptID: item.DingDeptID,
		}
		deptMap[item.ID] = deptNode

		// 如果是根节点 (PID == 0)，加入根节点列表
		if item.Pid == 0 {
			rootNodes = append(rootNodes, deptNode)
		}
	}

	// 现在根据 PID 将每个节点加入到其父节点的 Child 列表
	for _, node := range deptMap {
		if node.PID != 0 { // 排除根节点
			if parent, exists := deptMap[node.PID]; exists {
				parent.Child = append(parent.Child, node)
			}
		}
	}

	// 返回根节点列表
	return rootNodes, nil
}

func PrintDeptNodes(deptNodes []*DeptNode, indent string) {
	for _, node := range deptNodes {
		// 打印当前节点的信息
		fmt.Printf("%sID: %d, Name: %s, PID: %d, DingDeptID: %s\n", indent, node.ID, node.Name, node.PID, node.DingDeptID)

		// 如果有子节点，则递归打印子节点
		if len(node.Child) > 0 {
			PrintDeptNodes(node.Child, indent+"  ") // 增加缩进
		}
	}
}

// 深度优先遍历，寻找目标节点并收集从根节点到目标节点的链路和子节点 ID
func dfs(deptNodes []*DeptNode, targetID int32, path []int32, parentPath *[]int32, childrenPath *[]int32) bool {
	for _, node := range deptNodes {
		// 跟踪当前路径
		newPath := append(path, node.ID)

		// 如果找到目标节点，收集其路径和所有子节点的 ID
		if node.ID == targetID {
			*parentPath = newPath
			searchChild(node, childrenPath)
			return true
		}

		// 递归调用 DFS，继续遍历子节点
		if dfs(node.Child, targetID, newPath, parentPath, childrenPath) {
			return true
		}
	}

	return false
}

// 收集节点及其所有子节点的 ID
func searchChild(node *DeptNode, children *[]int32) {
	// 当前节点的 ID
	*children = append(*children, node.ID)
	// 遍历子节点
	for _, child := range node.Child {
		searchChild(child, children)
	}
}

// 获取目标节点的链路及其所有子节点的id
func (*UserDept) FindDeptPath(deptNodes []*DeptNode, targetID int32) (parent, children []int32, err error) {
	// parent 从根节点到目标节点的路径
	// children 目标节点及其所有子节点的id
	found := dfs(deptNodes, targetID, []int32{}, &parent, &children)
	if !found {
		return nil, nil, fmt.Errorf("节点 ID %d 未找到", targetID)
	}

	return parent, children, nil
}
