package response

import memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"

type RedemptionOrderDetail struct {
	Order memberModel.RedemptionOrder `json:"order"`
}
