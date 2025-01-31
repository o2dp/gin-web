package request

import (
	"gin-web/pkg/response"
)

// 获取消息列表结构体
type MessageRequestStruct struct {
	ToUserId          uint   `json:"toUserId"`
	Title             string `json:"title" form:"title"`
	Content           string `json:"content" form:"content"`
	Type              *uint  `json:"type" form:"type"`
	Status            *uint  `json:"status" form:"status"`
	response.PageInfo        // 分页参数
}

// 推送消息结构体
type PushMessageRequestStruct struct {
	FromUserId       uint
	Type             *ReqUint `json:"type" form:"type" validate:"required"`
	ToUserIds        []uint   `json:"toUserIds" form:"toUserIds"`
	ToRoleIds        []uint   `json:"toRoleIds" form:"toRoleIds"`
	Title            string   `json:"title" form:"title" validate:"required"`
	Content          string   `json:"content" form:"content" validate:"required"`
	IdempotenceToken string   `json:"idempotenceToken" form:"idempotenceToken"`
}

// 翻译需要校验的字段名称
func (s PushMessageRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Type"] = "消息类型"
	m["Title"] = "消息标题"
	m["Content"] = "消息内容"
	return m
}
