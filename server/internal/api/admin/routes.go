package admin

import "github.com/gin-gonic/gin"

func RegisterAdminRoutes(r *gin.RouterGroup, h *AdminHandler) {
	r.GET("/pending-users", h.ListPending)
	r.POST("/users/:id/approve", h.Approve)
	r.POST("/users/:id/suspend", h.Suspend)
	r.DELETE("/users/:id", h.DeleteUser)
}
