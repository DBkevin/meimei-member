package member

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const legacyDefaultProductCategory = "默认分类"

type LegacyMigrationSummary struct {
	Members      int
	Accounts     int
	Products     int
	Orders       int
	Transactions int
}

type legacyMember struct {
	ID          uint           `gorm:"column:id"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
	OpenID      *string        `gorm:"column:openid"`
	UnionID     *string        `gorm:"column:unionid"`
	Mobile      *string        `gorm:"column:mobile"`
	Nickname    *string        `gorm:"column:nickname"`
	AvatarURL   *string        `gorm:"column:avatar_url"`
	RealName    *string        `gorm:"column:real_name"`
	MemberLevel *string        `gorm:"column:member_level"`
	Status      *string        `gorm:"column:status"`
}

func (legacyMember) TableName() string {
	return "member_members"
}

type legacyPointAccount struct {
	ID                uint           `gorm:"column:id"`
	CreatedAt         time.Time      `gorm:"column:created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at"`
	MemberID          uint           `gorm:"column:member_id"`
	AvailablePoints   int64          `gorm:"column:available_points"`
	FrozenPoints      int64          `gorm:"column:frozen_points"`
	TotalEarnedPoints int64          `gorm:"column:total_earned_points"`
	TotalUsedPoints   int64          `gorm:"column:total_used_points"`
}

func (legacyPointAccount) TableName() string {
	return "member_point_accounts"
}

type legacyPointLog struct {
	ID           uint           `gorm:"column:id"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
	MemberID     uint           `gorm:"column:member_id"`
	ChangeType   *string        `gorm:"column:change_type"`
	ChangePoints int64          `gorm:"column:change_points"`
	BeforePoints int64          `gorm:"column:before_points"`
	AfterPoints  int64          `gorm:"column:after_points"`
	SourceType   *string        `gorm:"column:source_type"`
	SourceID     uint           `gorm:"column:source_id"`
	Remark       *string        `gorm:"column:remark"`
	OperatorID   uint           `gorm:"column:operator_id"`
}

func (legacyPointLog) TableName() string {
	return "member_point_logs"
}

type legacyPointProduct struct {
	ID             uint           `gorm:"column:id"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
	Name           *string        `gorm:"column:name"`
	CoverImage     *string        `gorm:"column:cover_image"`
	Description    *string        `gorm:"column:description"`
	PointsPrice    int64          `gorm:"column:points_price"`
	Stock          int64          `gorm:"column:stock"`
	LimitPerMember int64          `gorm:"column:limit_per_member"`
	Status         *string        `gorm:"column:status"`
	Sort           int64          `gorm:"column:sort"`
}

func (legacyPointProduct) TableName() string {
	return "member_point_goods"
}

type legacyExchangeOrder struct {
	ID         uint           `gorm:"column:id"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
	OrderNo    *string        `gorm:"column:order_no"`
	MemberID   uint           `gorm:"column:member_id"`
	GoodsID    uint           `gorm:"column:goods_id"`
	PointsCost int64          `gorm:"column:points_cost"`
	Status     *string        `gorm:"column:status"`
	VerifyCode *string        `gorm:"column:verify_code"`
	VerifiedAt *time.Time     `gorm:"column:verified_at"`
	OperatorID uint           `gorm:"column:operator_id"`
	Remark     *string        `gorm:"column:remark"`
}

func (legacyExchangeOrder) TableName() string {
	return "member_exchange_orders"
}

func MigrateLegacyMemberData(db *gorm.DB) (summary LegacyMigrationSummary, err error) {
	if db == nil {
		return summary, fmt.Errorf("database is nil")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			&memberModel.Member{},
			&memberModel.PointAccount{},
			&memberModel.PointTransaction{},
			&memberModel.PointProduct{},
			&memberModel.RedemptionOrder{},
		); err != nil {
			return err
		}

		if !hasAnyLegacyMemberTable(tx) {
			return nil
		}

		memberMap, migratedMembers, err := migrateLegacyMembers(tx)
		if err != nil {
			return err
		}
		summary.Members = migratedMembers

		productMap, migratedProducts, err := migrateLegacyProducts(tx)
		if err != nil {
			return err
		}
		summary.Products = migratedProducts

		migratedAccounts, err := migrateLegacyAccounts(tx)
		if err != nil {
			return err
		}
		summary.Accounts = migratedAccounts

		if err := ensureMemberAccounts(tx, memberMap); err != nil {
			return err
		}

		accountMap, err := loadPointAccountMap(tx)
		if err != nil {
			return err
		}

		migratedOrders, err := migrateLegacyOrders(tx, memberMap, productMap)
		if err != nil {
			return err
		}
		summary.Orders = migratedOrders

		migratedTransactions, err := migrateLegacyTransactions(tx, accountMap)
		if err != nil {
			return err
		}
		summary.Transactions = migratedTransactions

		return nil
	})
	return summary, err
}

func hasAnyLegacyMemberTable(tx *gorm.DB) bool {
	for _, table := range []string{
		legacyMember{}.TableName(),
		legacyPointAccount{}.TableName(),
		legacyPointLog{}.TableName(),
		legacyPointProduct{}.TableName(),
		legacyExchangeOrder{}.TableName(),
	} {
		if tx.Migrator().HasTable(table) {
			return true
		}
	}
	return false
}

func migrateLegacyMembers(tx *gorm.DB) (map[uint]memberModel.Member, int, error) {
	if !tx.Migrator().HasTable(legacyMember{}.TableName()) {
		return map[uint]memberModel.Member{}, 0, nil
	}

	var legacyMembers []legacyMember
	if err := tx.Unscoped().Order("id asc").Find(&legacyMembers).Error; err != nil {
		return nil, 0, err
	}
	if len(legacyMembers) == 0 {
		return map[uint]memberModel.Member{}, 0, nil
	}

	members := make([]memberModel.Member, 0, len(legacyMembers))
	for _, item := range legacyMembers {
		members = append(members, memberModel.Member{
			BaseModel: memberModel.BaseModel{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: item.DeletedAt,
			},
			Name:   buildLegacyMemberName(item),
			Phone:  cleanLegacyString(item.Mobile),
			Level:  normalizeLegacyMemberLevel(item.MemberLevel),
			Status: normalizeLegacyMemberStatus(item.Status),
			Remark: buildLegacyMemberRemark(item),
		})
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"name",
			"phone",
			"gender",
			"birthday",
			"source",
			"level",
			"status",
			"remark",
		}),
	}).Create(&members).Error; err != nil {
		return nil, 0, err
	}

	memberMap := make(map[uint]memberModel.Member, len(members))
	for _, item := range members {
		memberMap[item.ID] = item
	}
	return memberMap, len(members), nil
}

func migrateLegacyAccounts(tx *gorm.DB) (int, error) {
	if !tx.Migrator().HasTable(legacyPointAccount{}.TableName()) {
		return 0, nil
	}

	var legacyAccounts []legacyPointAccount
	if err := tx.Unscoped().Order("id asc").Find(&legacyAccounts).Error; err != nil {
		return 0, err
	}
	if len(legacyAccounts) == 0 {
		return 0, nil
	}

	accounts := make([]memberModel.PointAccount, 0, len(legacyAccounts))
	for _, item := range legacyAccounts {
		accounts = append(accounts, memberModel.PointAccount{
			BaseModel: memberModel.BaseModel{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: item.DeletedAt,
			},
			MemberID:     item.MemberID,
			Balance:      item.AvailablePoints,
			TotalEarned:  item.TotalEarnedPoints,
			TotalSpent:   item.TotalUsedPoints,
			FrozenPoints: item.FrozenPoints,
		})
	}

	err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "member_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"balance",
			"total_earned",
			"total_spent",
			"frozen_points",
		}),
	}).Create(&accounts).Error
	return len(accounts), err
}

func ensureMemberAccounts(tx *gorm.DB, memberMap map[uint]memberModel.Member) error {
	if len(memberMap) == 0 {
		return nil
	}

	memberIDs := make([]uint, 0, len(memberMap))
	for memberID := range memberMap {
		memberIDs = append(memberIDs, memberID)
	}

	var existing []memberModel.PointAccount
	if err := tx.Select("member_id").Where("member_id IN ?", memberIDs).Find(&existing).Error; err != nil {
		return err
	}

	existingMap := make(map[uint]struct{}, len(existing))
	for _, item := range existing {
		existingMap[item.MemberID] = struct{}{}
	}

	missing := make([]memberModel.PointAccount, 0)
	for memberID := range memberMap {
		if _, ok := existingMap[memberID]; ok {
			continue
		}
		missing = append(missing, memberModel.PointAccount{
			MemberID: memberID,
		})
	}

	if len(missing) == 0 {
		return nil
	}
	return tx.Create(&missing).Error
}

func loadPointAccountMap(tx *gorm.DB) (map[uint]memberModel.PointAccount, error) {
	var accounts []memberModel.PointAccount
	if err := tx.Unscoped().Find(&accounts).Error; err != nil {
		return nil, err
	}

	accountMap := make(map[uint]memberModel.PointAccount, len(accounts))
	for _, item := range accounts {
		accountMap[item.MemberID] = item
	}
	return accountMap, nil
}

func migrateLegacyProducts(tx *gorm.DB) (map[uint]memberModel.PointProduct, int, error) {
	if !tx.Migrator().HasTable(legacyPointProduct{}.TableName()) {
		return map[uint]memberModel.PointProduct{}, 0, nil
	}

	var legacyProducts []legacyPointProduct
	if err := tx.Unscoped().Order("id asc").Find(&legacyProducts).Error; err != nil {
		return nil, 0, err
	}
	if len(legacyProducts) == 0 {
		return map[uint]memberModel.PointProduct{}, 0, nil
	}

	products := make([]memberModel.PointProduct, 0, len(legacyProducts))
	for _, item := range legacyProducts {
		products = append(products, memberModel.PointProduct{
			BaseModel: memberModel.BaseModel{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: item.DeletedAt,
			},
			Name:        cleanLegacyString(item.Name),
			CoverURL:    cleanLegacyString(item.CoverImage),
			Category:    legacyDefaultProductCategory,
			PointsPrice: item.PointsPrice,
			Stock:       item.Stock,
			Status:      normalizeLegacyProductStatus(item.Status),
			Sort:        int(item.Sort),
			Description: buildLegacyProductDescription(item),
		})
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"name",
			"cover_url",
			"category",
			"points_price",
			"stock",
			"status",
			"sort",
			"description",
		}),
	}).Create(&products).Error; err != nil {
		return nil, 0, err
	}

	productMap := make(map[uint]memberModel.PointProduct, len(products))
	for _, item := range products {
		productMap[item.ID] = item
	}
	return productMap, len(products), nil
}

func migrateLegacyOrders(tx *gorm.DB, memberMap map[uint]memberModel.Member, productMap map[uint]memberModel.PointProduct) (int, error) {
	if !tx.Migrator().HasTable(legacyExchangeOrder{}.TableName()) {
		return 0, nil
	}

	var legacyOrders []legacyExchangeOrder
	if err := tx.Unscoped().Order("id asc").Find(&legacyOrders).Error; err != nil {
		return 0, err
	}
	if len(legacyOrders) == 0 {
		return 0, nil
	}

	orders := make([]memberModel.RedemptionOrder, 0, len(legacyOrders))
	for _, item := range legacyOrders {
		member := memberMap[item.MemberID]
		product := productMap[item.GoodsID]
		orders = append(orders, memberModel.RedemptionOrder{
			BaseModel: memberModel.BaseModel{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: item.DeletedAt,
			},
			OrderNo:       fallbackLegacyOrderNo(item.OrderNo, item.ID),
			MemberID:      item.MemberID,
			ProductID:     item.GoodsID,
			ProductName:   product.Name,
			Quantity:      1,
			UnitPoints:    item.PointsCost,
			TotalPoints:   item.PointsCost,
			Status:        normalizeLegacyOrderStatus(item.Status),
			ReceiverName:  member.Name,
			ReceiverPhone: member.Phone,
			Remark:        buildLegacyOrderRemark(item),
		})
	}

	err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"order_no",
			"member_id",
			"product_id",
			"product_name",
			"quantity",
			"unit_points",
			"total_points",
			"status",
			"receiver_name",
			"receiver_phone",
			"remark",
		}),
	}).Create(&orders).Error
	return len(orders), err
}

func migrateLegacyTransactions(tx *gorm.DB, accountMap map[uint]memberModel.PointAccount) (int, error) {
	if !tx.Migrator().HasTable(legacyPointLog{}.TableName()) {
		return 0, nil
	}

	var legacyLogs []legacyPointLog
	if err := tx.Unscoped().Order("id asc").Find(&legacyLogs).Error; err != nil {
		return 0, err
	}
	if len(legacyLogs) == 0 {
		return 0, nil
	}

	transactions := make([]memberModel.PointTransaction, 0, len(legacyLogs))
	for _, item := range legacyLogs {
		account, ok := accountMap[item.MemberID]
		if !ok {
			return 0, fmt.Errorf("missing point account for legacy member %d", item.MemberID)
		}

		transactionType, refType := normalizeLegacyPointLogType(item)
		transactions = append(transactions, memberModel.PointTransaction{
			BaseModel: memberModel.BaseModel{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: item.DeletedAt,
			},
			MemberID:      item.MemberID,
			AccountID:     account.ID,
			Type:          transactionType,
			Points:        absolutePoints(item.ChangePoints),
			BeforeBalance: item.BeforePoints,
			AfterBalance:  item.AfterPoints,
			RefType:       refType,
			RefID:         item.SourceID,
			Operator:      formatOperator(item.OperatorID),
			Remark:        cleanLegacyString(item.Remark),
		})
	}

	err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"member_id",
			"account_id",
			"type",
			"points",
			"before_balance",
			"after_balance",
			"ref_type",
			"ref_id",
			"operator",
			"remark",
		}),
	}).Create(&transactions).Error
	return len(transactions), err
}

func buildLegacyMemberName(item legacyMember) string {
	for _, value := range []*string{item.RealName, item.Nickname, item.Mobile} {
		if cleaned := cleanLegacyString(value); cleaned != "" {
			return cleaned
		}
	}
	return fmt.Sprintf("会员%d", item.ID)
}

func buildLegacyMemberRemark(item legacyMember) string {
	parts := make([]string, 0, 4)
	if nickname := cleanLegacyString(item.Nickname); nickname != "" && nickname != cleanLegacyString(item.RealName) {
		parts = append(parts, "legacy_nickname="+nickname)
	}
	if openID := cleanLegacyString(item.OpenID); openID != "" {
		parts = append(parts, "legacy_openid="+openID)
	}
	if unionID := cleanLegacyString(item.UnionID); unionID != "" {
		parts = append(parts, "legacy_unionid="+unionID)
	}
	if avatarURL := cleanLegacyString(item.AvatarURL); avatarURL != "" {
		parts = append(parts, "legacy_avatar_url="+avatarURL)
	}

	remark := strings.Join(parts, "; ")
	if len(remark) > 500 {
		return remark[:500]
	}
	return remark
}

func buildLegacyProductDescription(item legacyPointProduct) string {
	description := cleanLegacyString(item.Description)
	if item.LimitPerMember <= 0 {
		return description
	}

	legacyLimitRemark := fmt.Sprintf("[legacy] 每人限兑数量: %d", item.LimitPerMember)
	if description == "" {
		return legacyLimitRemark
	}
	return description + "\n\n" + legacyLimitRemark
}

func buildLegacyOrderRemark(item legacyExchangeOrder) string {
	parts := make([]string, 0, 3)
	if remark := cleanLegacyString(item.Remark); remark != "" {
		parts = append(parts, remark)
	}
	if verifyCode := cleanLegacyString(item.VerifyCode); verifyCode != "" {
		parts = append(parts, "legacy_verify_code="+verifyCode)
	}
	if item.VerifiedAt != nil && !item.VerifiedAt.IsZero() {
		parts = append(parts, "legacy_verified_at="+item.VerifiedAt.Format(time.RFC3339))
	}
	return strings.Join(parts, "; ")
}

func fallbackLegacyOrderNo(orderNo *string, orderID uint) string {
	if value := cleanLegacyString(orderNo); value != "" {
		return value
	}
	return fmt.Sprintf("LEGACY%06d", orderID)
}

func normalizeLegacyMemberLevel(level *string) string {
	if value := cleanLegacyString(level); value != "" {
		return value
	}
	return memberModel.MemberLevelStandard
}

func normalizeLegacyMemberStatus(status *string) int {
	switch strings.ToLower(cleanLegacyString(status)) {
	case "", "enabled", "enable", "active":
		return memberModel.MemberStatusEnabled
	default:
		return memberModel.MemberStatusDisabled
	}
}

func normalizeLegacyProductStatus(status *string) int {
	switch strings.ToLower(cleanLegacyString(status)) {
	case "on_sale", "on-sale", "enabled", "active":
		return memberModel.PointProductStatusOnSale
	default:
		return memberModel.PointProductStatusOffSale
	}
}

func normalizeLegacyOrderStatus(status *string) int {
	switch strings.ToLower(cleanLegacyString(status)) {
	case "completed", "verified", "used", "success":
		return memberModel.RedemptionOrderStatusCompleted
	case "cancelled", "canceled", "refunded", "closed", "void":
		return memberModel.RedemptionOrderStatusCancelled
	default:
		return memberModel.RedemptionOrderStatusPending
	}
}

func normalizeLegacyPointLogType(item legacyPointLog) (string, string) {
	changeType := strings.ToLower(cleanLegacyString(item.ChangeType))
	sourceType := strings.ToLower(cleanLegacyString(item.SourceType))

	switch changeType {
	case "adjust_add":
		return memberModel.PointTransactionTypeAdjust, memberModel.PointRefTypeManualAdjustAdd
	case "adjust_sub":
		return memberModel.PointTransactionTypeAdjust, memberModel.PointRefTypeManualAdjustSub
	case "use", "spend", "consume":
		return memberModel.PointTransactionTypeSpend, memberModel.PointRefTypeRedemptionOrder
	case "refund":
		return memberModel.PointTransactionTypeRefund, memberModel.PointRefTypeRedemptionOrderVoid
	}

	switch sourceType {
	case "exchange_order":
		return memberModel.PointTransactionTypeSpend, memberModel.PointRefTypeRedemptionOrder
	case "exchange_order_void":
		return memberModel.PointTransactionTypeRefund, memberModel.PointRefTypeRedemptionOrderVoid
	case "manual":
		if item.ChangePoints < 0 {
			return memberModel.PointTransactionTypeAdjust, memberModel.PointRefTypeManualAdjustSub
		}
		return memberModel.PointTransactionTypeAdjust, memberModel.PointRefTypeManualAdjustAdd
	}

	if item.ChangePoints < 0 {
		return memberModel.PointTransactionTypeSpend, memberModel.PointRefTypeRedemptionOrder
	}
	return memberModel.PointTransactionTypeEarn, ""
}

func cleanLegacyString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func absolutePoints(points int64) int64 {
	if points < 0 {
		return -points
	}
	return points
}

func LoadLegacyMigrationSummaryText(summary LegacyMigrationSummary) string {
	return strings.Join([]string{
		"members=" + strconv.Itoa(summary.Members),
		"accounts=" + strconv.Itoa(summary.Accounts),
		"products=" + strconv.Itoa(summary.Products),
		"orders=" + strconv.Itoa(summary.Orders),
		"transactions=" + strconv.Itoa(summary.Transactions),
	}, ", ")
}
