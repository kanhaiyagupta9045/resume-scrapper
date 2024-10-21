package main

import (
	"sync"

	"github.com/google/generative-ai-go/genai"
)

type Client struct {
	client *genai.Client
	mu     sync.Mutex
}
