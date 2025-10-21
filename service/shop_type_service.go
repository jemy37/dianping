package service

import (
	"context"
	"dianping/dao"
	"dianping/utils"
)

// GetShopTypeList 获取商铺类型列表
func GetShopTypeList(ctx context.Context) *utils.Result {
	// 从缓存中获取
	shopTypes, err := dao.GetShopTypeListCache(ctx, dao.Redis)
	if err == nil {
		return utils.SuccessResultWithData(shopTypes)
	}

	// 从数据库中获取
	shopTypes, err = dao.GetShopTypeList(ctx, dao.DB)
	if err != nil {
		return utils.ErrorResult("查询失败: " + err.Error())
	}
	// 缓存到redis
	if err = dao.SetShopTypeListCache(ctx, dao.Redis, shopTypes); err != nil {
		println("警告: 设置 ShopType 缓存失败: " + err.Error())
	}

	return utils.SuccessResultWithData(shopTypes)
}
