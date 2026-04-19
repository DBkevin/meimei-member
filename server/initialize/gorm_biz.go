package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberService "github.com/flipped-aurora/gin-vue-admin/server/service/member"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		&memberModel.Member{},
		&memberModel.PointAccount{},
		&memberModel.PointTransaction{},
		&memberModel.PointProduct{},
		&memberModel.RedemptionOrder{},
	)
	if err != nil {
		return err
	}
	return memberService.Bootstrap()
}
