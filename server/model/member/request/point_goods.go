package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type PointGoodsSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Status  string `json:"status" form:"status"`
}

type UpdateGoodsStatusReq struct {
	ID     uint   `json:"id" form:"id" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
}

type UpdateGoodsStockReq struct {
	ID    uint  `json:"id" form:"id" binding:"required"`
	Stock int64 `json:"stock" form:"stock" binding:"required"`
}
