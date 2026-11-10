package handler

import (
	"blog-back/internal/model"
	"blog-back/internal/response"
	"blog-back/internal/store"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

func bindJSON(c *gin.Context, target interface{}, limit int64) bool {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(target)
	if err == nil {
		var extra interface{}
		if e := decoder.Decode(&extra); e != io.EOF {
			if e == nil {
				err = errors.New("multiple JSON values")
			} else {
				err = e
			}
		}
	}
	if err != nil {
		var large *http.MaxBytesError
		if errors.As(err, &large) {
			response.Fail(c, response.CodeTooLarge, "请求内容过大", nil)
		} else {
			response.Fail(c, response.CodeParamInvalid, "请求必须是有效的 JSON 对象", nil)
		}
		return false
	}
	return true
}
func invalid(c *gin.Context, message string) bool {
	response.Fail(c, response.CodeParamInvalid, message, nil)
	return false
}
func validID(id int) bool { return id > 0 && int64(id) <= 2147483647 }
func validMediaURL(raw string) bool {
	if raw == "" {
		return true
	}
	if utf8.RuneCountInString(raw) > 512 || strings.ContainsAny(raw, "\r\n\t\\") {
		return false
	}
	if strings.HasPrefix(raw, "/") {
		return !strings.HasPrefix(raw, "//")
	}
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Host != "" && u.User == nil
}
func validateArticle(c *gin.Context, art *model.Article, tags []int) bool {
	art.Title = strings.TrimSpace(art.Title)
	if art.Title == "" || utf8.RuneCountInString(art.Title) > 255 {
		return invalid(c, "文章标题必填且最多255字")
	}
	if strings.TrimSpace(art.Content) == "" || len(art.Content) > 1024*1024 {
		return invalid(c, "正文不能为空且不能超过1MB")
	}
	if art.CoverImage != nil && !validMediaURL(*art.CoverImage) {
		return invalid(c, "封面地址必须为站内路径或 HTTP(S) 地址，最多512字")
	}
	if len(tags) == 0 || len(tags) > 20 {
		return invalid(c, "请选择1至20个标签")
	}
	seen := map[int]bool{}
	for _, id := range tags {
		if !validID(id) || seen[id] {
			return invalid(c, "标签 ID 必须有效且不能重复")
		}
		seen[id] = true
	}
	return true
}
func validateTag(c *gin.Context, tag *model.Tag) bool {
	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" || utf8.RuneCountInString(tag.Name) > 64 || utf8.RuneCountInString(tag.Description) > 255 {
		return invalid(c, "标签名称必填且最多64字，描述最多255字")
	}
	return true
}
func validateUser(c *gin.Context, user *model.User) bool {
	user.Username = strings.TrimSpace(user.Username)
	if user.Username == "" || utf8.RuneCountInString(user.Username) > 64 || utf8.RuneCountInString(user.Signature) > 255 {
		return invalid(c, "昵称必填且最多64字，签名最多255字")
	}
	if !validMediaURL(user.Avatar) {
		return invalid(c, "头像地址必须为站内路径或 HTTP(S) 地址，最多512字")
	}
	return true
}
func storeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		response.Fail(c, response.CodeNotFound, "请求的记录不存在", nil)
	case errors.Is(err, store.ErrInvalidTag):
		response.Fail(c, response.CodeParamInvalid, "所选标签已不存在，请刷新后重新选择", nil)
	case errors.Is(err, store.ErrInvalidAuthor):
		response.Fail(c, response.CodeParamInvalid, "博主资料不存在，请先保存资料", nil)
	case errors.Is(err, store.ErrConflict):
		response.Fail(c, response.CodeConflict, "记录仍被文章引用，请先处理关联文章", nil)
	default:
		response.Internal(c, err)
	}
}
