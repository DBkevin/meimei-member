package response

import memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"

type MemberDetail struct {
	Member  memberModel.Member       `json:"member"`
	Account memberModel.PointAccount `json:"account"`
}

type MemberOption struct {
	ID          uint   `json:"id"`
	Label       string `json:"label"`
	Mobile      string `json:"mobile"`
	Nickname    string `json:"nickname"`
	RealName    string `json:"realName"`
	MemberLevel string `json:"memberLevel"`
}
