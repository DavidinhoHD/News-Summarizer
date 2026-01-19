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
// Load API key
 err := godotenv.Load()
 if err != nil {
 	fmt.Println("Error loading .env file")
 	os.Exit(1)
}
openrouter_key := os.Getenv("openrouter_key")

	req := openrouter.Request{
		Model: "z-ai/glm-4.5-air:free",
		Messages: []openrouter.Message{
			{
				Role:	"system",
				Content: "you are a freandly assistant that answares the users questions",
			},
			{
				Role:	 "user",
				Content: "hello",
			},

		},
	}


	_, err = openrouter.MakeOpenRouterRequest(req, openrouter_key)
	if err != nil {
		fmt.Println(err)
	}


}
