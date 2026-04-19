package member

type RedemptionOrder struct {
	BaseModel
	OrderNo       string       `json:"orderNo" form:"orderNo" gorm:"column:order_no;size:64;uniqueIndex:idx_mm_redemption_order_no;comment:订单号"`
	MemberID      uint         `json:"memberId" form:"memberId" gorm:"column:member_id;index;comment:会员ID"`
	ProductID     uint         `json:"productId" form:"productId" gorm:"column:product_id;index;comment:商品ID"`
	ProductName   string       `json:"productName" form:"productName" gorm:"column:product_name;size:128;comment:商品名称"`
	Quantity      int64        `json:"quantity" form:"quantity" gorm:"column:quantity;default:1;comment:数量"`
	UnitPoints    int64        `json:"unitPoints" form:"unitPoints" gorm:"column:unit_points;comment:单件积分"`
	TotalPoints   int64        `json:"totalPoints" form:"totalPoints" gorm:"column:total_points;comment:总积分"`
	Status        int          `json:"status" form:"status" gorm:"column:status;default:1;index;comment:状态 1待处理 2已完成 3已取消"`
	ReceiverName  string       `json:"receiverName" form:"receiverName" gorm:"column:receiver_name;size:64;comment:收货人"`
	ReceiverPhone string       `json:"receiverPhone" form:"receiverPhone" gorm:"column:receiver_phone;size:32;comment:收货电话"`
	Remark        string       `json:"remark" form:"remark" gorm:"column:remark;size:255;comment:备注"`
	Member        Member       `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	Product       PointProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (RedemptionOrder) TableName() string {
	return "mm_redemption_orders"
}
