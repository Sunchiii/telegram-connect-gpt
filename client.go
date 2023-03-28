package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	METHOD = "POST"
)

type Topic struct {
	Model       string    `json:"model"`
	Temperature string    `json:"temperature"`
	Messages    []Message `json:"messages"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func Ask(payload Topic) string {
	client := &http.Client{}
	Payload, _ := json.Marshal(payload)
	newPayload := strings.NewReader(string(Payload))
	req, err := http.NewRequest(METHOD, os.Getenv("CHATGPT_API"), newPayload)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", os.Getenv("CHATGPT_TOKEN"))
	req.Header.Add("OpenAI-Organization", os.Getenv("ORGANIZATION"))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)

}
