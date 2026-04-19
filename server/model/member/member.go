package member

import "time"

type Member struct {
	BaseModel
	Name     string     `json:"name" form:"name" gorm:"column:name;size:64;comment:会员姓名"`
	Phone    string     `json:"phone" form:"phone" gorm:"column:phone;size:32;uniqueIndex:idx_mm_member_phone;comment:手机号"`
	Gender   string     `json:"gender" form:"gender" gorm:"column:gender;size:16;comment:性别"`
	Birthday *time.Time `json:"birthday,omitempty" form:"birthday" gorm:"column:birthday;type:date;comment:生日"`
	Source   string     `json:"source" form:"source" gorm:"column:source;size:64;comment:来源渠道"`
	Level    string     `json:"level" form:"level" gorm:"column:level;size:32;default:standard;comment:会员等级"`
	Status   int        `json:"status" form:"status" gorm:"column:status;default:1;index;comment:状态 1启用 2禁用"`
	Remark   string     `json:"remark" form:"remark" gorm:"column:remark;size:500;comment:备注"`
}

func (Member) TableName() string {
	return "mm_members"
}
