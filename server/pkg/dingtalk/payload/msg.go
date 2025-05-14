package payload

import (
	"godocms/pkg/dingtalk/payload/message"
	"sort"
	"strconv"
	"strings"
)

// CorpConvMessage 工作通知
type CorpConvMessage struct {
	// 发送消息时使用的微应用的ID。
	AgentId int `json:"agent_id" validate:"required"`

	// 接收者的企业内部用户的userId列表，最大用户列表长度100。
	UserIds []string `json:"-" validate:"lte=100"`

	UserIdList string `json:"userid_list,omitempty"`

	// 接收者的部门id列表，最大列表长度20。接收者是部门Id下包括子部门下的所有用户。
	DeptIdList string `json:"dept_id_list,omitempty"`

	DeptIds []int `json:"-" validate:"lte=20"`

	// 是否发送给企业全部用户。当设置为false时必须指定userid_list或dept_id_list其中一个参数的值。
	All *bool `json:"to_all_user,omitempty"`

	Msg message.Message `json:"msg,omitempty" validate:"required"`
}

type CorpConvMessageResponse struct {
	Response

	TaskId int `json:"task_id"`
}

type corpConvMessageBuilder struct {
	cm *CorpConvMessage
}

func NewCorpConvMessage(msg message.Message) *corpConvMessageBuilder {
	return &corpConvMessageBuilder{cm: &CorpConvMessage{Msg: msg}}
}

func (sb *corpConvMessageBuilder) SetAgentId(agentId int) *corpConvMessageBuilder {
	sb.cm.AgentId = agentId
	return sb
}

func (sb *corpConvMessageBuilder) SetUserIds(userId string, userIds ...string) *corpConvMessageBuilder {
	sb.cm.UserIds = append(userIds, userId)
	return sb
}

func (sb *corpConvMessageBuilder) SetUsers(userIds []string) *corpConvMessageBuilder {
	sb.cm.UserIds = userIds
	return sb
}

func (sb *corpConvMessageBuilder) SetDeptIds(deptId int, deptIds ...int) *corpConvMessageBuilder {
	sb.cm.DeptIds = append(deptIds, deptId)
	return sb
}

func (sb *corpConvMessageBuilder) SetAllUser(all bool) *corpConvMessageBuilder {
	sb.cm.All = &all
	return sb
}

func (sb *corpConvMessageBuilder) Build() *CorpConvMessage {
	cm := sb.cm
	cm.DeptIds = removeIntDuplicates(cm.DeptIds)
	cm.DeptIdList = strings.Join(removeIntDuplicatesToString(cm.DeptIds), ",")

	cm.UserIds = removeStringDuplicates(cm.UserIds)
	cm.UserIdList = strings.Join(cm.UserIds, ",")
	return cm
}

// removeIntDuplicates 去除重复的item
func removeIntDuplicates(item []int) (ids []int) {
	if len(item) <= 0 {
		return ids
	}
	sort.Ints(item)
	for i, id := range item {
		if (i >= 1 && id == item[i-1]) || id <= 0 {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

// removeIntDuplicatesToString 去除重复的item
func removeIntDuplicatesToString(item []int) (ids []string) {
	if len(item) <= 0 {
		return ids
	}
	sort.Ints(item)
	for i, id := range item {
		if (i >= 1 && id == item[i-1]) || id <= 0 {
			continue
		}
		ids = append(ids, strconv.Itoa(id))
	}
	return ids
}

// removeStringDuplicates 去除重复的item
func removeStringDuplicates(item []string) (ids []string) {
	if len(item) <= 0 {
		return ids
	}
	sort.Strings(item)
	for i, id := range item {
		if (i >= 1 && id == item[i-1]) || len(id) <= 0 {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}
