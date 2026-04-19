package member

import (
	"errors"
	"strings"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"gorm.io/gorm"
)

type PointAccountService struct{}

func (s *PointAccountService) GetPointAccountByMemberID(memberID uint) (account memberModel.PointAccount, err error) {
	err = bizDB().Transaction(func(tx *gorm.DB) error {
		member, loadErr := loadMember(tx, memberID)
		if loadErr != nil {
			return loadErr
		}

		queryErr := tx.Where("member_id = ?", memberID).First(&account).Error
		if queryErr == nil {
			account.Member = member
			return nil
		}
		if !errors.Is(queryErr, gorm.ErrRecordNotFound) {
			return queryErr
		}

		if createErr := tx.Create(&memberModel.PointAccount{MemberID: memberID}).Error; createErr != nil {
			queryErr = tx.Where("member_id = ?", memberID).First(&account).Error
			if queryErr == nil {
				account.Member = member
				return nil
			}
			return createErr
		}

		if queryErr = tx.Where("member_id = ?", memberID).First(&account).Error; queryErr != nil {
			return queryErr
		}
		account.Member = member
		return nil
	})
	return
}

func (s *PointAccountService) GetPointAccountList(info memberReq.PointAccountSearch) (list []memberModel.PointAccount, total int64, err error) {
	db := bizDB().Model(&memberModel.Member{})
	if info.MemberID > 0 {
		db = db.Where(memberModel.Member{}.TableName()+".id = ?", info.MemberID)
	}
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where(memberModel.Member{}.TableName()+".name LIKE ? OR "+memberModel.Member{}.TableName()+".phone LIKE ? OR "+memberModel.Member{}.TableName()+".source LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err = db.Count(&total).Error; err != nil {
		return
	}

	var members []memberModel.Member
	if err = db.Scopes(info.Paginate()).Order(memberModel.Member{}.TableName() + ".created_at desc").Find(&members).Error; err != nil {
		return
	}
	if len(members) == 0 {
		return
	}

	memberIDs := make([]uint, 0, len(members))
	for _, item := range members {
		memberIDs = append(memberIDs, item.ID)
	}

	var accounts []memberModel.PointAccount
	if err = bizDB().Where("member_id IN ?", memberIDs).Find(&accounts).Error; err != nil {
		return
	}

	accountMap := make(map[uint]memberModel.PointAccount, len(accounts))
	for _, item := range accounts {
		accountMap[item.MemberID] = item
	}

	list = make([]memberModel.PointAccount, 0, len(members))
	for _, item := range members {
		account, ok := accountMap[item.ID]
		if !ok {
			account = memberModel.PointAccount{MemberID: item.ID}
		}
		account.Member = item
		list = append(list, account)
	}
	return
}

func (s *PointAccountService) ManualAddPoints(req memberReq.AdjustPointsReq, operatorID uint) error {
	return s.manualAdjust(req, operatorID, pointActionAdjustAdd, memberModel.PointRefTypeManualAdjustAdd)
}

func (s *PointAccountService) ManualSubPoints(req memberReq.AdjustPointsReq, operatorID uint) error {
	return s.manualAdjust(req, operatorID, pointActionAdjustSub, memberModel.PointRefTypeManualAdjustSub)
}

func (s *PointAccountService) manualAdjust(req memberReq.AdjustPointsReq, operatorID uint, action string, refType string) error {
	return bizDB().Transaction(func(tx *gorm.DB) error {
		if _, err := loadMember(tx, req.MemberID); err != nil {
			return err
		}
		change, err := applyPointChange(tx, req.MemberID, action, req.Points)
		if err != nil {
			return err
		}
		return recordPointTransaction(tx, req.MemberID, change.Account.ID, action, req.Points, change.BeforeBalance, change.AfterBalance, refType, 0, formatOperator(operatorID), req.Remark)
	})
}
