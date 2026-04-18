package response

import memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"

type PointAccountDetail struct {
	Account memberModel.PointAccount `json:"account"`
}
