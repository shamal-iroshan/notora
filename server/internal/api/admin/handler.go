package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/shamal-iroshan/notora/internal/repository"
)

type AdminHandler struct {
	UserRepo *repository.UserRepository
}

func NewAdminHandler(repo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{UserRepo: repo}
}

func (h *AdminHandler) ListPending(ctx *gin.Context) {
	list, _ := h.UserRepo.ListPending()
	ctx.JSON(200, gin.H{"pending_users": list})
}

func (h *AdminHandler) Approve(ctx *gin.Context) {
	id, ok := toInt64Strict(ctx, "id")
	if !ok {
		return
	}
	h.UserRepo.Approve(id)
	ctx.JSON(200, gin.H{"status": "approved"})
}

func (h *AdminHandler) Suspend(ctx *gin.Context) {
	id, ok := toInt64Strict(ctx, "id")
	if !ok {
		return
	}
	h.UserRepo.Suspend(id)
	ctx.JSON(200, gin.H{"status": "suspended"})
}

func (h *AdminHandler) DeleteUser(ctx *gin.Context) {
	id, ok := toInt64Strict(ctx, "id")
	if !ok {
		return
	}
	h.UserRepo.DeleteUser(id)
	ctx.JSON(200, gin.H{"status": "deleted"})
}
