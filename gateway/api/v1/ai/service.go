package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type AiService struct {
	CommServiceURL string
}

func NewAiService() *AiService {
	commURL := os.Getenv("COMMUNICATION_SERVICE_URL")
	if commURL == "" {
		commURL = "http://localhost:8081/internal/ai/stream"
	}

	return &AiService{CommServiceURL: commURL}
}

func (s *AiService) ChatStream(reqDTO ChatRequest, write io.Writer) error {
	jsonData, err := json.Marshal(reqDTO)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.CommServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(write, resp.Body)

	return err
}
