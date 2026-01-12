package openrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type Request struct {
	Model 	string		`json:"model"`
	Messages []Message	`json:"messages"`
	Stream 	bool		`json:"stream"`
}

type OpenRouterResponse struct {
    ID       string   `json:"id"`
    Provider string   `json:"provider"`
    Model    string   `json:"model"`
    Object   string   `json:"object"`
    Created  int64    `json:"created"`
    Choices  []Choice `json:"choices"`
    Usage    Usage    `json:"usage"`
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

func MakeOpenRouterRequest(r Request, apiKey string) (*http.Response, error) {
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	// create Http request
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// set request Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()


	r_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(r_body))


	return resp, nil
}
