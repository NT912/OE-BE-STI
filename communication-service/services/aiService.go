package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"communication-service/configs"
)

const OpenAIEndpoint = "https://api.openai.com/v1/chat/completions"

type AiService struct{}

func NewAiService() *AiService {
	return &AiService{}
}

func (s *AiService) GetAIStream(prompt string, writer io.Writer) error {
	requestBody := OpenAIRequest{
		Model: "gpt-5",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("Error marshalling openai request: %w", err)
	}

	req, err := http.NewRequest("POST", OpenAIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error creating openai request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+configs.Env.OpenAIAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error performing openai request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("Openai api returned an error: %s - %s", resp.Status, string(body))
		log.Println(errorMsg)
		fmt.Fprint(writer, errorMsg)

		return errors.New(errorMsg)
	}

	_, err = io.Copy(writer, resp.Body)
	return err
}
