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
		if err := ensureSuperAdminAuthorityButtons(tx, menus); err != nil {
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
		{ApiGroup: "会员管理", Method: "PUT", Path: "/member/updateMemberStatus", Description: "更新会员状态"},
		{ApiGroup: "会员管理", Method: "GET", Path: "/member/getMemberPointAccount", Description: "获取会员积分账户"},
		{ApiGroup: "会员管理", Method: "GET", Path: "/member/getMemberOptions", Description: "获取会员选项"},

		{ApiGroup: "积分账户", Method: "GET", Path: "/pointAccount/findPointAccount", Description: "查询积分账户"},
		{ApiGroup: "积分账户", Method: "GET", Path: "/pointAccount/getPointAccountList", Description: "获取积分账户列表"},
		{ApiGroup: "积分账户", Method: "POST", Path: "/pointAccount/manualAddPoints", Description: "手工增加积分"},
		{ApiGroup: "积分账户", Method: "POST", Path: "/pointAccount/manualSubPoints", Description: "手工扣减积分"},

		{ApiGroup: "积分流水", Method: "GET", Path: "/pointTransaction/getPointTransactionList", Description: "获取积分流水列表"},

		{ApiGroup: "积分商品", Method: "POST", Path: "/pointProduct/createPointProduct", Description: "创建积分商品"},
		{ApiGroup: "积分商品", Method: "DELETE", Path: "/pointProduct/deletePointProduct", Description: "删除积分商品"},
		{ApiGroup: "积分商品", Method: "PUT", Path: "/pointProduct/updatePointProduct", Description: "更新积分商品"},
		{ApiGroup: "积分商品", Method: "GET", Path: "/pointProduct/findPointProduct", Description: "查询积分商品详情"},
		{ApiGroup: "积分商品", Method: "GET", Path: "/pointProduct/getPointProductList", Description: "获取积分商品列表"},
		{ApiGroup: "积分商品", Method: "PUT", Path: "/pointProduct/updatePointProductStatus", Description: "上下架积分商品"},
		{ApiGroup: "积分商品", Method: "GET", Path: "/pointProduct/getPointProductOptions", Description: "获取积分商品选项"},

		{ApiGroup: "兑换订单", Method: "POST", Path: "/redemptionOrder/createRedemptionOrder", Description: "创建兑换订单"},
		{ApiGroup: "兑换订单", Method: "GET", Path: "/redemptionOrder/findRedemptionOrder", Description: "查询兑换订单详情"},
		{ApiGroup: "兑换订单", Method: "GET", Path: "/redemptionOrder/getRedemptionOrderList", Description: "获取兑换订单列表"},
		{ApiGroup: "兑换订单", Method: "POST", Path: "/redemptionOrder/completeRedemptionOrder", Description: "完成兑换订单"},
		{ApiGroup: "兑换订单", Method: "POST", Path: "/redemptionOrder/cancelRedemptionOrder", Description: "取消兑换订单"},
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
		Buttons   []systemModel.SysBaseMenuBtn
	}{
		{
			Name:      "memberList",
			Path:      "members",
			Component: "view/member/member/index.vue",
			Sort:      1,
			Title:     "会员管理",
			Icon:      "user",
			Buttons: []systemModel.SysBaseMenuBtn{
				{Name: "add", Desc: "新增会员"},
				{Name: "edit", Desc: "编辑会员"},
				{Name: "status", Desc: "切换会员状态"},
				{Name: "account", Desc: "查看积分账户"},
				{Name: "delete", Desc: "删除会员"},
			},
		},
		{
			Name:      "pointAccount",
			Path:      "point-accounts",
			Component: "view/member/account/index.vue",
			Sort:      2,
			Title:     "积分账户",
			Icon:      "coin",
			Buttons: []systemModel.SysBaseMenuBtn{
				{Name: "adjustAdd", Desc: "手工增加积分"},
				{Name: "adjustSub", Desc: "手工扣减积分"},
			},
		},
		{Name: "pointTransaction", Path: "point-transactions", Component: "view/member/log/index.vue", Sort: 3, Title: "积分流水", Icon: "tickets"},
		{
			Name:      "pointProduct",
			Path:      "point-products",
			Component: "view/member/goods/index.vue",
			Sort:      4,
			Title:     "积分商品",
			Icon:      "goods-filled",
			Buttons: []systemModel.SysBaseMenuBtn{
				{Name: "add", Desc: "新增积分商品"},
				{Name: "edit", Desc: "编辑积分商品"},
				{Name: "status", Desc: "上下架积分商品"},
				{Name: "delete", Desc: "删除积分商品"},
			},
		},
		{
			Name:      "redemptionOrder",
			Path:      "redemption-orders",
			Component: "view/member/order/index.vue",
			Sort:      5,
			Title:     "兑换订单",
			Icon:      "memo",
			Buttons: []systemModel.SysBaseMenuBtn{
				{Name: "add", Desc: "新建兑换订单"},
				{Name: "info", Desc: "查看订单详情"},
				{Name: "complete", Desc: "完成兑换订单"},
				{Name: "cancel", Desc: "取消兑换订单"},
			},
		},
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
		if upsertErr = ensureMenuButtons(tx, menu.ID, item.Buttons); upsertErr != nil {
			return nil, upsertErr
		}
		result = append(result, menu)
	}
	return result, nil
}

func ensureMenuButtons(tx *gorm.DB, menuID uint, buttons []systemModel.SysBaseMenuBtn) error {
	if len(buttons) == 0 {
		return nil
	}

	for _, item := range buttons {
		var button systemModel.SysBaseMenuBtn
		err := tx.Where("sys_base_menu_id = ? AND name = ?", menuID, item.Name).First(&button).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item.SysBaseMenuID = menuID
			if err = tx.Create(&item).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if err = tx.Model(&button).Update("desc", item.Desc).Error; err != nil {
			return err
		}
	}
	return nil
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

func ensureSuperAdminAuthorityButtons(tx *gorm.DB, menus []systemModel.SysBaseMenu) error {
	const superAdminAuthorityID = 888

	menuIDs := make([]uint, 0, len(menus))
	for _, item := range menus {
		if item.ParentId == 0 {
			continue
		}
		menuIDs = append(menuIDs, item.ID)
	}
	if len(menuIDs) == 0 {
		return nil
	}

	var existing []systemModel.SysAuthorityBtn
	if err := tx.Where("authority_id = ? AND sys_menu_id IN ?", superAdminAuthorityID, menuIDs).Find(&existing).Error; err != nil {
		return err
	}
	existingMap := make(map[uint]map[uint]struct{}, len(existing))
	for _, item := range existing {
		if existingMap[item.SysMenuID] == nil {
			existingMap[item.SysMenuID] = make(map[uint]struct{})
		}
		existingMap[item.SysMenuID][item.SysBaseMenuBtnID] = struct{}{}
	}

	missing := make([]systemModel.SysAuthorityBtn, 0)
	for _, menu := range menus {
		if menu.ParentId == 0 {
			continue
		}
		var buttons []systemModel.SysBaseMenuBtn
		if err := tx.Where("sys_base_menu_id = ?", menu.ID).Find(&buttons).Error; err != nil {
			return err
		}
		for _, button := range buttons {
			if _, ok := existingMap[menu.ID][button.ID]; ok {
				continue
			}
			missing = append(missing, systemModel.SysAuthorityBtn{
				AuthorityId:      superAdminAuthorityID,
				SysMenuID:        menu.ID,
				SysBaseMenuBtnID: button.ID,
			})
		}
	}

	if len(missing) == 0 {
		return nil
	}
	return tx.Create(&missing).Error
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

		{Ptype: "p", V0: "888", V1: "/pointAccount/findPointAccount", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/pointAccount/getPointAccountList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/pointAccount/manualAddPoints", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/pointAccount/manualSubPoints", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/pointTransaction/getPointTransactionList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/pointProduct/createPointProduct", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/pointProduct/deletePointProduct", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/pointProduct/updatePointProduct", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/pointProduct/findPointProduct", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/pointProduct/getPointProductList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/pointProduct/updatePointProductStatus", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/pointProduct/getPointProductOptions", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/redemptionOrder/createRedemptionOrder", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/redemptionOrder/findRedemptionOrder", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/redemptionOrder/getRedemptionOrderList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/redemptionOrder/completeRedemptionOrder", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/redemptionOrder/cancelRedemptionOrder", V2: "POST"},
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
