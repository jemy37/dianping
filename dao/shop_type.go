package dao

import (
	"context"
	"encoding/json"
	"errors"
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

// GetShopTypeByName 根据 name 查询 ShopType，找不到返回 (nil, nil)
func GetShopTypeByName(ctx context.Context, db *gorm.DB, name string) (*models.ShopType, error) {
	var st models.ShopType
	err := db.WithContext(ctx).Where("name = ?", name).First(&st).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &st, nil
}

// UpdateShopType 更新 ShopType，通常用于回写 Sort 字段
func UpdateShopType(ctx context.Context, db *gorm.DB, st *models.ShopType) error {
	return db.WithContext(ctx).Model(&models.ShopType{}).Where("id = ?", st.ID).Updates(st).Error
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
