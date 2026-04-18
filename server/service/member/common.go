package member

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

func buildOrderNo() string {
	return fmt.Sprintf("ME%s%04d", time.Now().Format("20060102150405"), randomNumber(10000))
}

func buildVerifyCode() string {
	return fmt.Sprintf("%06d", randomNumber(1000000))
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

func getOrCreateAccountForUpdate(tx *gorm.DB, memberID uint) (memberModel.PointAccount, error) {
	var account memberModel.PointAccount
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("member_id = ?", memberID).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		account = memberModel.PointAccount{MemberID: memberID}
		if err = tx.Create(&account).Error; err != nil {
			return account, err
		}
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", account.ID).First(&account).Error
	}
	return account, err
}

func adjustPoints(tx *gorm.DB, memberID uint, changeType string, points int64, sourceType string, sourceID uint, remark string, operatorID uint) (memberModel.PointAccount, error) {
	var account memberModel.PointAccount
	if points <= 0 {
		return account, errors.New("积分必须大于0")
	}
	account, err := getOrCreateAccountForUpdate(tx, memberID)
	if err != nil {
		return account, err
	}

	before := account.AvailablePoints
	after := before
	earned := account.TotalEarnedPoints
	used := account.TotalUsedPoints

	switch changeType {
	case memberModel.PointChangeTypeEarn, memberModel.PointChangeTypeAdjustAdd:
		after += points
		earned += points
	case memberModel.PointChangeTypeUse, memberModel.PointChangeTypeAdjustSub:
		if before < points {
			return account, errors.New("会员可用积分不足")
		}
		after -= points
		used += points
	case memberModel.PointChangeTypeRefund:
		after += points
		if used >= points {
			used -= points
		} else {
			used = 0
		}
	default:
		return account, errors.New("不支持的积分变动类型")
	}

	updates := map[string]interface{}{
		"available_points":    after,
		"total_earned_points": earned,
		"total_used_points":   used,
	}
	if err = tx.Model(&memberModel.PointAccount{}).Where("id = ?", account.ID).Updates(updates).Error; err != nil {
		return account, err
	}
	account.AvailablePoints = after
	account.TotalEarnedPoints = earned
	account.TotalUsedPoints = used

	pointLog := memberModel.PointLog{
		MemberID:     memberID,
		ChangeType:   changeType,
		ChangePoints: points,
		BeforePoints: before,
		AfterPoints:  after,
		SourceType:   sourceType,
		SourceID:     sourceID,
		Remark:       remark,
		OperatorID:   operatorID,
	}
	if err = tx.Create(&pointLog).Error; err != nil {
		return account, err
	}

	return account, nil
}

func bizDB() *gorm.DB {
	return global.GVA_DB
}
