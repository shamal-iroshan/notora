package share

import (
	"github.com/gin-gonic/gin"
)

// Protected routes (requires JWT)
func RegisterProtectedShareRoutes(r *gin.RouterGroup, handler *ShareHandler) {
	r.POST("/notes/:id/share", handler.CreateShare)
	r.DELETE("/notes/:id/share", handler.DisableShare)
}

// Public routes (no JWT)
func RegisterPublicShareRoutes(r *gin.RouterGroup, handler *ShareHandler) {
	r.GET("/share/:token", handler.PublicGet)
}
