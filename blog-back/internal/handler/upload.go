package handler

import (
	"blog-back/internal/response"
	"blog-back/internal/upload"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	svc *upload.Service
}

func NewUploadHandler(svc *upload.Service) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) Image(c *gin.Context) {
	// Bound the complete multipart request before parsing any file.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, upload.MaxImageBytes+1024*1024)
	defer func() {
		if c.Request.MultipartForm != nil {
			_ = c.Request.MultipartForm.RemoveAll()
		}
	}()

	file, err := c.FormFile("file")
	if err != nil {
		file, err = c.FormFile("file[]")
	}
	if err != nil {
		var large *http.MaxBytesError
		if errors.As(err, &large) {
			response.Fail(c, response.CodeTooLarge, "图片请求过大，单张图片不能超过10MB", nil)
		} else {
			response.Fail(c, response.CodeUploadError, "请选择有效的图片文件上传", nil)
		}
		return
	}
	// 简单大小限制（10MB）
	if file.Size > upload.MaxImageBytes {
		response.Fail(c, response.CodeUploadError, "图片大小不能超过 10MB", nil)
		return
	}
	url, err := h.svc.SaveImage(file)
	if err != nil {
		var invalid upload.ValidationError
		if errors.As(err, &invalid) {
			response.Fail(c, response.CodeUploadError, invalid.Error(), nil)
		} else {
			response.Internal(c, err)
		}
		return
	}
	response.OK(c, gin.H{"imageUrl": url})
}
