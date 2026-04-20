package member

import (
	"strconv"

	commonRes "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

// DashboardApi 会员积分数据概览API
type DashboardApi struct{}

// GetDashboardSummary
// @Tags      Dashboard
// @Summary   获取数据概览
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Success   200  {object}  commonRes.Response{data=DashboardSummary}  "获取成功"
// @Router    /member/dashboard/summary [get]
func (a *DashboardApi) GetDashboardSummary(c *gin.Context) {
	summary, err := dashboardService.GetDashboardSummary()
	if err != nil {
		commonRes.FailWithMessage("获取数据概览失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(summary, "获取成功", c)
}

// GetRecentTransactions
// @Tags      Dashboard
// @Summary   获取最近积分流水
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     limit  query    int  false  "数量限制，默认10"
// @Success   200  {object}  commonRes.Response{data=[]PointTransaction}  "获取成功"
// @Router    /member/dashboard/recentTransactions [get]
func (a *DashboardApi) GetRecentTransactions(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	transactions, err := dashboardService.GetRecentTransactions(limit)
	if err != nil {
		commonRes.FailWithMessage("获取最近积分流水失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(transactions, "获取成功", c)
}

// GetRecentOrders
// @Tags      Dashboard
// @Summary   获取最近兑换订单
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     limit  query    int  false  "数量限制，默认10"
// @Success   200  {object}  commonRes.Response{data=[]RedemptionOrder}  "获取成功"
// @Router    /member/dashboard/recentOrders [get]
func (a *DashboardApi) GetRecentOrders(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	orders, err := dashboardService.GetRecentOrders(limit)
	if err != nil {
		commonRes.FailWithMessage("获取最近兑换订单失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(orders, "获取成功", c)
}

// GetLowStockProducts
// @Tags      Dashboard
// @Summary   获取库存不足商品
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     limit  query    int  false  "数量限制，默认10"
// @Success   200  {object}  commonRes.Response{data=[]PointProduct}  "获取成功"
// @Router    /member/dashboard/lowStockProducts [get]
func (a *DashboardApi) GetLowStockProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	products, err := dashboardService.GetLowStockProducts(limit)
	if err != nil {
		commonRes.FailWithMessage("获取库存不足商品失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(products, "获取成功", c)
}
