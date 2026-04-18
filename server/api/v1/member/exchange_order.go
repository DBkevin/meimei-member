package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExchangeOrderApi struct{}

func (a *ExchangeOrderApi) CreateExchangeOrder(c *gin.Context) {
	var req memberReq.CreateExchangeOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := exchangeOrderService.CreateExchangeOrder(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("创建兑换订单失败", zap.Error(err))
		response.FailWithMessage("创建兑换订单失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建兑换订单成功", c)
}

func (a *ExchangeOrderApi) FindExchangeOrder(c *gin.Context) {
	var info request.GetById
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order, err := exchangeOrderService.GetExchangeOrder(info.Uint())
	if err != nil {
		global.GVA_LOG.Error("查询兑换订单失败", zap.Error(err))
		response.FailWithMessage("查询兑换订单失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(order, "查询兑换订单成功", c)
}

func (a *ExchangeOrderApi) GetExchangeOrderList(c *gin.Context) {
	var pageInfo memberReq.ExchangeOrderSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := exchangeOrderService.GetExchangeOrderList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取兑换订单列表失败", zap.Error(err))
		response.FailWithMessage("获取兑换订单列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取兑换订单列表成功", c)
}

func (a *ExchangeOrderApi) VerifyExchangeOrder(c *gin.Context) {
	var req memberReq.OperateExchangeOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := exchangeOrderService.VerifyExchangeOrder(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("核销兑换订单失败", zap.Error(err))
		response.FailWithMessage("核销兑换订单失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("核销兑换订单成功", c)
}

func (a *ExchangeOrderApi) CancelExchangeOrder(c *gin.Context) {
	var req memberReq.OperateExchangeOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := exchangeOrderService.CancelExchangeOrder(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("取消兑换订单失败", zap.Error(err))
		response.FailWithMessage("取消兑换订单失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("取消兑换订单成功", c)
}

func (a *ExchangeOrderApi) RefundExchangeOrder(c *gin.Context) {
	var req memberReq.OperateExchangeOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := exchangeOrderService.RefundExchangeOrder(req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("退款兑换订单失败", zap.Error(err))
		response.FailWithMessage("退款兑换订单失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("退款兑换订单成功", c)
}
