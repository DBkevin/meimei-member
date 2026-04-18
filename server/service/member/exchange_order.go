package member

import (
	"errors"
	"strings"
	"time"

	memberModel "github.com/flipped-aurora/gin-vue-admin/server/model/member"
	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExchangeOrderService struct{}

func (s *ExchangeOrderService) CreateExchangeOrder(req memberReq.CreateExchangeOrderReq, operatorID uint) error {
	return bizDB().Transaction(func(tx *gorm.DB) error {
		mem, err := loadMember(tx, req.MemberID)
		if err != nil {
			return err
		}
		if mem.Status != memberModel.MemberStatusEnabled {
			return errors.New("会员已被禁用，无法兑换")
		}

		var goods memberModel.PointGoods
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.GoodsID).First(&goods).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("积分商品不存在")
			}
			return err
		}
		if goods.Status != memberModel.GoodsStatusOnSale {
			return errors.New("商品未上架，无法兑换")
		}
		if goods.Stock <= 0 {
			return errors.New("商品库存不足")
		}
		if goods.LimitPerMember > 0 {
			var usedCount int64
			if err = tx.Model(&memberModel.ExchangeOrder{}).
				Where("member_id = ? AND goods_id = ? AND status IN ?", req.MemberID, req.GoodsID, []string{memberModel.OrderStatusPending, memberModel.OrderStatusCompleted}).
				Count(&usedCount).Error; err != nil {
				return err
			}
			if usedCount >= goods.LimitPerMember {
				return errors.New("已达到该商品每人限兑数量")
			}
		}

		order := memberModel.ExchangeOrder{
			OrderNo:    buildOrderNo(),
			MemberID:   req.MemberID,
			GoodsID:    req.GoodsID,
			PointsCost: goods.PointsPrice,
			Status:     memberModel.OrderStatusPending,
			VerifyCode: buildVerifyCode(),
			OperatorID: operatorID,
			Remark:     strings.TrimSpace(req.Remark),
		}
		if err = tx.Create(&order).Error; err != nil {
			return err
		}

		if _, err = adjustPoints(tx, req.MemberID, memberModel.PointChangeTypeUse, goods.PointsPrice, memberModel.PointSourceTypeExchangeOrder, order.ID, order.Remark, operatorID); err != nil {
			return err
		}

		result := tx.Model(&memberModel.PointGoods{}).
			Where("id = ? AND stock >= ?", goods.ID, int64(1)).
			UpdateColumn("stock", gorm.Expr("stock - ?", int64(1)))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errors.New("商品库存不足")
		}
		return nil
	})
}

func (s *ExchangeOrderService) GetExchangeOrder(id uint) (order memberModel.ExchangeOrder, err error) {
	err = bizDB().Preload("Member").Preload("Goods").Where("id = ?", id).First(&order).Error
	return
}

func (s *ExchangeOrderService) GetExchangeOrderList(info memberReq.ExchangeOrderSearch) (list []memberModel.ExchangeOrder, total int64, err error) {
	db := bizDB().Model(&memberModel.ExchangeOrder{}).
		Joins("LEFT JOIN " + memberModel.Member{}.TableName() + " ON " + memberModel.Member{}.TableName() + ".id = " + memberModel.ExchangeOrder{}.TableName() + ".member_id").
		Joins("LEFT JOIN " + memberModel.PointGoods{}.TableName() + " ON " + memberModel.PointGoods{}.TableName() + ".id = " + memberModel.ExchangeOrder{}.TableName() + ".goods_id")

	if info.MemberID > 0 {
		db = db.Where(memberModel.ExchangeOrder{}.TableName()+".member_id = ?", info.MemberID)
	}
	if info.GoodsID > 0 {
		db = db.Where(memberModel.ExchangeOrder{}.TableName()+".goods_id = ?", info.GoodsID)
	}
	if info.Status != "" {
		db = db.Where(memberModel.ExchangeOrder{}.TableName()+".status = ?", info.Status)
	}
	if keyword := strings.TrimSpace(info.Keyword); keyword != "" {
		db = db.Where(memberModel.ExchangeOrder{}.TableName()+".order_no LIKE ? OR "+memberModel.ExchangeOrder{}.TableName()+".verify_code LIKE ? OR "+memberModel.Member{}.TableName()+".mobile LIKE ? OR "+memberModel.Member{}.TableName()+".nickname LIKE ? OR "+memberModel.Member{}.TableName()+".real_name LIKE ? OR "+memberModel.PointGoods{}.TableName()+".name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Preload("Member").Preload("Goods").Scopes(info.Paginate()).Order(memberModel.ExchangeOrder{}.TableName() + ".created_at desc").Find(&list).Error
	return
}

func (s *ExchangeOrderService) VerifyExchangeOrder(req memberReq.OperateExchangeOrderReq, operatorID uint) error {
	return bizDB().Transaction(func(tx *gorm.DB) error {
		order, err := s.lockOrder(tx, req.ID)
		if err != nil {
			return err
		}
		if order.Status != memberModel.OrderStatusPending {
			return errors.New("只有待核销订单可以核销")
		}
		now := time.Now()
		updates := map[string]interface{}{
			"status":      memberModel.OrderStatusCompleted,
			"verified_at": &now,
			"operator_id": operatorID,
			"remark":      mergeRemark(order.Remark, req.Remark),
		}
		return tx.Model(&memberModel.ExchangeOrder{}).Where("id = ?", order.ID).Updates(updates).Error
	})
}

func (s *ExchangeOrderService) CancelExchangeOrder(req memberReq.OperateExchangeOrderReq, operatorID uint) error {
	return s.rollbackExchangeOrder(req, operatorID, memberModel.OrderStatusPending, memberModel.OrderStatusCancelled)
}

func (s *ExchangeOrderService) RefundExchangeOrder(req memberReq.OperateExchangeOrderReq, operatorID uint) error {
	return s.rollbackExchangeOrder(req, operatorID, memberModel.OrderStatusCompleted, memberModel.OrderStatusRefunded)
}

func (s *ExchangeOrderService) rollbackExchangeOrder(req memberReq.OperateExchangeOrderReq, operatorID uint, fromStatus string, toStatus string) error {
	return bizDB().Transaction(func(tx *gorm.DB) error {
		order, err := s.lockOrder(tx, req.ID)
		if err != nil {
			return err
		}
		if order.Status != fromStatus {
			if toStatus == memberModel.OrderStatusCancelled {
				return errors.New("只有待核销订单可以取消")
			}
			return errors.New("只有已完成订单可以退款")
		}

		var goods memberModel.PointGoods
		if err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", order.GoodsID).First(&goods).Error; err != nil {
			return err
		}

		if _, err = adjustPoints(tx, order.MemberID, memberModel.PointChangeTypeRefund, order.PointsCost, memberModel.PointSourceTypeExchangeOrderVoid, order.ID, req.Remark, operatorID); err != nil {
			return err
		}

		if err = tx.Model(&memberModel.PointGoods{}).Where("id = ?", goods.ID).UpdateColumn("stock", gorm.Expr("stock + ?", 1)).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"status":      toStatus,
			"operator_id": operatorID,
			"remark":      mergeRemark(order.Remark, req.Remark),
		}
		if toStatus == memberModel.OrderStatusRefunded {
			updates["verified_at"] = nil
		}
		return tx.Model(&memberModel.ExchangeOrder{}).Where("id = ?", order.ID).Updates(updates).Error
	})
}

func (s *ExchangeOrderService) lockOrder(tx *gorm.DB, id uint) (memberModel.ExchangeOrder, error) {
	var order memberModel.ExchangeOrder
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return order, errors.New("兑换订单不存在")
	}
	return order, err
}
