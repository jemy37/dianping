package handler

import (
	"dianping/service"
	"dianping/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateBlog 创建博客
// EN: Create a blog post
func CreateBlog(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Images  string `json:"images"`
		ShopId  uint   `json:"shopId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result := service.CreateBlog(c.Request.Context(), userID.(uint), req.Title, req.Content, req.Images, req.ShopId)
	utils.Response(c, result)
}

// LikeBlog 点赞博客
// EN: Like a blog post
func LikeBlog(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	blogIdStr := c.Param("id")
	blogId, err := strconv.ParseUint(blogIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的博客ID")
		return
	}

	result := service.LikeBlog(c.Request.Context(), userID.(uint), uint(blogId))
	utils.Response(c, result)
}

// GetBlogList 获取博客列表
// EN: Get blog list (paginated)
func GetBlogList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetBlogList(c.Request.Context(), page, size)
	utils.Response(c, result)
}

// GetBlogOfFollow 获取关注用户的博客列表
// EN: Get blogs from followed users (inbox/feed)
func GetBlogOfFollow(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	lastId, _ := strconv.Atoi(c.DefaultQuery("lastId", "0"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "10"))
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))

	result := service.GetBlogOfFollow(c.Request.Context(), userID.(uint), lastId, offset, count)
	utils.Response(c, result)
}

// GetBlogById 根据ID获取博客
// EN: Get a blog by ID
func GetBlogById(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "无效的博客ID")
        return
    }

    // 可未登录访问：若未登录则按未点赞处理
    // EN: Allow unauthenticated access; if no user, treat as not liked
    var uid uint = 0
    if userID, exists := c.Get("userID"); exists {
        uid = userID.(uint)
    }

    result := service.GetBlogById(c.Request.Context(), uint(id), uid)
    utils.Response(c, result)
}

// GetHotBlogList 获取热门博客列表
// EN: Get hot blogs sorted by likes
func GetHotBlogList(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

    // 可未登录访问：若未登录则按未点赞处理
    // EN: Allow unauthenticated; if no user, treat as not liked
    var uid uint = 0
    if userID, exists := c.Get("userID"); exists {
        uid = userID.(uint)
    }

    result := service.GetHotBlogList(c.Request.Context(), page, size, uid)
    utils.Response(c, result)
}

// GetMyBlogList 获取我的博客列表
// EN: Get my blogs
func GetMyBlogList(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
        return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetMyBlogList(c.Request.Context(), userID.(uint), page, size)
	utils.Response(c, result)
}

// GetBlogOfShop 获取指定商铺的博客列表
// EN: Get blogs for a specific shop
func GetBlogOfShop(c *gin.Context) {
    shopIdStr := c.Param("id")
    sid, err := strconv.ParseUint(shopIdStr, 10, 32)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "无效的商铺ID")
        return
    }
    page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

    // 可未登录访问：若未登录则按未点赞处理
    var uid uint = 0
    if userID, exists := c.Get("userID"); exists {
        uid = userID.(uint)
    }

    result := service.GetBlogsByShop(c.Request.Context(), uint(sid), page, size, uid)
    utils.Response(c, result)
}
