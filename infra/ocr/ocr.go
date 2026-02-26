package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const url = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

type ocrclient struct {
	apiKey string
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
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
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

func Init(apiKey string) {
	client = &ocrclient{apiKey: apiKey}
}

func OCR(imgUrl, prompt string) (res string, err error) {
	// 1. 构造请求数据
	reqBody := ChatRequest{
		Model: "qwen-vl-ocr-2025-11-20",
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "image_url",
						ImageURL: &ImageURL{
							URL: imgUrl,
						},
					},
					{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	// 2. 创建 HTTP Request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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
