package member

import (
	"testing"
	"time"

	adapter "github.com/casbin/gorm-adapter/v3"
	systemModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type bootstrapTestAPI struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Path        string
	Description string
	ApiGroup    string `gorm:"column:api_group"`
	Method      string
}

func (bootstrapTestAPI) TableName() string {
	return "sys_apis"
}

type bootstrapTestAuthority struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	AuthorityId   uint           `gorm:"primaryKey;column:authority_id"`
	AuthorityName string         `gorm:"column:authority_name"`
	ParentId      *uint          `gorm:"column:parent_id"`
	DefaultRouter string         `gorm:"column:default_router"`
}

func (bootstrapTestAuthority) TableName() string {
	return "sys_authorities"
}

type bootstrapTestBaseMenu struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	MenuLevel      uint           `gorm:"column:menu_level"`
	ParentId       uint           `gorm:"column:parent_id"`
	Path           string
	Name           string
	Hidden         bool
	Component      string
	Sort           int
	ActiveName     string `gorm:"column:active_name"`
	KeepAlive      bool   `gorm:"column:keep_alive"`
	DefaultMenu    bool   `gorm:"column:default_menu"`
	Title          string
	Icon           string
	CloseTab       bool   `gorm:"column:close_tab"`
	TransitionType string `gorm:"column:transition_type"`
}

func (bootstrapTestBaseMenu) TableName() string {
	return "sys_base_menus"
}

type bootstrapTestBaseMenuParameter struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	SysBaseMenuID uint           `gorm:"column:sys_base_menu_id"`
	Type          string
	Key           string
	Value         string
}

func (bootstrapTestBaseMenuParameter) TableName() string {
	return "sys_base_menu_parameters"
}

type bootstrapTestBaseMenuBtn struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Name          string
	Desc          string
	SysBaseMenuID uint `gorm:"column:sys_base_menu_id"`
}

func (bootstrapTestBaseMenuBtn) TableName() string {
	return "sys_base_menu_btns"
}

type bootstrapTestAuthorityMenu struct {
	MenuId      string `gorm:"column:sys_base_menu_id"`
	AuthorityId string `gorm:"column:sys_authority_authority_id"`
}

func (bootstrapTestAuthorityMenu) TableName() string {
	return "sys_authority_menus"
}

type bootstrapTestAuthorityBtn struct {
	AuthorityId      uint `gorm:"column:authority_id"`
	SysMenuID        uint `gorm:"column:sys_menu_id"`
	SysBaseMenuBtnID uint `gorm:"column:sys_base_menu_btn_id"`
}

func (bootstrapTestAuthorityBtn) TableName() string {
	return "sys_authority_btns"
}

func setupBootstrapTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := setupMemberTestDB(t)
	require.NoError(t, db.AutoMigrate(
		&bootstrapTestAPI{},
		&bootstrapTestAuthority{},
		&bootstrapTestBaseMenu{},
		&bootstrapTestBaseMenuParameter{},
		&bootstrapTestBaseMenuBtn{},
		&bootstrapTestAuthorityMenu{},
		&bootstrapTestAuthorityBtn{},
		&adapter.CasbinRule{},
	))

	resetBootstrapTestTables(t, db)
	t.Cleanup(func() {
		resetBootstrapTestTables(t, db)
	})

	return db
}

func resetBootstrapTestTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	require.NoError(t, db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error)

	for _, model := range []interface{}{
		&bootstrapTestAuthorityBtn{},
		&bootstrapTestAuthorityMenu{},
		&bootstrapTestBaseMenuBtn{},
		&bootstrapTestBaseMenuParameter{},
		&bootstrapTestBaseMenu{},
		&bootstrapTestAPI{},
		&adapter.CasbinRule{},
		&bootstrapTestAuthority{},
	} {
		require.NoError(t, db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(model).Error)
	}

	require.NoError(t, db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error)
}

func TestBootstrapSeedsMemberMenusButtonsAndAuthorityGrants(t *testing.T) {
	db := setupBootstrapTestDB(t)

	require.NoError(t, db.Create(&bootstrapTestAuthority{
		AuthorityId:   888,
		AuthorityName: "超级管理员",
		DefaultRouter: "dashboard",
	}).Error)

	require.NoError(t, Bootstrap())
	require.NoError(t, Bootstrap())

	var apiCount int64
	require.NoError(t, db.Model(&systemModel.SysApi{}).Count(&apiCount).Error)
	require.EqualValues(t, 25, apiCount)

	var parent systemModel.SysBaseMenu
	require.NoError(t, db.Where("name = ?", "medicalMember").First(&parent).Error)
	require.EqualValues(t, 0, parent.ParentId)

	var children []systemModel.SysBaseMenu
	require.NoError(t, db.Where("parent_id = ?", parent.ID).Order("sort asc").Find(&children).Error)
	require.Len(t, children, 5)

	expectedButtons := map[string][]string{
		"memberList":       {"add", "edit", "status", "account", "delete"},
		"pointAccount":     {"adjustAdd", "adjustSub"},
		"pointTransaction": nil,
		"pointProduct":     {"add", "edit", "status", "delete"},
		"redemptionOrder":  {"add", "info", "complete", "cancel"},
	}

	for _, menu := range children {
		expected, ok := expectedButtons[menu.Name]
		require.Truef(t, ok, "发现未预期菜单: %s", menu.Name)

		var buttons []systemModel.SysBaseMenuBtn
		require.NoError(t, db.Where("sys_base_menu_id = ?", menu.ID).Order("id asc").Find(&buttons).Error)
		require.Len(t, buttons, len(expected))

		for index, button := range buttons {
			require.Equal(t, expected[index], button.Name)

			var grantCount int64
			require.NoError(t, db.Model(&systemModel.SysAuthorityBtn{}).
				Where("authority_id = ? AND sys_menu_id = ? AND sys_base_menu_btn_id = ?", 888, menu.ID, button.ID).
				Count(&grantCount).Error)
			require.EqualValues(t, 1, grantCount)
		}
	}

	var authorityMenuCount int64
	require.NoError(t, db.Model(&systemModel.SysAuthorityMenu{}).
		Where("sys_authority_authority_id = ?", "888").
		Count(&authorityMenuCount).Error)
	require.EqualValues(t, 6, authorityMenuCount)

	var authorityButtonCount int64
	require.NoError(t, db.Model(&systemModel.SysAuthorityBtn{}).
		Where("authority_id = ?", 888).
		Count(&authorityButtonCount).Error)
	require.EqualValues(t, 15, authorityButtonCount)

	var casbinRuleCount int64
	require.NoError(t, db.Model(&adapter.CasbinRule{}).
		Where("ptype = ? AND v0 = ?", "p", "888").
		Count(&casbinRuleCount).Error)
	require.EqualValues(t, 25, casbinRuleCount)
}
