package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type AiService struct{}

func NewAiService() *AiService {
	return &AiService{}
}

func (s *AiService) GetAIStream(writer io.Writer) error {
	responseWords := []string{"Đây ", "là ", "câu ", "trả ", "lời ", "được ", "stream ", "từ ", "AI ", "thật ", "sự."}

	flusher, ok := writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	for _, word := range responseWords {
		_, err := fmt.Fprint(writer, word)
		if err != nil {
			return err
		}
		flusher.Flush()
		log.Printf("Sent Chunk: %s", word)
		time.Sleep(200 * time.Millisecond)
	}
	log.Println("Finished streaming AI response.")
	return nil
}
