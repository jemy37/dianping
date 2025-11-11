package handler

import (
	"dianping/service"
	"dianping/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册
// EN: Register a new user
func UserRegister(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		NickName string `json:"nickName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result := service.UserRegister(req.Phone, req.Code, req.Password, req.NickName)
	utils.Response(c, result)
}

// UserLogin 用户登录
// EN: Login with phone and code
func UserLogin(c *gin.Context) {
	// 这里只能支持验证码登录
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 判断手机号格式是否正确
	if ok := utils.IsPhoneValid(req.Phone); !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, "手机号格式不正确")
		return
	}

	// 判断验证码格式是否正确
	if utils.IsCodeInvalid(req.Code) {
		utils.ErrorResponse(c, http.StatusBadRequest, "验证码格式不正确")
		return
	}

	result := service.UserLogin(req.Phone, req.Code)
	utils.Response(c, result)
}

// UserPasswordLogin 用户密码登录（手机号或昵称）
// EN: Login with password using phone or nickname
func UserPasswordLogin(c *gin.Context) {
    var req struct {
        Phone    string `json:"phone"`
        NickName string `json:"nickName"`
        Password string `json:"password" binding:"required,min=6"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
        return
    }
    if req.Phone == "" && req.NickName == "" {
        utils.ErrorResponse(c, http.StatusBadRequest, "需要提供手机号或昵称")
        return
    }
    var result *utils.Result
    if req.NickName != "" {
        result = service.UserPasswordLogin(req.NickName, req.Password, true)
    } else {
        // 若提供 phone，校验手机号格式
        if ok := utils.IsPhoneValid(req.Phone); !ok {
            utils.ErrorResponse(c, http.StatusBadRequest, "手机号格式不正确")
            return
        }
        result = service.UserPasswordLogin(req.Phone, req.Password, false)
    }
    utils.Response(c, result)
}

// GetUserInfo 获取用户信息
// EN: Get current user info
func GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	result := service.GetUserInfo(userID.(uint))
	utils.Response(c, result)
}

// UpdateUserInfo 更新用户信息
// EN: Update current user profile
func UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		NickName string `json:"nickName"`
		Icon     string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result := service.UpdateUserInfo(userID.(uint), req.NickName, req.Icon)
	utils.Response(c, result)
}

// UserLogout 用户登出
// EN: Logout (stateless JWT; client discards token)
func UserLogout(c *gin.Context) {
	// TODO: 实现登出逻辑，清除token等
	utils.SuccessResponse(c, "登出成功")
}

// SendCode 发送验证码
// EN: Send login verification code
func SendCode(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "手机号不能为空")
		return
	}

	if !utils.IsPhoneValid(phone) {
		utils.ErrorResponse(c, http.StatusBadRequest, "手机号格式不正确")
		return
	}
	result := service.SendCode(phone)
	utils.Response(c, result)
}

// Sign 用户签到
// EN: Daily user sign-in
func Sign(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	result := service.Sign(c.Request.Context(), userID.(uint))
	utils.Response(c, result)
}

// CheckSign 获取用户签到状态
func CheckSign(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	month := c.Query("month")
	if month == "" {
		// 以当前月为准
		month = time.Now().Format("2006-01")
	}

	result := service.CheckSign(c.Request.Context(), userID.(uint), month)
	utils.Response(c, result)
}
