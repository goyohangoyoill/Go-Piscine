package main

type Success struct {
	Success bool   `json:"success"`
	Error   bool   `json:"error"`
	Content string `json:"content"`
}
