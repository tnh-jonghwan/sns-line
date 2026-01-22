package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Env struct {
	Kid                  string
	ChannelId            string
	LineApiPrefix        string
	LineAccessToken      string
	InstagramVerifyToken string
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
		lineAccessToken := os.Getenv("LINE_ACCESS_TOKEN")
		instagramVerifyToken := os.Getenv("INSTAGRAM_VERIFY_TOKEN")

		if kid == "" || channelId == "" || lineApiPrefix == "" || lineAccessToken == "" {
			log.Fatal("Environment variables KID, CHANNEL_ID, LINE_API_PREFIX, and LINE_ACCESS_TOKEN must be set")
		}

		if instagramVerifyToken == "" {
			log.Println("Warning: INSTAGRAM_VERIFY_TOKEN not set")
		}

		instance = &Env{
			Kid:                  kid,
			ChannelId:            channelId,
			LineApiPrefix:        lineApiPrefix,
			LineAccessToken:      lineAccessToken,
			InstagramVerifyToken: instagramVerifyToken,
		}
	})

	return instance
}
