package member

import (
	"time"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
)

type DashboardService struct{}

// DashboardSummary 数据概览统计
type DashboardSummary struct {
	// 会员统计
	TotalMembers    int64 `json:"totalMembers"`
	EnabledMembers  int64 `json:"enabledMembers"`
	DisabledMembers int64 `json:"disabledMembers"`
	TodayNewMembers int64 `json:"todayNewMembers"`

	// 积分统计
	TotalPointsBalance  int64 `json:"totalPointsBalance"`
	TotalPointsIssued   int64 `json:"totalPointsIssued"`
	TotalPointsConsumed int64 `json:"totalPointsConsumed"`
	TodayPointsIssued   int64 `json:"todayPointsIssued"`
	TodayPointsConsumed int64 `json:"todayPointsConsumed"`

	// 商品统计
	TotalProducts      int64 `json:"totalProducts"`
	OnSaleProducts     int64 `json:"onSaleProducts"`
	OffSaleProducts    int64 `json:"offSaleProducts"`
	OutOfStockProducts int64 `json:"outOfStockProducts"`

	// 订单统计
	TotalOrders     int64 `json:"totalOrders"`
	PendingOrders   int64 `json:"pendingOrders"`
	CompletedOrders int64 `json:"completedOrders"`
	CancelledOrders int64 `json:"cancelledOrders"`
	TodayNewOrders  int64 `json:"todayNewOrders"`
}

// GetDashboardSummary 获取数据概览
func (s *DashboardService) GetDashboardSummary() (DashboardSummary, error) {
	var summary DashboardSummary
	db := bizDB()

	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.ParseInLocation("2006-01-02 00:00:00", today+" 00:00:00", time.Local)

	// 会员统计
	db.Model(&memberModel.Member{}).Count(&summary.TotalMembers)
	db.Model(&memberModel.Member{}).Where("status = ?", memberModel.MemberStatusEnabled).Count(&summary.EnabledMembers)
	db.Model(&memberModel.Member{}).Where("status = ?", memberModel.MemberStatusDisabled).Count(&summary.DisabledMembers)
	db.Model(&memberModel.Member{}).Where("created_at >= ?", todayStart).Count(&summary.TodayNewMembers)

	// 积分统计 - 从所有账户汇总
	var accounts []memberModel.PointAccount
	db.Find(&accounts)
	for _, acc := range accounts {
		summary.TotalPointsBalance += acc.Balance
		summary.TotalPointsIssued += acc.TotalEarned
		summary.TotalPointsConsumed += acc.TotalSpent
	}

	// 今日积分发放和消耗
	var todayTransactions []memberModel.PointTransaction
	db.Where("created_at >= ?", todayStart).Find(&todayTransactions)
	for _, tx := range todayTransactions {
		if tx.Type == memberModel.PointTransactionTypeEarn || tx.Type == memberModel.PointTransactionTypeAdjust {
			summary.TodayPointsIssued += tx.Points
		} else if tx.Type == memberModel.PointTransactionTypeSpend {
			summary.TodayPointsConsumed += tx.Points
		}
	}

	// 商品统计
	db.Model(&memberModel.PointProduct{}).Count(&summary.TotalProducts)
	db.Model(&memberModel.PointProduct{}).Where("status = ?", memberModel.PointProductStatusOnSale).Count(&summary.OnSaleProducts)
	db.Model(&memberModel.PointProduct{}).Where("status = ?", memberModel.PointProductStatusOffSale).Count(&summary.OffSaleProducts)
	db.Model(&memberModel.PointProduct{}).Where("stock <= ?", 0).Count(&summary.OutOfStockProducts)

	// 订单统计
	db.Model(&memberModel.RedemptionOrder{}).Count(&summary.TotalOrders)
	db.Model(&memberModel.RedemptionOrder{}).Where("status = ?", memberModel.RedemptionOrderStatusPending).Count(&summary.PendingOrders)
	db.Model(&memberModel.RedemptionOrder{}).Where("status = ?", memberModel.RedemptionOrderStatusCompleted).Count(&summary.CompletedOrders)
	db.Model(&memberModel.RedemptionOrder{}).Where("status = ?", memberModel.RedemptionOrderStatusCancelled).Count(&summary.CancelledOrders)
	db.Model(&memberModel.RedemptionOrder{}).Where("created_at >= ?", todayStart).Count(&summary.TodayNewOrders)

	return summary, nil
}

// GetRecentTransactions 获取最近积分流水
func (s *DashboardService) GetRecentTransactions(limit int) ([]memberModel.PointTransaction, error) {
	if limit <= 0 {
		limit = 10
	}
	var transactions []memberModel.PointTransaction
	err := bizDB().Order("created_at desc").Limit(limit).Find(&transactions).Error
	return transactions, err
}

// GetRecentOrders 获取最近兑换订单
func (s *DashboardService) GetRecentOrders(limit int) ([]memberModel.RedemptionOrder, error) {
	if limit <= 0 {
		limit = 10
	}
	var orders []memberModel.RedemptionOrder
	err := bizDB().Order("created_at desc").Limit(limit).Find(&orders).Error
	return orders, err
}

// GetLowStockProducts 获取库存不足商品
func (s *DashboardService) GetLowStockProducts(limit int) ([]memberModel.PointProduct, error) {
	if limit <= 0 {
		limit = 10
	}
	var products []memberModel.PointProduct
	err := bizDB().Where("stock <= ?", 0).Order("updated_at desc").Limit(limit).Find(&products).Error
	return products, err
}
