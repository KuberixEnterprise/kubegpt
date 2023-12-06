package cache

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	CacheEnabled bool                `json:"cacheEnabled"`
	CacheData    map[string][]string `json:"cacheData"`
}

type Cache struct {
	Data map[string]CacheItem
}
type CacheItem struct {
	Message   string
	Timestamp time.Time
}

func NewCache() *Cache {
	return &Cache{
		Data: make(map[string]CacheItem),
	}
}

func (c *Cache) CacheAdd(key, value string) {
	c.Data[key] = CacheItem{
		Message:   value,
		Timestamp: time.Now(),
	}
}

// CacheUpdate 캐시 업데이트 (기존 값 대체)
func (c *Cache) CacheUpdate(key, value string) {
	if item, exists := c.Data[key]; exists {
		item.Timestamp = time.Now() // 타임스탬프 갱신
		c.Data[key] = item
	}
}

// 설정 파일에 변경된 캐시 데이터 저장
func (c *Cache) LoadCacheFromFile(filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.Data = make(map[string]CacheItem)
			return nil
		}
		return err
	}
	cacheData := make(map[string]CacheItem)
	if err := json.Unmarshal(fileData, &cacheData); err != nil {
		return err
	}

	c.Data = cacheData
	return nil
}

func (c *Cache) SaveCacheToFile(filePath string) error {
	dataBytes, err := json.Marshal(c.Data)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, dataBytes, 0644)
}

func (c Cache) DuplicateEvent(key, value string) bool {
	if item, exists := c.Data[key]; exists {
		return item.Message == value
	}
	return false
}

func (c *Cache) Cleanup(currentEvents []string) {
	for key := range c.Data {
		if !contains(currentEvents, key) {
			delete(c.Data, key)
		}
	}
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
