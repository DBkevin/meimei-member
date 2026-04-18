package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PointGoodsRouter struct{}

func (r *PointGoodsRouter) InitPointGoodsRouter(Router *gin.RouterGroup) {
	pointGoodsRouter := Router.Group("pointGoods").Use(middleware.OperationRecord())
	pointGoodsRouterWithoutRecord := Router.Group("pointGoods")
	{
		pointGoodsRouter.POST("createPointGoods", pointGoodsApi.CreatePointGoods)
		pointGoodsRouter.DELETE("deletePointGoods", pointGoodsApi.DeletePointGoods)
		pointGoodsRouter.PUT("updatePointGoods", pointGoodsApi.UpdatePointGoods)
		pointGoodsRouter.PUT("updatePointGoodsStatus", pointGoodsApi.UpdatePointGoodsStatus)
		pointGoodsRouter.PUT("updatePointGoodsStock", pointGoodsApi.UpdatePointGoodsStock)
	}
	{
		pointGoodsRouterWithoutRecord.GET("findPointGoods", pointGoodsApi.FindPointGoods)
		pointGoodsRouterWithoutRecord.GET("getPointGoodsList", pointGoodsApi.GetPointGoodsList)
		pointGoodsRouterWithoutRecord.GET("getPointGoodsOptions", pointGoodsApi.GetPointGoodsOptions)
	}
}
