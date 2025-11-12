package dao

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
)

const jwtBlacklistKeyPrefix = "auth:jwt:blacklist:"
const jwtTokenStringKeyPrefix = "auth:jwt:blacklist:tk:"

// BlacklistTokenString 将 token 字符串加入黑名单（使用 SHA256 摘要作为 key）
// EN: Blacklist the raw token string using SHA256 digest as key
func BlacklistTokenString(ctx context.Context, token string, ttl time.Duration) error {
    if Redis == nil {
        return fmt.Errorf("redis client not initialized")
    }
    if token == "" {
        return fmt.Errorf("empty token")
    }
    key := jwtTokenStringKeyPrefix + sha256Hex(token)
    return Redis.Set(ctx, key, 1, ttl).Err()
}

// IsTokenStringBlacklisted 检查 token 字符串是否在黑名单
// EN: Check whether the token string is blacklisted
func IsTokenStringBlacklisted(ctx context.Context, token string) (bool, error) {
    if Redis == nil {
        return false, fmt.Errorf("redis client not initialized")
    }
    if token == "" {
        return false, nil
    }
    key := jwtTokenStringKeyPrefix + sha256Hex(token)
    n, err := Redis.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return n > 0, nil
}

// sha256Hex 计算字符串的 SHA256 十六进制
func sha256Hex(s string) string {
    h := sha256.New()
    h.Write([]byte(s))
    return hex.EncodeToString(h.Sum(nil))
}

