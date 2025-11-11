package dao

import (
	"context"
	"dianping/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// CreateBlog 创建博客
// EN: Create blog record
func CreateBlog(ctx context.Context, blog *models.Blog) error {
	return DB.WithContext(ctx).Create(blog).Error
}

// GetBlogByID 根据ID获取博客
// EN: Get blog by ID
func GetBlogByID(ctx context.Context, id uint) (*models.Blog, error) {
	var blog models.Blog
	err := DB.WithContext(ctx).First(&blog, id).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// GetBlogList 获取博客列表（分页）
// EN: Get blog list with pagination
func GetBlogList(ctx context.Context, offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	// 获取总数
	if err := DB.WithContext(ctx).Model(&models.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := DB.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at desc").Find(&blogs).Error
	return blogs, total, err
}

// GetHotBlogList 获取热门博客列表（按点赞数排序）
// EN: Get hot blogs sorted by likes
func GetHotBlogList(ctx context.Context, offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	// 获取总数
	if err := DB.WithContext(ctx).Model(&models.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按点赞数排序
	err := DB.WithContext(ctx).Offset(offset).Limit(limit).Order("liked desc, created_at desc").Find(&blogs).Error
	return blogs, total, err
}

// GetMyBlogList 获取用户的博客列表
// EN: Get blogs created by the user
func GetMyBlogList(ctx context.Context, userID uint, offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	// 获取总数
	if err := DB.WithContext(ctx).Model(&models.Blog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := DB.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at desc").Find(&blogs).Error
	return blogs, total, err
}

// GetBlogLike 检查用户是否已点赞博客
// EN: Get like record for user-blog pair
func GetBlogLike(ctx context.Context, userID, blogID uint) (*models.BlogLike, error) {
	var like models.BlogLike
	err := DB.WithContext(ctx).Where("user_id = ? AND blog_id = ?", userID, blogID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// CreateBlogLike 创建博客点赞
// EN: Create a like record
func CreateBlogLike(ctx context.Context, like *models.BlogLike) error {
	return DB.WithContext(ctx).Create(like).Error
}

// DeleteBlogLike 删除博客点赞
// EN: Delete a like record
func DeleteBlogLike(ctx context.Context, like *models.BlogLike) error {
	return DB.WithContext(ctx).Delete(like).Error
}

// DeleteBlogLikeByUser 删除指定用户对指定博客的点赞记录（数据库）
// EN: Delete like record for user-blog pair
func DeleteBlogLikeByUser(ctx context.Context, userID, blogID uint) error {
	return DB.WithContext(ctx).Where("user_id = ? AND blog_id = ?", userID, blogID).Delete(&models.BlogLike{}).Error
}

// IncrementBlogLiked 增加博客点赞数
// EN: Increment blog like counter
func IncrementBlogLiked(ctx context.Context, blogID uint) error {
	return DB.WithContext(ctx).Model(&models.Blog{}).Where("id = ?", blogID).UpdateColumn("liked", DB.Raw("liked + 1")).Error
}

// DecrementBlogLiked 减少博客点赞数
// EN: Decrement blog like counter
func DecrementBlogLiked(ctx context.Context, blogID uint) error {
	return DB.WithContext(ctx).Model(&models.Blog{}).Where("id = ?", blogID).UpdateColumn("liked", DB.Raw("liked - 1")).Error
}

// GetBlogByIDs 根据博客 ID 获取博客详情
// EN: Batch get blogs by IDs
func GetBlogByIDs(ctx context.Context, blogIDs []uint) ([]models.Blog, error) {
	var blogs []models.Blog
	err := DB.WithContext(ctx).Where("id IN ?", blogIDs).Find(&blogs).Error
	return blogs, err
}

// ======= redis 相关操作 =========

const (
	// 博客点赞集合的键名格式：blog_like:%d
	blogLikeKey = "blog:liked:"
	feedKey     = "feed:"
)

// IsLikedMember 检查用户是否已点赞博客（使用 Redis SortedSet）
// EN: Check like existence via Sorted Set; fallback to DB then backfill Redis
func IsLikedMember(ctx context.Context, rds *redis.Client, userID, blogID uint) (bool, error) {
	// 使用 SortedSet 检查用户是否在点赞集合中，使用 ZRank 查找 member 是否存在
	key := blogLikeKey + strconv.Itoa(int(blogID))
	member := strconv.Itoa(int(userID))

	// 1) 尝试从 Redis 获取
	_, err := rds.ZRank(ctx, key, member).Result()
	if err == nil {
		return true, nil
	}
	// 如果是其他 Redis 错误，返回错误（网络/权限等）
	if err != redis.Nil {
		return false, err
	}

	// 2) redis 没有命中 -> 回退到数据库查询（cache-aside）
	if like, derr := GetBlogLike(ctx, userID, blogID); derr == nil && like != nil {
		// DB 命中，尝试回写到 Redis（best-effort）
		_ = rds.ZAdd(ctx, key, &redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: member,
		}).Err()
		return true, nil
	} else if derr == nil {
		// DB 返回 nil like 且无 error（不太可能，但处理为未命中）
		return false, nil
	} else {
		// 若 DB 返回记录未找到，返回 false；否则返回错误
		if derr == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, derr
	}
}

func RemoveLikedMember(ctx context.Context, rds *redis.Client, userID, blogID uint) error {
    // 从 SortedSet 中移除用户
    // EN: Remove user from the blog's like zset
    return rds.ZRem(ctx, blogLikeKey+strconv.Itoa(int(blogID)), strconv.Itoa(int(userID))).Err()
}

// GetBlogsByShop 获取指定商铺的博客列表（分页）
// EN: List blogs by shop ID with pagination
func GetBlogsByShop(ctx context.Context, shopID uint, offset, limit int) ([]models.Blog, int64, error) {
    var blogs []models.Blog
    var total int64

    if err := DB.WithContext(ctx).Model(&models.Blog{}).Where("shop_id = ?", shopID).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := DB.WithContext(ctx).Where("shop_id = ?", shopID).Offset(offset).Limit(limit).Order("created_at desc").Find(&blogs).Error
    return blogs, total, err
}

func SaveLikedMember(ctx context.Context, rds *redis.Client, userID, blogID uint) error {
    // 向 SortedSet 中添加用户，分数设为当前时间戳
    // EN: Add user to zset with current timestamp as score
    return rds.ZAdd(ctx, blogLikeKey+strconv.Itoa(int(blogID)), &redis.Z{
        Score:  float64(time.Now().Unix()),
        Member: strconv.Itoa(int(userID)),
    }).Err()
}

func GetTopKBloglikedMember(ctx context.Context, rds *redis.Client, blogID uint, k int) ([]string, error) {
    // 从 SortedSet 中获取点赞数最多的 k 个用户
    // EN: Get top-K likers from zset
    return rds.ZRevRange(ctx, blogLikeKey+strconv.Itoa(int(blogID)), 0, int64(k-1)).Result()
}

func FeedToUserRedis(ctx context.Context, rds *redis.Client, userID uint, blogID uint) error {
    // 向用户的 feed 队列中插入id，使用 zset 存储，分数设为当前时间戳
    // EN: Push blog ID to user's feed zset
    return rds.ZAdd(ctx, feedKey+strconv.Itoa(int(userID)), &redis.Z{
        Score:  float64(time.Now().Unix()),
        Member: strconv.Itoa(int(blogID)),
    }).Err()
}

func GetFeedFromUserRedis(ctx context.Context, rds *redis.Client, userID uint, lastId, offset, count int) ([]uint, int64, int, error) {
    // 从用户的 feed 队列中获取博客 ID
    // EN: Read IDs from user's feed zset with score window
    result, err := rds.ZRevRangeByScoreWithScores(ctx, feedKey+strconv.Itoa(int(userID)), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.Itoa(int(lastId)),
		Offset: int64(offset),
		Count:  int64(count),
	}).Result()
	if err != nil {
		return nil, 0, 0, err
	}

	if len(result) == 0 {
		return nil, 0, 0, nil
	}

	var blogIds []uint
	minTime := result[len(result)-1].Score
	offset = 0
	for _, z := range result {
		blogId, _ := strconv.Atoi(z.Member.(string))
		blogIds = append(blogIds, uint(blogId))

		if z.Score == minTime {
			offset++
		}
	}

	return blogIds, int64(minTime), int(offset), nil
}
