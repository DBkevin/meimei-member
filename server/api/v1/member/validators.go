package member

import (
	"errors"
	"fmt"
	"regexp"

	memberReq "github.com/flipped-aurora/gin-vue-admin/server/model/member/request"
)

// ValidateMemberInput 验证会员输入参数
func ValidateMemberInput(input memberReq.MemberBaseInput) error {
	if input.Name == "" || len(input.Name) > 100 {
		return errors.New("会员名称不能为空且不超过100个字符")
	}

	if input.Phone != "" {
		matched, _ := regexp.MatchString(`^1[0-9]{10}$`, input.Phone)
		if !matched {
			return errors.New("手机号格式不正确")
		}
	}

	if input.Gender != "" && input.Gender != "M" && input.Gender != "F" && input.Gender != "U" {
		return errors.New("性别值无效，应为 M/F/U")
	}

	if input.Source != "" && len(input.Source) > 50 {
		return errors.New("来源不超过50个字符")
	}

	if input.Level != "" && len(input.Level) > 20 {
		return errors.New("会员等级不超过20个字符")
	}

	if input.Status != 1 && input.Status != 2 {
		return errors.New("会员状态值无效，应为 1启用 / 2禁用")
	}

	return nil
}

// ValidatePointProductInput 验证积分商品输入参数
func ValidatePointProductInput(input memberReq.PointProductBaseInput) error {
	if input.Name == "" || len(input.Name) > 200 {
		return errors.New("商品名称不能为空且不超过200个字符")
	}

	if input.Category == "" || len(input.Category) > 100 {
		return errors.New("商品分类不能为空且不超过100个字符")
	}

	if input.PointsPrice <= 0 || input.PointsPrice > 1000000 {
		return errors.New("积分兑换价格必须大于0且不超过1000000")
	}

	if input.Stock < 0 {
		return errors.New("商品库存不能为负数")
	}

	if input.Status != 1 && input.Status != 2 {
		return errors.New("商品状态值无效，应为 1上架 / 2下架")
	}

	if input.Sort < 0 || input.Sort > 1000 {
		return errors.New("排序值应在 0-1000 之间")
	}

	return nil
}

// ValidateRedemptionOrderInput 验证兑换订单输入参数
func ValidateRedemptionOrderInput(input memberReq.CreateRedemptionOrderReq) error {
	if input.MemberID == 0 {
		return errors.New("会员ID不能为空")
	}

	if input.ProductID == 0 {
		return errors.New("商品ID不能为空")
	}

	if input.Quantity <= 0 || input.Quantity > 1000 {
		return errors.New("兑换数量必须大于0且不超过1000")
	}

	if len(input.ReceiverName) > 50 {
		return errors.New("收货人名称不超过50个字符")
	}

	if len(input.ReceiverPhone) > 20 {
		return errors.New("收货人电话不超过20个字符")
	}

	if len(input.Remark) > 500 {
		return errors.New("备注不超过500个字符")
	}

	return nil
}

// ValidateAdjustPointsInput 验证调整积分输入参数
func ValidateAdjustPointsInput(input memberReq.AdjustPointsReq) error {
	if input.MemberID == 0 {
		return errors.New("会员ID不能为空")
	}

	if input.Points <= 0 || input.Points > 10000000 {
		return errors.New("调整积分数量必须大于0且不超过10000000")
	}

	if input.Remark == "" || len(input.Remark) > 500 {
		return errors.New("备注不能为空且不超过500个字符")
	}

	return nil
}

// ValidateIDInput 验证ID输入参数
func ValidateIDInput(id uint) error {
	if id == 0 {
		return fmt.Errorf("ID不能为空或为0")
	}
	return nil
}
