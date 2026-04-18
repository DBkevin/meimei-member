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

		createErr := tx.Create(&memberModel.PointAccount{MemberID: memberID}).Error
		if createErr != nil {
			retryErr := tx.Where("member_id = ?", memberID).First(&account).Error
			if retryErr == nil {
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
	return account, nil
}

func (s *PointAccountService) GetPointAccountList(info memberReq.PointAccountSearch) (list []memberModel.PointAccount, total int64, err error) {
	db := bizDB().Model(&memberModel.Member{})
	if info.MemberID > 0 {
		db = db.Where(memberModel.Member{}.TableName()+".id = ?", info.MemberID)
	}
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where(memberModel.Member{}.TableName()+".mobile LIKE ? OR "+memberModel.Member{}.TableName()+".nickname LIKE ? OR "+memberModel.Member{}.TableName()+".real_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var members []memberModel.Member
	err = db.Scopes(info.Paginate()).Order(memberModel.Member{}.TableName() + ".created_at desc").Find(&members).Error
	if err != nil || len(members) == 0 {
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
			account = s.buildDefaultAccount(item)
		} else {
			account.Member = item
		}
		list = append(list, account)
	}
	return
}

func (s *PointAccountService) ManualAddPoints(req memberReq.AdjustPointsReq, operatorID uint) error {
	return bizDB().Transaction(func(tx *gorm.DB) error {
		if _, err := loadMember(tx, req.MemberID); err != nil {
			return err
		}
		_, err := adjustPoints(tx, req.MemberID, memberModel.PointChangeTypeAdjustAdd, req.Points, memberModel.PointSourceTypeManual, 0, req.Remark, operatorID)
		return err
	})
}

func (s *PointAccountService) ManualSubPoints(req memberReq.AdjustPointsReq, operatorID uint) error {
	return bizDB().Transaction(func(tx *gorm.DB) error {
		if _, err := loadMember(tx, req.MemberID); err != nil {
			return err
		}
		_, err := adjustPoints(tx, req.MemberID, memberModel.PointChangeTypeAdjustSub, req.Points, memberModel.PointSourceTypeManual, 0, req.Remark, operatorID)
		return err
	})
}

func (s *PointAccountService) buildDefaultAccount(member memberModel.Member) memberModel.PointAccount {
	return memberModel.PointAccount{
		BaseModel: memberModel.BaseModel{
			CreatedAt: member.CreatedAt,
			UpdatedAt: member.UpdatedAt,
		},
		MemberID: member.ID,
		Member:   member,
	}
}
