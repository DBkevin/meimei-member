package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PointProductRouter struct{}

func (r *PointProductRouter) InitPointProductRouter(Router *gin.RouterGroup) {
	pointProductRouter := Router.Group("pointProduct").Use(middleware.OperationRecord())
	pointProductRouterWithoutRecord := Router.Group("pointProduct")
	{
		pointProductRouter.POST("createPointProduct", pointProductApi.CreatePointProduct)
		pointProductRouter.DELETE("deletePointProduct", pointProductApi.DeletePointProduct)
		pointProductRouter.PUT("updatePointProduct", pointProductApi.UpdatePointProduct)
		pointProductRouter.PUT("updatePointProductStatus", pointProductApi.UpdatePointProductStatus)
	}
	{
		pointProductRouterWithoutRecord.GET("findPointProduct", pointProductApi.FindPointProduct)
		pointProductRouterWithoutRecord.GET("getPointProductList", pointProductApi.GetPointProductList)
		pointProductRouterWithoutRecord.GET("getPointProductOptions", pointProductApi.GetPointProductOptions)
	}
}
