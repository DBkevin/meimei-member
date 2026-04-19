package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type RedemptionOrderSearch struct {
	commonReq.PageInfo
	Keyword   string `json:"keyword" form:"keyword"`
	MemberID  uint   `json:"memberId" form:"memberId"`
	ProductID uint   `json:"productId" form:"productId"`
	Status    int    `json:"status" form:"status"`
}

type CreateRedemptionOrderReq struct {
	MemberID      uint   `json:"memberId" form:"memberId" binding:"required"`
	ProductID     uint   `json:"productId" form:"productId" binding:"required"`
	Quantity      int64  `json:"quantity" form:"quantity" binding:"required"`
	ReceiverName  string `json:"receiverName" form:"receiverName"`
	ReceiverPhone string `json:"receiverPhone" form:"receiverPhone"`
	Remark        string `json:"remark" form:"remark"`
}

type OperateRedemptionOrderReq struct {
	ID     uint   `json:"id" form:"id" binding:"required"`
	Remark string `json:"remark" form:"remark"`
}
