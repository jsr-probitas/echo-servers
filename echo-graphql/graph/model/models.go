package model

type Message struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type EchoResult struct {
	Message *string `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}
