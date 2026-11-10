package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

// AdminSessions stores opaque eight-hour sessions. Restart/logout revokes them.
// The administrator password is never accepted as an API bearer token.
type AdminSessions struct {
	password string
	mu       sync.Mutex
	sessions map[string]time.Time
	now      func() time.Time
}

func NewAdminSessions(password string) *AdminSessions {
	return &AdminSessions{password: strings.TrimSpace(password), sessions: make(map[string]time.Time), now: time.Now}
}

func (s *AdminSessions) Login(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	if s.password == "" {
		c.AbortWithStatusJSON(503, gin.H{"code": 503, "msg": "管理员登录尚未配置"})
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
	var input struct {
		Password string `json:"password"`
	}
	if c.ShouldBindJSON(&input) != nil {
		c.AbortWithStatusJSON(400, gin.H{"code": 400, "msg": "登录参数无效"})
		return
	}
	want, got := sha256.Sum256([]byte(s.password)), sha256.Sum256([]byte(input.Password))
	if subtle.ConstantTimeCompare(want[:], got[:]) != 1 {
		c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "管理密钥错误"})
		return
	}
	var random [32]byte
	if _, err := rand.Read(random[:]); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"code": 500, "msg": "无法创建会话"})
		return
	}
	token := hex.EncodeToString(random[:])
	now := s.now()
	expires := now.Add(8 * time.Hour)
	s.mu.Lock()
	for key, expiry := range s.sessions {
		if !now.Before(expiry) {
			delete(s.sessions, key)
		}
	}
	s.sessions[token] = expires
	s.mu.Unlock()
	c.JSON(200, gin.H{"code": 200, "data": gin.H{"token": token, "expiresAt": expires}})
}

func sessionToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}

func (s *AdminSessions) Require() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		token := sessionToken(c)
		s.mu.Lock()
		expires, exists := s.sessions[token]
		valid := s.password != "" && exists && s.now().Before(expires)
		if !valid {
			delete(s.sessions, token)
		}
		s.mu.Unlock()
		if !valid {
			c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "请重新登录管理员账号"})
			return
		}
		c.Next()
	}
}

func (s *AdminSessions) Logout(c *gin.Context) {
	s.mu.Lock()
	delete(s.sessions, sessionToken(c))
	s.mu.Unlock()
	c.Header("Cache-Control", "no-store")
	c.JSON(200, gin.H{"code": 200})
}
