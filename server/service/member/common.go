package member

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	pointActionEarn      = "earn"
	pointActionSpend     = "spend"
	pointActionAdjustAdd = "adjust_add"
	pointActionAdjustSub = "adjust_sub"
	pointActionRefund    = "refund"
)

type pointChangeResult struct {
	Account       memberModel.PointAccount
	BeforeBalance int64
	AfterBalance  int64
}

func mergeRemark(origin, extra string) string {
	origin = strings.TrimSpace(origin)
	extra = strings.TrimSpace(extra)
	switch {
	case origin == "":
		return extra
	case extra == "":
		return origin
	default:
		return origin + "；" + extra
	}
}

func formatOperator(operatorID uint) string {
	if operatorID == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(operatorID), 10)
}

func buildOrderNo() string {
	return fmt.Sprintf("ME%s%04d", time.Now().Format("20060102150405"), randomNumber(10000))
}

func randomNumber(max int64) int64 {
	if max <= 0 {
		return 0
	}
	var n uint32
	if err := binaryReadRandom(&n); err != nil {
		return time.Now().UnixNano() % max
	}
	return int64(n) % max
}

func binaryReadRandom(target *uint32) error {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	*target = uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	return nil
}

func loadMember(tx *gorm.DB, memberID uint) (memberModel.Member, error) {
	var member memberModel.Member
	if err := tx.Where("id = ?", memberID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return member, errors.New("会员不存在")
		}
		return member, err
	}
	return member, nil
}

func loadPointAccountByMemberID(tx *gorm.DB, memberID uint, lock bool, unscoped bool) (memberModel.PointAccount, error) {
	var account memberModel.PointAccount
	query := tx
	if unscoped {
		query = query.Unscoped()
	}
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := query.Where("member_id = ?", memberID).First(&account).Error
	return account, err
}

func restoreDeletedPointAccount(tx *gorm.DB, memberID uint, lock bool) (memberModel.PointAccount, error) {
	account, err := loadPointAccountByMemberID(tx, memberID, lock, true)
	if err != nil {
		return account, err
	}
	if !account.DeletedAt.Valid {
		return account, nil
	}

	updateResult := tx.Unscoped().
		Model(&memberModel.PointAccount{}).
		Where("id = ?", account.ID).
		Updates(map[string]interface{}{"deleted_at": nil})
	if updateResult.Error != nil {
		return account, updateResult.Error
	}
	if updateResult.RowsAffected != 1 {
		return account, errors.New("恢复积分账户失败")
	}

	return loadPointAccountByMemberID(tx, memberID, lock, false)
}

func getOrCreateAccount(tx *gorm.DB, memberID uint, lock bool) (memberModel.PointAccount, error) {
	account, err := loadPointAccountByMemberID(tx, memberID, lock, false)
	if err == nil {
		return account, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return account, err
	}

	account, err = restoreDeletedPointAccount(tx, memberID, lock)
	if err == nil {
		return account, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return account, err
	}

	account = memberModel.PointAccount{MemberID: memberID}
	if err = tx.Create(&account).Error; err != nil {
		if restored, restoreErr := restoreDeletedPointAccount(tx, memberID, lock); restoreErr == nil {
			return restored, nil
		}
		if retryAccount, retryErr := loadPointAccountByMemberID(tx, memberID, lock, false); retryErr == nil {
			return retryAccount, nil
		}
		return account, err
	}

	return loadPointAccountByMemberID(tx, memberID, lock, false)
}

func getOrCreateAccountForUpdate(tx *gorm.DB, memberID uint) (memberModel.PointAccount, error) {
	return getOrCreateAccount(tx, memberID, true)
}

func applyPointChange(tx *gorm.DB, memberID uint, action string, points int64) (pointChangeResult, error) {
	var result pointChangeResult
	if points <= 0 {
		return result, errors.New("积分必须大于0")
	}
	account, err := getOrCreateAccountForUpdate(tx, memberID)
	if err != nil {
		return result, err
	}

	before := account.Balance
	after := before
	earned := account.TotalEarned
	spent := account.TotalSpent

	switch action {
	case pointActionEarn, pointActionAdjustAdd:
		after += points
		earned += points
	case pointActionSpend, pointActionAdjustSub:
		if before < points {
			return result, errors.New("会员积分余额不足")
		}
		after -= points
		spent += points
	case pointActionRefund:
		after += points
	default:
		return result, errors.New("不支持的积分变动类型")
	}

	updates := map[string]interface{}{
		"balance":       after,
		"total_earned":  earned,
		"total_spent":   spent,
		"frozen_points": account.FrozenPoints,
	}
	updateResult := tx.Model(&memberModel.PointAccount{}).Where("id = ?", account.ID).Updates(updates)
	if updateResult.Error != nil {
		return result, updateResult.Error
	}
	if updateResult.RowsAffected != 1 {
		return result, errors.New("更新积分账户失败")
	}
	account.Balance = after
	account.TotalEarned = earned
	account.TotalSpent = spent
	result = pointChangeResult{
		Account:       account,
		BeforeBalance: before,
		AfterBalance:  after,
	}
	return result, nil
}

func actionToTransactionType(action string) string {
	switch action {
	case pointActionEarn:
		return memberModel.PointTransactionTypeEarn
	case pointActionSpend:
		return memberModel.PointTransactionTypeSpend
	case pointActionAdjustAdd, pointActionAdjustSub:
		return memberModel.PointTransactionTypeAdjust
	case pointActionRefund:
		return memberModel.PointTransactionTypeRefund
	default:
		return memberModel.PointTransactionTypeAdjust
	}
}

func recordPointTransaction(tx *gorm.DB, memberID uint, accountID uint, action string, points int64, beforeBalance int64, afterBalance int64, refType string, refID uint, operator string, remark string) error {
	transaction := memberModel.PointTransaction{
		MemberID:      memberID,
		AccountID:     accountID,
		Type:          actionToTransactionType(action),
		Points:        points,
		BeforeBalance: beforeBalance,
		AfterBalance:  afterBalance,
		RefType:       refType,
		RefID:         refID,
		Operator:      strings.TrimSpace(operator),
		Remark:        strings.TrimSpace(remark),
	}
	return tx.Create(&transaction).Error
}

func bizDB() *gorm.DB {
	return global.GVA_DB
}
