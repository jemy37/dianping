package service

import (
	"context"
	"dianping/dao"
	"dianping/utils"
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

var shopTypeListLocalLock sync.Mutex

// GetShopTypeList 获取商铺类型列表，防止缓存击穿使用本地+分布式混合锁
func GetShopTypeList(ctx context.Context) *utils.Result {
	// 1) 先从缓存中获取 无竞争 缓存正常
	shopTypes, err := dao.GetShopTypeListCache(ctx, dao.Redis)
	if err == nil {
		return utils.SuccessResultWithData(shopTypes)
	}

	// 2) 使用本地 mutex 限制同进程并发
	shopTypeListLocalLock.Lock()
	defer shopTypeListLocalLock.Unlock()

	// 双重检查缓存 防止同进程内的多个协程重复获取锁，此时可能缓存已经被其他协程回写
	shopTypes, err = dao.GetShopTypeListCache(ctx, dao.Redis)
	if err == nil {
		return utils.SuccessResultWithData(shopTypes)
	}

	// 3) 尝试获取分布式锁，只有拿到锁的实例去查询数据库并回写缓存
	lockKey := "lock:shop_type:list"
	ok, lockValue := utils.TryLockWithTTL(ctx, dao.Redis, lockKey, 5*time.Second)
	if !ok {
		maxWait := 1000 * time.Millisecond
		base := 50 * time.Millisecond
		waited := time.Duration(0)
		for waited <= maxWait {
			shopTypes, err = dao.GetShopTypeListCache(ctx, dao.Redis)
			if err == nil {
				return utils.SuccessResultWithData(shopTypes)
			}
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
		// 兜底去 DB 直接读取
		shopTypes, err = dao.GetShopTypeList(ctx, dao.DB)
		if err != nil {
			return utils.ErrorResult("查询失败: " + err.Error())
		}
		return utils.SuccessResultWithData(shopTypes)
	}
	// 拿到分布式锁，确保释放
	defer func() {
		if ok := utils.UnLockSafe(ctx, dao.Redis, lockKey, lockValue); !ok {
			log.Printf("警告: 释放 shop_type:list 锁失败")
		}
	}()

	// 再次双重检查缓存，防止重复 DB 查询
	shopTypes, err = dao.GetShopTypeListCache(ctx, dao.Redis)
	if err == nil {
		return utils.SuccessResultWithData(shopTypes)
	}

	// 真正从数据库读取并回写缓存
	shopTypes, err = dao.GetShopTypeList(ctx, dao.DB)
	if err != nil {
		return utils.ErrorResult("查询失败: " + err.Error())
	}
	if err = dao.SetShopTypeListCache(ctx, dao.Redis, shopTypes); err != nil {
		log.Printf("警告: 设置 ShopType 缓存失败: %v", err)
	}

	return utils.SuccessResultWithData(shopTypes)
}
