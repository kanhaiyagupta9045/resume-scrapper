package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err.Error())
	}
}
func main() {

	router := gin.Default()
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("API_KEY")))

	if err != nil {
		log.Fatal(err)
	}
	c := &Client{
		client: client,
		mu:     sync.Mutex{},
	}

	router.POST("/parse-resume", c.Scrape())

	router.Run()
}
