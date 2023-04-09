package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/b0noi/go-utils/v2/gcp"
)

func Message(message string) (GptChatCompletionMessage, error) {
	apiKey, err := gcp.AccessSecretVersion("projects/243625208291/secrets/open-key/versions/1")
	if err != nil {
		return GptChatCompletionMessage{}, err
	}

	url := "https://api.openai.com/v1/chat/completions"

	requestBody, err := prepareGPT4RequestBody(message)
	if err != nil {
		return GptChatCompletionMessage{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return GptChatCompletionMessage{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GptChatCompletionMessage{}, err
	}
	defer resp.Body.Close()
	// var response map[string]interface{}
	var response GptChatCompletionMessage
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return GptChatCompletionMessage{}, err
	}
	return response, nil
}

func prepareGPT4RequestBody(message string) ([]byte, error) {
	gptMessages := make([]map[string]string, 1)
	gptMessages[0] = map[string]string{
		"role":    "user",
		"content": message,
	}

	// Marshal the request body for GPT-4
	requestBody, err := json.Marshal(map[string]interface{}{
		"messages":   gptMessages,
		"max_tokens": 2000,
		"n":          1,
		"model":      "gpt-4",
	})

	if err != nil {
		return nil, err
	}

	return requestBody, nil
}
