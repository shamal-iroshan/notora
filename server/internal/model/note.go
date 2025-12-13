package model

type Note struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	IsPinned   bool   `json:"is_pinned"`
	IsArchived bool   `json:"is_archived"`
	IsDeleted  bool   `json:"is_deleted"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
