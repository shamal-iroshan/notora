package encrypted

type CreateEncryptedNoteDTO struct {
	TitleCiphertext   string `json:"title"`
	ContentCiphertext string `json:"content"`
	TitleNonce        string `json:"title_nonce"`
	ContentNonce      string `json:"content_nonce"`
	NoteSalt          string `json:"note_salt"`
}

type UpdateEncryptedNoteDTO struct {
	TitleCiphertext   string `json:"title"`
	ContentCiphertext string `json:"content"`
	TitleNonce        string `json:"title_nonce"`
	ContentNonce      string `json:"content_nonce"`
	NoteSalt          string `json:"note_salt"`
}
