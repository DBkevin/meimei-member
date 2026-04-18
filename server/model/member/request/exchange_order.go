package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type ExchangeOrderSearch struct {
	request.PageInfo
	Keyword  string `json:"keyword" form:"keyword"`
	MemberID uint   `json:"memberId" form:"memberId"`
	GoodsID  uint   `json:"goodsId" form:"goodsId"`
	Status   string `json:"status" form:"status"`
}

type CreateExchangeOrderReq struct {
	MemberID uint   `json:"memberId" form:"memberId" binding:"required"`
	GoodsID  uint   `json:"goodsId" form:"goodsId" binding:"required"`
	Remark   string `json:"remark" form:"remark"`
}

type OperateExchangeOrderReq struct {
	ID     uint   `json:"id" form:"id" binding:"required"`
	Remark string `json:"remark" form:"remark"`
}
