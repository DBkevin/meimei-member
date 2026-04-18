package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExchangeOrderRouter struct{}

func (r *ExchangeOrderRouter) InitExchangeOrderRouter(Router *gin.RouterGroup) {
	exchangeOrderRouter := Router.Group("exchangeOrder").Use(middleware.OperationRecord())
	exchangeOrderRouterWithoutRecord := Router.Group("exchangeOrder")
	{
		exchangeOrderRouter.POST("createExchangeOrder", exchangeOrderApi.CreateExchangeOrder)
		exchangeOrderRouter.POST("verifyExchangeOrder", exchangeOrderApi.VerifyExchangeOrder)
		exchangeOrderRouter.POST("cancelExchangeOrder", exchangeOrderApi.CancelExchangeOrder)
		exchangeOrderRouter.POST("refundExchangeOrder", exchangeOrderApi.RefundExchangeOrder)
	}
	{
		exchangeOrderRouterWithoutRecord.GET("findExchangeOrder", exchangeOrderApi.FindExchangeOrder)
		exchangeOrderRouterWithoutRecord.GET("getExchangeOrderList", exchangeOrderApi.GetExchangeOrderList)
	}
}
