package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AIRatePerMinute int
	AIMaxConcurrent int
	AIDailyLimit    int
	AIMaxTokens     int
	TrustedProxies  []string

	Port             int
	DBDSN            string
	UploadPath       string
	UploadURLPrefix  string
	CORSOrigins      []string
	DashScopeAPIKey  string
	DashScopeModel   string
	DashScopeBaseURL string
}

func Load() Config {
	port, _ := strconv.Atoi(getEnv("SERVER_PORT", "9090"))
	origins := strings.Split(getEnv("CORS_ORIGINS", ""), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return Config{
		AIRatePerMinute: positiveInt("AI_RATE_PER_MINUTE", 6, 60),
		AIMaxConcurrent: positiveInt("AI_MAX_CONCURRENT", 3, 20),
		AIDailyLimit:    positiveInt("AI_DAILY_LIMIT", 500, 100000),
		AIMaxTokens:     positiveInt("AI_MAX_TOKENS", 1024, 4096),
		TrustedProxies:  strings.FieldsFunc(getEnv("TRUSTED_PROXIES", "127.0.0.1,::1"), func(r rune) bool { return r == ',' || r == ' ' }),

		Port:             port,
		DBDSN:            getEnv("DB_DSN", ""),
		UploadPath:       getEnv("UPLOAD_PATH", "./uploads"), // 默认使用相对路径，避免硬编码
		UploadURLPrefix:  getEnv("UPLOAD_URL_PREFIX", "/upload/images/"),
		CORSOrigins:      origins,
		DashScopeAPIKey:  getEnv("DASHSCOPE_API_KEY", ""),
		DashScopeModel:   getEnv("DASHSCOPE_MODEL", "qwen-plus,qwen-turbo"),
		DashScopeBaseURL: getEnv("DASHSCOPE_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Invalid/zero values never disable a protection.
func positiveInt(key string, fallback, maximum int) int {
	n, err := strconv.Atoi(os.Getenv(key))
	if err != nil || n <= 0 {
		return fallback
	}
	if n > maximum {
		return maximum
	}
	return n
}
