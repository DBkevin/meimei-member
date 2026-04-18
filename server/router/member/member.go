package member

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MemberRouter struct{}

func (r *MemberRouter) InitMemberRouter(Router *gin.RouterGroup) {
	memberRouter := Router.Group("member").Use(middleware.OperationRecord())
	memberRouterWithoutRecord := Router.Group("member")
	{
		memberRouter.POST("createMember", memberApi.CreateMember)
		memberRouter.DELETE("deleteMember", memberApi.DeleteMember)
		memberRouter.PUT("updateMember", memberApi.UpdateMember)
		memberRouter.PUT("updateMemberStatus", memberApi.UpdateMemberStatus)
	}
	{
		memberRouterWithoutRecord.GET("findMember", memberApi.FindMember)
		memberRouterWithoutRecord.GET("getMemberList", memberApi.GetMemberList)
		memberRouterWithoutRecord.GET("getMemberPointAccount", memberApi.GetMemberPointAccount)
		memberRouterWithoutRecord.GET("getMemberOptions", memberApi.GetMemberOptions)
	}
}
