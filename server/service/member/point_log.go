package member

import (
	"strings"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
)

type PointTransactionService struct{}

func (s *PointTransactionService) GetPointTransactionList(info memberReq.PointTransactionSearch) (list []memberModel.PointTransaction, total int64, err error) {
	db := bizDB().Model(&memberModel.PointTransaction{}).
		Joins("LEFT JOIN " + memberModel.Member{}.TableName() + " ON " + memberModel.Member{}.TableName() + ".id = " + memberModel.PointTransaction{}.TableName() + ".member_id")

	if info.MemberID > 0 {
		db = db.Where(memberModel.PointTransaction{}.TableName()+".member_id = ?", info.MemberID)
	}
	if info.Type != "" {
		db = db.Where(memberModel.PointTransaction{}.TableName()+".type = ?", info.Type)
	}
	if info.RefType != "" {
		db = db.Where(memberModel.PointTransaction{}.TableName()+".ref_type = ?", info.RefType)
	}
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where(memberModel.Member{}.TableName()+".name LIKE ? OR "+memberModel.Member{}.TableName()+".phone LIKE ? OR "+memberModel.PointTransaction{}.TableName()+".remark LIKE ? OR "+memberModel.PointTransaction{}.TableName()+".operator LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err = db.Count(&total).Error; err != nil {
		return
	}

	err = db.Preload("Member").Preload("Account").Scopes(info.Paginate()).Order(memberModel.PointTransaction{}.TableName() + ".created_at desc").Find(&list).Error
	return
}
