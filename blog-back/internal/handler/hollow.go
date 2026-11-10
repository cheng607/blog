package handler

import (
	"strings"
	"time"
	"unicode/utf8"

	"blog-back/internal/model"
	"blog-back/internal/response"
	"blog-back/internal/store"

	"github.com/gin-gonic/gin"
)

type HollowHandler struct {
	store *store.HollowStore
}

func NewHollowHandler(s *store.HollowStore) *HollowHandler {
	return &HollowHandler{store: s}
}

type hollowInput struct {
	Word string `json:"word"`
	Date string `json:"date"`
}

func (h *HollowHandler) List(c *gin.Context) {
	list, err := h.store.List()
	if err != nil {
		storeError(c, err)
		return
	}
	response.OK(c, list)
}

func (h *HollowHandler) Save(c *gin.Context) {
	var in hollowInput
	if !bindJSON(c, &in, 16<<10) {
		return
	}
	in.Word = strings.TrimSpace(in.Word)
	if in.Word == "" || utf8.RuneCountInString(in.Word) > 1000 {
		invalid(c, "树洞内容必填且最多1000字")
		return
	}
	item := model.Hollow{Word: in.Word}
	if in.Date != "" {
		layouts := []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05"}
		var parsed time.Time
		var err error
		for _, layout := range layouts {
			parsed, err = time.ParseInLocation(layout, in.Date, time.Local)
			if err == nil {
				item.Date = &parsed
				break
			}
		}
		if item.Date == nil {
			response.Fail(c, response.CodeParamInvalid, "\u65e5\u671f\u683c\u5f0f\u65e0\u6548", nil)
			return
		}
	}
	if err := h.store.Save(&item); err != nil {
		storeError(c, err)
		return
	}
	response.OK(c, item)
}
