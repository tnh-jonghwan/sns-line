package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Env struct {
	Kid           string
	ChannelId     string
	LineApiPrefix string
	AccessToken   string
}

var (
	instance *Env
	once     sync.Once
)

// 싱글톤 패턴
func GetEnv() *Env {
	once.Do(func() {
		// .env 파일 로드 (파일이 없어도 에러 무시)
		_ = godotenv.Load()

		kid := os.Getenv("KID")
		channelId := os.Getenv("CHANNEL_ID")
		lineApiPrefix := os.Getenv("LINE_API_PREFIX")
		accessToken := os.Getenv("ACCESS_TOKEN")

		if kid == "" || channelId == "" || lineApiPrefix == "" || accessToken == "" {
			log.Fatal("Environment variables KID, CHANNEL_ID, LINE_API_PREFIX, and ACCESS_TOKEN must be set")
		}

		instance = &Env{
			Kid:           kid,
			ChannelId:     channelId,
			LineApiPrefix: lineApiPrefix,
			AccessToken:   accessToken,
		}
	})

	return instance
}
