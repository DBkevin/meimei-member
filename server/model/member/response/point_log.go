package response

import memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"

type PointLogDetail struct {
	Log memberModel.PointLog `json:"log"`
}
