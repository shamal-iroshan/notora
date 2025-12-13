package encrypted

import "github.com/gin-gonic/gin"

func RegisterEncryptedNotesRoutes(r *gin.RouterGroup, h *EncryptedNotesHandler) {
	r.POST("/", h.Create)
	r.GET("/", h.List)
	r.GET("/:id", h.Get)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}
