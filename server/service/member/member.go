package member

import (
	"errors"
	"strings"
	"time"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	memberRes "github.com/flipped-aurora/gin-vue-admin/server/model/member/response"
	"gorm.io/gorm"
)

type MemberService struct{}

func (s *MemberService) CreateMember(req memberReq.CreateMemberReq) error {
	member, err := s.buildMember(req.MemberBaseInput)
	if err != nil {
		return err
	}

	return bizDB().Transaction(func(tx *gorm.DB) error {
		if err = s.ensurePhoneUnique(tx, member.Phone, 0); err != nil {
			return err
		}
		if err = tx.Create(&member).Error; err != nil {
			return err
		}
		return tx.Create(&memberModel.PointAccount{MemberID: member.ID}).Error
	})
}

func (s *MemberService) DeleteMember(id uint) error {
	return bizDB().Transaction(func(tx *gorm.DB) error {
		if _, err := loadMember(tx, id); err != nil {
			return err
		}

		var transactionCount int64
		if err := tx.Model(&memberModel.PointTransaction{}).Where("member_id = ?", id).Count(&transactionCount).Error; err != nil {
			return err
		}

		var orderCount int64
		if err := tx.Model(&memberModel.RedemptionOrder{}).Where("member_id = ?", id).Count(&orderCount).Error; err != nil {
			return err
		}

		if transactionCount > 0 || orderCount > 0 {
			return errors.New("该会员已有积分记录或兑换订单，不允许删除，请改为禁用。")
		}

		if err := tx.Where("member_id = ?", id).Delete(&memberModel.PointAccount{}).Error; err != nil {
			return err
		}
		return tx.Delete(&memberModel.Member{}, "id = ?", id).Error
	})
}

func (s *MemberService) UpdateMember(req memberReq.UpdateMemberReq) error {
	member, err := s.buildMember(req.MemberBaseInput)
	if err != nil {
		return err
	}
	member.ID = req.ID

	return bizDB().Transaction(func(tx *gorm.DB) error {
		if _, err = loadMember(tx, req.ID); err != nil {
			return err
		}
		if err = s.ensurePhoneUnique(tx, member.Phone, req.ID); err != nil {
			return err
		}
		updates := map[string]interface{}{
			"name":     member.Name,
			"phone":    member.Phone,
			"gender":   member.Gender,
			"birthday": member.Birthday,
			"source":   member.Source,
			"level":    member.Level,
			"status":   member.Status,
			"remark":   member.Remark,
		}
		return tx.Model(&memberModel.Member{}).Where("id = ?", req.ID).Updates(updates).Error
	})
}

func (s *MemberService) GetMember(id uint) (memberRes.MemberDetail, error) {
	var result memberRes.MemberDetail
	member, err := s.getMemberEntity(id)
	if err != nil {
		return result, err
	}

	account, err := (&PointAccountService{}).GetPointAccountByMemberID(id)
	if err != nil {
		return result, err
	}

	result.Member = member
	result.Account = account
	return result, nil
}

func (s *MemberService) GetMemberList(info memberReq.MemberSearch) (list []memberModel.Member, total int64, err error) {
	db := bizDB().Model(&memberModel.Member{})
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where("name LIKE ? OR phone LIKE ? OR source LIKE ? OR remark LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if name := strings.TrimSpace(info.Name); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if phone := strings.TrimSpace(info.Phone); phone != "" {
		db = db.Where("phone LIKE ?", "%"+phone+"%")
	}
	if source := strings.TrimSpace(info.Source); source != "" {
		db = db.Where("source = ?", source)
	}
	if level := strings.TrimSpace(info.Level); level != "" {
		db = db.Where("level = ?", level)
	}
	if info.Status > 0 {
		db = db.Where("status = ?", info.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(info.Paginate()).Order("created_at desc").Find(&list).Error
	return
}

func (s *MemberService) UpdateMemberStatus(req memberReq.UpdateMemberStatusReq) error {
	if !s.isValidStatus(req.Status) {
		return errors.New("会员状态不合法")
	}
	result := bizDB().Model(&memberModel.Member{}).Where("id = ?", req.ID).Update("status", req.Status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("会员不存在")
	}
	return nil
}

func (s *MemberService) GetMemberOptions(keyword string) (list []memberRes.MemberOption, err error) {
	db := bizDB().Model(&memberModel.Member{}).Where("status = ?", memberModel.MemberStatusEnabled)
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		db = db.Where("name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var members []memberModel.Member
	if err = db.Order("created_at desc").Limit(50).Find(&members).Error; err != nil {
		return nil, err
	}

	list = make([]memberRes.MemberOption, 0, len(members))
	for _, item := range members {
		label := item.Name
		if label == "" {
			label = item.Phone
		}
		list = append(list, memberRes.MemberOption{
			ID:    item.ID,
			Label: label,
			Name:  item.Name,
			Phone: item.Phone,
			Level: item.Level,
		})
	}
	return list, nil
}

func (s *MemberService) getMemberEntity(id uint) (memberModel.Member, error) {
	return loadMember(bizDB(), id)
}

func (s *MemberService) buildMember(input memberReq.MemberBaseInput) (memberModel.Member, error) {
	member := memberModel.Member{
		Name:   strings.TrimSpace(input.Name),
		Phone:  strings.TrimSpace(input.Phone),
		Gender: strings.TrimSpace(input.Gender),
		Source: strings.TrimSpace(input.Source),
		Level:  strings.TrimSpace(input.Level),
		Remark: strings.TrimSpace(input.Remark),
	}
	if member.Level == "" {
		member.Level = memberModel.MemberLevelStandard
	}
	if input.Status == 0 {
		member.Status = memberModel.MemberStatusEnabled
	} else {
		member.Status = input.Status
	}
	if !s.isValidStatus(member.Status) {
		return member, errors.New("会员状态不合法")
	}
	if member.Name == "" {
		return member, errors.New("会员姓名不能为空")
	}
	if member.Phone == "" {
		return member, errors.New("手机号不能为空")
	}
	if input.Birthday != "" {
		birthday, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(input.Birthday), time.Local)
		if err != nil {
			return member, errors.New("生日格式不正确，请使用 YYYY-MM-DD")
		}
		member.Birthday = &birthday
	}
	return member, nil
}

func (s *MemberService) ensurePhoneUnique(tx *gorm.DB, phone string, excludeID uint) error {
	var count int64
	db := tx.Model(&memberModel.Member{}).Where("phone = ?", phone)
	if excludeID > 0 {
		db = db.Where("id <> ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("手机号已存在")
	}
	return nil
}

func (s *MemberService) isValidStatus(status int) bool {
	return status == memberModel.MemberStatusEnabled || status == memberModel.MemberStatusDisabled
}
