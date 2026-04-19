package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type RedemptionOrderRouter struct{}

func (r *RedemptionOrderRouter) InitRedemptionOrderRouter(Router *gin.RouterGroup) {
	redemptionOrderRouter := Router.Group("redemptionOrder").Use(middleware.OperationRecord())
	redemptionOrderRouterWithoutRecord := Router.Group("redemptionOrder")
	{
		redemptionOrderRouter.POST("createRedemptionOrder", redemptionOrderApi.CreateRedemptionOrder)
		redemptionOrderRouter.POST("completeRedemptionOrder", redemptionOrderApi.CompleteRedemptionOrder)
		redemptionOrderRouter.POST("cancelRedemptionOrder", redemptionOrderApi.CancelRedemptionOrder)
	}
	{
		redemptionOrderRouterWithoutRecord.GET("findRedemptionOrder", redemptionOrderApi.FindRedemptionOrder)
		redemptionOrderRouterWithoutRecord.GET("getRedemptionOrderList", redemptionOrderApi.GetRedemptionOrderList)
	}
}
