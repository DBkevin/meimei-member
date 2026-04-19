package member

import (
	"errors"
	"strings"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
)

type PointProductService struct{}

func (s *PointProductService) CreatePointProduct(req memberReq.CreatePointProductReq) error {
	product := s.buildPointProduct(req.PointProductBaseInput)
	if err := s.validatePointProduct(product); err != nil {
		return err
	}
	return bizDB().Create(&product).Error
}

func (s *PointProductService) DeletePointProduct(id uint) error {
	var count int64
	if err := bizDB().Model(&memberModel.RedemptionOrder{}).Where("product_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该商品已有兑换订单，不允许删除，请改为下架。")
	}
	return bizDB().Delete(&memberModel.PointProduct{}, "id = ?", id).Error
}

func (s *PointProductService) UpdatePointProduct(req memberReq.UpdatePointProductReq) error {
	product := s.buildPointProduct(req.PointProductBaseInput)
	product.ID = req.ID
	if err := s.validatePointProduct(product); err != nil {
		return err
	}

	result := bizDB().Model(&memberModel.PointProduct{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"name":         product.Name,
		"cover_url":    product.CoverURL,
		"category":     product.Category,
		"points_price": product.PointsPrice,
		"stock":        product.Stock,
		"status":       product.Status,
		"sort":         product.Sort,
		"description":  product.Description,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("积分商品不存在")
	}
	return nil
}

func (s *PointProductService) GetPointProduct(id uint) (product memberModel.PointProduct, err error) {
	err = bizDB().Where("id = ?", id).First(&product).Error
	if err != nil {
		return product, err
	}
	return
}

func (s *PointProductService) GetPointProductList(info memberReq.PointProductSearch) (list []memberModel.PointProduct, total int64, err error) {
	db := bizDB().Model(&memberModel.PointProduct{})
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where("name LIKE ? OR category LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category := strings.TrimSpace(info.Category); category != "" {
		db = db.Where("category = ?", category)
	}
	if info.Status > 0 {
		db = db.Where("status = ?", info.Status)
	}

	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Scopes(info.Paginate()).Order("sort asc, created_at desc").Find(&list).Error
	return
}

func (s *PointProductService) UpdatePointProductStatus(req memberReq.UpdatePointProductStatusReq) error {
	if !s.isValidProductStatus(req.Status) {
		return errors.New("商品状态不合法")
	}
	result := bizDB().Model(&memberModel.PointProduct{}).Where("id = ?", req.ID).Update("status", req.Status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("积分商品不存在")
	}
	return nil
}

func (s *PointProductService) GetPointProductOptions(keyword string) (list []memberModel.PointProduct, err error) {
	db := bizDB().Model(&memberModel.PointProduct{}).Where("status = ?", memberModel.PointProductStatusOnSale)
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = db.Order("sort asc, created_at desc").Limit(50).Find(&list).Error
	return
}

func (s *PointProductService) buildPointProduct(input memberReq.PointProductBaseInput) memberModel.PointProduct {
	product := memberModel.PointProduct{
		Name:        strings.TrimSpace(input.Name),
		CoverURL:    strings.TrimSpace(input.CoverURL),
		Category:    strings.TrimSpace(input.Category),
		PointsPrice: input.PointsPrice,
		Stock:       input.Stock,
		Status:      input.Status,
		Sort:        input.Sort,
		Description: strings.TrimSpace(input.Description),
	}
	if product.Status == 0 {
		product.Status = memberModel.PointProductStatusOffSale
	}
	return product
}

func (s *PointProductService) validatePointProduct(product memberModel.PointProduct) error {
	if product.Name == "" {
		return errors.New("商品名称不能为空")
	}
	if product.PointsPrice <= 0 {
		return errors.New("兑换积分必须大于0")
	}
	if product.Stock < 0 {
		return errors.New("库存不能小于0")
	}
	if product.Sort < 0 {
		return errors.New("排序不能小于0")
	}
	if !s.isValidProductStatus(product.Status) {
		return errors.New("商品状态不合法")
	}
	return nil
}

func (s *PointProductService) isValidProductStatus(status int) bool {
	return status == memberModel.PointProductStatusOnSale || status == memberModel.PointProductStatusOffSale
}
