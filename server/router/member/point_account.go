package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PointAccountRouter struct{}

func (r *PointAccountRouter) InitPointAccountRouter(Router *gin.RouterGroup) {
	pointAccountRouter := Router.Group("pointAccount").Use(middleware.OperationRecord())
	pointAccountRouterWithoutRecord := Router.Group("pointAccount")
	{
		pointAccountRouter.POST("manualAddPoints", pointAccountApi.ManualAddPoints)
		pointAccountRouter.POST("manualSubPoints", pointAccountApi.ManualSubPoints)
	}
	{
		pointAccountRouterWithoutRecord.GET("findPointAccount", pointAccountApi.FindPointAccount)
		pointAccountRouterWithoutRecord.GET("getPointAccountList", pointAccountApi.GetPointAccountList)
	}
}
