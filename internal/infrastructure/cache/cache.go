package cache

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

// RenderCache は描画結果のキャッシュ
type RenderCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	maxSize int // キャッシュの最大エントリ数
}

// CacheEntry はキャッシュエントリ
type CacheEntry struct {
	Key    string
	Result []byte
}

// NewRenderCache は新しいキャッシュを作成する
func NewRenderCache(maxSize int) *RenderCache {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &RenderCache{
		entries: make(map[string]*CacheEntry),
		maxSize: maxSize,
	}
}

// ComputeKey はデータからキャッシュキーを計算する
func ComputeKey(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// Get はキャッシュからエントリを取得する
func (c *RenderCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.Result, true
}

// Put はキャッシュにエントリを追加する
func (c *RenderCache) Put(key string, result []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// キャッシュサイズ上限に達した場合、最も古いエントリを削除
	if len(c.entries) >= c.maxSize {
		// 簡易的に最初に見つかったエントリを削除
		for k := range c.entries {
			delete(c.entries, k)
			break
		}
	}

	c.entries[key] = &CacheEntry{
		Key:    key,
		Result: result,
	}
}

// Invalidate はキャッシュからエントリを削除する
func (c *RenderCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Clear はキャッシュを全てクリアする
func (c *RenderCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}

// Size はキャッシュ内のエントリ数を返す
func (c *RenderCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
