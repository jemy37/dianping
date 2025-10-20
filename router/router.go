package router

import (
	"dianping/handler"
	"dianping/utils"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 添加中间件
	r.Use(utils.CORSMiddleware())
	r.Use(utils.LoggerMiddleware())
	r.Use(utils.UVStatMiddleware()) // UV统计中间件

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关路由
		userGroup := api.Group("/user")
		{
			userGroup.POST("/code", handler.SendCode)                               //发送验证码√
			userGroup.POST("/register", handler.UserRegister)                       //用户注册√
			userGroup.POST("/login", handler.UserLogin)                             // 用户登录√
			userGroup.POST("/logout", handler.UserLogout)                           // 用户登出
			userGroup.GET("/me", utils.JWTMiddleware(), handler.GetUserInfo)        //获取个人信息√
			userGroup.PUT("/update", utils.JWTMiddleware(), handler.UpdateUserInfo) // 更新个人信息√
			userGroup.POST("/sign", utils.JWTMiddleware(), handler.Sign)            // 签到
		}

		// 商铺相关路由
		shopGroup := api.Group("/shop")
		{
			shopGroup.GET("/list", handler.GetShopList)                                 // 获取商铺列表√
			shopGroup.GET("/:id", handler.GetShopById)                                  // 通过商铺ID 获取商铺信息√
			shopGroup.GET("/of/type", handler.GetShopByType)                            // 根据类型获取商铺√
			shopGroup.GET("/of/name", handler.GetShopByName)                            // 根据名称搜索商铺√
			shopGroup.POST("", handler.SaveShop)                                        // 新增商铺
			shopGroup.PUT("/update", handler.UpdateShop)                                // 更新商铺√
			shopGroup.GET("/:id/nearby", utils.JWTMiddleware(), handler.GetNearbyShops) // 获取某个商铺附近的商铺
		}

		// 商铺类型相关路由
		shopTypeGroup := api.Group("/shop-type")
		{
			shopTypeGroup.GET("/list", handler.GetShopTypeList)
		}

		// 优惠券相关路由
		voucherGroup := api.Group("/voucher")
		{
			voucherGroup.GET("/list/:shopId", handler.GetVoucherList)   // 根据商铺ID信息获取优惠券列表
			voucherGroup.POST("", handler.AddVoucher)                   // TODO: 实现新增普通券功能
			voucherGroup.POST("/seckill", handler.AddSeckillVoucher)    // 新增秒杀券√
			voucherGroup.GET("/seckill/:id", handler.GetSeckillVoucher) // 获取秒杀券详情√
		}

		// 优惠券订单相关路由
		voucherOrderGroup := api.Group("/voucher-order")
		{
			voucherOrderGroup.POST("/seckill/:id", utils.JWTMiddleware(), handler.SeckillVoucher) // 秒杀优惠券√
		}

		// 博客相关路由
		blogGroup := api.Group("/blog")
		{
			blogGroup.POST("", utils.JWTMiddleware(), handler.CreateBlog)               // 创建博客√
			blogGroup.PUT("/like/:id", utils.JWTMiddleware(), handler.LikeBlog)         // 给博客点赞√
			blogGroup.GET("/hot", handler.GetHotBlogList)                               // 获取热门博客列表
			blogGroup.GET("/of/me", utils.JWTMiddleware(), handler.GetMyBlogList)       // 获取我的博客列表√
			blogGroup.GET("/:id", handler.GetBlogById)                                  // 通过ID获取博客
			blogGroup.GET("/of/follow", utils.JWTMiddleware(), handler.GetBlogOfFollow) // 获取关注用户的博客列表√
		}

		// 关注相关路由
		followGroup := api.Group("/follow")
		{
			followGroup.POST("/:id", utils.JWTMiddleware(), handler.Follow)
			followGroup.DELETE("/:id", utils.JWTMiddleware(), handler.Unfollow)
			followGroup.GET("/common/:id", utils.JWTMiddleware(), handler.GetCommonFollows)
		}

		// 统计相关路由
		statGroup := api.Group("/stat")
		{
			statGroup.GET("/uv/today", handler.GetTodayUV)     // 获取今日UV
			statGroup.GET("/uv/daily", handler.GetDailyUV)     // 获取指定日期UV
			statGroup.GET("/uv/range", handler.GetUVRange)     // 获取日期范围UV
			statGroup.GET("/uv/recent", handler.GetRecentUV)   // 获取最近N天UV
			statGroup.GET("/uv/summary", handler.GetUVSummary) // 获取UV统计摘要
		}

		pprofGroup := api.Group("/debug/pprof")
		{
			pprofGroup.GET("/", func(c *gin.Context) {
				pprof.Index(c.Writer, c.Request)
			})
			pprofGroup.GET("/cmdline", func(c *gin.Context) {
				pprof.Cmdline(c.Writer, c.Request)
			})
			pprofGroup.GET("/profile", func(c *gin.Context) {
				pprof.Profile(c.Writer, c.Request)
			})
			pprofGroup.POST("/symbol", func(c *gin.Context) {
				pprof.Symbol(c.Writer, c.Request)
			})
			pprofGroup.GET("/symbol", func(c *gin.Context) {
				pprof.Symbol(c.Writer, c.Request)
			})
			pprofGroup.GET("/trace", func(c *gin.Context) {
				pprof.Trace(c.Writer, c.Request)
			})
		}
	}

	// 健康检查
	r.GET("/health", handler.HealthCheck)

	return r
}
