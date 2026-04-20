package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	commonRes "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PointProductApi struct{}

// CreatePointProduct
// @Tags      PointProduct
// @Summary   创建积分商品
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberReq.CreatePointProductReq true  "积分商品信息"
// @Success   200   {object}  commonRes.Response{msg=string}  "创建积分商品成功"
// @Router    /pointProduct/createPointProduct [post]
func (a *PointProductApi) CreatePointProduct(c *gin.Context) {
	var req memberReq.CreatePointProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := ValidatePointProductInput(req.PointProductBaseInput); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointProductService.CreatePointProduct(req); err != nil {
		global.GVA_LOG.Error("创建积分商品失败", zap.Error(err))
		commonRes.FailWithMessage("创建积分商品失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithMessage("创建积分商品成功", c)
}

// DeletePointProduct
// @Tags      PointProduct
// @Summary   删除积分商品
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      commonReq.GetById             true  "积分商品ID"
// @Success   200   {object}  commonRes.Response{msg=string} "删除积分商品成功"
// @Router    /pointProduct/deletePointProduct [delete]
func (a *PointProductApi) DeletePointProduct(c *gin.Context) {
	var info commonReq.GetById
	if err := c.ShouldBindJSON(&info); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := ValidateIDInput(info.Uint()); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointProductService.DeletePointProduct(info.Uint()); err != nil {
		global.GVA_LOG.Error("删除积分商品失败", zap.Error(err))
		commonRes.FailWithMessage("删除积分商品失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithMessage("删除积分商品成功", c)
}

// UpdatePointProduct
// @Tags      PointProduct
// @Summary   更新积分商品
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberReq.UpdatePointProductReq true  "积分商品信息"
// @Success   200   {object}  commonRes.Response{msg=string}  "更新积分商品成功"
// @Router    /pointProduct/updatePointProduct [put]
func (a *PointProductApi) UpdatePointProduct(c *gin.Context) {
	var req memberReq.UpdatePointProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := ValidateIDInput(req.ID); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := ValidatePointProductInput(req.PointProductBaseInput); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointProductService.UpdatePointProduct(req); err != nil {
		global.GVA_LOG.Error("更新积分商品失败", zap.Error(err))
		commonRes.FailWithMessage("更新积分商品失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithMessage("更新积分商品成功", c)
}

// FindPointProduct
// @Tags      PointProduct
// @Summary   查询积分商品详情
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     id    query     int                true  "积分商品ID"
// @Success   200   {object}  commonRes.Response "查询积分商品成功"
// @Router    /pointProduct/findPointProduct [get]
func (a *PointProductApi) FindPointProduct(c *gin.Context) {
	var info commonReq.GetById
	if err := c.ShouldBindQuery(&info); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	product, err := pointProductService.GetPointProduct(info.Uint())
	if err != nil {
		global.GVA_LOG.Error("查询积分商品失败", zap.Error(err))
		commonRes.FailWithMessage("查询积分商品失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(product, "查询积分商品成功", c)
}

// GetPointProductList
// @Tags      PointProduct
// @Summary   分页获取积分商品列表
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  query     memberReq.PointProductSearch true  "分页与筛选条件"
// @Success   200   {object}  commonRes.Response{data=commonRes.PageResult,msg=string} "获取积分商品列表成功"
// @Router    /pointProduct/getPointProductList [get]
func (a *PointProductApi) GetPointProductList(c *gin.Context) {
	var pageInfo memberReq.PointProductSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pointProductService.GetPointProductList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取积分商品列表失败", zap.Error(err))
		commonRes.FailWithMessage("获取积分商品列表失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(commonRes.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取积分商品列表成功", c)
}

// UpdatePointProductStatus
// @Tags      PointProduct
// @Summary   更新积分商品状态
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberReq.UpdatePointProductStatusReq true  "积分商品ID与状态"
// @Success   200   {object}  commonRes.Response{msg=string} "更新积分商品状态成功"
// @Router    /pointProduct/updatePointProductStatus [put]
func (a *PointProductApi) UpdatePointProductStatus(c *gin.Context) {
	var req memberReq.UpdatePointProductStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	if err := pointProductService.UpdatePointProductStatus(req); err != nil {
		global.GVA_LOG.Error("更新积分商品状态失败", zap.Error(err))
		commonRes.FailWithMessage("更新积分商品状态失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithMessage("更新积分商品状态成功", c)
}

// GetPointProductOptions
// @Tags      PointProduct
// @Summary   获取积分商品选项
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     keyword  query     string             false  "关键字"
// @Success   200      {object}  commonRes.Response "获取积分商品选项成功"
// @Router    /pointProduct/getPointProductOptions [get]
func (a *PointProductApi) GetPointProductOptions(c *gin.Context) {
	list, err := pointProductService.GetPointProductOptions(c.Query("keyword"))
	if err != nil {
		global.GVA_LOG.Error("获取积分商品选项失败", zap.Error(err))
		commonRes.FailWithMessage("获取积分商品选项失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithData(gin.H{"list": list}, c)
}
