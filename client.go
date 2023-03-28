package main

var (
	METHOD = "POST"
)

type Topic struct {
	Model       string    `json:"model"`
	Temperature string    `json:"temperature"`
	Message     []Message `json:"message"`
}
type Message struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}
