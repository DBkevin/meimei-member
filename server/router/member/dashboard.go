package member

import (
	"github.com/gin-gonic/gin"
)

type DashboardRouter struct{}

func (r *DashboardRouter) InitDashboardRouter(Router *gin.RouterGroup) {
	dashboardRouter := Router.Group("dashboard").Use()
	{
		dashboardRouter.GET("summary", dashboardApi.GetDashboardSummary)
		dashboardRouter.GET("recentTransactions", dashboardApi.GetRecentTransactions)
		dashboardRouter.GET("recentOrders", dashboardApi.GetRecentOrders)
		dashboardRouter.GET("lowStockProducts", dashboardApi.GetLowStockProducts)
	}
}
