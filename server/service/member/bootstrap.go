package member

import (
	"errors"
	"strconv"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	systemModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

func Bootstrap() error {
	if global.GVA_DB == nil {
		return nil
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := ensureSysApis(tx); err != nil {
			return err
		}
		menus, err := ensureMenus(tx)
		if err != nil {
			return err
		}
		if err := ensureSuperAdminMenus(tx, menus); err != nil {
			return err
		}
		return ensureSuperAdminCasbin(tx)
	})
}

func ensureSysApis(tx *gorm.DB) error {
	apis := []systemModel.SysApi{
		{ApiGroup: "会员管理", Method: "POST", Path: "/member/createMember", Description: "创建会员"},
		{ApiGroup: "会员管理", Method: "DELETE", Path: "/member/deleteMember", Description: "删除会员"},
		{ApiGroup: "会员管理", Method: "PUT", Path: "/member/updateMember", Description: "更新会员"},
		{ApiGroup: "会员管理", Method: "GET", Path: "/member/findMember", Description: "查询会员详情"},
		{ApiGroup: "会员管理", Method: "GET", Path: "/member/getMemberList", Description: "获取会员列表"},
		{ApiGroup: "会员管理", Method: "PUT", Path: "/member/updateMemberStatus", Description: "启用禁用会员"},
		{ApiGroup: "会员管理", Method: "GET", Path: "/member/getMemberPointAccount", Description: "获取会员积分账户"},
		{ApiGroup: "会员管理", Method: "GET", Path: "/member/getMemberOptions", Description: "获取会员选项"},

		{ApiGroup: "积分账户", Method: "GET", Path: "/pointAccount/getPointAccountList", Description: "获取积分账户列表"},
		{ApiGroup: "积分账户", Method: "POST", Path: "/pointAccount/manualAddPoints", Description: "手工增加积分"},
		{ApiGroup: "积分账户", Method: "POST", Path: "/pointAccount/manualSubPoints", Description: "手工扣减积分"},

		{ApiGroup: "积分流水", Method: "GET", Path: "/pointLog/getPointLogList", Description: "获取积分流水列表"},

		{ApiGroup: "积分商品", Method: "POST", Path: "/pointGoods/createPointGoods", Description: "创建积分商品"},
		{ApiGroup: "积分商品", Method: "DELETE", Path: "/pointGoods/deletePointGoods", Description: "删除积分商品"},
		{ApiGroup: "积分商品", Method: "PUT", Path: "/pointGoods/updatePointGoods", Description: "更新积分商品"},
		{ApiGroup: "积分商品", Method: "GET", Path: "/pointGoods/findPointGoods", Description: "查询积分商品详情"},
		{ApiGroup: "积分商品", Method: "GET", Path: "/pointGoods/getPointGoodsList", Description: "获取积分商品列表"},
		{ApiGroup: "积分商品", Method: "PUT", Path: "/pointGoods/updatePointGoodsStatus", Description: "上下架积分商品"},
		{ApiGroup: "积分商品", Method: "PUT", Path: "/pointGoods/updatePointGoodsStock", Description: "更新积分商品库存"},
		{ApiGroup: "积分商品", Method: "GET", Path: "/pointGoods/getPointGoodsOptions", Description: "获取积分商品选项"},

		{ApiGroup: "兑换订单", Method: "POST", Path: "/exchangeOrder/createExchangeOrder", Description: "创建兑换订单"},
		{ApiGroup: "兑换订单", Method: "GET", Path: "/exchangeOrder/findExchangeOrder", Description: "查询兑换订单详情"},
		{ApiGroup: "兑换订单", Method: "GET", Path: "/exchangeOrder/getExchangeOrderList", Description: "获取兑换订单列表"},
		{ApiGroup: "兑换订单", Method: "POST", Path: "/exchangeOrder/verifyExchangeOrder", Description: "核销兑换订单"},
		{ApiGroup: "兑换订单", Method: "POST", Path: "/exchangeOrder/cancelExchangeOrder", Description: "取消兑换订单"},
		{ApiGroup: "兑换订单", Method: "POST", Path: "/exchangeOrder/refundExchangeOrder", Description: "退款兑换订单"},
	}

	for _, item := range apis {
		var api systemModel.SysApi
		err := tx.Where("path = ? AND method = ?", item.Path, item.Method).First(&api).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&item).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if err = tx.Model(&api).Updates(map[string]interface{}{
			"api_group":   item.ApiGroup,
			"description": item.Description,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func ensureMenus(tx *gorm.DB) ([]systemModel.SysBaseMenu, error) {
	parent, err := upsertMenu(tx, "medicalMember", map[string]interface{}{
		"menu_level": 0,
		"parent_id":  0,
		"path":       "medical-member",
		"hidden":     false,
		"component":  "view/routerHolder.vue",
		"sort":       20,
		"title":      "医美会员系统",
		"icon":       "tickets",
		"keep_alive": false,
	}, systemModel.SysBaseMenu{Name: "medicalMember"})
	if err != nil {
		return nil, err
	}

	children := []struct {
		Name      string
		Path      string
		Component string
		Sort      int
		Title     string
		Icon      string
	}{
		{Name: "memberList", Path: "members", Component: "view/member/member/index.vue", Sort: 1, Title: "会员管理", Icon: "user"},
		{Name: "pointAccount", Path: "point-accounts", Component: "view/member/account/index.vue", Sort: 2, Title: "积分账户", Icon: "coin"},
		{Name: "pointLog", Path: "point-logs", Component: "view/member/log/index.vue", Sort: 3, Title: "积分流水", Icon: "tickets"},
		{Name: "pointGoods", Path: "point-goods", Component: "view/member/goods/index.vue", Sort: 4, Title: "积分商品", Icon: "goods-filled"},
		{Name: "exchangeOrder", Path: "exchange-orders", Component: "view/member/order/index.vue", Sort: 5, Title: "兑换订单", Icon: "memo"},
	}

	result := []systemModel.SysBaseMenu{parent}
	for _, item := range children {
		menu, upsertErr := upsertMenu(tx, item.Name, map[string]interface{}{
			"menu_level": 1,
			"parent_id":  parent.ID,
			"path":       item.Path,
			"hidden":     false,
			"component":  item.Component,
			"sort":       item.Sort,
			"title":      item.Title,
			"icon":       item.Icon,
			"keep_alive": true,
		}, systemModel.SysBaseMenu{Name: item.Name})
		if upsertErr != nil {
			return nil, upsertErr
		}
		result = append(result, menu)
	}
	return result, nil
}

func upsertMenu(tx *gorm.DB, name string, updates map[string]interface{}, defaults systemModel.SysBaseMenu) (systemModel.SysBaseMenu, error) {
	var menu systemModel.SysBaseMenu
	err := tx.Where("name = ?", name).First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		defaults.Name = name
		if err = tx.Create(&defaults).Error; err != nil {
			return menu, err
		}
		menu = defaults
	} else if err != nil {
		return menu, err
	}
	if err = tx.Model(&menu).Updates(updates).Error; err != nil {
		return menu, err
	}
	err = tx.Where("id = ?", menu.ID).First(&menu).Error
	return menu, err
}

func ensureSuperAdminMenus(tx *gorm.DB, menus []systemModel.SysBaseMenu) error {
	var auth systemModel.SysAuthority
	if err := tx.Where("authority_id = ?", 888).First(&auth).Error; err != nil {
		return err
	}

	var existingRelations []systemModel.SysAuthorityMenu
	if err := tx.Where("sys_authority_authority_id = ?", auth.AuthorityId).Find(&existingRelations).Error; err != nil {
		return err
	}

	existingMenuIDs := make(map[string]struct{}, len(existingRelations))
	for _, item := range existingRelations {
		existingMenuIDs[item.MenuId] = struct{}{}
	}

	authorityID := strconv.FormatUint(uint64(auth.AuthorityId), 10)
	missingRelations := make([]systemModel.SysAuthorityMenu, 0, len(menus))
	for _, item := range menus {
		menuID := strconv.FormatUint(uint64(item.ID), 10)
		if _, ok := existingMenuIDs[menuID]; ok {
			continue
		}
		missingRelations = append(missingRelations, systemModel.SysAuthorityMenu{
			MenuId:      menuID,
			AuthorityId: authorityID,
		})
	}

	if len(missingRelations) == 0 {
		return nil
	}

	return tx.Create(&missingRelations).Error
}

func ensureSuperAdminCasbin(tx *gorm.DB) error {
	rules := []adapter.CasbinRule{
		{Ptype: "p", V0: "888", V1: "/member/createMember", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/member/deleteMember", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/member/updateMember", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/member/findMember", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/member/getMemberList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/member/updateMemberStatus", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/member/getMemberPointAccount", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/member/getMemberOptions", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/pointAccount/getPointAccountList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/pointAccount/manualAddPoints", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/pointAccount/manualSubPoints", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/pointLog/getPointLogList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/pointGoods/createPointGoods", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/pointGoods/deletePointGoods", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/pointGoods/updatePointGoods", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/pointGoods/findPointGoods", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/pointGoods/getPointGoodsList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/pointGoods/updatePointGoodsStatus", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/pointGoods/updatePointGoodsStock", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/pointGoods/getPointGoodsOptions", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/exchangeOrder/createExchangeOrder", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/exchangeOrder/findExchangeOrder", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/exchangeOrder/getExchangeOrderList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/exchangeOrder/verifyExchangeOrder", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/exchangeOrder/cancelExchangeOrder", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/exchangeOrder/refundExchangeOrder", V2: "POST"},
	}

	for _, item := range rules {
		var rule adapter.CasbinRule
		err := tx.Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", item.Ptype, item.V0, item.V1, item.V2).First(&rule).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&item).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}
