package member

import (
	"strings"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"gorm.io/gorm"
)

type PointAccountService struct{}

func (s *PointAccountService) GetPointAccountByMemberID(memberID uint) (account memberModel.PointAccount, err error) {
	err = bizDB().Preload("Member").Where("member_id = ?", memberID).First(&account).Error
	return
}

func (s *PointAccountService) GetPointAccountList(info memberReq.PointAccountSearch) (list []memberModel.PointAccount, total int64, err error) {
	db := bizDB().Model(&memberModel.PointAccount{}).
		Joins("LEFT JOIN " + memberModel.Member{}.TableName() + " ON " + memberModel.Member{}.TableName() + ".id = " + memberModel.PointAccount{}.TableName() + ".member_id")
	if info.MemberID > 0 {
		db = db.Where(memberModel.PointAccount{}.TableName()+".member_id = ?", info.MemberID)
	}
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where(memberModel.Member{}.TableName()+".mobile LIKE ? OR "+memberModel.Member{}.TableName()+".nickname LIKE ? OR "+memberModel.Member{}.TableName()+".real_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Preload("Member").Scopes(info.Paginate()).Order(memberModel.PointAccount{}.TableName() + ".created_at desc").Find(&list).Error
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
