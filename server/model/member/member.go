package member

type Member struct {
	BaseModel
	OpenID      string `json:"openid" form:"openid" gorm:"column:openid;size:128;index:idx_member_openid,unique;comment:微信openid"`
	UnionID     string `json:"unionid" form:"unionid" gorm:"column:unionid;size:128;index;comment:微信unionid"`
	Mobile      string `json:"mobile" form:"mobile" gorm:"column:mobile;size:32;uniqueIndex:idx_member_mobile;comment:手机号"`
	Nickname    string `json:"nickname" form:"nickname" gorm:"column:nickname;size:128;comment:昵称"`
	AvatarURL   string `json:"avatarUrl" form:"avatarUrl" gorm:"column:avatar_url;size:255;comment:头像地址"`
	RealName    string `json:"realName" form:"realName" gorm:"column:real_name;size:64;comment:真实姓名"`
	MemberLevel string `json:"memberLevel" form:"memberLevel" gorm:"column:member_level;size:32;default:standard;comment:会员等级"`
	Status      string `json:"status" form:"status" gorm:"column:status;size:32;default:enabled;index;comment:会员状态"`
}

func (Member) TableName() string {
	return "member_members"
}
