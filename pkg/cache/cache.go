package cache

import (
	"encoding/json"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	CacheEnabled bool                `json:"cacheEnabled"`
	CacheData    map[string][]string `json:"cacheData"`
}

func initViper() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetDefault("cacheEnabled", true)
	return viper.ReadInConfig()
}

func cacheAdd(cacheData map[string][]string, key, value string) {
	if _, exists := cacheData[key]; !exists {
		cacheData[key] = []string{}
	}
	cacheData[key] = append(cacheData[key], value)
}

// 캐시 삭제
func cacheDelete(cacheData map[string][]string, key string) {
	delete(cacheData, key)
}

// 캐시 업데이트 (기존 값 대체)
func cacheUpdate(cacheData map[string][]string, key, value string) {
	cacheData[key] = []string{value}
}

// 설정 파일에 변경된 캐시 데이터 저장
func loadCacheFromFile(filePath string) (map[string][]string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string][]string), nil
		}
		return nil, err
	}
	cacheData := make(map[string][]string)
	if err := json.Unmarshal(fileData, &cacheData); err != nil {
		return nil, err
	}
	return cacheData, nil
}

func saveCacheToFile(cacheData map[string][]string, filePath string) error {
	dataBytes, err := json.Marshal(cacheData)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, dataBytes, 0644)
}
