package member

type PointAccount struct {
	BaseModel
	MemberID     uint   `json:"memberId" form:"memberId" gorm:"column:member_id;uniqueIndex:idx_mm_point_account_member;comment:会员ID"`
	Balance      int64  `json:"balance" form:"balance" gorm:"column:balance;default:0;comment:当前可用积分"`
	TotalEarned  int64  `json:"totalEarned" form:"totalEarned" gorm:"column:total_earned;default:0;comment:累计获得积分"`
	TotalSpent   int64  `json:"totalSpent" form:"totalSpent" gorm:"column:total_spent;default:0;comment:累计消耗积分"`
	FrozenPoints int64  `json:"frozenPoints" form:"frozenPoints" gorm:"column:frozen_points;default:0;comment:冻结积分"`
	Member       Member `json:"member,omitempty" gorm:"foreignKey:MemberID"`
}

func (PointAccount) TableName() string {
	return "mm_point_accounts"
}
