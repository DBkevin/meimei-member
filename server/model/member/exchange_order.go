package member

import "time"

type ExchangeOrder struct {
	BaseModel
	OrderNo    string     `json:"orderNo" form:"orderNo" gorm:"column:order_no;size:64;uniqueIndex;comment:订单号"`
	MemberID   uint       `json:"memberId" form:"memberId" gorm:"column:member_id;index;comment:会员ID"`
	GoodsID    uint       `json:"goodsId" form:"goodsId" gorm:"column:goods_id;index;comment:商品ID"`
	PointsCost int64      `json:"pointsCost" form:"pointsCost" gorm:"column:points_cost;comment:积分消耗"`
	Status     string     `json:"status" form:"status" gorm:"column:status;size:32;default:pending;index;comment:订单状态"`
	VerifyCode string     `json:"verifyCode" form:"verifyCode" gorm:"column:verify_code;size:32;index;comment:核销码"`
	VerifiedAt *time.Time `json:"verifiedAt" form:"verifiedAt" gorm:"column:verified_at;comment:核销时间"`
	OperatorID uint       `json:"operatorId" form:"operatorId" gorm:"column:operator_id;index;comment:操作人ID"`
	Remark     string     `json:"remark" form:"remark" gorm:"column:remark;size:255;comment:备注"`
	Member     Member     `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	Goods      PointGoods `json:"goods,omitempty" gorm:"foreignKey:GoodsID"`
}

func (ExchangeOrder) TableName() string {
	return "member_exchange_orders"
}
