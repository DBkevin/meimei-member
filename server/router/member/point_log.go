package member

import "github.com/gin-gonic/gin"

type PointTransactionRouter struct{}

func (r *PointTransactionRouter) InitPointTransactionRouter(Router *gin.RouterGroup) {
	pointTransactionRouterWithoutRecord := Router.Group("pointTransaction")
	{
		pointTransactionRouterWithoutRecord.GET("getPointTransactionList", pointTransactionApi.GetPointTransactionList)
	}
}
