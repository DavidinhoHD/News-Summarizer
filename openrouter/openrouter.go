package openrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ReasoningEffort string
const (
	ReasoningEffortNone ReasoningEffort = "none"
	ReasoningEffortLow ReasoningEffort = "low"
	ReasoningEffortMedium ReasoningEffort = "medium"
	ReasoningEffortHigh ReasoningEffort = "high"
)

const (
	baseURL = 	"https://openrouter.ai/api/v1/chat/completions"
	defaultTimeout = 30 * time.Second
)


type Request struct {
	Model 		string		`json:"model"`
	Messages 	[]Message	`json:"messages"`
	Stream 		bool		`json:"stream,omitempty"`
	Temperature float64   	`json:"temperature,omitempty"`
    MaxTokens   int       	`json:"max_tokens,omitempty"`
	Reasoning   Reasoning 	`json:"reasoning,omitempty"`
}

type Reasoning struct {
	Effort ReasoningEffort
}


type OpenRouterResponse struct {
    ID       string   `json:"id"`
    Provider string   `json:"provider"`
    Model    string   `json:"model"`
    Object   string   `json:"object"`
    Created  int64    `json:"created"`
    Choices  []Choice `json:"choices"`
    Usage    Usage    `json:"usage"`
    Error	 *APIError `json:"error,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Type string `json:"type"`
	Code string `json:"code"`
}

type Choice struct {
    Logprobs           interface{}     `json:"logprobs"`
    FinishReason       string          `json:"finish_reason"`
    NativeFinishReason string          `json:"native_finish_reason"`
    Index              int             `json:"index"`
    Message            Message         `json:"message"`
}

type Message struct {
    Role             string            `json:"role"`
    Content          string            `json:"content"`
    Refusal          interface{}       `json:"refusal"`
    Reasoning        string            `json:"reasoning,omitempty"`
    ReasoningDetails []ReasoningDetail `json:"reasoning_details,omitempty"`
}

type ReasoningDetail struct {
    Format string `json:"format"`
    Index  int    `json:"index"`
    Type   string `json:"type"`
    Text   string `json:"text"`
}

type Usage struct {
    PromptTokens            int                  `json:"prompt_tokens"`
    CompletionTokens        int                  `json:"completion_tokens"`
    TotalTokens             int                  `json:"total_tokens"`
    Cost                    float64              `json:"cost"`
    IsByok                  bool                 `json:"is_byok"`
    PromptTokensDetails     PromptTokensDetails  `json:"prompt_tokens_details"`
    CostDetails             CostDetails          `json:"cost_details"`
    CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
}

type PromptTokensDetails struct {
    CachedTokens int `json:"cached_tokens"`
    AudioTokens  int `json:"audio_tokens"`
    VideoTokens  int `json:"video_tokens"`
}

type CostDetails struct {
    UpstreamInferenceCost            interface{} `json:"upstream_inference_cost"`
    UpstreamInferencePromptCost      float64     `json:"upstream_inference_prompt_cost"`
    UpstreamInferenceCompletionsCost float64     `json:"upstream_inference_completions_cost"`
}

type CompletionTokensDetails struct {
    ReasoningTokens int `json:"reasoning_tokens"`
    ImageTokens     int `json:"image_tokens"`
}


// resoning effort validator
func (r *ReasoningEffort) isValid() bool {
	switch *r {
		case ReasoningEffortNone:
			return true
		case ReasoningEffortLow:
			return true
		case ReasoningEffortMedium:
			return true
		case ReasoningEffortHigh:
			return true
		default:
			return false
	}
}


func MakeOpenRouterRequest(r Request, apiKey string) (OpenRouterResponse, error) {
	var response OpenRouterResponse

	if apiKey == "" {
		return OpenRouterResponse{}, fmt.Errorf("API key is required")
	}

	// validate reasoning effort if not empty
	if r.Reasoning.Effort != "" && !r.Reasoning.Effort.isValid() {
		return OpenRouterResponse{}, fmt.Errorf("invalid reasoning effort")
	}

	// Marshal request body
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return OpenRouterResponse{}, err
	}

	// create Http request
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return OpenRouterResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	// set request Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return OpenRouterResponse{}, err
	}
	defer resp.Body.Close()


	r_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OpenRouterResponse{}, err
	}

	// turn JSON response to OpenRouterResponse
	err = json.Unmarshal(r_body, &response)
	if err != nil {
		return OpenRouterResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	// check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return OpenRouterResponse{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}


	// check if content is empty
	if response.Choices[0].Message.Content == "" {
		return OpenRouterResponse{}, fmt.Errorf("no content returned")
	}


	return response, nil
}
