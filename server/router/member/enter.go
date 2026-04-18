package member

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	MemberRouter
	PointAccountRouter
	PointLogRouter
	PointGoodsRouter
	ExchangeOrderRouter
}

var (
	memberApi        = api.ApiGroupApp.MemberApiGroup.MemberApi
	pointAccountApi  = api.ApiGroupApp.MemberApiGroup.PointAccountApi
	pointLogApi      = api.ApiGroupApp.MemberApiGroup.PointLogApi
	pointGoodsApi    = api.ApiGroupApp.MemberApiGroup.PointGoodsApi
	exchangeOrderApi = api.ApiGroupApp.MemberApiGroup.ExchangeOrderApi
)
