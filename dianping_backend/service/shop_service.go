package service

import (
	"context"
	"dianping/dao"
	"dianping/models"
	"dianping/utils"
	"fmt"
	"log"
	"math/rand/v2"
	"net/url"
	"time"
)

var shopSF utils.SingleflightGroup

// GetShopById 根据ID获取商铺
func GetShopById(ctx context.Context, id uint) *utils.Result {
	// 1. 布隆过滤器检查，防止缓存击穿
	flag, err := utils.CheckIDExistsWithRedis(ctx, dao.Redis, "shop", id)
	if err != nil {
		log.Fatalf("检查布隆过滤器失败: %v", err)
	}
	if !flag {
		// 布隆过滤器判断商铺不存在，直接返回
		return utils.ErrorResult("商铺不存在")
	}

	// 2. 先读取带逻辑过期的缓存
	if cached, fresh, err := dao.GetShopCacheByIdWithLogicalExpire(ctx, dao.Redis, id); err == nil && cached != nil {
		if fresh {
			return utils.SuccessResultWithData(cached)
		}
		go tryRebuildShopCache(context.Background(), id)
		return utils.SuccessResultWithData(cached)
	}

	// 3. 使用 singleflight 合并并发重建
	key := fmt.Sprintf("shop:%d", id)
	val, err := shopSF.Do(key, func() (interface{}, error) {
		if cached, fresh, err := dao.GetShopCacheByIdWithLogicalExpire(ctx, dao.Redis, id); err == nil && cached != nil {
			if fresh {
				return cached, nil
			}
		}
		return loadAndFillShopCache(ctx, id)
	})
	if err != nil {
		return utils.ErrorResult("查询失败: " + err.Error())
	}
	return utils.SuccessResultWithData(val)
}

// UpdateShopById 根据ID更新商铺
func UpdateShopById(ctx context.Context, shop *models.Shop) *utils.Result {

	// 0. 启动事务
	tx := dao.DB.Begin()
	defer func() { // 捕获异常
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 更新数据库
	err := dao.UpdateShop(ctx, tx, shop)

	// 2. 更新失败
	if err != nil {
		tx.Rollback()
		return utils.ErrorResult("更新失败: " + err.Error())
	}

	// 3. 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return utils.ErrorResult("更新失败: " + err.Error())
	}

	// 4. 事务成功后删除缓存（最终一致性）
	err = dao.DelShopCacheById(ctx, dao.Redis, shop.ID)
	if err != nil {
		// 记录日志但不影响业务结果
		log.Printf("警告: 删除缓存失败，商铺ID=%d, 错误=%v", shop.ID, err)
	}

	// 5. 返回结果
	return utils.SuccessResult("更新成功")
}

// GetShopList 获取商铺列表
func GetShopList(page, size int) *utils.Result {
	var shops []models.Shop
	var total int64

	offset := (page - 1) * size

	// 获取总数
	dao.DB.Model(&models.Shop{}).Count(&total)

	// 分页查询
	err := dao.DB.Offset(offset).Limit(size).Find(&shops).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  shops,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetShopByType 根据类型获取商铺
func GetShopByType(typeId uint, page, size int) *utils.Result {
	var shops []models.Shop
	var total int64

	offset := (page - 1) * size

	// 获取总数
	dao.DB.Model(&models.Shop{}).Where("type_id = ?", typeId).Count(&total)

	// 分页查询
	err := dao.DB.Where("type_id = ?", typeId).Offset(offset).Limit(size).Find(&shops).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  shops,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetShopByName 根据名称搜索商铺
func GetShopByName(name string, page, size int) *utils.Result {
	var shops []models.Shop
	var total int64

	offset := (page - 1) * size

	query := dao.DB.Model(&models.Shop{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	err := query.Offset(offset).Limit(size).Find(&shops).Error
	if err != nil {
		return utils.ErrorResult("查询失败")
	}

	return utils.SuccessResultWithData(map[string]interface{}{
		"list":  shops,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetNearbyShops 获取某个店铺的附近某个距离的所有点
func GetNearbyShops(ctx context.Context, shopId uint, radius float64, count int) *utils.Result {
	// 1. 查询店铺
	shop, err := dao.GetShopById(ctx, dao.DB, shopId)
	if err != nil {
		return utils.ErrorResult("查询店铺失败: " + err.Error())
	}

	// 2. 查询附近的同类型商铺
	shopIds, err := dao.GetNearbyShops(ctx, dao.Redis, shop, radius, "km", count)
	if err != nil {
		return utils.ErrorResult("查询附近商铺失败: " + err.Error())
	}

	// 3. 返回结果
	return utils.SuccessResultWithData(shopIds)
}

func tryRebuildShopCache(ctx context.Context, id uint) {
	lockKey := fmt.Sprintf("lock:shop:rebuild:%d", id)
	ok, lockVal := utils.TryLockWithTTL(ctx, dao.Redis, lockKey, 5*time.Second)
	if !ok {
		return
	}
	defer utils.UnLockSafe(ctx, dao.Redis, lockKey, lockVal)
	if _, err := loadAndFillShopCache(ctx, id); err != nil {
		log.Printf("异步重建商铺缓存失败: %v", err)
	}
}

func loadAndFillShopCache(ctx context.Context, id uint) (*models.Shop, error) {
	lockKey := fmt.Sprintf("lock:shop:%d", id)
	ok, lockVal := utils.TryLockWithTTL(ctx, dao.Redis, lockKey, 5*time.Second)
	if !ok {
		base := 50 * time.Millisecond
		waited := time.Duration(0)
		for i := 0; i < 5 && waited < 500*time.Millisecond; i++ {
			time.Sleep(base)
			waited += base
			base *= 2
			if base > 200*time.Millisecond {
				base = 200 * time.Millisecond
			}
			if cached, _, err := dao.GetShopCacheByIdWithLogicalExpire(ctx, dao.Redis, id); err == nil && cached != nil {
				return cached, nil
			}
		}
		var retry bool
		retry, lockVal = utils.TryLockWithTTL(ctx, dao.Redis, lockKey, 3*time.Second)
		if !retry {
			if cached, _, err := dao.GetShopCacheByIdWithLogicalExpire(ctx, dao.Redis, id); err == nil && cached != nil {
				return cached, nil
			}
			return nil, fmt.Errorf("服务繁忙")
		}
		ok = true
	}
	defer func() {
		if ok {
			utils.UnLockSafe(ctx, dao.Redis, lockKey, lockVal)
		}
	}()

	shop, err := dao.GetShopById(ctx, dao.DB, id)
	if err != nil {
		return nil, err
	}

	ttlMinutes := 30 + (rand.IntN(10) - 5)
	if ttlMinutes < 5 {
		ttlMinutes = 5
	}
	expireAt := time.Now().Add(time.Duration(ttlMinutes) * time.Minute).Unix()
	realTTL := 24 * time.Hour
	if err := dao.SetShopCacheByIdWithLogicalExpire(ctx, dao.Redis, id, shop, expireAt, realTTL); err != nil {
		log.Printf("设置逻辑过期缓存失败: %v", err)
	}
	return shop, nil
}

// CreateShopWithType 创建类型商铺

func CreateShopWithType(ctx context.Context, shop *models.Shop, typeIcon string) *utils.Result {
	//使用本地 mutex 限制同进程并发
	shopTypeListLocalLock.Lock()
	defer shopTypeListLocalLock.Unlock()

	// 使用按 name 的分布式锁，避免并发创建重复的 ShopType
	typeName := shop.Name
	//通过escape 转义特殊字符
	lockKey := fmt.Sprintf("lock:shop_type:name:%s", url.QueryEscape(typeName))
	ok, lockValue := utils.TryLockWithTTL(ctx, dao.Redis, lockKey, 5*time.Second)
	if !ok {
		maxWait := 1000 * time.Millisecond
		base := 50 * time.Millisecond
		waited := time.Duration(0)
		for waited <= maxWait {

			//尝试再拿锁
			if ok2, lockValue2 := utils.TryLockWithTTL(ctx, dao.Redis, lockKey, 5*time.Second); ok2 {
				lockValue = lockValue2
				break
			}
			// 等待指数退避+随机抖动
			sleep := base + time.Duration(rand.IntN(50))*time.Millisecond
			time.Sleep(sleep)
			waited += sleep
			base *= 2
			if base > 500*time.Millisecond {
				base = 500 * time.Millisecond
			}
		}
		// 未拿到锁：短等待后再次检查数据库
		time.Sleep(50 * time.Millisecond)
		if existing, err := dao.GetShopByNameAndAddress(ctx, dao.DB, shop.Name, shop.Address); err == nil && existing != nil {
			return utils.SuccessResultWithData(existing)

		} else {
			return utils.ErrorResult("操作超时，请重试")
		}
	} else {
		// 拿到锁，确保释放
		defer func() {
			if ok := utils.UnLockSafe(ctx, dao.Redis, lockKey, lockValue); !ok {
				log.Printf("警告: 释放 shop_type:name 锁失败: %s", typeName)
			}
		}()
	}

	// 再次检查 DB，以防并发已创建
	if existing, err := dao.GetShopByNameAndAddress(ctx, dao.DB, shop.Name, shop.Address); err == nil && existing != nil {

		return utils.SuccessResultWithData(existing)
	}

	// 没有则创建shop
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := dao.CreateShop(ctx, tx, shop); err != nil {
		tx.Rollback()
		return utils.ErrorResult("创建商铺失败: " + err.Error())
	}

	st := &models.ShopType{
		Name:   typeName,
		Icon:   typeIcon,
		TypeID: shop.TypeID,
	}
	if err := dao.CreateShopType(ctx, tx, st); err != nil {
		tx.Rollback()
		return utils.ErrorResult("创建商铺类型失败: " + err.Error())
	}
	// 用回写ID的方式设置 Sort = int(st.ID)
	st.Sort = int(st.ID)
	st.ShopId = shop.ID
	if err := dao.UpdateShopType(ctx, tx, st); err != nil {
		tx.Rollback()
		return utils.ErrorResult("更新 ShopType Sort 失败: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return utils.ErrorResult("提交事务失败: " + err.Error())
	}

	go func(id uint) {
		bf := utils.CreateShopBloomFilter(dao.Redis)
		if _, err := bf.AddID(context.Background(), id); err != nil {
			log.Printf("添加商铺ID到Bloom失败: %v", err)
		}
	}(shop.ID)

	return utils.SuccessResultWithData(shop.ID)
}
