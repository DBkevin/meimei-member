package member

type PointLog struct {
	BaseModel
	MemberID     uint   `json:"memberId" form:"memberId" gorm:"column:member_id;index;comment:会员ID"`
	ChangeType   string `json:"changeType" form:"changeType" gorm:"column:change_type;size:32;index;comment:变动类型"`
	ChangePoints int64  `json:"changePoints" form:"changePoints" gorm:"column:change_points;comment:变动积分"`
	BeforePoints int64  `json:"beforePoints" form:"beforePoints" gorm:"column:before_points;comment:变动前积分"`
	AfterPoints  int64  `json:"afterPoints" form:"afterPoints" gorm:"column:after_points;comment:变动后积分"`
	SourceType   string `json:"sourceType" form:"sourceType" gorm:"column:source_type;size:32;index;comment:来源类型"`
	SourceID     uint   `json:"sourceId" form:"sourceId" gorm:"column:source_id;index;comment:来源ID"`
	Remark       string `json:"remark" form:"remark" gorm:"column:remark;size:255;comment:备注"`
	OperatorID   uint   `json:"operatorId" form:"operatorId" gorm:"column:operator_id;index;comment:操作人ID"`
	Member       Member `json:"member,omitempty" gorm:"foreignKey:MemberID"`
}

func (PointLog) TableName() string {
	return "member_point_logs"
}
