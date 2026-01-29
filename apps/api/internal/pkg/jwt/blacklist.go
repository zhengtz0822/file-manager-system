package jwt

import (
	"time"
)

// Blacklist Token 黑名单接口
// 可以使用 Redis 或内存实现
type Blacklist interface {
	Add(token string, expiration time.Duration) error
	Exists(token string) (bool, error)
	Remove(token string) error
}

// MemoryBlacklist 内存黑名单实现（适合单机）
type MemoryBlacklist struct {
	tokens map[string]time.Time
}

// NewMemoryBlacklist 创建内存黑名单
func NewMemoryBlacklist() *MemoryBlacklist {
	bl := &MemoryBlacklist{
		tokens: make(map[string]time.Time),
	}
	// 启动清理协程
	go bl.cleanup()
	return bl
}

// Add 添加 token 到黑名单
func (m *MemoryBlacklist) Add(token string, expiration time.Duration) error {
	m.tokens[token] = time.Now().Add(expiration)
	return nil
}

// Exists 检查 token 是否在黑名单
func (m *MemoryBlacklist) Exists(token string) (bool, error) {
	exp, ok := m.tokens[token]
	if !ok {
		return false, nil
	}
	// 检查是否过期
	if time.Now().After(exp) {
		delete(m.tokens, token)
		return false, nil
	}
	return true, nil
}

// Remove 从黑名单移除
func (m *MemoryBlacklist) Remove(token string) error {
	delete(m.tokens, token)
	return nil
}

// cleanup 定期清理过期的 token
func (m *MemoryBlacklist) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for token, exp := range m.tokens {
			if now.After(exp) {
				delete(m.tokens, token)
			}
		}
	}
}

// RevokeToken 撤销 token
func RevokeToken(token string) error {
	// 解析 token 获取过期时间
	claims, err := ParseToken(token)
	if err != nil {
		return err
	}

	expiration := claims.ExpiresAt.Sub(time.Now())
	if expiration < 0 {
		return nil // token 已经过期
	}

	// 使用全局黑名单（需要在主程序初始化）
	return blacklist.Add(token, expiration)
}

// IsRevoked 检查 token 是否被撤销
func IsRevoked(token string) bool {
	exists, _ := blacklist.Exists(token)
	return exists
}

var blacklist Blacklist

// InitBlacklist 初始化黑名单
func InitBlacklist(bl Blacklist) {
	blacklist = bl
}

func init() {
	// 默认使用内存黑名单
	blacklist = NewMemoryBlacklist()
}
