package main

import (
	"fmt"
	"os"

	"DavidinhoHD/News-Summarizer/openrouter"

	"github.com/joho/godotenv"
)

var (
	systemPromt = "you are a freandly assistant that answares the users questions"
)


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
				Content: systemPromt,
			},
			{
				Role:	 "user",
				Content: "hello",
			},

		},
	}


	resp, err := openrouter.MakeOpenRouterRequest(req, openrouter_key)
	if err != nil {
		fmt.Println(err)
	}

	m := resp.Choices[0].Message.Content
	fmt.Println(m)

}
