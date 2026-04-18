package member

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	MemberApi
	PointAccountApi
	PointLogApi
	PointGoodsApi
	ExchangeOrderApi
}

var (
	memberService        = service.ServiceGroupApp.MemberServiceGroup.MemberService
	pointAccountService  = service.ServiceGroupApp.MemberServiceGroup.PointAccountService
	pointLogService      = service.ServiceGroupApp.MemberServiceGroup.PointLogService
	pointGoodsService    = service.ServiceGroupApp.MemberServiceGroup.PointGoodsService
	exchangeOrderService = service.ServiceGroupApp.MemberServiceGroup.ExchangeOrderService
)
