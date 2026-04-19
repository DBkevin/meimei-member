package member

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type memberTestFileConfig struct {
	Mysql config.Mysql `yaml:"mysql"`
}

func setupMemberTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	mysqlConfig, dsn := loadMemberTestDBConfig(t)
	if mysqlConfig != nil {
		ensureMemberTestDatabase(t, *mysqlConfig)
	}

	db, err := gorm.Open(gormmysql.New(gormmysql.Config{
		DSN:               dsn,
		DefaultStringSize: 191,
	}), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(4)

	err = db.AutoMigrate(
		&memberModel.Member{},
		&memberModel.PointAccount{},
		&memberModel.PointTransaction{},
		&memberModel.PointProduct{},
		&memberModel.RedemptionOrder{},
	)
	require.NoError(t, err)

	resetMemberTestTables(t, db)

	oldDB := global.GVA_DB
	global.GVA_DB = db
	t.Cleanup(func() {
		resetMemberTestTables(t, db)
		global.GVA_DB = oldDB
		_ = sqlDB.Close()
	})

	return db
}

func loadMemberTestDBConfig(t *testing.T) (*config.Mysql, string) {
	t.Helper()

	if dsn := strings.TrimSpace(os.Getenv("MEMBER_TEST_DSN")); dsn != "" {
		return nil, dsn
	}

	configPath := strings.TrimSpace(os.Getenv("MEMBER_TEST_CONFIG"))
	if configPath == "" {
		_, filename, _, ok := runtime.Caller(0)
		require.True(t, ok, "无法定位会员测试文件路径")
		configPath = filepath.Join(filepath.Dir(filename), "..", "..", "config.local.yaml")
	}

	content, err := os.ReadFile(configPath)
	require.NoError(t, err, "读取会员测试配置失败: %s", configPath)

	var cfg memberTestFileConfig
	err = yaml.Unmarshal(content, &cfg)
	require.NoError(t, err, "解析会员测试配置失败: %s", configPath)

	dbName := strings.TrimSpace(os.Getenv("MEMBER_TEST_DB_NAME"))
	if dbName == "" {
		require.NotEmpty(t, cfg.Mysql.Dbname, "config.local.yaml 中缺少 mysql.db-name")
		dbName = cfg.Mysql.Dbname + "_codex_test"
	}
	cfg.Mysql.Dbname = dbName
	if cfg.Mysql.Config == "" {
		cfg.Mysql.Config = "charset=utf8mb4&parseTime=True&loc=Local"
	}

	return &cfg.Mysql, cfg.Mysql.Dsn()
}

func ensureMemberTestDatabase(t *testing.T, mysqlConfig config.Mysql) {
	t.Helper()

	adminDSN := mysqlConfig.Username + ":" + mysqlConfig.Password + "@tcp(" + mysqlConfig.Path + ":" + mysqlConfig.Port + ")/?" + mysqlConfig.Config
	sqlDB, err := sql.Open("mysql", adminDSN)
	require.NoError(t, err)
	defer func() {
		_ = sqlDB.Close()
	}()

	_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS `" + mysqlConfig.Dbname + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	require.NoError(t, err)
}

func resetMemberTestTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	require.NoError(t, db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error)
	for _, table := range []string{
		memberModel.RedemptionOrder{}.TableName(),
		memberModel.PointTransaction{}.TableName(),
		memberModel.PointAccount{}.TableName(),
		memberModel.PointProduct{}.TableName(),
		memberModel.Member{}.TableName(),
	} {
		require.NoError(t, db.Exec("TRUNCATE TABLE `"+table+"`").Error)
	}
	require.NoError(t, db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error)
}

func createTestMember(t *testing.T, req memberReq.CreateMemberReq) memberModel.Member {
	t.Helper()
	service := &MemberService{}
	require.NoError(t, service.CreateMember(req))

	var member memberModel.Member
	err := bizDB().Where("phone = ?", req.Phone).First(&member).Error
	require.NoError(t, err)
	return member
}

func createTestProduct(t *testing.T, req memberReq.CreatePointProductReq) memberModel.PointProduct {
	t.Helper()
	service := &PointProductService{}
	require.NoError(t, service.CreatePointProduct(req))

	var product memberModel.PointProduct
	err := bizDB().Where("name = ?", req.Name).First(&product).Error
	require.NoError(t, err)
	return product
}

func TestMemberServiceCreateMemberCreatesAccount(t *testing.T) {
	setupMemberTestDB(t)

	member := createTestMember(t, memberReq.CreateMemberReq{
		MemberBaseInput: memberReq.MemberBaseInput{
			Name:   "张三",
			Phone:  "13800000001",
			Gender: "female",
			Source: "douyin",
			Level:  "standard",
			Status: memberModel.MemberStatusEnabled,
			Remark: "首位会员",
		},
	})

	var account memberModel.PointAccount
	err := bizDB().Where("member_id = ?", member.ID).First(&account).Error
	require.NoError(t, err)
	require.Equal(t, int64(0), account.Balance)
	require.Equal(t, int64(0), account.TotalEarned)
	require.Equal(t, int64(0), account.TotalSpent)

	err = (&MemberService{}).CreateMember(memberReq.CreateMemberReq{
		MemberBaseInput: memberReq.MemberBaseInput{
			Name:  "李四",
			Phone: "13800000001",
		},
	})
	require.EqualError(t, err, "手机号已存在")
}

func TestPointAccountServiceGetByMemberIDRebuildsMissingAccount(t *testing.T) {
	setupMemberTestDB(t)

	member := createTestMember(t, memberReq.CreateMemberReq{
		MemberBaseInput: memberReq.MemberBaseInput{
			Name:  "王五",
			Phone: "13800000002",
		},
	})

	require.NoError(t, bizDB().Where("member_id = ?", member.ID).Delete(&memberModel.PointAccount{}).Error)

	account, err := (&PointAccountService{}).GetPointAccountByMemberID(member.ID)
	require.NoError(t, err)
	require.Equal(t, member.ID, account.MemberID)
	require.Equal(t, member.Name, account.Member.Name)

	var count int64
	err = bizDB().Model(&memberModel.PointAccount{}).Where("member_id = ?", member.ID).Count(&count).Error
	require.NoError(t, err)
	require.EqualValues(t, 1, count)
}

func TestPointAccountServiceManualAdjustWritesTransactions(t *testing.T) {
	setupMemberTestDB(t)

	member := createTestMember(t, memberReq.CreateMemberReq{
		MemberBaseInput: memberReq.MemberBaseInput{
			Name:  "赵六",
			Phone: "13800000003",
		},
	})

	accountService := &PointAccountService{}
	require.NoError(t, accountService.ManualAddPoints(memberReq.AdjustPointsReq{
		MemberID: member.ID,
		Points:   50,
		Remark:   "首充赠送",
	}, 1001))

	require.NoError(t, accountService.ManualSubPoints(memberReq.AdjustPointsReq{
		MemberID: member.ID,
		Points:   20,
		Remark:   "后台扣减",
	}, 1002))

	err := accountService.ManualSubPoints(memberReq.AdjustPointsReq{
		MemberID: member.ID,
		Points:   100,
		Remark:   "超额扣减",
	}, 1003)
	require.EqualError(t, err, "会员积分余额不足")

	account, err := accountService.GetPointAccountByMemberID(member.ID)
	require.NoError(t, err)
	require.Equal(t, int64(30), account.Balance)
	require.Equal(t, int64(50), account.TotalEarned)
	require.Equal(t, int64(20), account.TotalSpent)

	var transactions []memberModel.PointTransaction
	err = bizDB().Order("id asc").Find(&transactions).Error
	require.NoError(t, err)
	require.Len(t, transactions, 2)
	require.Equal(t, memberModel.PointTransactionTypeAdjust, transactions[0].Type)
	require.Equal(t, memberModel.PointRefTypeManualAdjustAdd, transactions[0].RefType)
	require.Equal(t, int64(0), transactions[0].BeforeBalance)
	require.Equal(t, int64(50), transactions[0].AfterBalance)
	require.Equal(t, "1001", transactions[0].Operator)
	require.Equal(t, memberModel.PointRefTypeManualAdjustSub, transactions[1].RefType)
	require.Equal(t, int64(50), transactions[1].BeforeBalance)
	require.Equal(t, int64(30), transactions[1].AfterBalance)
	require.Equal(t, "1002", transactions[1].Operator)
}

func TestRedemptionOrderServiceCreateCancelAndComplete(t *testing.T) {
	setupMemberTestDB(t)

	member := createTestMember(t, memberReq.CreateMemberReq{
		MemberBaseInput: memberReq.MemberBaseInput{
			Name:  "孙七",
			Phone: "13800000004",
		},
	})

	require.NoError(t, (&PointAccountService{}).ManualAddPoints(memberReq.AdjustPointsReq{
		MemberID: member.ID,
		Points:   200,
		Remark:   "初始化积分",
	}, 1001))

	product := createTestProduct(t, memberReq.CreatePointProductReq{
		PointProductBaseInput: memberReq.PointProductBaseInput{
			Name:        "水光护理",
			Category:    "护理",
			PointsPrice: 60,
			Stock:       5,
			Status:      memberModel.PointProductStatusOnSale,
			Sort:        1,
			Description: "基础水光护理",
		},
	})

	orderService := &RedemptionOrderService{}
	require.NoError(t, orderService.CreateRedemptionOrder(memberReq.CreateRedemptionOrderReq{
		MemberID:      member.ID,
		ProductID:     product.ID,
		Quantity:      2,
		ReceiverName:  "孙七",
		ReceiverPhone: "13800000004",
		Remark:        "首次兑换",
	}, 2001))

	var firstOrder memberModel.RedemptionOrder
	err := bizDB().Where("member_id = ?", member.ID).Order("id asc").First(&firstOrder).Error
	require.NoError(t, err)
	require.Equal(t, memberModel.RedemptionOrderStatusPending, firstOrder.Status)
	require.Equal(t, int64(120), firstOrder.TotalPoints)

	account, err := (&PointAccountService{}).GetPointAccountByMemberID(member.ID)
	require.NoError(t, err)
	require.Equal(t, int64(80), account.Balance)

	var updatedProduct memberModel.PointProduct
	err = bizDB().Where("id = ?", product.ID).First(&updatedProduct).Error
	require.NoError(t, err)
	require.Equal(t, int64(3), updatedProduct.Stock)

	require.NoError(t, orderService.CancelRedemptionOrder(memberReq.OperateRedemptionOrderReq{
		ID:     firstOrder.ID,
		Remark: "客户取消",
	}, 2002))

	account, err = (&PointAccountService{}).GetPointAccountByMemberID(member.ID)
	require.NoError(t, err)
	require.Equal(t, int64(200), account.Balance)

	err = bizDB().Where("id = ?", product.ID).First(&updatedProduct).Error
	require.NoError(t, err)
	require.Equal(t, int64(5), updatedProduct.Stock)

	err = bizDB().Where("id = ?", firstOrder.ID).First(&firstOrder).Error
	require.NoError(t, err)
	require.Equal(t, memberModel.RedemptionOrderStatusCancelled, firstOrder.Status)

	require.NoError(t, orderService.CreateRedemptionOrder(memberReq.CreateRedemptionOrderReq{
		MemberID:      member.ID,
		ProductID:     product.ID,
		Quantity:      1,
		ReceiverName:  "孙七",
		ReceiverPhone: "13800000004",
		Remark:        "再次兑换",
	}, 2003))

	var secondOrder memberModel.RedemptionOrder
	err = bizDB().Where("member_id = ?", member.ID).Order("id desc").First(&secondOrder).Error
	require.NoError(t, err)

	require.NoError(t, orderService.CompleteRedemptionOrder(memberReq.OperateRedemptionOrderReq{
		ID:     secondOrder.ID,
		Remark: "门店核销完成",
	}, 2004))

	err = bizDB().Where("id = ?", secondOrder.ID).First(&secondOrder).Error
	require.NoError(t, err)
	require.Equal(t, memberModel.RedemptionOrderStatusCompleted, secondOrder.Status)

	account, err = (&PointAccountService{}).GetPointAccountByMemberID(member.ID)
	require.NoError(t, err)
	require.Equal(t, int64(140), account.Balance)

	err = bizDB().Where("id = ?", product.ID).First(&updatedProduct).Error
	require.NoError(t, err)
	require.Equal(t, int64(4), updatedProduct.Stock)

	var transactionCount int64
	err = bizDB().Model(&memberModel.PointTransaction{}).Count(&transactionCount).Error
	require.NoError(t, err)
	require.EqualValues(t, 4, transactionCount)
}

func TestMemberServiceDeleteMemberRejectsPointHistory(t *testing.T) {
	setupMemberTestDB(t)

	member := createTestMember(t, memberReq.CreateMemberReq{
		MemberBaseInput: memberReq.MemberBaseInput{
			Name:  "周八",
			Phone: "13800000005",
		},
	})

	require.NoError(t, (&PointAccountService{}).ManualAddPoints(memberReq.AdjustPointsReq{
		MemberID: member.ID,
		Points:   20,
		Remark:   "开卡赠送",
	}, 3001))

	err := (&MemberService{}).DeleteMember(member.ID)
	require.EqualError(t, err, "该会员已有积分记录或兑换订单，不允许删除，请改为禁用。")

	var count int64
	err = bizDB().Model(&memberModel.Member{}).Where("id = ?", member.ID).Count(&count).Error
	require.NoError(t, err)
	require.EqualValues(t, 1, count)
}

func TestRedemptionOrderServiceCreateRedemptionOrderRejectsDisabledMemberAndInsufficientStock(t *testing.T) {
	setupMemberTestDB(t)

	disabledMember := createTestMember(t, memberReq.CreateMemberReq{
		MemberBaseInput: memberReq.MemberBaseInput{
			Name:   "吴九",
			Phone:  "13800000006",
			Status: memberModel.MemberStatusDisabled,
		},
	})

	enabledMember := createTestMember(t, memberReq.CreateMemberReq{
		MemberBaseInput: memberReq.MemberBaseInput{
			Name:  "郑十",
			Phone: "13800000007",
		},
	})

	require.NoError(t, (&PointAccountService{}).ManualAddPoints(memberReq.AdjustPointsReq{
		MemberID: enabledMember.ID,
		Points:   500,
		Remark:   "测试积分",
	}, 3002))

	product := createTestProduct(t, memberReq.CreatePointProductReq{
		PointProductBaseInput: memberReq.PointProductBaseInput{
			Name:        "热玛吉护理",
			Category:    "项目",
			PointsPrice: 120,
			Stock:       1,
			Status:      memberModel.PointProductStatusOnSale,
			Sort:        2,
		},
	})

	orderService := &RedemptionOrderService{}
	err := orderService.CreateRedemptionOrder(memberReq.CreateRedemptionOrderReq{
		MemberID:      disabledMember.ID,
		ProductID:     product.ID,
		Quantity:      1,
		ReceiverName:  "吴九",
		ReceiverPhone: "13800000006",
	}, 3003)
	require.EqualError(t, err, "该会员已禁用，不能创建兑换订单。")

	err = orderService.CreateRedemptionOrder(memberReq.CreateRedemptionOrderReq{
		MemberID:      enabledMember.ID,
		ProductID:     product.ID,
		Quantity:      2,
		ReceiverName:  "郑十",
		ReceiverPhone: "13800000007",
	}, 3004)
	require.EqualError(t, err, "商品库存不足，无法兑换")
}
