package member

import "github.com/gin-gonic/gin"

type PointLogRouter struct{}

func (r *PointLogRouter) InitPointLogRouter(Router *gin.RouterGroup) {
	pointLogRouterWithoutRecord := Router.Group("pointLog")
	{
		pointLogRouterWithoutRecord.GET("getPointLogList", pointLogApi.GetPointLogList)
	}
}
