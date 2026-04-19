package member

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	MemberApi
	PointAccountApi
	PointTransactionApi
	PointProductApi
	RedemptionOrderApi
}

var (
	memberService           = service.ServiceGroupApp.MemberServiceGroup.MemberService
	pointAccountService     = service.ServiceGroupApp.MemberServiceGroup.PointAccountService
	pointTransactionService = service.ServiceGroupApp.MemberServiceGroup.PointTransactionService
	pointProductService     = service.ServiceGroupApp.MemberServiceGroup.PointProductService
	redemptionOrderService  = service.ServiceGroupApp.MemberServiceGroup.RedemptionOrderService
)
