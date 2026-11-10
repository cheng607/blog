package handler

import (
	"log"
	"strconv"

	"blog-back/internal/model"
	"blog-back/internal/response"
	"blog-back/internal/store"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	store *store.ArticleStore
}

func NewArticleHandler(s *store.ArticleStore) *ArticleHandler {
	return &ArticleHandler{store: s}
}

func (h *ArticleHandler) Save(c *gin.Context) {
	var payload model.ArticlePayload
	if !bindJSON(c, &payload, 2<<20) {
		return
	}
	tagIDs := payload.ResolvedTagIDs()
	if !validID(payload.UserID) {
		invalid(c, "作者 ID 无效")
		return
	}
	art := &model.Article{Title: payload.Title, Content: payload.Content, CoverImage: payload.CoverImage, UserID: payload.UserID}
	if !validateArticle(c, art, tagIDs) {
		return
	}
	id, err := h.store.Save(art, tagIDs)
	if err != nil {
		storeError(c, err)
		return
	}
	art.ID = id
	// A read failure after a committed write must not encourage duplicate publishing.
	if saved, err := h.store.ByIDWithTags(id); err == nil && saved != nil {
		art = saved
	} else if err != nil {
		log.Printf("article read after save: %v", err)
	}
	response.OK(c, art)
}

func (h *ArticleHandler) List(c *gin.Context) { h.paged(c, 0, 0) }

func (h *ArticleHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	art, err := h.store.ByIDWithTags(id)
	if err != nil {
		storeError(c, err)
		return
	}
	if art == nil {
		response.Fail(c, response.CodeArticleNotExist, "\u6587\u7ae0\u4e0d\u5b58\u5728", nil)
		return
	}
	response.OK(c, art)
}

func (h *ArticleHandler) ByUser(c *gin.Context) {
	id, ok := parseIDParam(c, "userId")
	if !ok {
		return
	}
	h.paged(c, 0, id)
}

func (h *ArticleHandler) ByTag(c *gin.Context) {
	id, ok := parseIDParam(c, "tagId")
	if !ok {
		return
	}
	h.paged(c, id, 0)
}

func (h *ArticleHandler) Update(c *gin.Context) {
	var art model.Article
	if !bindJSON(c, &art, 2<<20) {
		return
	}
	if !validID(art.ID) {
		invalid(c, "文章 ID 无效")
		return
	}
	if !validateArticle(c, &art, art.TagIDs) {
		return
	}
	if err := h.store.Update(&art, art.TagIDs); err != nil {
		storeError(c, err)
		return
	}
	if saved, err := h.store.ByIDWithTags(art.ID); err == nil && saved != nil {
		response.OK(c, saved)
		return
	} else if err != nil {
		log.Printf("article read after update: %v", err)
	}
	response.OK(c, art)
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.store.Delete(id); err != nil {
		storeError(c, err)
		return
	}
	response.OKEmpty(c)
}

func (h *ArticleHandler) View(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	n, err := h.store.IncrementViews(id)
	if err != nil {
		storeError(c, err)
		return
	}
	response.OK(c, n)
}

func parseID(c *gin.Context) (int, bool) {
	return parseIDParam(c, "id")
}

func parseIDParam(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
	if err != nil || !validID(id) {
		response.Fail(c, response.CodeParamInvalid, "\u53c2\u6570\u65e0\u6548", nil)
		return 0, false
	}
	return id, true
}
