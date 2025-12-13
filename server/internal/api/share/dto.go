package share

// ----- RESPONSE DTOs -----

type ShareResponse struct {
	ShareURL string `json:"share_url"`
}

type DisableShareResponse struct {
	Status string `json:"status"`
}

type PublicSharedNoteResponse struct {
	Note interface{} `json:"note"` // Could replace with a concrete Note model later
}
