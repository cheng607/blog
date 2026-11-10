package handler

import (
	"blog-back/internal/model"
	"blog-back/internal/response"
	"blog-back/internal/store"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	store *store.TagStore
}

func NewTagHandler(s *store.TagStore) *TagHandler {
	return &TagHandler{store: s}
}

func (h *TagHandler) Save(c *gin.Context) {
	var t model.Tag
	if !bindJSON(c, &t, 16<<10) {
		return
	}
	if !validateTag(c, &t) {
		return
	}
	id, err := h.store.Save(&t)
	if err != nil {
		storeError(c, err)
		return
	}
	t.ID = id
	response.OK(c, t)
}

func (h *TagHandler) List(c *gin.Context) {
	list, err := h.store.ListWithCount()
	if err != nil {
		storeError(c, err)
		return
	}
	response.OK(c, list)
}

func (h *TagHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	t, err := h.store.ByIDWithCount(id)
	if err != nil {
		storeError(c, err)
		return
	}
	if t == nil {
		response.Fail(c, response.CodeNotFound, "请求的记录不存在", nil)
		return
	}
	response.OK(c, t)
}

func (h *TagHandler) Update(c *gin.Context) {
	var t model.Tag
	if !bindJSON(c, &t, 16<<10) {
		return
	}
	if !validID(t.ID) {
		invalid(c, "记录 ID 无效")
		return
	}
	if !validateTag(c, &t) {
		return
	}
	if err := h.store.Update(&t); err != nil {
		storeError(c, err)
		return
	}
	response.OK(c, t)
}

func (h *TagHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.store.Delete(id); err != nil {
		storeError(c, err)
		return
	}
	response.OKEmpty(c)
}
