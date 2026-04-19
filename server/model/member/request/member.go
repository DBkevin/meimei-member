package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type MemberBaseInput struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Gender   string `json:"gender" form:"gender"`
	Birthday string `json:"birthday" form:"birthday"`
	Source   string `json:"source" form:"source"`
	Level    string `json:"level" form:"level"`
	Status   int    `json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
}

type CreateMemberReq struct {
	MemberBaseInput
}

type UpdateMemberReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
	MemberBaseInput
}

type MemberSearch struct {
	commonReq.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Name    string `json:"name" form:"name"`
	Phone   string `json:"phone" form:"phone"`
	Source  string `json:"source" form:"source"`
	Level   string `json:"level" form:"level"`
	Status  int    `json:"status" form:"status"`
}

type UpdateMemberStatusReq struct {
	ID     uint `json:"id" form:"id" binding:"required"`
	Status int  `json:"status" form:"status" binding:"required"`
}

type MemberOptionsReq struct {
	Keyword string `json:"keyword" form:"keyword"`
}
