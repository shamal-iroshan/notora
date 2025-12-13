package notes

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNoteFlagsRequest struct {
	IsPinned   *bool `json:"is_pinned"`
	IsArchived *bool `json:"is_archived"`
	IsDeleted  *bool `json:"is_deleted"`
}

type SearchNotesRequest struct {
	Query string `json:"query" binding:"required"`
}

type SearchRequest struct {
	Query string `json:"query" binding:"required"`
}
