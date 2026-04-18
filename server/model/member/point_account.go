package member

type PointAccount struct {
	BaseModel
	MemberID          uint   `json:"memberId" form:"memberId" gorm:"column:member_id;uniqueIndex;comment:会员ID"`
	AvailablePoints   int64  `json:"availablePoints" form:"availablePoints" gorm:"column:available_points;default:0;comment:可用积分"`
	FrozenPoints      int64  `json:"frozenPoints" form:"frozenPoints" gorm:"column:frozen_points;default:0;comment:冻结积分"`
	TotalEarnedPoints int64  `json:"totalEarnedPoints" form:"totalEarnedPoints" gorm:"column:total_earned_points;default:0;comment:累计获得积分"`
	TotalUsedPoints   int64  `json:"totalUsedPoints" form:"totalUsedPoints" gorm:"column:total_used_points;default:0;comment:累计使用积分"`
	Member            Member `json:"member,omitempty" gorm:"foreignKey:MemberID"`
}

func (PointAccount) TableName() string {
	return "member_point_accounts"
}
