package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
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

func (s *AiService) ChatStream(reqDTO ChatRequest, writer io.Writer) error {
	jsonData, err := json.Marshal(reqDTO)
	if err != nil {
		return err
	}

	// Set timeout cho việc đợi response
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

	buf := make([]byte, 1024)

	flusher, ok := writer.(http.Flusher)
	if !ok {
		log.Println("Warning: ResponseWriter does not support flushing.")
	}
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			log.Printf("Gateway received chunk: %s", string(buf[:n]))
			if _, writeErr := writer.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
			if ok {
				flusher.Flush()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return err
}
