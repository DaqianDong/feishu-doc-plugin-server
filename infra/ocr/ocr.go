package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ocrclient struct {
	apiKey string
	model  string
	url    string
}

// Content 结构体用于处理多模态内容（文本或图片URL）
type Content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type ChatRequest struct {
	Model           string    `json:"model"`
	ReasoningFormat string    `json:"reasoning_format"`
	Messages        []Message `json:"messages"`
	Temperature     float32   `json:"temperature"`
	TopK            float32   `json:"top_k"`
	TopP            float32   `json:"top_p"`
}

// 响应结构体（简化版）
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

var client *ocrclient

func Init(url, model, apiKey string) {
	client = &ocrclient{url: url, model: model, apiKey: apiKey}
}

func OCR(imgUrl, prompt string) (res string, err error) {
	// 1. 构造请求数据
	reqBody := ChatRequest{
		Model: client.model,
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: prompt,
					},
					{
						Type: "image_url",
						ImageURL: &ImageURL{
							URL: imgUrl,
						},
					},
				},
			},
		},
		Temperature:     0,
		TopK:            1,
		TopP:            1,
		ReasoningFormat: "none",
	}

	jsonData, _ := json.Marshal(reqBody)

	// 2. 创建 HTTP Request
	req, _ := http.NewRequest("POST", client.url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 3. 发送请求 (http.Client 默认支持 --location 的重定向逻辑)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 4. 读取并解析结果
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API 返回错误 (状态码 %d): %s\n", resp.StatusCode, string(body))
		return
	}

	var chatResp ChatResponse
	if err = json.Unmarshal(body, &chatResp); err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return
	}

	// 打印 OCR 识别出的 JSON 结果
	if len(chatResp.Choices) > 0 {
		res = chatResp.Choices[0].Message.Content
		fmt.Println("识别结果: " + res)
	}
	return
}
