package dao

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"dianping/models"
)

func GetShopTypeList(ctx context.Context, db *gorm.DB) ([]*models.ShopType, error) {
	var shopTypes []*models.ShopType
	err := db.WithContext(ctx).Model(&models.ShopType{}).Order("sort").Find(&shopTypes).Error
	if err != nil {
		return nil, err
	}
	return shopTypes, nil
}

// ===========缓存相关=============

const (
	ShopTypeCache = "cache:shop_type"
)

func GetShopTypeListCache(ctx context.Context, rds *redis.Client) ([]*models.ShopType, error) {
	key := ShopTypeCache
	str, err := rds.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var shopTypes []*models.ShopType
	if err := json.Unmarshal([]byte(str), &shopTypes); err != nil {
		return nil, err
	}
	return shopTypes, nil
}

func SetShopTypeListCache(ctx context.Context, rds *redis.Client, shopTypes []*models.ShopType) error {
	// 设置一小时的过期时间
	b, err := json.Marshal(shopTypes)
	if err != nil {
		return err
	}
	return rds.Set(ctx, ShopTypeCache, b, time.Hour).Err()
}
