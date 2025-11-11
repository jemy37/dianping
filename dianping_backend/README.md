
## 简介 (Chinese)

一个基于 Gin + GORM + Redis 的点评系统示例，涵盖用户、商铺、博客、关注、优惠券（含秒杀）等功能。项目实现了缓存击穿/穿透/雪崩的防护、Redis GEO 搜索、Redis Bloom 布隆过滤器、Redis Stream 消费者用于异步下单等实践。

### 项目结构

```
dianping/
├── api-test/                 # HTTP 测试脚本
├── config/                   # 配置
│   ├── config.go             # 配置解析
│   └── application.yaml      # 默认配置
├── dao/                      # 数据访问层 (DB/Redis)
├── handler/                  # 请求处理层 (Controller)
├── models/                   # 数据模型
├── router/                   # 路由
├── script/                   # Lua 脚本 (秒杀)
├── service/                  # 业务层
├── utils/                    # 工具库/中间件
├── main.go                   # 启动入口
└── README.md                 # 项目说明
```

### 技术栈

- Web: Gin
- 数据库: MySQL + GORM
- 缓存/消息: Redis (含 GEO、HyperLogLog、Stream)
- 防护: RedisBloom 布隆过滤器（redis/redis-stack 镜像）
- 认证: JWT
- 配置: YAML

### 功能模块概览

- 用户：注册、验证码登录、信息查询/更新、签到（位图）
- 商铺：详情、分页、按类型、名称搜索、附近搜索（GEO）、创建/更新
- 博客：创建、点赞、热门列表、我的列表、关注人动态（Feed）
- 关注：关注/取关/共同关注
- 优惠券：普通券创建/查询；秒杀券创建/查询/下单（Lua+Stream）

### 快速开始

1) 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+（推荐 redis/redis-stack:latest 以启用 RedisBloom 模块）

使用 Docker 启动 Redis Stack（包含 Bloom/JSON/TimeSeries 等模块）:

```
docker run -d --name redis-stack -p 6379:6379 redis/redis-stack:latest
```

2) 安装依赖

```
go mod tidy
```

3) 配置修改

编辑 `config/application.yaml` 设置数据库、Redis、JWT 等配置。生产环境请覆盖 `jwt.secret` 且勿将明文密码提交版本库。

4) 初始化与运行

```
go run main.go
```

服务默认监听 `http://localhost:8080`，程序会自动进行 GORM 迁移，并初始化 Redis 连接、Bloom 过滤器、Stream 消费者、商铺 GEO 索引缓存等。

### 主要接口（节选）

 - `POST /api/user/code?phone=` 发送验证码
  - `POST /api/user/register` 注册
  - `POST /api/user/login` 登录
  - `POST /api/user/login/password` 密码登录（支持手机号或昵称）
  - `GET /api/user/me` 当前用户信息（鉴权）
  - `PUT /api/user/update` 更新信息（鉴权）
- 商铺：
- 博客：
  - `POST /api/blog` 创建（鉴权）
  - `PUT /api/blog/like/:id` 点赞（鉴权）
  - `GET /api/blog` 全部博客（分页，可未登录）
  - `GET /api/blog/hot` 热门（可未登录）
  - `GET /api/blog/of/me` 我的（鉴权）
  - `GET /api/blog/of/shop/:id` 按商铺获取博客（可未登录）
  - `GET /api/blog/:id` 详情（可未登录）
  - `GET /api/blog/of/follow` 关注动态（鉴权）
- 优惠券：
- 流式下单：Lua 校验 + Redis Stream 消费者组处理订单
- 指标统计：HyperLogLog UV 统计中间件
### 前端（React + Vite）
本仓库的同级目录下提供了示例前端：`../dianping-frontend`
- 技术栈：React 18 + Vite + React Router + Axios
- 开发端口：`http://localhost:5173`（通过 Vite 代理转发 `/api` 到后端 `http://localhost:8080`）
  - 若后端端口变更，请修改前端 `vite.config.ts` 中的 `server.proxy['/api'].target`
运行步骤：
```
cd ../dianping-frontend
npm install
npm run dev
```
前端示例包含：
- 首页：全部商铺列表（按类型 sort 排序），点击进入详情
- 类型页：按固定 typeId（1/2/5）展示“烧烤/咖啡/火锅”，点击查看该类型商铺
- 商铺详情：商铺信息、优惠券、与本店相关的博客列表，并支持登录后在此页发博文
- 博客：默认展示全部博客；可切换“热门”按点赞数排序；支持查看博客详情
- 登录：支持“手机号+验证码”或“手机号/昵称 + 密码”两种方式


### 已知改进点 / TODO

- 登出目前为前端丢弃 token，可选实现 token 黑名单
- 接口鉴权一致性已修正（博客热门/详情支持未登录访问）
- 生产建议：
  - 移除配置中的默认 secret，使用环境变量或密管
  - 完善日志/追踪/告警；Stream 消费失败重试与死信队列
  - 为 Lua/Stream/DB 操作补充更细的监控指标

---

## Overview (English)

This is a demo “dianping” system built with Gin, GORM and Redis. It implements users, shops, blogs, follow, vouchers (including seckill). It demonstrates cache protection (breakdown/penetration), Redis GEO search, Redis Bloom filters, and Redis Stream consumers for async ordering.

### Structure

See the tree above. Key layers: router, handler, service, dao, models, utils, scripts.

### Stack

- Gin (HTTP), GORM (ORM)
- MySQL, Redis (GEO/HyperLogLog/Stream)
- RedisBloom (via redis/redis-stack) for Bloom filters
- JWT authentication, YAML config

### Features

- Users: register, login via code, profile, daily sign-in (bitmap)
- Shops: detail, pagination, by type/name, nearby via GEO, create/update
- Blogs: create, like, hot list, mine, follow feed
- Follow: follow/unfollow/common-follows
- Vouchers: normal voucher create/list; seckill create/detail/purchase (Lua + Stream)

### Quick Start

1) Requirements: Go 1.21+, MySQL 8+, Redis 6+ (prefer redis/redis-stack for Bloom)

2) Dependencies: `go mod tidy`

3) Config: edit `config/application.yaml` to set DB/Redis/JWT. Use a strong JWT secret in production.

4) Run: `go run main.go`

The app bootstraps DB migrations, Redis clients, Bloom filters, Stream consumers and GEO caches.

### Selected APIs

- Users:
  - `POST /api/user/code?phone=` Send login code
  - `POST /api/user/register` Register
  - `POST /api/user/login` Login with phone + code
  - `POST /api/user/login/password` Login with password (phone or nickname)
  - `GET /api/user/me` Current user (auth)
  - `PUT /api/user/update` Update profile (auth)
- Shops:
  - `GET /api/shop/list` List
  - `GET /api/shop/:id` Detail
  - `GET /api/shop/of/type` By type (pagination)
  - `GET /api/shop/of/name` Search by name
  - `POST /api/shop/createShop` Create
  - `PUT /api/shop/update` Update
  - `GET /api/shop/:id/nearby` Nearby shops (auth)
- Blogs:
  - `POST /api/blog` Create (auth)
  - `PUT /api/blog/like/:id` Like/unlike (auth)
  - `GET /api/blog` All blogs (pagination)
  - `GET /api/blog/hot` Hot blogs (by likes)
  - `GET /api/blog/of/me` My blogs (auth)
  - `GET /api/blog/of/shop/:id` Blogs of a shop
  - `GET /api/blog/:id` Blog detail
  - `GET /api/blog/of/follow` Follow feed (auth)
- Vouchers:
  - `POST /api/voucher` Create normal voucher
  - `GET /api/voucher/list/:shopId` List vouchers of shop
  - `POST /api/voucher/seckill` Create seckill voucher
  - `GET /api/voucher/seckill/:id` Seckill voucher detail
  - `POST /api/voucher-order/seckill/:id` Join seckill (auth)

### Frontend (React + Vite)

There is a sample frontend at the sibling folder: `../dianping-frontend`.

- Stack: React 18, Vite, React Router, Axios
- Dev server: `http://localhost:5173` with proxy for `/api` → `http://localhost:8080`
  - If backend port differs, change `vite.config.ts` `server.proxy['/api'].target`.

Run:

```
cd ../dianping-frontend
npm install
npm run dev
```

Included pages:
- Home: all shops, sorted by their type's sort value; click to shop detail
- Types: fixed mapping by `typeId` (1=BBQ, 2=Cafe, 5=Hotpot), click to list shops of that type
- Shop detail: shop info, vouchers, related blogs; post a blog for the shop when logged in
- Blogs: list all blogs with toggle to “Hot”; blog detail route
- Login: supports phone+code, and phone/nickname + password

### Highlights

- Clear layering and unified responses
- Robust caching and Bloom filters to mitigate DB pressure
- GEO-based “nearby shops”
- Seckill with Lua checks and Redis Stream consumer group
- HLL-based UV middleware

### Notes / TODO

- Logout is stateless; implement token blacklist if needed
- Prefer secrets via env/secret manager
- Add retries/DLQ for Stream processing in production
