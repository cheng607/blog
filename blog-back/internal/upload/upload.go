package upload

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MaxImageBytes = 10 * 1024 * 1024
const maxImagePixels = 16_000_000

type ValidationError string
func (e ValidationError) Error() string { return string(e) }

type Service struct{ BasePath, URLPrefix string }

func NewService(basePath, urlPrefix string) *Service {
	return &Service{BasePath: basePath, URLPrefix: strings.TrimRight(urlPrefix, "/") + "/"}
}

func (s *Service) SaveImage(file *multipart.FileHeader) (string, error) {
	if file == nil || file.Size == 0 {
		return "", ValidationError("上传文件不能为空")
	}
	if file.Size > MaxImageBytes {
		return "", ValidationError("图片大小不能超过10MB")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", ValidationError("仅支持JPG、JPEG、PNG图片")
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	data, err := io.ReadAll(io.LimitReader(src, MaxImageBytes+1))
	if err != nil {
		return "", err
	}
	if len(data) > MaxImageBytes {
		return "", ValidationError("图片大小不能超过10MB")
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", ValidationError("文件不是有效图片")
	}
	if (format != "png" && format != "jpeg") || (ext == ".png") != (format == "png") {
		return "", ValidationError("图片内容与扩展名不一致")
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || int64(cfg.Width)*int64(cfg.Height) > maxImagePixels {
		return "", ValidationError("图片分辨率过大，最多支持1600万像素")
	}
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", ValidationError("图片已损坏或不完整")
	}
	// Re-encode to exclude appended payloads and unrelated metadata.
	var clean bytes.Buffer
	if format == "png" {
		err = png.Encode(&clean, decoded)
	} else {
		err = jpeg.Encode(&clean, decoded, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return "", ValidationError("图片处理失败")
	}
	if clean.Len() > MaxImageBytes {
		return "", ValidationError("处理后的图片超过10MB，请缩小图片后重试")
	}
	dateDir := time.Now().Format("20060102")
	saveDir := filepath.Join(s.BasePath, "images", dateDir)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", err
	}
	out, err := os.CreateTemp(saveDir, "image-*"+ext)
	if err != nil {
		return "", err
	}
	name := out.Name()
	_, writeErr := out.Write(clean.Bytes())
	closeErr := out.Close()
	if writeErr != nil || closeErr != nil {
		_ = os.Remove(name)
		return "", fmt.Errorf("图片保存失败")
	}
	return s.URLPrefix + dateDir + "/" + filepath.Base(name), nil
}
