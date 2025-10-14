package dao

import (
	"dianping/models"
)

// GetAllVoucherIDs 获取所有优惠券ID
func GetAllVoucherIDs() ([]uint, error) {
	var ids []uint
	err := DB.Model(&models.Voucher{}).Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}
