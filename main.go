package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	systemPrompt = `You are a helpful assistant that generates search queries based on user questions. Only generate one search query.`
	userQuestion = `What is the recent news in physics this week`
	contentSummaryQuerry = `You are a helpful assistant that briefly summarizes the content of a webpage. Summarize the users input.`
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
	exaKey := os.Getenv("EXA_KEY")
	llmAPI := os.Getenv("LLM_API")



	// generate search query based on user question
	searchQuery := ollamaRequest(llmAPI, fmt.Sprintf("%s\n%s", systemPrompt, userQuestion))
	cleandSearchQuery := removeThinkingTags(searchQuery.Response)

	// search for articles 
	articles := exaSearchRequest(exaKey, cleandSearchQuery)
	var results []string
	for _, res:= range articles.Results {
		results = append(results, res.URL)
	}

	// get content of articles
	contentResponse := exaGetContent(exaKey, results)
	var content []string
	for _, res := range contentResponse.Results {
		content = append(content, res.Text)
	}

	// summarize content via LLM
	summaryResponse := ollamaRequest(llmAPI, fmt.Sprintf("%s\n%s", contentSummaryQuerry, strings.Join(content, "\n")))
	cleandSummaryResponse := removeThinkingTags(summaryResponse.Response)

	fmt.Println(cleandSummaryResponse)
}