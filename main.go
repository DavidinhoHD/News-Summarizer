package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"DavidinhoHD/News-Summarizer/openrouter"

	"github.com/joho/godotenv"
)

var (
	systemPrompt = `You are a helpful assistant that generates search queries based on user questions. Only generate one search query.`
	//userQuestion = `What is the recent news in physics this week`
	contentSummaryQuery = `You are a helpful assistant that briefly summarizes the content of a webpage. Summarize the users input.`
)

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your news topic: ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func main() {

	req := openrouter.Request{
		Model: "z-ai/glm-4.5-air:free",
		Message: []openrouter.Message{
			{
				Role:	"system",
				Content: "you are a freandly assistant that answares the users questions",
			},
			{
				Role:	 "User",
				Content: "hello",
			},

		},
	}


	 err := godotenv.Load()
	 if err != nil {
	 	fmt.Println("Error loading .env file")
	 	os.Exit(1)
	}
	openrouter_key := os.Getenv("openrouter_key")

	_, err = openrouter.MakeOpenRouterRequest(req, openrouter_key)
	if err != nil {
		fmt.Println(err)
	}

	// exaKey := os.Getenv("EXA_KEY")
	// llmAPI := os.Getenv("LLM_API")

	// userQuestion := getUserInput()


	// // generate search query based on user question
	// searchQuery := ollamaRequest(llmAPI, fmt.Sprintf("%s\n%s", systemPrompt, userQuestion))
	// cleanedSearchQuery := removeThinkingTags(searchQuery.Response)

	// // search for articles
	// articles := exaSearchRequest(exaKey, cleanedSearchQuery)
	// var results []string
	// for _, res:= range articles.Results {
	// 	results = append(results, res.URL)
	// }

	// // get content of articles
	// contentResponse := exaGetContent(exaKey, results)
	// var content []string
	// for _, res := range contentResponse.Results {
	// 	content = append(content, res.Text)
	// }

	// // summarize content via LLM
	// summaryResponse := ollamaRequest(llmAPI, fmt.Sprintf("%s\n%s", contentSummaryQuery, strings.Join(content, "\n")))
	// cleanedSummaryResponse := removeThinkingTags(summaryResponse.Response)

	// fmt.Println(cleanedSummaryResponse)
}
