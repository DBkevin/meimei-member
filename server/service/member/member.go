package member

import (
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	memberRes "github.com/flipped-aurora/gin-vue-admin/server/model/member/response"
	"gorm.io/gorm"
)

type MemberService struct{}

func (s *MemberService) CreateMember(info *memberModel.Member) error {
	info.OpenID = strings.TrimSpace(info.OpenID)
	info.UnionID = strings.TrimSpace(info.UnionID)
	info.Mobile = strings.TrimSpace(info.Mobile)
	info.Nickname = strings.TrimSpace(info.Nickname)
	info.RealName = strings.TrimSpace(info.RealName)
	info.AvatarURL = strings.TrimSpace(info.AvatarURL)
	if info.MemberLevel == "" {
		info.MemberLevel = "standard"
	}
	if info.Status == "" {
		info.Status = memberModel.MemberStatusEnabled
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := s.ensureMemberUnique(tx, info.OpenID, info.Mobile, 0); err != nil {
			return err
		}
		if err := tx.Create(info).Error; err != nil {
			return err
		}
		return tx.Create(&memberModel.PointAccount{MemberID: info.ID}).Error
	})
}

func (s *MemberService) DeleteMember(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var flowCount int64
		if err := tx.Model(&memberModel.PointLog{}).Where("member_id = ?", id).Count(&flowCount).Error; err != nil {
			return err
		}
		if flowCount > 0 {
			return errors.New("会员已有积分流水，无法删除")
		}

		var orderCount int64
		if err := tx.Model(&memberModel.ExchangeOrder{}).Where("member_id = ?", id).Count(&orderCount).Error; err != nil {
			return err
		}
		if orderCount > 0 {
			return errors.New("会员已有兑换订单，无法删除")
		}

		if err := tx.Where("member_id = ?", id).Delete(&memberModel.PointAccount{}).Error; err != nil {
			return err
		}
		return tx.Delete(&memberModel.Member{}, "id = ?", id).Error
	})
}

func (s *MemberService) UpdateMember(info *memberModel.Member) error {
	info.OpenID = strings.TrimSpace(info.OpenID)
	info.UnionID = strings.TrimSpace(info.UnionID)
	info.Mobile = strings.TrimSpace(info.Mobile)
	info.Nickname = strings.TrimSpace(info.Nickname)
	info.RealName = strings.TrimSpace(info.RealName)
	info.AvatarURL = strings.TrimSpace(info.AvatarURL)
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := s.ensureMemberUnique(tx, info.OpenID, info.Mobile, info.ID); err != nil {
			return err
		}
		updates := map[string]interface{}{
			"openid":       info.OpenID,
			"unionid":      info.UnionID,
			"mobile":       info.Mobile,
			"nickname":     info.Nickname,
			"avatar_url":   info.AvatarURL,
			"real_name":    info.RealName,
			"member_level": info.MemberLevel,
			"status":       info.Status,
		}
		return tx.Model(&memberModel.Member{}).Where("id = ?", info.ID).Updates(updates).Error
	})
}

func (s *MemberService) GetMember(id uint) (memberRes.MemberDetail, error) {
	var result memberRes.MemberDetail
	err := global.GVA_DB.Where("id = ?", id).First(&result.Member).Error
	if err != nil {
		return result, err
	}
	_ = global.GVA_DB.Where("member_id = ?", id).First(&result.Account).Error
	return result, nil
}

func (s *MemberService) GetMemberList(info memberReq.MemberSearch) (list []memberModel.Member, total int64, err error) {
	db := global.GVA_DB.Model(&memberModel.Member{})
	if info.Mobile != "" {
		db = db.Where("mobile LIKE ?", "%"+strings.TrimSpace(info.Mobile)+"%")
	}
	if info.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+strings.TrimSpace(info.Nickname)+"%")
	}
	if info.RealName != "" {
		db = db.Where("real_name LIKE ?", "%"+strings.TrimSpace(info.RealName)+"%")
	}
	if info.MemberLevel != "" {
		db = db.Where("member_level = ?", info.MemberLevel)
	}
	if info.Status != "" {
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
	if req.Status != memberModel.MemberStatusEnabled && req.Status != memberModel.MemberStatusDisabled {
		return errors.New("会员状态不合法")
	}
	return global.GVA_DB.Model(&memberModel.Member{}).Where("id = ?", req.ID).Update("status", req.Status).Error
}

func (s *MemberService) GetMemberOptions(keyword string) (list []memberRes.MemberOption, err error) {
	db := global.GVA_DB.Model(&memberModel.Member{}).Where("status = ?", memberModel.MemberStatusEnabled)
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		db = db.Where("mobile LIKE ? OR nickname LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	var members []memberModel.Member
	err = db.Order("created_at desc").Limit(50).Find(&members).Error
	if err != nil {
		return
	}
	list = make([]memberRes.MemberOption, 0, len(members))
	for _, item := range members {
		label := item.Mobile
		if item.RealName != "" {
			label = item.RealName
		} else if item.Nickname != "" {
			label = item.Nickname
		}
		list = append(list, memberRes.MemberOption{
			ID:          item.ID,
			Label:       label,
			Mobile:      item.Mobile,
			Nickname:    item.Nickname,
			RealName:    item.RealName,
			MemberLevel: item.MemberLevel,
		})
	}
	return
}

func (s *MemberService) ensureMemberUnique(tx *gorm.DB, openID, mobile string, excludeID uint) error {
	if openID != "" {
		var count int64
		db := tx.Model(&memberModel.Member{}).Where("openid = ?", openID)
		if excludeID > 0 {
			db = db.Where("id <> ?", excludeID)
		}
		if err := db.Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("openid 已存在")
		}
	}
	if mobile != "" {
		var count int64
		db := tx.Model(&memberModel.Member{}).Where("mobile = ?", mobile)
		if excludeID > 0 {
			db = db.Where("id <> ?", excludeID)
		}
		if err := db.Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("手机号已存在")
		}
	}
	return nil
}
