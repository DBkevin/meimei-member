package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PointLogApi struct{}

func (a *PointLogApi) GetPointLogList(c *gin.Context) {
	var pageInfo memberReq.PointLogSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pointLogService.GetPointLogList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取积分流水列表失败", zap.Error(err))
		response.FailWithMessage("获取积分流水列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取积分流水列表成功", c)
}
