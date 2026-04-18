package response

import memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"

type PointGoodsDetail struct {
	Goods memberModel.PointGoods `json:"goods"`
}
