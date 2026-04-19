package response

import memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"

type PointTransactionDetail struct {
	Transaction memberModel.PointTransaction `json:"transaction"`
}
