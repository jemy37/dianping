package handler

import (
	"dianping/service"
	"dianping/utils"

	"github.com/gin-gonic/gin"
)

// GetShopTypeList 获取商铺类型列表
func GetShopTypeList(c *gin.Context) {
	result := service.GetShopTypeList(c.Request.Context())
	utils.Response(c, result)
}
