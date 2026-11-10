package handler

import (
	"blog-back/internal/response"
	"blog-back/internal/store"
	"github.com/gin-gonic/gin"
	"strconv"
)

func parsePagination(c *gin.Context) (int, int, bool) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 || page > 1000000 {
		invalid(c, "页码必须为1至1000000的整数")
		return 0, 0, false
	}
	size, err := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	if err != nil || size < 1 || size > 50 {
		invalid(c, "每页条数必须为1至50的整数")
		return 0, 0, false
	}
	return page, size, true
}
func (h *ArticleHandler) paged(c *gin.Context, tagID, userID int) {
	page, size, ok := parsePagination(c)
	if !ok {
		return
	}
	result, err := h.store.Page(c.Request.Context(), page, size, tagID, userID)
	if err != nil {
		storeError(c, err)
		return
	}
	response.OK(c, result)
}
func SiteSummary(s *store.SummaryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := s.Get(c.Request.Context())
		if err != nil {
			storeError(c, err)
			return
		}
		response.OK(c, result)
	}
}
