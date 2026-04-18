package member

type PointGoods struct {
	BaseModel
	Name           string `json:"name" form:"name" gorm:"column:name;size:128;index;comment:商品名称"`
	CoverImage     string `json:"coverImage" form:"coverImage" gorm:"column:cover_image;size:255;comment:封面图"`
	Description    string `json:"description" form:"description" gorm:"column:description;type:text;comment:商品描述"`
	PointsPrice    int64  `json:"pointsPrice" form:"pointsPrice" gorm:"column:points_price;comment:积分价格"`
	Stock          int64  `json:"stock" form:"stock" gorm:"column:stock;default:0;comment:库存"`
	LimitPerMember int64  `json:"limitPerMember" form:"limitPerMember" gorm:"column:limit_per_member;default:0;comment:每人限兑数量"`
	Status         string `json:"status" form:"status" gorm:"column:status;size:32;default:on_sale;index;comment:商品状态"`
	Sort           int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序"`
}

func (PointGoods) TableName() string {
	return "member_point_goods"
}
