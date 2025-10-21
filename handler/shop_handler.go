package handler

import (
	"dianping/models"
	"dianping/service"
	"dianping/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetShopById 根据ID获取商铺
func GetShopById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的商铺ID")
		return
	}

	result := service.GetShopById(c.Request.Context(), uint(id))
	utils.Response(c, result)
}

// GetShopList 获取商铺列表
func GetShopList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetShopList(page, size)
	utils.Response(c, result)
}

// GetShopByType 根据类型获取商铺
func GetShopByType(c *gin.Context) {
	typeIdStr := c.Query("typeId")
	if typeIdStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "类型ID不能为空")
		return
	}

	typeId, err := strconv.ParseUint(typeIdStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的类型ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetShopByType(uint(typeId), page, size)
	utils.Response(c, result)
}

// GetShopByName 根据名称搜索商铺
func GetShopByName(c *gin.Context) {
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result := service.GetShopByName(name, page, size)
	utils.Response(c, result)
}

// SaveShop 新增商铺
func SaveShop(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		TypeID uint   `json:"typeId"`
		// TypeName  string  `json:"typeName"`
		TypeIcon  string  `json:"typeIcon"`
		Images    string  `json:"images"`
		Area      string  `json:"area"`
		Address   string  `json:"address"`
		X         float64 `json:"x"`
		Y         float64 `json:"y"`
		AvgPrice  int     `json:"avgPrice"`
		OpenHours string  `json:"openHours"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数校验失败: "+err.Error())
		return
	}

	shop := &models.Shop{
		Name:      req.Name,
		TypeID:    req.TypeID,
		Images:    req.Images,
		Area:      req.Area,
		Address:   req.Address,
		X:         req.X,
		Y:         req.Y,
		AvgPrice:  req.AvgPrice,
		OpenHours: req.OpenHours,
	}

	result := service.CreateShopWithType(c.Request.Context(), shop, req.TypeIcon)
	utils.Response(c, result)
}

// UpdateShop 更新商铺信息
func UpdateShop(c *gin.Context) {
	// 1. 参数校验
	var shop models.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数校验失败")
		return
	}

	// 2. 更新商铺
	result := service.UpdateShopById(c.Request.Context(), &shop)
	utils.Response(c, result)
}

// GetNearbyShops 获取某个店铺的附近某个距离的所有点
func GetNearbyShops(c *gin.Context) {
	// 1. 参数校验
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的商铺ID")
		return
	}

	radius, err := strconv.ParseFloat(c.DefaultQuery("radius", "1.0"), 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的半径")
		return
	}

	count, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的数量")
		return
	}

	// 2. 查询附近的商铺
	result := service.GetNearbyShops(c.Request.Context(), uint(id), radius, count)
	utils.Response(c, result)
}
