package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	systemPrompt = `You are a helpful assistant that generates search queries based on user questions. Only generate one search query.`
	userQuestion = `What is the recent news in physics this week`
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
	exaKey := os.Getenv("EXA_KEY")
	llmKey := os.Getenv("LLM_KEY")

	// fmt.Println(llmKey)
	ollamaRequest(llmKey)
	exaRequest(exaKey, userQuestion)
}