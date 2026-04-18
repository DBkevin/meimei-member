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

type MemberApi struct{}

// CreateMember
// @Tags      Member
// @Summary   创建会员
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberModel.Member            true  "会员信息"
// @Success   200   {object}  response.Response{msg=string} "创建会员成功"
// @Router    /member/createMember [post]
func (a *MemberApi) CreateMember(c *gin.Context) {
	var member memberModel.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := memberService.CreateMember(&member); err != nil {
		global.GVA_LOG.Error("创建会员失败", zap.Error(err))
		response.FailWithMessage("创建会员失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建会员成功", c)
}

// DeleteMember
// @Tags      Member
// @Summary   删除会员
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById              true  "会员ID"
// @Success   200   {object}  response.Response{msg=string} "删除会员成功"
// @Router    /member/deleteMember [delete]
func (a *MemberApi) DeleteMember(c *gin.Context) {
	var info request.GetById
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := memberService.DeleteMember(info.Uint()); err != nil {
		global.GVA_LOG.Error("删除会员失败", zap.Error(err))
		response.FailWithMessage("删除会员失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除会员成功", c)
}

// UpdateMember
// @Tags      Member
// @Summary   更新会员
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberModel.Member            true  "会员信息"
// @Success   200   {object}  response.Response{msg=string} "更新会员成功"
// @Router    /member/updateMember [put]
func (a *MemberApi) UpdateMember(c *gin.Context) {
	var member memberModel.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if member.ID == 0 {
		response.FailWithMessage("会员ID不能为空", c)
		return
	}
	if err := memberService.UpdateMember(&member); err != nil {
		global.GVA_LOG.Error("更新会员失败", zap.Error(err))
		response.FailWithMessage("更新会员失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新会员成功", c)
}

// FindMember
// @Tags      Member
// @Summary   查询会员详情
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     id    query     int                        true  "会员ID"
// @Success   200   {object}  response.Response          "查询会员成功"
// @Router    /member/findMember [get]
func (a *MemberApi) FindMember(c *gin.Context) {
	var info request.GetById
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := memberService.GetMember(info.Uint())
	if err != nil {
		global.GVA_LOG.Error("查询会员失败", zap.Error(err))
		response.FailWithMessage("查询会员失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "查询会员成功", c)
}

// GetMemberList
// @Tags      Member
// @Summary   分页获取会员列表
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  query     memberReq.MemberSearch      true  "分页与筛选条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string} "获取会员列表成功"
// @Router    /member/getMemberList [get]
func (a *MemberApi) GetMemberList(c *gin.Context) {
	var pageInfo memberReq.MemberSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := memberService.GetMemberList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取会员列表失败", zap.Error(err))
		response.FailWithMessage("获取会员列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取会员列表成功", c)
}

// UpdateMemberStatus
// @Tags      Member
// @Summary   更新会员状态
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     data  body      memberReq.UpdateMemberStatusReq true  "会员ID与状态"
// @Success   200   {object}  response.Response{msg=string} "更新会员状态成功"
// @Router    /member/updateMemberStatus [put]
func (a *MemberApi) UpdateMemberStatus(c *gin.Context) {
	var req memberReq.UpdateMemberStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := memberService.UpdateMemberStatus(req); err != nil {
		global.GVA_LOG.Error("更新会员状态失败", zap.Error(err))
		response.FailWithMessage("更新会员状态失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新会员状态成功", c)
}

// GetMemberPointAccount
// @Tags      Member
// @Summary   获取会员积分账户
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     id    query     int                        true  "会员ID"
// @Success   200   {object}  response.Response          "获取会员积分账户成功"
// @Router    /member/getMemberPointAccount [get]
func (a *MemberApi) GetMemberPointAccount(c *gin.Context) {
	var info request.GetById
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	account, err := pointAccountService.GetPointAccountByMemberID(info.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取会员积分账户失败", zap.Error(err))
		response.FailWithMessage("获取会员积分账户失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(account, "获取会员积分账户成功", c)
}

// GetMemberOptions
// @Tags      Member
// @Summary   获取会员选项
// @Security  ApiKeyAuth
// @Accept    application/json
// @Produce   application/json
// @Param     keyword  query     string                     false  "关键字"
// @Success   200      {object}  response.Response          "获取会员选项成功"
// @Router    /member/getMemberOptions [get]
func (a *MemberApi) GetMemberOptions(c *gin.Context) {
	var req memberReq.MemberOptionsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := memberService.GetMemberOptions(req.Keyword)
	if err != nil {
		global.GVA_LOG.Error("获取会员选项失败", zap.Error(err))
		response.FailWithMessage("获取会员选项失败:"+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"list": list}, c)
}
