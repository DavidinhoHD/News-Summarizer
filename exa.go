package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ExaRequest struct {
	Query string `json:"query"`
	Text  bool   `json:"text"`
	StartPublishedDate time.Time `json:"startPublishingDate"`
}

type ExaResponse struct {
	RequestID          string `json:"requestId"`
	ResolvedSearchType string `json:"resolvedSearchType"`
	Results            []struct {
		Title           string    `json:"title"`
		URL             string    `json:"url"`
		PublishedDate   time.Time `json:"publishedDate"`
		Author          string    `json:"author"`
		Score           float64   `json:"score"`
		ID              string    `json:"id"`
		Image           string    `json:"image"`
		Favicon         string    `json:"favicon"`
		Text            string    `json:"text"`
		Highlights      []string  `json:"highlights"`
		HighlightScores []float64 `json:"highlightScores"`
		Summary         string    `json:"summary"`
		Subpages        []struct {
			ID              string    `json:"id"`
			URL             string    `json:"url"`
			Title           string    `json:"title"`
			Author          string    `json:"author"`
			PublishedDate   time.Time `json:"publishedDate"`
			Text            string    `json:"text"`
			Summary         string    `json:"summary"`
			Highlights      []string  `json:"highlights"`
			HighlightScores []float64 `json:"highlightScores"`
		} `json:"subpages"`
		Extras struct {
			Links []any `json:"links"`
		} `json:"extras"`
	} `json:"results"`
	SearchType  string `json:"searchType"`
	Context     string `json:"context"`
	CostDollars struct {
		Total     float64 `json:"total"`
		BreakDown []struct {
			Search    float64 `json:"search"`
			Contents  int     `json:"contents"`
			Breakdown struct {
				KeywordSearch    int     `json:"keywordSearch"`
				NeuralSearch     float64 `json:"neuralSearch"`
				ContentText      int     `json:"contentText"`
				ContentHighlight int     `json:"contentHighlight"`
				ContentSummary   int     `json:"contentSummary"`
			} `json:"breakdown"`
		} `json:"breakDown"`
		PerRequestPrices struct {
			NeuralSearch125Results      float64 `json:"neuralSearch_1_25_results"`
			NeuralSearch26100Results    float64 `json:"neuralSearch_26_100_results"`
			NeuralSearch100PlusResults  int     `json:"neuralSearch_100_plus_results"`
			KeywordSearch1100Results    float64 `json:"keywordSearch_1_100_results"`
			KeywordSearch100PlusResults int     `json:"keywordSearch_100_plus_results"`
		} `json:"perRequestPrices"`
		PerPagePrices struct {
			ContentText      float64 `json:"contentText"`
			ContentHighlight float64 `json:"contentHighlight"`
			ContentSummary   float64 `json:"contentSummary"`
		} `json:"perPagePrices"`
	} `json:"costDollars"`
}

type ExaGetContentRequest struct {
	Urls []string `json:"urls"`
	Text bool     `json:"text"`
}

type ExaGetContentResponse struct {
	RequestID string `json:"requestId"`
	Results   []struct {
		Title           string    `json:"title"`
		URL             string    `json:"url"`
		PublishedDate   time.Time `json:"publishedDate"`
		Author          string    `json:"author"`
		Score           float64   `json:"score"`
		ID              string    `json:"id"`
		Image           string    `json:"image"`
		Favicon         string    `json:"favicon"`
		Text            string    `json:"text"`
		Highlights      []string  `json:"highlights"`
		HighlightScores []float64 `json:"highlightScores"`
		Summary         string    `json:"summary"`
		Subpages        []struct {
			ID              string    `json:"id"`
			URL             string    `json:"url"`
			Title           string    `json:"title"`
			Author          string    `json:"author"`
			PublishedDate   time.Time `json:"publishedDate"`
			Text            string    `json:"text"`
			Summary         string    `json:"summary"`
			Highlights      []string  `json:"highlights"`
			HighlightScores []float64 `json:"highlightScores"`
		} `json:"subpages"`
		Extras struct {
			Links []any `json:"links"`
		} `json:"extras"`
	} `json:"results"`
	Context  string `json:"context"`
	Statuses []struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Error  struct {
			Tag            string `json:"tag"`
			HTTPStatusCode int    `json:"httpStatusCode"`
		} `json:"error"`
	} `json:"statuses"`
	CostDollars struct {
		Total     float64 `json:"total"`
		BreakDown []struct {
			Search    float64 `json:"search"`
			Contents  int     `json:"contents"`
			Breakdown struct {
				KeywordSearch    int     `json:"keywordSearch"`
				NeuralSearch     float64 `json:"neuralSearch"`
				ContentText      int     `json:"contentText"`
				ContentHighlight int     `json:"contentHighlight"`
				ContentSummary   int     `json:"contentSummary"`
			} `json:"breakdown"`
		} `json:"breakDown"`
		PerRequestPrices struct {
			NeuralSearch125Results      float64 `json:"neuralSearch_1_25_results"`
			NeuralSearch26100Results    float64 `json:"neuralSearch_26_100_results"`
			NeuralSearch100PlusResults  int     `json:"neuralSearch_100_plus_results"`
			KeywordSearch1100Results    float64 `json:"keywordSearch_1_100_results"`
			KeywordSearch100PlusResults int     `json:"keywordSearch_100_plus_results"`
		} `json:"perRequestPrices"`
		PerPagePrices struct {
			ContentText      float64 `json:"contentText"`
			ContentHighlight float64 `json:"contentHighlight"`
			ContentSummary   float64 `json:"contentSummary"`
		} `json:"perPagePrices"`
	} `json:"costDollars"`
}

func exaRequest(exaKey string, q string) ExaResponse {
	// create JSON request body
	reqData := ExaRequest{
		Query: q,
		Text:  true,
		// set start publish date to current day -7
		StartPublishedDate: time.Now().AddDate(0, 0, -7),
	}

	jsonBodyBytes, err := json.Marshal(reqData)
	if err != nil {
		fmt.Println("Error marshalling JSON")
		os.Exit(1)
	}

	// create HTTP request
	req, err := http.NewRequest("POST", "https://api.exa.ai/search", bytes.NewBuffer(jsonBodyBytes))
	if err != nil {
		fmt.Println("Error creating request", err)
		os.Exit(1)
	}
	req.Header.Set("x-api-key", exaKey)
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
	var response ExaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON")
		os.Exit(1)
	}

	return response
}


func exaGetContent(exaKey string, urls []string) ExaGetContentResponse {
	// create JSON request body
	reqData := ExaGetContentRequest{
		Urls: urls,
		Text: true,
	}

	jsonBodyBytes, err := json.Marshal(reqData)
	if err != nil {
		fmt.Println("Error marshalling JSON")
		os.Exit(1)
	}

	// create HTTP request
	req, err := http.NewRequest("POST", "https://api.exa.ai/contents", bytes.NewBuffer(jsonBodyBytes))
	if err != nil {
		fmt.Println("Error creating request", err)
		os.Exit(1)
	}
	req.Header.Set("x-api-key", exaKey)
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
	var response ExaGetContentResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON")
		os.Exit(1)
	}

	return response
}