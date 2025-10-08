package ai

type ChatRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}
