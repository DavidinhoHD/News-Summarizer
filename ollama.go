package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

// remove <thinking> tags from LLM response
func removeThinkingTags(response string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleaned := re.ReplaceAllString(response, "")
	return strings.TrimSpace(cleaned)
}

// make a request to an LLM and return the response
func ollamaRequest(llmKey string) OllamaResponse {
	// create JSON request body
	jsonBody := OllamaRequest{
		Model:  "qwen3:8b",
		Prompt: systemPrompt + userQuestion,
		Stream: false,
	}
	jsonBodyBytes, err := json.Marshal(jsonBody)
	if err != nil {
		fmt.Println("Error marshalling JSON")
		os.Exit(1)
	}

	// create HTTP request
	req, err := http.NewRequest("POST", llmKey, bytes.NewBuffer(jsonBodyBytes))
	if err != nil {
		fmt.Println("Error creating request", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error making request", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body")
		os.Exit(1)
	}

	// parse response body as JSON
	var response OllamaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON")
		os.Exit(1)
	}

	return response
}
