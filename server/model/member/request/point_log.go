package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type PointTransactionSearch struct {
	commonReq.PageInfo
	MemberID uint   `json:"memberId" form:"memberId"`
	Keyword  string `json:"keyword" form:"keyword"`
	Type     string `json:"type" form:"type"`
	RefType  string `json:"refType" form:"refType"`
}
