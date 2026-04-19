package response

import memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"

type PointProductDetail struct {
	Product memberModel.PointProduct `json:"product"`
}
