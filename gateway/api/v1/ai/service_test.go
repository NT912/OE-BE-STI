package ai

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAiService_ChatStream(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseWords := []string{"Hello", " ", "World", "!"}
		for _, word := range responseWords {
			fmt.Fprint(w, word)
			w.(http.Flusher).Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}))

	defer mockServer.Close()

	service := NewAiService()
	service.CommServiceURL = mockServer.URL

	request := ChatRequest{Prompt: "test"}
	var outputBuffer bytes.Buffer
	err := service.ChatStream(request, &outputBuffer)
	if err != nil {
		t.Fatalf("ChatStream returned an unexpected error: %v", err)
	}

	expectedResponse := "Hello World!"
	if outputBuffer.String() != expectedResponse {
		t.Errorf("Expected response '%s', but got '%s'", expectedResponse, outputBuffer.String())
	}

	t.Logf("Test successful, recived: %s", outputBuffer.String())
}
