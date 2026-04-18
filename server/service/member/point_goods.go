package member

import (
	"errors"
	"strings"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
)

type PointGoodsService struct{}

func (s *PointGoodsService) CreatePointGoods(info *memberModel.PointGoods) error {
	s.normalizeGoods(info)
	return bizDB().Create(info).Error
}

func (s *PointGoodsService) DeletePointGoods(id uint) error {
	var count int64
	if err := bizDB().Model(&memberModel.ExchangeOrder{}).Where("goods_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("商品已有兑换订单，无法删除")
	}
	return bizDB().Delete(&memberModel.PointGoods{}, "id = ?", id).Error
}

func (s *PointGoodsService) UpdatePointGoods(info *memberModel.PointGoods) error {
	s.normalizeGoods(info)
	updates := map[string]interface{}{
		"name":             info.Name,
		"cover_image":      info.CoverImage,
		"description":      info.Description,
		"points_price":     info.PointsPrice,
		"stock":            info.Stock,
		"limit_per_member": info.LimitPerMember,
		"status":           info.Status,
		"sort":             info.Sort,
	}
	return bizDB().Model(&memberModel.PointGoods{}).Where("id = ?", info.ID).Updates(updates).Error
}

func (s *PointGoodsService) GetPointGoods(id uint) (goods memberModel.PointGoods, err error) {
	err = bizDB().Where("id = ?", id).First(&goods).Error
	return
}

func (s *PointGoodsService) GetPointGoodsList(info memberReq.PointGoodsSearch) (list []memberModel.PointGoods, total int64, err error) {
	db := bizDB().Model(&memberModel.PointGoods{})
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(info.Paginate()).Order("sort asc, created_at desc").Find(&list).Error
	return
}

func (s *PointGoodsService) UpdatePointGoodsStatus(req memberReq.UpdateGoodsStatusReq) error {
	if req.Status != memberModel.GoodsStatusOnSale && req.Status != memberModel.GoodsStatusOffSale {
		return errors.New("商品状态不合法")
	}
	return bizDB().Model(&memberModel.PointGoods{}).Where("id = ?", req.ID).Update("status", req.Status).Error
}

func (s *PointGoodsService) UpdatePointGoodsStock(req memberReq.UpdateGoodsStockReq) error {
	if req.Stock < 0 {
		return errors.New("库存不能小于0")
	}
	return bizDB().Model(&memberModel.PointGoods{}).Where("id = ?", req.ID).Update("stock", req.Stock).Error
}

func (s *PointGoodsService) GetPointGoodsOptions(keyword string) (list []memberModel.PointGoods, err error) {
	db := bizDB().Model(&memberModel.PointGoods{}).Where("status = ?", memberModel.GoodsStatusOnSale)
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = db.Order("sort asc, created_at desc").Limit(50).Find(&list).Error
	return
}

func (s *PointGoodsService) normalizeGoods(info *memberModel.PointGoods) {
	info.Name = strings.TrimSpace(info.Name)
	info.CoverImage = strings.TrimSpace(info.CoverImage)
	info.Description = strings.TrimSpace(info.Description)
	if info.Status == "" {
		info.Status = memberModel.GoodsStatusOnSale
	}
}
