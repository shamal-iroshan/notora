package model

type EncryptedNoteMetadata struct {
	ID              int64  `json:"id"`
	TitleCiphertext string `json:"title"`
	TitleNonce      string `json:"title_nonce"`
	NoteSalt        string `json:"note_salt"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type EncryptedNoteResponse struct {
	ID                int64  `json:"id"`
	TitleCiphertext   string `json:"title"`
	ContentCiphertext string `json:"content"`
	TitleNonce        string `json:"title_nonce"`
	ContentNonce      string `json:"content_nonce"`
	NoteSalt          string `json:"note_salt"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type CreateEncryptedNoteInput struct {
	TitleCiphertext   string
	ContentCiphertext string
	TitleNonce        string
	ContentNonce      string
	NoteSalt          string
}

type UpdateEncryptedNoteInput struct {
	TitleCiphertext   string
	ContentCiphertext string
	TitleNonce        string
	ContentNonce      string
	NoteSalt          string
}
