package member

import (
	"errors"
	"strings"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RedemptionOrderService struct{}

func (s *RedemptionOrderService) CreateRedemptionOrder(req memberReq.CreateRedemptionOrderReq, operatorID uint) error {
	req.Remark = strings.TrimSpace(req.Remark)
	req.ReceiverName = strings.TrimSpace(req.ReceiverName)
	req.ReceiverPhone = strings.TrimSpace(req.ReceiverPhone)
	if req.Quantity <= 0 {
		return errors.New("兑换数量必须大于0")
	}

	return bizDB().Transaction(func(tx *gorm.DB) error {
		member, err := loadMember(tx, req.MemberID)
		if err != nil {
			return err
		}
		if member.Status != memberModel.MemberStatusEnabled {
			return errors.New("该会员已禁用，不能创建兑换订单。")
		}

		var product memberModel.PointProduct
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.ProductID).First(&product).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("积分商品不存在")
			}
			return err
		}
		if product.Status != memberModel.PointProductStatusOnSale {
			return errors.New("商品已下架，无法兑换")
		}
		if product.PointsPrice <= 0 {
			return errors.New("商品积分价格异常，无法兑换")
		}
		if product.Stock < req.Quantity {
			return errors.New("商品库存不足，无法兑换")
		}

		totalPoints := product.PointsPrice * req.Quantity
		change, err := applyPointChange(tx, req.MemberID, pointActionSpend, totalPoints)
		if err != nil {
			return err
		}

		stockResult := tx.Model(&memberModel.PointProduct{}).
			Where("id = ? AND stock >= ?", product.ID, req.Quantity).
			UpdateColumn("stock", gorm.Expr("stock - ?", req.Quantity))
		if stockResult.Error != nil {
			return stockResult.Error
		}
		if stockResult.RowsAffected != 1 {
			return errors.New("商品库存不足")
		}

		order := memberModel.RedemptionOrder{
			OrderNo:       buildOrderNo(),
			MemberID:      req.MemberID,
			ProductID:     product.ID,
			ProductName:   product.Name,
			Quantity:      req.Quantity,
			UnitPoints:    product.PointsPrice,
			TotalPoints:   totalPoints,
			Status:        memberModel.RedemptionOrderStatusPending,
			ReceiverName:  req.ReceiverName,
			ReceiverPhone: req.ReceiverPhone,
			Remark:        req.Remark,
		}
		if err = tx.Create(&order).Error; err != nil {
			return err
		}

		return recordPointTransaction(tx, req.MemberID, change.Account.ID, pointActionSpend, totalPoints, change.BeforeBalance, change.AfterBalance, memberModel.PointRefTypeRedemptionOrder, order.ID, formatOperator(operatorID), req.Remark)
	})
}

func (s *RedemptionOrderService) GetRedemptionOrder(id uint) (order memberModel.RedemptionOrder, err error) {
	err = bizDB().Preload("Member").Preload("Product").Where("id = ?", id).First(&order).Error
	return
}

func (s *RedemptionOrderService) GetRedemptionOrderList(info memberReq.RedemptionOrderSearch) (list []memberModel.RedemptionOrder, total int64, err error) {
	db := bizDB().Model(&memberModel.RedemptionOrder{}).
		Joins("LEFT JOIN " + memberModel.Member{}.TableName() + " ON " + memberModel.Member{}.TableName() + ".id = " + memberModel.RedemptionOrder{}.TableName() + ".member_id").
		Joins("LEFT JOIN " + memberModel.PointProduct{}.TableName() + " ON " + memberModel.PointProduct{}.TableName() + ".id = " + memberModel.RedemptionOrder{}.TableName() + ".product_id")

	if info.MemberID > 0 {
		db = db.Where(memberModel.RedemptionOrder{}.TableName()+".member_id = ?", info.MemberID)
	}
	if info.ProductID > 0 {
		db = db.Where(memberModel.RedemptionOrder{}.TableName()+".product_id = ?", info.ProductID)
	}
	if info.Status > 0 {
		db = db.Where(memberModel.RedemptionOrder{}.TableName()+".status = ?", info.Status)
	}
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where(memberModel.RedemptionOrder{}.TableName()+".order_no LIKE ? OR "+memberModel.RedemptionOrder{}.TableName()+".receiver_name LIKE ? OR "+memberModel.RedemptionOrder{}.TableName()+".receiver_phone LIKE ? OR "+memberModel.RedemptionOrder{}.TableName()+".product_name LIKE ? OR "+memberModel.Member{}.TableName()+".name LIKE ? OR "+memberModel.Member{}.TableName()+".phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Preload("Member").Preload("Product").Scopes(info.Paginate()).Order(memberModel.RedemptionOrder{}.TableName() + ".created_at desc").Find(&list).Error
	return
}

func (s *RedemptionOrderService) CompleteRedemptionOrder(req memberReq.OperateRedemptionOrderReq, operatorID uint) error {
	_ = operatorID
	return bizDB().Transaction(func(tx *gorm.DB) error {
		order, err := s.lockOrder(tx, req.ID)
		if err != nil {
			return err
		}
		if order.Status != memberModel.RedemptionOrderStatusPending {
			return errors.New("只有待处理订单可以完成")
		}

		updateResult := tx.Model(&memberModel.RedemptionOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"status": memberModel.RedemptionOrderStatusCompleted,
			"remark": mergeRemark(order.Remark, req.Remark),
		})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected != 1 {
			return errors.New("更新订单状态失败")
		}
		return nil
	})
}

func (s *RedemptionOrderService) CancelRedemptionOrder(req memberReq.OperateRedemptionOrderReq, operatorID uint) error {
	return bizDB().Transaction(func(tx *gorm.DB) error {
		order, err := s.lockOrder(tx, req.ID)
		if err != nil {
			return err
		}
		if order.Status != memberModel.RedemptionOrderStatusPending {
			return errors.New("只有待处理订单可以取消")
		}

		change, err := applyPointChange(tx, order.MemberID, pointActionRefund, order.TotalPoints)
		if err != nil {
			return err
		}

		stockResult := tx.Model(&memberModel.PointProduct{}).
			Where("id = ?", order.ProductID).
			UpdateColumn("stock", gorm.Expr("stock + ?", order.Quantity))
		if stockResult.Error != nil {
			return stockResult.Error
		}
		if stockResult.RowsAffected != 1 {
			return errors.New("恢复商品库存失败")
		}

		updateResult := tx.Model(&memberModel.RedemptionOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"status": memberModel.RedemptionOrderStatusCancelled,
			"remark": mergeRemark(order.Remark, req.Remark),
		})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected != 1 {
			return errors.New("更新订单状态失败")
		}

		return recordPointTransaction(tx, order.MemberID, change.Account.ID, pointActionRefund, order.TotalPoints, change.BeforeBalance, change.AfterBalance, memberModel.PointRefTypeRedemptionOrderVoid, order.ID, formatOperator(operatorID), req.Remark)
	})
}

func (s *RedemptionOrderService) lockOrder(tx *gorm.DB, id uint) (memberModel.RedemptionOrder, error) {
	var order memberModel.RedemptionOrder
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return order, errors.New("兑换订单不存在")
	}
	return order, err
}
