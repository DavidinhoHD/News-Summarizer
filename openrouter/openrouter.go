package openrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type Request struct {
	Model 	string
	Message []Message
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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

	//client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println("request error")
	// 	return nil, err
	// }
	// defer resp.Body.Close()

	resp := &http.Response{}

	r_body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(r_body))


	return resp, nil
}
