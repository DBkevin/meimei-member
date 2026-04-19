package member

const (
	MemberLevelStandard = "standard"

	MemberStatusEnabled  = 1
	MemberStatusDisabled = 2

	PointTransactionTypeEarn   = "earn"
	PointTransactionTypeSpend  = "spend"
	PointTransactionTypeAdjust = "adjust"
	PointTransactionTypeRefund = "refund"

	PointRefTypeManualAdjustAdd     = "manual_adjust_add"
	PointRefTypeManualAdjustSub     = "manual_adjust_sub"
	PointRefTypeRedemptionOrder     = "redemption_order"
	PointRefTypeRedemptionOrderVoid = "redemption_order_cancel"

	PointProductStatusOnSale  = 1
	PointProductStatusOffSale = 2

	RedemptionOrderStatusPending   = 1
	RedemptionOrderStatusCompleted = 2
	RedemptionOrderStatusCancelled = 3
)
