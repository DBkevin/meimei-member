package member

type PointTransaction struct {
	BaseModel
	MemberID      uint         `json:"memberId" form:"memberId" gorm:"column:member_id;index;comment:会员ID"`
	AccountID     uint         `json:"accountId" form:"accountId" gorm:"column:account_id;index;comment:积分账户ID"`
	Type          string       `json:"type" form:"type" gorm:"column:type;size:32;index;comment:流水类型"`
	Points        int64        `json:"points" form:"points" gorm:"column:points;comment:积分数量"`
	BeforeBalance int64        `json:"beforeBalance" form:"beforeBalance" gorm:"column:before_balance;comment:变动前余额"`
	AfterBalance  int64        `json:"afterBalance" form:"afterBalance" gorm:"column:after_balance;comment:变动后余额"`
	RefType       string       `json:"refType" form:"refType" gorm:"column:ref_type;size:64;index;comment:来源类型"`
	RefID         uint         `json:"refId" form:"refId" gorm:"column:ref_id;index;comment:来源ID"`
	Operator      string       `json:"operator" form:"operator" gorm:"column:operator;size:64;comment:操作人"`
	Remark        string       `json:"remark" form:"remark" gorm:"column:remark;size:255;comment:备注"`
	Member        Member       `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	Account       PointAccount `json:"account,omitempty" gorm:"foreignKey:AccountID"`
}

func (PointTransaction) TableName() string {
	return "mm_point_transactions"
}
