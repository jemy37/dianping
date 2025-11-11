package service

import (
	"context"
	"dianping/dao"
	"dianping/models"
	"dianping/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// SeckillVoucher 秒杀优惠券
// EN: Seckill purchase entry. Runs Lua for stock/user checks and publishes to Redis Stream.
func SeckillVoucher(ctx context.Context, userId, voucherId uint) *utils.Result {
	// 从文件当中加载脚本（缓存已读内容）
	script, err := os.ReadFile("script/seckill.lua")
	if err != nil {
		log.Printf("读取秒杀脚本失败: %v", err)
		return utils.ErrorResult("系统错误")
	}
	scriptStr := string(script)

	// 生成临时 orderId 并传入 Lua 脚本以便 xadd 中包含 id 字段
	var orderId string
	if idWorker == nil {
		idWorker = utils.NewRedisIdWorker(dao.Redis, 16)
	}
	if id, err := idWorker.NextId(ctx, "order"); err == nil {
		orderId = strconv.FormatInt(id, 10)
	} else {
		log.Printf("生成orderId失败: %v", err)
		orderId = ""
		//test
		println(orderId)
	}

	// 1. 执行Lua脚本
	result := dao.Redis.Eval(ctx, scriptStr, []string{}, strconv.Itoa(int(voucherId)), strconv.Itoa(int(userId)), orderId)
	if result.Err() != nil {
		log.Printf("执行秒杀脚本失败: %v", result.Err())
		return utils.ErrorResult("系统错误")
	}

	// 2. 判断结果是否为 0，0的时候有资格完成
	r, err := result.Int()
	if err != nil {
		log.Printf("获取秒杀脚本返回值失败: %v", err)
		return utils.ErrorResult("系统错误")
	}
	if r != 0 {
		if r == 1 {
			return utils.ErrorResult("库存不足")
		}
		return utils.ErrorResult("不能重复购买")
	}

	// 3. 已经加入到消息队列了

	// 4. 返回订单ID（这里可以生成一个临时ID或者返回成功信息）
	return utils.SuccessResultWithData("秒杀成功，订单处理中...")
}

// StreamOrderInfo Redis Stream中的订单信息结构体
// EN: Order info payload structure stored in Redis Stream
type StreamOrderInfo struct {
	UserID    string `json:"userId"`
	VoucherID string `json:"voucherId"`
	OrderID   string `json:"id"`
}

// Stream消费者相关配置
// EN: Redis Stream consumer configuration
var (
	streamKey     = "stream.orders"     // Stream名称
	groupName     = "order-group"       // 消费者组名称
	consumerCount = 3                   // 消费者数量
	streamOnce    sync.Once             // 确保Stream只初始化一次
	stopChan      = make(chan struct{}) // 停止信号
	wg            sync.WaitGroup        // 等待组，用于优雅关闭
	idWorker      *utils.RedisIdWorker
)

// InitStreamConsumer 初始化Redis Stream消费者
// EN: Initialize Redis Stream consumers (group + workers)
func InitStreamConsumer() error {
	var initErr error
	streamOnce.Do(func() {
		ctx := context.Background()

		// 1. 检查Stream是否存在，如果不存在则创建
		exists, err := checkStreamExists(ctx, streamKey)
		if err != nil {
			initErr = fmt.Errorf("检查Stream失败: %v", err)
			return
		}

		if !exists {
			// 创建一个空的Stream（通过添加临时消息然后删除）
			result := dao.Redis.XAdd(ctx, &redis.XAddArgs{
				Stream: streamKey,
				ID:     "*",
				Values: map[string]interface{}{"init": "temp"},
			})
			if result.Err() != nil {
				initErr = fmt.Errorf("创建Stream失败: %v", result.Err())
				return
			}
			// 删除临时消息
			dao.Redis.XDel(ctx, streamKey, result.Val())
		}

		// 2. 创建消费者组（如果不存在）
		err = dao.Redis.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()
		if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
			initErr = fmt.Errorf("创建消费者组失败: %v", err)
			return
		}

		// 3. 启动消费者
		for i := 0; i < consumerCount; i++ {
			consumerName := fmt.Sprintf("consumer-%d", i)
			wg.Add(1)
			go streamConsumer(consumerName, i)
		}

		log.Printf("Redis Stream消费者初始化完成，Stream: %s, 消费者组: %s, 消费者数量: %d",
			streamKey, groupName, consumerCount)
	})

	return initErr
}

// checkStreamExists 检查Stream是否存在
// EN: Check if the Stream key exists in Redis
func checkStreamExists(ctx context.Context, streamKey string) (bool, error) {
	result := dao.Redis.Exists(ctx, streamKey)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val() > 0, nil
}

// streamConsumer Stream消费者worker
// EN: Worker loop that reads and processes Stream messages
func streamConsumer(consumerName string, workerID int) {
	defer wg.Done()

	log.Printf("Stream消费者 %s (Worker %d) 启动", consumerName, workerID)

	ctx := context.Background()

	for {
		select {
		case <-stopChan:
			log.Printf("Stream消费者 %s (Worker %d) 收到停止信号，正在退出", consumerName, workerID)
			return
		default:
			// 从Stream中读取消息
			messages, err := readStreamMessages(ctx, consumerName)
			if err != nil {
				log.Printf("消费者 %s 读取消息失败: %v", consumerName, err)
				time.Sleep(time.Second * 2) // 出错时等待2秒再重试
				continue
			}

			// 处理每条消息
			for _, msg := range messages {
				err := processStreamMessage(ctx, msg, consumerName)
				if err != nil {
					log.Printf("消费者 %s 处理消息失败: msgID=%s, error=%v",
						consumerName, msg.ID, err)
				} else {
					log.Printf("消费者 %s 成功处理消息: msgID=%s", consumerName, msg.ID)
					// 确认消息已处理
					dao.Redis.XAck(ctx, streamKey, groupName, msg.ID)
				}
			}

			// 如果没有消息，短暂休眠
			if len(messages) == 0 {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

// readStreamMessages 从Stream中读取消息
// EN: Read pending first, then new messages from the Stream
func readStreamMessages(ctx context.Context, consumerName string) ([]redis.XMessage, error) {
	// 首先尝试读取pending消息（之前未确认的消息）
	pendingResult := dao.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamKey, "0"}, // "0"表示读取pending消息
		Count:    10,
		Block:    0, // 不阻塞
	})

	if pendingResult.Err() == nil && len(pendingResult.Val()) > 0 && len(pendingResult.Val()[0].Messages) > 0 {
		return pendingResult.Val()[0].Messages, nil
	}

	// 如果没有pending消息，读取新消息
	result := dao.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamKey, ">"}, // ">"表示读取新消息
		Count:    10,
		Block:    time.Second * 1, // 阻塞1秒
	})

	if result.Err() != nil {
		// 阻塞读取超时等情况，go-redis 可能返回 redis.Nil，视为无消息而不是错误
		if result.Err() == redis.Nil {
			return []redis.XMessage{}, nil
		}
		return nil, result.Err()
	}

	if len(result.Val()) > 0 && len(result.Val()[0].Messages) > 0 {
		return result.Val()[0].Messages, nil
	}

	return []redis.XMessage{}, nil
}

// processStreamMessage 处理单条Stream消息
// EN: Parse and dispatch a single Stream message
func processStreamMessage(ctx context.Context, msg redis.XMessage, consumerName string) error {
	// 解析消息内容
	orderInfo, err := parseOrderMessage(msg)
	if err != nil {
		return fmt.Errorf("解析消息失败: %v", err)
	}

	// 转换字符串ID为uint
	userID, err := strconv.ParseUint(orderInfo.UserID, 10, 32)
	if err != nil {
		return fmt.Errorf("解析用户ID失败: %v", err)
	}

	voucherID, err := strconv.ParseUint(orderInfo.VoucherID, 10, 32)
	if err != nil {
		return fmt.Errorf("解析优惠券ID失败: %v", err)
	}

	// 处理订单
	return processStreamOrder(ctx, uint(userID), uint(voucherID), orderInfo.OrderID)
}

// parseOrderMessage 解析订单消息
// EN: Parse order fields from Stream message values
func parseOrderMessage(msg redis.XMessage) (*StreamOrderInfo, error) {
	orderInfo := &StreamOrderInfo{}

	// 从消息中提取字段
	if userID, ok := msg.Values["userId"].(string); ok {
		orderInfo.UserID = userID
	} else {
		return nil, fmt.Errorf("消息中缺少userId字段")
	}

	if voucherID, ok := msg.Values["voucherId"].(string); ok {
		orderInfo.VoucherID = voucherID
	} else {
		return nil, fmt.Errorf("消息中缺少voucherId字段")
	}

	if orderID, ok := msg.Values["orderId"].(string); ok {
		orderInfo.OrderID = orderID
	} else {
		return nil, fmt.Errorf("消息中缺少orderId字段")
	}

	return orderInfo, nil
}

// processStreamOrder 处理Stream中的订单
// EN: Transactionally check idempotency, decrement stock and create order
func processStreamOrder(ctx context.Context, userID, voucherID uint, orderID string) error {
	// 开始数据库事务
	tx := dao.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("订单处理发生panic: %v", r)
		}
	}()

	// 1) 检查用户是否已经有秒杀券订单（防止重复）
	exists, err := dao.CheckSeckillVoucherOrderExists(ctx, tx, userID, voucherID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("检查重复订单失败: %v", err)
	}
	if exists {
		// 用户已下单，直接回滚事务并返回 nil（视为已处理）
		tx.Rollback()
		log.Printf("用户已存在秒杀订单: userID=%d, voucherID=%d", userID, voucherID)
		return nil
	}

	// 2) 在事务内扣减秒杀券库存（保证与订单创建在同一事务）
	// Use raw SQL expression for atomic decrement
	result := tx.Model(&models.SeckillVoucher{}).
		Where("voucher_id = ? AND stock >= ?", voucherID, 1).
		UpdateColumn("stock", gorm.Expr("stock - ?", 1))
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("更新库存失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("库存不足")
	}

	// 3) 同步扣减关联的普通券库存（tb_voucher）——保证与上面操作在同一事务中
	vResult := tx.Model(&models.Voucher{}).
		Where("id = ? AND stock >= ?", voucherID, 1).
		UpdateColumn("stock", gorm.Expr("stock - ?", 1))
	if vResult.Error != nil {
		tx.Rollback()
		return fmt.Errorf("更新关联券库存失败: %v", vResult.Error)
	}
	if vResult.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("关联券库存不足")
	}

	// 创建订单
	now := time.Now()
	order := &models.VoucherOrder{
		UserID:      userID,
		VoucherID:   voucherID,
		PayType:     1,
		Status:      1,
		CreateTime:  &now,
		VoucherType: 2, // 秒杀券类型
	}

	// 如果Lua脚本提供了订单ID，可以使用它
	if orderID != "" {
		// 尝试将 string orderID 解析为无符号整数并保存到模型的 OrderID 字段
		if id64, err := strconv.ParseUint(orderID, 10, 64); err == nil {
			order.OrderID = uint(id64)
		} else {
			// 解析失败则记录警告，但不阻止下单
			log.Printf("警告: 解析 orderId (%s) 为整数失败: %v", orderID, err)
		}
	}

	// 创建订单记录
	if err := dao.CreateVoucherOrder(ctx, tx, order); err != nil {
		tx.Rollback()
		return fmt.Errorf("创建订单失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("提交事务失败: %v", err)
	}

	// 记录创建成功：包含 DB 自增主键和原始 orderId 字符串（如果有）
	log.Printf("成功创建订单: userID=%d, voucherID=%d, dbID=%d, orderID=%s",
		userID, voucherID, order.ID, orderID)

	return nil
}

// StopStreamConsumers 停止所有Stream消费者（用于优雅关闭）
// EN: Gracefully stop all Stream consumer workers
func StopStreamConsumers() {
	log.Println("正在停止Stream消费者...")
	close(stopChan)
	wg.Wait()
	log.Println("所有Stream消费者已停止")
}

// GetStreamInfo 获取Stream状态信息（用于监控）
// EN: Get Stream/groups/consumers info for monitoring
func GetStreamInfo() (map[string]interface{}, error) {
	ctx := context.Background()

	// 获取Stream基本信息
	streamInfo, err := dao.Redis.XInfoStream(ctx, streamKey).Result()
	if err != nil {
		return nil, fmt.Errorf("获取Stream信息失败: %v", err)
	}

	// 获取消费者组信息
	groupInfo, err := dao.Redis.XInfoGroups(ctx, streamKey).Result()
	if err != nil {
		return nil, fmt.Errorf("获取消费者组信息失败: %v", err)
	}

	// 获取消费者信息
	consumerInfo, err := dao.Redis.XInfoConsumers(ctx, streamKey, groupName).Result()
	if err != nil {
		return nil, fmt.Errorf("获取消费者信息失败: %v", err)
	}

	return map[string]interface{}{
		"stream":    streamInfo,
		"groups":    groupInfo,
		"consumers": consumerInfo,
	}, nil
}
