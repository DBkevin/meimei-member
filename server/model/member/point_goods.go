package member

type PointProduct struct {
	BaseModel
	Name        string `json:"name" form:"name" gorm:"column:name;size:128;index;comment:商品名称"`
	CoverURL    string `json:"coverUrl" form:"coverUrl" gorm:"column:cover_url;size:255;comment:封面图"`
	Category    string `json:"category" form:"category" gorm:"column:category;size:64;index;comment:分类"`
	PointsPrice int64  `json:"pointsPrice" form:"pointsPrice" gorm:"column:points_price;comment:兑换所需积分"`
	Stock       int64  `json:"stock" form:"stock" gorm:"column:stock;default:0;comment:库存"`
	Status      int    `json:"status" form:"status" gorm:"column:status;default:2;index;comment:状态 1上架 2下架"`
	Sort        int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序"`
	Description string `json:"description" form:"description" gorm:"column:description;type:text;comment:说明"`
}

func (PointProduct) TableName() string {
	return "mm_point_products"
}
