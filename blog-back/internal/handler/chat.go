package handler

import (
	"blog-back/internal/ailimit"
	"blog-back/internal/config"
	"blog-back/internal/response"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const chatTotalTimeout = 65 * time.Second
const chatAttemptTimeout = 30 * time.Second
const maxChatAttempts = 2

type ChatHandler struct {
	apiKey                       string
	models                       []string
	baseURL                      string
	client                       *http.Client
	limiter                      *ailimit.Limiter
	maxTokens                    int
	totalTimeout, attemptTimeout time.Duration
}

func NewChatHandler(cfg config.Config) *ChatHandler {
	models := parseModels(cfg.DashScopeModel)
	if len(models) == 0 {
		models = []string{"qwen-plus", "qwen-turbo"}
	}
	base := strings.TrimRight(strings.TrimSpace(cfg.DashScopeBaseURL), "/")
	if base == "" {
		base = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}
	rate, concurrent, daily, tokens := cfg.AIRatePerMinute, cfg.AIMaxConcurrent, cfg.AIDailyLimit, cfg.AIMaxTokens
	if rate <= 0 {
		rate = 6
	}
	if concurrent <= 0 {
		concurrent = 3
	}
	if daily <= 0 {
		daily = 500
	}
	if tokens <= 0 {
		tokens = 1024
	}
	return &ChatHandler{apiKey: strings.TrimSpace(cfg.DashScopeAPIKey), models: models, baseURL: base,
		client: &http.Client{}, maxTokens: tokens,
		limiter:      ailimit.New(rate, concurrent, daily, filepath.Join(cfg.UploadPath, ".ai-quota.json")),
		totalTimeout: chatTotalTimeout, attemptTimeout: chatAttemptTimeout}
}
func parseModels(raw string) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, m := range strings.Split(raw, ",") {
		m = strings.TrimSpace(m)
		if m != "" && !seen[m] {
			out = append(out, m)
			seen[m] = true
		}
	}
	return out
}

type chatBody struct {
	Message string `json:"message"`
}
type dashScopeRequest struct {
	Model     string             `json:"model"`
	Messages  []dashScopeMessage `json:"messages"`
	MaxTokens int                `json:"max_tokens"`
}
type dashScopeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type dashScopeResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error json.RawMessage `json:"error"`
}

func (h *ChatHandler) Status(c *gin.Context) { response.OK(c, gin.H{"configured": h.apiKey != ""}) }
func (h *ChatHandler) Chat(c *gin.Context) {
	if h.apiKey == "" {
		response.HTTPError(c, 503, "AI 暂时不可用，请稍后再试")
		return
	}
	var body chatBody
	if !bindJSON(c, &body, 16*1024) {
		return
	}
	msg := strings.TrimSpace(body.Message)
	if msg == "" || !utf8.ValidString(msg) || utf8.RuneCountInString(msg) > 2000 {
		response.Fail(c, response.CodeParamInvalid, "问题不能为空且不能超过2000字", nil)
		return
	}
	release, retry, message := h.limiter.Acquire(c.ClientIP())
	if release == nil {
		c.Header("Retry-After", strconv.Itoa(retry))
		response.HTTPError(c, 429, message)
		return
	}
	defer release()
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.totalTimeout)
	defer cancel()
	var last error
	for i, model := range h.models {
		if i >= maxChatAttempts || ctx.Err() != nil {
			break
		}
		if err := h.limiter.ReserveDaily(); err != nil {
			if errors.Is(err, ailimit.ErrDailyLimit) {
				c.Header("Retry-After", strconv.Itoa(h.limiter.DailyRetryAfter()))
				response.HTTPError(c, 429, "今日 AI 额度已用完，请明天再来")
			} else {
				log.Printf("AI quota storage error: %v", err)
				response.HTTPError(c, 503, "AI 暂时不可用，请稍后再试")
			}
			return
		}
		var content string
		var retryable bool
		content, last, retryable = h.callModel(ctx, model, msg)
		if last == nil {
			response.OK(c, gin.H{"content": content, "model": model})
			return
		}
		// Do not log provider bodies, prompts or credentials.
		log.Printf("AI upstream failure model=%s: %v", model, last)
		if !retryable {
			break
		}
	}
	if c.Request.Context().Err() != nil {
		return
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(last, context.DeadlineExceeded) {
		response.HTTPError(c, 504, "AI 回答超时，请稍后重试")
		return
	}
	response.HTTPError(c, 502, "AI 服务暂时无法回答，请稍后重试")
}

func (h *ChatHandler) callModel(parent context.Context, model, msg string) (string, error, bool) {
	ctx, cancel := context.WithTimeout(parent, h.attemptTimeout)
	defer cancel()
	payload, _ := json.Marshal(dashScopeRequest{Model: model, MaxTokens: h.maxTokens, Messages: []dashScopeMessage{
		{Role: "system", Content: "你是博客站点的 AI 助手，用简洁中文回答问题。"}, {Role: "user", Content: msg},
	}})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", errors.New("invalid upstream URL"), false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	resp, err := h.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err(), parent.Err() == nil
		}
		return "", errors.New("upstream network error"), true
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		retryable := resp.StatusCode == 408 || resp.StatusCode == 429 || resp.StatusCode >= 500
		return "", fmt.Errorf("upstream HTTP %d", resp.StatusCode), retryable
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, (1<<20)+1))
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err(), parent.Err() == nil
		}
		return "", errors.New("upstream read error"), true
	}
	if len(raw) > 1<<20 {
		return "", errors.New("upstream response too large"), false
	}
	var result dashScopeResponse
	if json.Unmarshal(raw, &result) != nil {
		return "", errors.New("invalid upstream JSON"), true
	}
	if len(result.Error) > 0 && string(result.Error) != "null" {
		return "", errors.New("upstream API error"), false
	}
	if len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		return "", errors.New("empty upstream response"), true
	}
	return strings.TrimSpace(result.Choices[0].Message.Content), nil, false
}
