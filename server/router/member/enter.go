package member

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	MemberRouter
	PointAccountRouter
	PointTransactionRouter
	PointProductRouter
	RedemptionOrderRouter
	DashboardRouter
}

var (
	memberApi           = api.ApiGroupApp.MemberApiGroup.MemberApi
	pointAccountApi     = api.ApiGroupApp.MemberApiGroup.PointAccountApi
	pointTransactionApi = api.ApiGroupApp.MemberApiGroup.PointTransactionApi
	pointProductApi     = api.ApiGroupApp.MemberApiGroup.PointProductApi
	redemptionOrderApi  = api.ApiGroupApp.MemberApiGroup.RedemptionOrderApi
	dashboardApi        = api.ApiGroupApp.MemberApiGroup.DashboardApi
)
