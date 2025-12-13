package share

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shamal-iroshan/notora/internal/config"
	"github.com/shamal-iroshan/notora/internal/service"
)

// ShareHandler handles sharing-related API endpoints.
type ShareHandler struct {
	Service *service.ShareService
	CFG     *config.Config
}

// Constructor for ShareHandler
func NewShareHandler(service *service.ShareService, cfg *config.Config) *ShareHandler {
	return &ShareHandler{
		Service: service,
		CFG:     cfg,
	}
}

// Helper to convert string → int64 safely
func toInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

// -------------------------------------------------------------
// POST /api/notes/:id/share  (Protected)
// Creates a shareable link for a note
// -------------------------------------------------------------
func (h *ShareHandler) CreateShare(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	token, err := h.Service.CreateShare(userID, noteID)
	if err != nil {
		ctx.JSON(403, gin.H{"error": "not allowed"})
		return
	}

	shareURL := fmt.Sprintf("%s/api/share/%s", h.CFG.AppBaseURL, token)

	ctx.JSON(200, ShareResponse{
		ShareURL: shareURL,
	})
}

// -------------------------------------------------------------
// DELETE /api/notes/:id/share  (Protected)
// Disables sharable link for the note
// -------------------------------------------------------------
func (h *ShareHandler) DisableShare(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	if err := h.Service.DisableShare(userID, noteID); err != nil {
		ctx.JSON(403, gin.H{"error": "not allowed"})
		return
	}

	ctx.JSON(200, DisableShareResponse{
		Status: "share_disabled",
	})
}

// -------------------------------------------------------------
// GET /api/share/:token  (Public)
// Public endpoint that fetches a shared note
// -------------------------------------------------------------
func (h *ShareHandler) PublicGet(ctx *gin.Context) {
	token := ctx.Param("token")

	note, err := h.Service.GetSharedNote(token)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}

	ctx.JSON(200, PublicSharedNoteResponse{
		Note: note,
	})
}
