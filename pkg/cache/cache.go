package cache

import (
	"encoding/json"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	CacheEnabled bool                `json:"cacheEnabled"`
	CacheData    map[string][]string `json:"cacheData"`
}

type Cache struct {
	Data map[string]CacheItem
}
type CacheItem struct {
	Message    string
	Answer     string
	Timestamp  time.Time
	ErrorTime  time.Time
	ErrorCount int
}

func NewCache() *Cache {
	return &Cache{
		Data: make(map[string]CacheItem),
	}
}

func (c *Cache) CacheAdd(key, value string, count int) {
	c.Data[key] = CacheItem{
		Message:    value,
		Timestamp:  time.Now(),
		ErrorTime:  time.Now(),
		ErrorCount: count,
	}
	log.Printf("캐시 추가: %v\n", key)
}

// CacheUpdate 캐시 업데이트 (기존 값 대체)
func (c *Cache) CacheTimeUpdate(key string) {
	if item, exists := c.Data[key]; exists {
		item.Timestamp = time.Now() // 타임스탬프 갱신
		c.Data[key] = item
		log.Printf("error 20분 경과: %v\n", key)
	}
}

func (c *Cache) CacheErrorTimeUpdate(key string, count int) {
	if item, exists := c.Data[key]; exists {
		if item.ErrorCount != count {
			item.ErrorCount = count
			item.ErrorTime = time.Now()
			c.Data[key] = item
			log.Printf("에러가 발생하여 캐시 업데이트: %v\n", key)
		}
	}
}

func (c *Cache) CacheGPTUpdate(key string, value string) {
	if item, exists := c.Data[key]; exists {
		answer := value
		item.Answer = answer
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
	// 빈 파일인 경우 데이터 저장 후 리턴
	if len(fileData) == 0 {
		c.Data = make(map[string]CacheItem)
		return nil
	}

	cacheData := make(map[string]CacheItem)

	if err := json.Unmarshal(fileData, &cacheData); err != nil {
		return err
	}

	c.Data = cacheData
	//log.Printf("캐시 파일 읽기: %v\n", c.Data)
	return nil
}

func (c *Cache) SaveCacheToFile(filePath string) error {
	dataBytes, err := json.Marshal(c.Data)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, dataBytes, 0644)
}

func (c *Cache) DuplicateEvent(key, value string) bool {
	if item, exists := c.Data[key]; exists {
		return item.Message == value
	}
	return false
}

func (c *Cache) Cleanup(currentEvents []string) {
	for key := range c.Data {
		if !contains(currentEvents, key) {
			delete(c.Data, key)
			log.Printf("에러가 없으므로 캐시 삭제: %v\n", key)
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
