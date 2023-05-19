package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var accessToken string = os.Getenv("OPENAI_ACCESS_TOKEN")
	var question string

	fmt.Println("Question: ")
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _ := reader.ReadString('\n')
		if len(strings.TrimSpace(line)) == 0 {
			break
		}

		var lines []string
		lines = append(lines, line)

		for _, line := range lines {
			question += line
		}
	}

	question = strings.ReplaceAll(question, "\n", "\\n")
	question = strings.ReplaceAll(question, "\"", "\\\"")

	URL := "https://api.openai.com/v1/chat/completions"
	METHOD := "POST"

	PAYLOAD := strings.NewReader(`{
    "model": "gpt-3.5-turbo",
    "messages": [
        {
            "role": "user",
            "content": "` + question + `"
        }
    ],
		"stream": true,
    "temperature": 0.8
}`)

	client := &http.Client{}
	req, err := http.NewRequest(METHOD, URL, PAYLOAD)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Answer: ")

	stream := bufio.NewReader(res.Body)
	for {
		line, err := stream.ReadString('\n')

		type ChatCompletion struct {
			ID      string `json:"id"`
			Object  string `json:"object"`
			Created int    `json:"created"`
			Model   string `json:"model"`
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
				Index        int    `json:"index"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}

		line = strings.ReplaceAll(line, "data: ", "")

		var chatCompletion ChatCompletion = ChatCompletion{}
		json.Unmarshal([]byte(line), &chatCompletion)

		if len(chatCompletion.ID) > 0 {
			fmt.Print(chatCompletion.Choices[0].Delta.Content)
		}

		if err != nil {
			break
		}
	}
	res.Body.Close()

	fmt.Println()
}
