package encrypted

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shamal-iroshan/notora/internal/model"
	"github.com/shamal-iroshan/notora/internal/service"
)

type EncryptedNotesHandler struct {
	Service *service.EncryptedNotesService
}

func NewEncryptedNotesHandler(s *service.EncryptedNotesService) *EncryptedNotesHandler {
	return &EncryptedNotesHandler{Service: s}
}

func toInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

// POST /api/encrypted-notes
func (h *EncryptedNotesHandler) Create(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	var dto CreateEncryptedNoteDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

	input := model.CreateEncryptedNoteInput{
		TitleCiphertext:   dto.TitleCiphertext,
		ContentCiphertext: dto.ContentCiphertext,
		TitleNonce:        dto.TitleNonce,
		ContentNonce:      dto.ContentNonce,
		NoteSalt:          dto.NoteSalt,
	}

	id, err := h.Service.Create(userID, input)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "db error"})
		return
	}

	ctx.JSON(200, gin.H{"id": id, "status": "created"})
}

// GET /api/encrypted-notes
func (h *EncryptedNotesHandler) List(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	notes, err := h.Service.List(userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "db error"})
		return
	}

	ctx.JSON(200, notes)
}

// GET /api/encrypted-notes/:id
func (h *EncryptedNotesHandler) Get(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	note, err := h.Service.Get(userID, noteID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}

	ctx.JSON(200, note)
}

// PUT /api/encrypted-notes/:id
func (h *EncryptedNotesHandler) Update(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	var dto UpdateEncryptedNoteDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

	input := model.UpdateEncryptedNoteInput{
		TitleCiphertext:   dto.TitleCiphertext,
		ContentCiphertext: dto.ContentCiphertext,
		TitleNonce:        dto.TitleNonce,
		ContentNonce:      dto.ContentNonce,
		NoteSalt:          dto.NoteSalt,
	}

	if err := h.Service.Update(userID, noteID, input); err != nil {
		ctx.JSON(500, gin.H{"error": "db error"})
		return
	}

	ctx.JSON(200, gin.H{"status": "updated"})
}

// DELETE /api/encrypted-notes/:id
func (h *EncryptedNotesHandler) Delete(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	if err := h.Service.Delete(userID, noteID); err != nil {
		ctx.JSON(500, gin.H{"error": "db error"})
		return
	}

	ctx.JSON(200, gin.H{"status": "deleted"})
}
