package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type PointProductBaseInput struct {
	Name        string `json:"name" form:"name" binding:"required"`
	CoverURL    string `json:"coverUrl" form:"coverUrl"`
	Category    string `json:"category" form:"category"`
	PointsPrice int64  `json:"pointsPrice" form:"pointsPrice" binding:"required"`
	Stock       int64  `json:"stock" form:"stock" binding:"required"`
	Status      int    `json:"status" form:"status"`
	Sort        int    `json:"sort" form:"sort"`
	Description string `json:"description" form:"description"`
}

type CreatePointProductReq struct {
	PointProductBaseInput
}

type UpdatePointProductReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
	PointProductBaseInput
}

type PointProductSearch struct {
	commonReq.PageInfo
	Keyword  string `json:"keyword" form:"keyword"`
	Category string `json:"category" form:"category"`
	Status   int    `json:"status" form:"status"`
}

type UpdatePointProductStatusReq struct {
	ID     uint `json:"id" form:"id" binding:"required"`
	Status int  `json:"status" form:"status" binding:"required"`
}

type PointProductOptionsReq struct {
	Keyword string `json:"keyword" form:"keyword"`
}
