package member

import (
	"strings"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
)

type PointLogService struct{}

func (s *PointLogService) GetPointLogList(info memberReq.PointLogSearch) (list []memberModel.PointLog, total int64, err error) {
	db := bizDB().Model(&memberModel.PointLog{}).
		Joins("LEFT JOIN " + memberModel.Member{}.TableName() + " ON " + memberModel.Member{}.TableName() + ".id = " + memberModel.PointLog{}.TableName() + ".member_id")

	if info.MemberID > 0 {
		db = db.Where(memberModel.PointLog{}.TableName()+".member_id = ?", info.MemberID)
	}
	if info.ChangeType != "" {
		db = db.Where(memberModel.PointLog{}.TableName()+".change_type = ?", info.ChangeType)
	}
	if info.SourceType != "" {
		db = db.Where(memberModel.PointLog{}.TableName()+".source_type = ?", info.SourceType)
	}
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where(memberModel.Member{}.TableName()+".mobile LIKE ? OR "+memberModel.Member{}.TableName()+".nickname LIKE ? OR "+memberModel.Member{}.TableName()+".real_name LIKE ? OR "+memberModel.PointLog{}.TableName()+".remark LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Preload("Member").Scopes(info.Paginate()).Order(memberModel.PointLog{}.TableName() + ".created_at desc").Find(&list).Error
	return
}
