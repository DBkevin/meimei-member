package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type PointLogSearch struct {
	request.PageInfo
	MemberID   uint   `json:"memberId" form:"memberId"`
	Keyword    string `json:"keyword" form:"keyword"`
	ChangeType string `json:"changeType" form:"changeType"`
	SourceType string `json:"sourceType" form:"sourceType"`
}
