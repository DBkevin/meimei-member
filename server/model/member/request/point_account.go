package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type PointAccountSearch struct {
	request.PageInfo
	MemberID uint   `json:"memberId" form:"memberId"`
	Keyword  string `json:"keyword" form:"keyword"`
}

type AdjustPointsReq struct {
	MemberID uint   `json:"memberId" form:"memberId" binding:"required"`
	Points   int64  `json:"points" form:"points" binding:"required"`
	Remark   string `json:"remark" form:"remark"`
}
