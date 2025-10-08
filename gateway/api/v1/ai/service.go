package ai

type AiService struct {
	CommServiceURL string
}

func NewAiService() *AiService {
	return &AiService{}
}
