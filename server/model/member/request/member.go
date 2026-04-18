package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type MemberSearch struct {
	request.PageInfo
	Mobile      string `json:"mobile" form:"mobile"`
	Nickname    string `json:"nickname" form:"nickname"`
	RealName    string `json:"realName" form:"realName"`
	MemberLevel string `json:"memberLevel" form:"memberLevel"`
	Status      string `json:"status" form:"status"`
}

type UpdateMemberStatusReq struct {
	ID     uint   `json:"id" form:"id" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
}

type MemberOptionsReq struct {
	Keyword string `json:"keyword" form:"keyword"`
}
