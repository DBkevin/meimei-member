package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonRes "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PointAccountApi struct{}

// FindPointAccount
// @Tags      PointAccount
// @Summary   查询会员积分账户
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     memberId  query     int                true  "会员ID"
// @Success   200       {object}  commonRes.Response "查询积分账户成功"
// @Router    /pointAccount/findPointAccount [get]
func (a *PointAccountApi) FindPointAccount(c *gin.Context) {
	var req memberReq.GetPointAccountReq
	if err := c.ShouldBindQuery(&req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	account, err := pointAccountService.GetPointAccountByMemberID(req.MemberID)
	if err != nil {
		global.GVA_LOG.Error("查询积分账户失败", zap.Error(err))
		commonRes.FailWithMessage("查询积分账户失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(account, "查询积分账户成功", c)
}

// GetPointAccountList
// @Tags      PointAccount
// @Summary   分页获取积分账户列表
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  query     memberReq.PointAccountSearch true  "分页与筛选条件"
// @Success   200   {object}  commonRes.Response{data=commonRes.PageResult,msg=string} "获取积分账户列表成功"
// @Router    /pointAccount/getPointAccountList [get]
func (a *PointAccountApi) GetPointAccountList(c *gin.Context) {
	var pageInfo memberReq.PointAccountSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pointAccountService.GetPointAccountList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取积分账户列表失败", zap.Error(err))
		commonRes.FailWithMessage("获取积分账户列表失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(commonRes.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取积分账户列表成功", c)
}

// ManualAddPoints
// @Tags      PointAccount
// @Summary   手工增加积分
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberReq.AdjustPointsReq    true  "会员ID与积分数量"
// @Success   200   {object}  commonRes.Response{msg=string} "手工增加积分成功"
// @Router    /pointAccount/manualAddPoints [post]
func (a *PointAccountApi) ManualAddPoints(c *gin.Context) {
	var req memberReq.AdjustPointsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := ValidateAdjustPointsInput(req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointAccountService.ManualAddPoints(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("手工增加积分失败", zap.Error(err))
		commonRes.FailWithMessage("手工增加积分失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithMessage("手工增加积分成功", c)
}

// ManualSubPoints
// @Tags      PointAccount
// @Summary   手工扣减积分
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberReq.AdjustPointsReq    true  "会员ID与积分数量"
// @Success   200   {object}  commonRes.Response{msg=string} "手工扣减积分成功"
// @Router    /pointAccount/manualSubPoints [post]
func (a *PointAccountApi) ManualSubPoints(c *gin.Context) {
	var req memberReq.AdjustPointsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := ValidateAdjustPointsInput(req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointAccountService.ManualSubPoints(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("手工扣减积分失败", zap.Error(err))
		commonRes.FailWithMessage("手工扣减积分失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithMessage("手工扣减积分成功", c)
}
