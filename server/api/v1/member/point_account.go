package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PointAccountApi struct{}

func (a *PointAccountApi) GetPointAccountList(c *gin.Context) {
	var pageInfo memberReq.PointAccountSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pointAccountService.GetPointAccountList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取积分账户列表失败", zap.Error(err))
		response.FailWithMessage("获取积分账户列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取积分账户列表成功", c)
}

func (a *PointAccountApi) ManualAddPoints(c *gin.Context) {
	var req memberReq.AdjustPointsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointAccountService.ManualAddPoints(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("手工增加积分失败", zap.Error(err))
		response.FailWithMessage("手工增加积分失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("手工增加积分成功", c)
}

func (a *PointAccountApi) ManualSubPoints(c *gin.Context) {
	var req memberReq.AdjustPointsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointAccountService.ManualSubPoints(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("手工扣减积分失败", zap.Error(err))
		response.FailWithMessage("手工扣减积分失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("手工扣减积分成功", c)
}
