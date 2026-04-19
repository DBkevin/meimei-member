package member

import (
	"testing"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupLegacyMigrationTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := setupMemberTestDB(t)
	resetLegacyMigrationTables(t, db)
	createLegacyMigrationTables(t, db)
	t.Cleanup(func() {
		resetLegacyMigrationTables(t, db)
	})
	return db
}

func createLegacyMigrationTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	statements := []string{
		`CREATE TABLE member_members (
			id bigint unsigned NOT NULL AUTO_INCREMENT,
			created_at datetime(3) DEFAULT NULL,
			updated_at datetime(3) DEFAULT NULL,
			deleted_at datetime(3) DEFAULT NULL,
			openid varchar(128) DEFAULT NULL,
			unionid varchar(128) DEFAULT NULL,
			mobile varchar(32) DEFAULT NULL,
			nickname varchar(128) DEFAULT NULL,
			avatar_url varchar(255) DEFAULT NULL,
			real_name varchar(64) DEFAULT NULL,
			member_level varchar(32) DEFAULT 'standard',
			status varchar(32) DEFAULT 'enabled',
			PRIMARY KEY (id),
			UNIQUE KEY idx_member_mobile (mobile)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE member_point_accounts (
			id bigint unsigned NOT NULL AUTO_INCREMENT,
			created_at datetime(3) DEFAULT NULL,
			updated_at datetime(3) DEFAULT NULL,
			deleted_at datetime(3) DEFAULT NULL,
			member_id bigint unsigned DEFAULT NULL,
			available_points bigint DEFAULT 0,
			frozen_points bigint DEFAULT 0,
			total_earned_points bigint DEFAULT 0,
			total_used_points bigint DEFAULT 0,
			PRIMARY KEY (id),
			UNIQUE KEY idx_member_point_accounts_member_id (member_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE member_point_logs (
			id bigint unsigned NOT NULL AUTO_INCREMENT,
			created_at datetime(3) DEFAULT NULL,
			updated_at datetime(3) DEFAULT NULL,
			deleted_at datetime(3) DEFAULT NULL,
			member_id bigint unsigned DEFAULT NULL,
			change_type varchar(32) DEFAULT NULL,
			change_points bigint DEFAULT NULL,
			before_points bigint DEFAULT NULL,
			after_points bigint DEFAULT NULL,
			source_type varchar(32) DEFAULT NULL,
			source_id bigint unsigned DEFAULT NULL,
			remark varchar(255) DEFAULT NULL,
			operator_id bigint unsigned DEFAULT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE member_point_goods (
			id bigint unsigned NOT NULL AUTO_INCREMENT,
			created_at datetime(3) DEFAULT NULL,
			updated_at datetime(3) DEFAULT NULL,
			deleted_at datetime(3) DEFAULT NULL,
			name varchar(128) DEFAULT NULL,
			cover_image varchar(255) DEFAULT NULL,
			description text,
			points_price bigint DEFAULT NULL,
			stock bigint DEFAULT 0,
			limit_per_member bigint DEFAULT 0,
			status varchar(32) DEFAULT 'on_sale',
			sort bigint DEFAULT 0,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE member_exchange_orders (
			id bigint unsigned NOT NULL AUTO_INCREMENT,
			created_at datetime(3) DEFAULT NULL,
			updated_at datetime(3) DEFAULT NULL,
			deleted_at datetime(3) DEFAULT NULL,
			order_no varchar(64) DEFAULT NULL,
			member_id bigint unsigned DEFAULT NULL,
			goods_id bigint unsigned DEFAULT NULL,
			points_cost bigint DEFAULT NULL,
			status varchar(32) DEFAULT 'pending',
			verify_code varchar(32) DEFAULT NULL,
			verified_at datetime(3) DEFAULT NULL,
			operator_id bigint unsigned DEFAULT NULL,
			remark varchar(255) DEFAULT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY idx_member_exchange_orders_order_no (order_no)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}

	for _, statement := range statements {
		require.NoError(t, db.Exec(statement).Error)
	}
}

func resetLegacyMigrationTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	require.NoError(t, db.Exec("DROP TABLE IF EXISTS member_exchange_orders").Error)
	require.NoError(t, db.Exec("DROP TABLE IF EXISTS member_point_logs").Error)
	require.NoError(t, db.Exec("DROP TABLE IF EXISTS member_point_accounts").Error)
	require.NoError(t, db.Exec("DROP TABLE IF EXISTS member_point_goods").Error)
	require.NoError(t, db.Exec("DROP TABLE IF EXISTS member_members").Error)
}

func TestMigrateLegacyMemberData(t *testing.T) {
	db := setupLegacyMigrationTestDB(t)

	require.NoError(t, db.Exec(`
		INSERT INTO member_members (id, created_at, updated_at, openid, unionid, mobile, nickname, avatar_url, real_name, member_level, status)
		VALUES
			(1, '2026-04-18 10:00:00', '2026-04-18 10:00:00', 'openid-1', 'unionid-1', '13800000001', '小美', 'https://cdn.example.com/avatar.png', '张三', 'premium', 'enabled'),
			(2, '2026-04-18 11:00:00', '2026-04-18 11:00:00', NULL, NULL, '13800000002', '旧昵称', NULL, '', 'standard', 'disabled')
	`).Error)

	require.NoError(t, db.Exec(`
		INSERT INTO member_point_accounts (id, created_at, updated_at, member_id, available_points, frozen_points, total_earned_points, total_used_points)
		VALUES (11, '2026-04-18 10:05:00', '2026-04-18 10:05:00', 1, 10, 2, 12, 2)
	`).Error)

	require.NoError(t, db.Exec(`
		INSERT INTO member_point_goods (id, created_at, updated_at, name, cover_image, description, points_price, stock, limit_per_member, status, sort)
		VALUES (21, '2026-04-18 12:00:00', '2026-04-18 12:00:00', '玻尿酸补水', 'https://cdn.example.com/goods.png', '补水项目', 99, 8, 2, 'on_sale', 3)
	`).Error)

	require.NoError(t, db.Exec(`
		INSERT INTO member_exchange_orders (id, created_at, updated_at, order_no, member_id, goods_id, points_cost, status, verify_code, operator_id, remark)
		VALUES (31, '2026-04-18 13:00:00', '2026-04-18 13:00:00', 'MELEGACY001', 1, 21, 99, 'refunded', 'CODE31', 1001, '老订单备注')
	`).Error)

	require.NoError(t, db.Exec(`
		INSERT INTO member_point_logs (id, created_at, updated_at, member_id, change_type, change_points, before_points, after_points, source_type, source_id, operator_id, remark)
		VALUES
			(41, '2026-04-18 10:06:00', '2026-04-18 10:06:00', 1, 'adjust_add', 10, 0, 10, 'manual', 0, 1001, '手工补分'),
			(42, '2026-04-18 13:01:00', '2026-04-18 13:01:00', 1, 'use', 99, 10, -89, 'exchange_order', 31, 1001, '兑换扣减'),
			(43, '2026-04-18 13:02:00', '2026-04-18 13:02:00', 1, 'refund', 99, -89, 10, 'exchange_order_void', 31, 1001, '退款返还')
	`).Error)

	summary, err := MigrateLegacyMemberData(db)
	require.NoError(t, err)
	require.Equal(t, 2, summary.Members)
	require.Equal(t, 1, summary.Accounts)
	require.Equal(t, 1, summary.Products)
	require.Equal(t, 1, summary.Orders)
	require.Equal(t, 3, summary.Transactions)

	var member memberModel.Member
	require.NoError(t, db.Where("id = ?", 1).First(&member).Error)
	require.Equal(t, "张三", member.Name)
	require.Equal(t, "13800000001", member.Phone)
	require.Equal(t, "premium", member.Level)
	require.Equal(t, memberModel.MemberStatusEnabled, member.Status)
	require.Contains(t, member.Remark, "legacy_openid=openid-1")

	member = memberModel.Member{}
	require.NoError(t, db.Where("id = ?", 2).First(&member).Error)
	require.Equal(t, "旧昵称", member.Name)
	require.Equal(t, memberModel.MemberStatusDisabled, member.Status)

	var accounts []memberModel.PointAccount
	require.NoError(t, db.Order("member_id asc").Find(&accounts).Error)
	require.Len(t, accounts, 2)
	require.Equal(t, uint(1), accounts[0].MemberID)
	require.Equal(t, int64(10), accounts[0].Balance)
	require.Equal(t, int64(2), accounts[0].FrozenPoints)
	require.Equal(t, uint(2), accounts[1].MemberID)
	require.Equal(t, int64(0), accounts[1].Balance)

	var product memberModel.PointProduct
	require.NoError(t, db.Where("id = ?", 21).First(&product).Error)
	require.Equal(t, legacyDefaultProductCategory, product.Category)
	require.Equal(t, memberModel.PointProductStatusOnSale, product.Status)
	require.Contains(t, product.Description, "每人限兑数量: 2")

	var order memberModel.RedemptionOrder
	require.NoError(t, db.Where("id = ?", 31).First(&order).Error)
	require.Equal(t, uint(21), order.ProductID)
	require.Equal(t, "玻尿酸补水", order.ProductName)
	require.Equal(t, int64(1), order.Quantity)
	require.Equal(t, int64(99), order.UnitPoints)
	require.Equal(t, int64(99), order.TotalPoints)
	require.Equal(t, memberModel.RedemptionOrderStatusCancelled, order.Status)
	require.Equal(t, "张三", order.ReceiverName)
	require.Equal(t, "13800000001", order.ReceiverPhone)
	require.Contains(t, order.Remark, "legacy_verify_code=CODE31")

	var transactions []memberModel.PointTransaction
	require.NoError(t, db.Order("id asc").Find(&transactions).Error)
	require.Len(t, transactions, 3)
	require.Equal(t, accounts[0].ID, transactions[0].AccountID)
	require.Equal(t, memberModel.PointTransactionTypeAdjust, transactions[0].Type)
	require.Equal(t, memberModel.PointRefTypeManualAdjustAdd, transactions[0].RefType)
	require.Equal(t, int64(10), transactions[0].Points)
	require.Equal(t, memberModel.PointTransactionTypeSpend, transactions[1].Type)
	require.Equal(t, memberModel.PointRefTypeRedemptionOrder, transactions[1].RefType)
	require.Equal(t, uint(31), transactions[1].RefID)
	require.Equal(t, memberModel.PointTransactionTypeRefund, transactions[2].Type)
	require.Equal(t, memberModel.PointRefTypeRedemptionOrderVoid, transactions[2].RefType)
}
