package member

const (
	MemberStatusEnabled  = "enabled"
	MemberStatusDisabled = "disabled"

	PointChangeTypeEarn      = "earn"
	PointChangeTypeUse       = "use"
	PointChangeTypeAdjustAdd = "adjust_add"
	PointChangeTypeAdjustSub = "adjust_sub"
	PointChangeTypeRefund    = "refund"

	PointSourceTypeManual            = "manual"
	PointSourceTypeExchangeOrder     = "exchange_order"
	PointSourceTypeExchangeOrderVoid = "exchange_order_void"

	GoodsStatusOnSale  = "on_sale"
	GoodsStatusOffSale = "off_sale"

	OrderStatusPending   = "pending"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
	OrderStatusRefunded  = "refunded"
)
