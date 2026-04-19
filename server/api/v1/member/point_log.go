package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonRes "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PointTransactionApi struct{}

// GetPointTransactionList
// @Tags      PointTransaction
// @Summary   分页获取积分流水列表
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  query     memberReq.PointTransactionSearch true  "分页与筛选条件"
// @Success   200   {object}  commonRes.Response{data=commonRes.PageResult,msg=string} "获取积分流水列表成功"
// @Router    /pointTransaction/getPointTransactionList [get]
func (a *PointTransactionApi) GetPointTransactionList(c *gin.Context) {
	var pageInfo memberReq.PointTransactionSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		commonRes.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pointTransactionService.GetPointTransactionList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取积分流水列表失败", zap.Error(err))
		commonRes.FailWithMessage("获取积分流水列表失败:"+err.Error(), c)
		return
	}
	commonRes.OkWithDetailed(commonRes.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取积分流水列表成功", c)
}
