package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
)

func (cl *Client) Scrape() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("resume")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "please attach the resume"})
			return
		}
		fmt.Println(file.Filename)
		fileReader, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not open file"})
			return
		}
		defer fileReader.Close()

		cl.mu.Lock()
		defer cl.mu.Unlock()

		genaifile, err := cl.client.UploadFile(context.Background(), strings.TrimSuffix(file.Filename, ".pdf"), fileReader, nil)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer cl.client.DeleteFile(context.Background(), genaifile.Name)
		model := cl.client.GenerativeModel("gemini-1.5-flash")

		resp, err := model.GenerateContent(context.Background(),
			genai.Text("Extract the sections (education, skills, experience) from my attached resume and return them in a concise JSON format."),
			genai.FileData{URI: genaifile.URI})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		prompt_response := printresponse(resp)
		c.JSON(http.StatusOK, prompt_response)

	}
}
func printresponse(resp *genai.GenerateContentResponse) *genai.Part {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				return &part
			}
		}
	}
	return nil
}
