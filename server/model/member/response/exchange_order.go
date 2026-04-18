package response

import memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"

type ExchangeOrderDetail struct {
	Order memberModel.ExchangeOrder `json:"order"`
}
