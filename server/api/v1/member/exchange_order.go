package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	commonRes "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RedemptionOrderApi struct{}

// CreateRedemptionOrder
// @Tags      RedemptionOrder
// @Summary   创建兑换订单
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberReq.CreateRedemptionOrderReq true  "会员ID、商品ID和兑换信息"
// @Success   200   {object}  commonRes.Response{msg=string}     "创建兑换订单成功"
// @Router    /redemptionOrder/createRedemptionOrder [post]
func (a *RedemptionOrderApi) CreateRedemptionOrder(c *gin.Context) {
	var req memberReq.CreateRedemptionOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := ValidateRedemptionOrderInput(req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := redemptionOrderService.CreateRedemptionOrder(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("创建兑换订单失败", zap.Error(err))
		commonRes.FailWithMessage("创建兑换订单失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithMessage("创建兑换订单成功", c)
}

// FindRedemptionOrder
// @Tags      RedemptionOrder
// @Summary   查询兑换订单详情
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     id    query     int                true  "兑换订单ID"
// @Success   200   {object}  commonRes.Response "查询兑换订单成功"
// @Router    /redemptionOrder/findRedemptionOrder [get]
func (a *RedemptionOrderApi) FindRedemptionOrder(c *gin.Context) {
	var info commonReq.GetById
	if err := c.ShouldBindQuery(&info); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	order, err := redemptionOrderService.GetRedemptionOrder(info.Uint())
	if err != nil {
		global.GVA_LOG.Error("查询兑换订单失败", zap.Error(err))
		commonRes.FailWithMessage("查询兑换订单失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(order, "查询兑换订单成功", c)
}

// GetRedemptionOrderList
// @Tags      RedemptionOrder
// @Summary   分页获取兑换订单列表
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  query     memberReq.RedemptionOrderSearch true  "分页与筛选条件"
// @Success   200   {object}  commonRes.Response{data=commonRes.PageResult,msg=string} "获取兑换订单列表成功"
// @Router    /redemptionOrder/getRedemptionOrderList [get]
func (a *RedemptionOrderApi) GetRedemptionOrderList(c *gin.Context) {
	var pageInfo memberReq.RedemptionOrderSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := redemptionOrderService.GetRedemptionOrderList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取兑换订单列表失败", zap.Error(err))
		commonRes.FailWithMessage("获取兑换订单列表失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(commonRes.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取兑换订单列表成功", c)
}

// CompleteRedemptionOrder
// @Tags      RedemptionOrder
// @Summary   完成兑换订单
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberReq.OperateRedemptionOrderReq true  "订单ID与备注"
// @Success   200   {object}  commonRes.Response{msg=string}      "完成兑换订单成功"
// @Router    /redemptionOrder/completeRedemptionOrder [post]
func (a *RedemptionOrderApi) CompleteRedemptionOrder(c *gin.Context) {
	var req memberReq.OperateRedemptionOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := ValidateIDInput(req.ID); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := redemptionOrderService.CompleteRedemptionOrder(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("完成兑换订单失败", zap.Error(err))
		commonRes.FailWithMessage("完成兑换订单失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithMessage("完成兑换订单成功", c)
}

// CancelRedemptionOrder
// @Tags      RedemptionOrder
// @Summary   取消兑换订单
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberReq.OperateRedemptionOrderReq true  "订单ID与备注"
// @Success   200   {object}  commonRes.Response{msg=string}      "取消兑换订单成功"
// @Router    /redemptionOrder/cancelRedemptionOrder [post]
func (a *RedemptionOrderApi) CancelRedemptionOrder(c *gin.Context) {
	var req memberReq.OperateRedemptionOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := ValidateIDInput(req.ID); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := redemptionOrderService.CancelRedemptionOrder(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("取消兑换订单失败", zap.Error(err))
		commonRes.FailWithMessage("取消兑换订单失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithMessage("取消兑换订单成功", c)
}
