package notes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shamal-iroshan/notora/internal/service"
)

type NoteHandler struct {
	Service *service.NoteService
}

func NewNoteHandler(service *service.NoteService) *NoteHandler {
	return &NoteHandler{Service: service}
}

// Create
func (h *NoteHandler) Create(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	var req CreateNoteRequest
	if ctx.ShouldBindJSON(&req) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	id, err := h.Service.Create(userID, req.Title, req.Content)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "could not create"})
		return
	}

	ctx.JSON(201, gin.H{"id": id})
}

// Get one
func (h *NoteHandler) Get(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	note, err := h.Service.Get(userID, noteID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}

	ctx.JSON(200, gin.H{"note": note})
}

// Get all
func (h *NoteHandler) GetAll(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	rows, err := h.Service.GetAll(userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed"})
		return
	}
	defer rows.Close()

	var notes []interface{}

	for rows.Next() {
		var (
			id                        int64
			title, content            string
			pinned, archived, deleted int
			createdAt, updatedAt      string
		)

		rows.Scan(&id, &title, &content, &pinned, &archived, &deleted, &createdAt, &updatedAt)

		notes = append(notes, gin.H{
			"id": id, "title": title, "content": content,
			"is_pinned":   pinned == 1,
			"is_archived": archived == 1,
			"is_deleted":  deleted == 1,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}

	ctx.JSON(200, gin.H{"notes": notes})
}

// Update text
func (h *NoteHandler) Update(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	var req UpdateNoteRequest
	if ctx.ShouldBindJSON(&req) != nil {
		ctx.JSON(400, gin.H{"error": "invalid"})
		return
	}

	if err := h.Service.Update(userID, noteID, req.Title, req.Content); err != nil {
		ctx.JSON(500, gin.H{"error": "failed"})
		return
	}

	ctx.JSON(200, gin.H{"status": "updated"})
}

// Update flags (pin/archive/delete)
func (h *NoteHandler) UpdateFlags(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	var req UpdateNoteFlagsRequest
	if ctx.ShouldBindJSON(&req) != nil {
		ctx.JSON(400, gin.H{"error": "invalid"})
		return
	}

	if err := h.Service.UpdateFlags(userID, noteID, req.IsPinned, req.IsArchived, req.IsDeleted); err != nil {
		ctx.JSON(500, gin.H{"error": "failed"})
		return
	}

	ctx.JSON(200, gin.H{"status": "updated"})
}

// Delete permanently
func (h *NoteHandler) DeleteForever(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	if err := h.Service.DeleteForever(userID, noteID); err != nil {
		ctx.JSON(500, gin.H{"error": "failed"})
		return
	}

	ctx.JSON(200, gin.H{"status": "deleted"})
}

func toInt64(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func (h *NoteHandler) Duplicate(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	noteID := toInt64(ctx.Param("id"))

	newID, err := h.Service.Duplicate(userID, noteID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "duplicate failed"})
		return
	}

	ctx.JSON(201, gin.H{"id": newID})
}

func (h *NoteHandler) Metadata(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	rows, err := h.Service.Metadata(userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed"})
		return
	}
	defer rows.Close()

	var notes []interface{}

	for rows.Next() {
		var (
			id                        int64
			title                     string
			updatedAt                 string
			pinned, archived, deleted int
		)

		rows.Scan(&id, &title, &updatedAt, &pinned, &archived, &deleted)

		notes = append(notes, gin.H{
			"id":          id,
			"title":       title,
			"updated_at":  updatedAt,
			"is_pinned":   pinned == 1,
			"is_archived": archived == 1,
			"is_deleted":  deleted == 1,
		})
	}

	ctx.JSON(200, gin.H{"notes": notes})
}

// Search handles POST /notes/search
func (h *NoteHandler) Search(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	var req SearchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "query is required"})
		return
	}

	notes, err := h.Service.Search(userID, req.Query)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "search failed"})
		return
	}

	ctx.JSON(200, gin.H{"results": notes})
}
