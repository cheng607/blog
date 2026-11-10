package handler

import (
	"blog-back/internal/model"
	"blog-back/internal/response"
	"blog-back/internal/store"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	store *store.UserStore
}

func NewUserHandler(s *store.UserStore) *UserHandler {
	return &UserHandler{store: s}
}

func (h *UserHandler) List(c *gin.Context) {
	list, err := h.store.List()
	if err != nil {
		storeError(c, err)
		return
	}
	response.OK(c, list)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	u, err := h.store.ByID(id)
	if err != nil {
		storeError(c, err)
		return
	}
	if u == nil {
		response.Fail(c, response.CodeNotFound, "请求的记录不存在", nil)
		return
	}
	response.OK(c, u)
}

func (h *UserHandler) Save(c *gin.Context) {
	var u model.User
	if !bindJSON(c, &u, 16<<10) {
		return
	}
	if !validateUser(c, &u) {
		return
	}
	id, err := h.store.Save(&u)
	if err != nil {
		storeError(c, err)
		return
	}
	u.ID = id
	response.OK(c, u)
}

func (h *UserHandler) Update(c *gin.Context) {
	var u model.User
	if !bindJSON(c, &u, 16<<10) {
		return
	}
	if !validID(u.ID) {
		invalid(c, "记录 ID 无效")
		return
	}
	if !validateUser(c, &u) {
		return
	}
	_, err := h.store.Update(&u)
	if err != nil {
		storeError(c, err)
		return
	}
	response.OK(c, u)
}

func (h *UserHandler) Delete(c *gin.Context) {
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
