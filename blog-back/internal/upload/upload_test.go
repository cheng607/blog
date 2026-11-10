package upload

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fileHeader(t *testing.T, name string, data []byte) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = part.Write(data); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err = req.ParseMultipartForm(12 << 20); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = req.MultipartForm.RemoveAll() })
	return req.MultipartForm.File["file"][0]
}

func TestImageValidationAndStorage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	var pngData, jpegData bytes.Buffer
	if err := png.Encode(&pngData, img); err != nil {
		t.Fatal(err)
	}
	if err := jpeg.Encode(&jpegData, img, nil); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name  string
		data  []byte
		valid bool
	}{
		{"photo.png", pngData.Bytes(), true},
		{"photo.JPG", jpegData.Bytes(), true},
		{"photo.jpeg", jpegData.Bytes(), true},
		{"fake.png", []byte("<script>alert(1)</script>"), false},
		{"wrong.jpg", pngData.Bytes(), false},
		{"broken.png", pngData.Bytes()[:40], false},
		{"photo.gif", pngData.Bytes(), false},
		{"empty.png", nil, false},
		{"large.png", make([]byte, MaxImageBytes+1), false},
		{"appended.png", append(append([]byte{}, pngData.Bytes()...), []byte("<script>payload</script>")...), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			url, err := NewService(dir, "/upload/images/").SaveImage(fileHeader(t, tc.name, tc.data))
			if !tc.valid {
				if err == nil {
					t.Fatal("invalid image accepted")
				}
				entries, _ := os.ReadDir(dir)
				if len(entries) != 0 {
					t.Fatal("invalid image written to disk")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(dir, "images", strings.TrimPrefix(url, "/upload/images/")))
			if err != nil {
				t.Fatal(err)
			}
			if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(data, []byte("<script>")) {
				t.Fatal("appended payload retained")
			}
		})
	}
}
