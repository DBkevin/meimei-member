package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type PointAccountSearch struct {
	commonReq.PageInfo
	MemberID uint   `json:"memberId" form:"memberId"`
	Keyword  string `json:"keyword" form:"keyword"`
}

type GetPointAccountReq struct {
	MemberID uint `json:"memberId" form:"memberId" binding:"required"`
}

type AdjustPointsReq struct {
	MemberID uint   `json:"memberId" form:"memberId" binding:"required"`
	Points   int64  `json:"points" form:"points" binding:"required"`
	Remark   string `json:"remark" form:"remark"`
}
