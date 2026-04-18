package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PointGoodsApi struct{}

func (a *PointGoodsApi) CreatePointGoods(c *gin.Context) {
	var goods memberModel.PointGoods
	if err := c.ShouldBindJSON(&goods); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointGoodsService.CreatePointGoods(&goods); err != nil {
		global.GVA_LOG.Error("创建积分商品失败", zap.Error(err))
		response.FailWithMessage("创建积分商品失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建积分商品成功", c)
}

func (a *PointGoodsApi) DeletePointGoods(c *gin.Context) {
	var info request.GetById
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointGoodsService.DeletePointGoods(info.Uint()); err != nil {
		global.GVA_LOG.Error("删除积分商品失败", zap.Error(err))
		response.FailWithMessage("删除积分商品失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除积分商品成功", c)
}

func (a *PointGoodsApi) UpdatePointGoods(c *gin.Context) {
	var goods memberModel.PointGoods
	if err := c.ShouldBindJSON(&goods); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if goods.ID == 0 {
		response.FailWithMessage("商品ID不能为空", c)
		return
	}
	if err := pointGoodsService.UpdatePointGoods(&goods); err != nil {
		global.GVA_LOG.Error("更新积分商品失败", zap.Error(err))
		response.FailWithMessage("更新积分商品失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新积分商品成功", c)
}

func (a *PointGoodsApi) FindPointGoods(c *gin.Context) {
	var info request.GetById
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	goods, err := pointGoodsService.GetPointGoods(info.Uint())
	if err != nil {
		global.GVA_LOG.Error("查询积分商品失败", zap.Error(err))
		response.FailWithMessage("查询积分商品失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(goods, "查询积分商品成功", c)
}

func (a *PointGoodsApi) GetPointGoodsList(c *gin.Context) {
	var pageInfo memberReq.PointGoodsSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pointGoodsService.GetPointGoodsList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取积分商品列表失败", zap.Error(err))
		response.FailWithMessage("获取积分商品列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取积分商品列表成功", c)
}

func (a *PointGoodsApi) UpdatePointGoodsStatus(c *gin.Context) {
	var req memberReq.UpdateGoodsStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointGoodsService.UpdatePointGoodsStatus(req); err != nil {
		global.GVA_LOG.Error("更新积分商品状态失败", zap.Error(err))
		response.FailWithMessage("更新积分商品状态失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新积分商品状态成功", c)
}

func (a *PointGoodsApi) UpdatePointGoodsStock(c *gin.Context) {
	var req memberReq.UpdateGoodsStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointGoodsService.UpdatePointGoodsStock(req); err != nil {
		global.GVA_LOG.Error("更新积分商品库存失败", zap.Error(err))
		response.FailWithMessage("更新积分商品库存失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新积分商品库存成功", c)
}

func (a *PointGoodsApi) GetPointGoodsOptions(c *gin.Context) {
	list, err := pointGoodsService.GetPointGoodsOptions(c.Query("keyword"))
	if err != nil {
		global.GVA_LOG.Error("获取积分商品选项失败", zap.Error(err))
		response.FailWithMessage("获取积分商品选项失败:"+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"list": list}, c)
}
