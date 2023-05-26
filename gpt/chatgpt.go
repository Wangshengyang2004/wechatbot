package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ylsislove/wechatbot/config"
	"github.com/ylsislove/wechatbot/pkg/logger"
)

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Index        int            `json:"index"`
	FinishReason string         `json:"finish_reason"`
	Message      ChatGPTMessage `json:"message"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string           `json:"model"`
	Messages         []ChatGPTMessage `json:"messages"`
	MaxTokens        uint             `json:"max_tokens"`
	Temperature      float64          `json:"temperature"`
	TopP             int              `json:"top_p"`
	FrequencyPenalty int              `json:"frequency_penalty"`
	PresencePenalty  int              `json:"presence_penalty"`
	Stop             []string         `json:"stop"`
	User             string           `json:"user"`
}

type ChatGPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
// -d '{"model": "gpt-3.5-turbo", "messages": [{"role":"system", "content":"You are a assistant"}, {"role":"user", "content": "give me good song"}], "temperature": 0, "max_tokens": 7}'
func ChatCompletions(msg string) (string, error) {
	cfg := config.LoadConfig()
	requestBody := ChatGPTRequestBody{
		Model: cfg.Model,
		Messages: []ChatGPTMessage{
			{Role: "system", Content: "你是一个可以聊天，写作，编程的机器人，如果你收到了中文问题，请用中文回复；英文问题则用英文."},
			{Role: "user", Content: msg},
		},
		MaxTokens:        cfg.MaxTokens,
		Temperature:      cfg.Temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("request gpt json string : %v", string(requestData)))
	req, err := http.NewRequest("POST", cfg.BaseUrl+"chat/completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}
	// fmt.Println(cfg.BaseUrl + "chat/completions")

	apiKey := config.LoadConfig().ApiKey
	proxy := config.LoadConfig().Proxy
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	// client := &http.Client{Timeout: cfg.RequestTimeout * time.Second}

	var client *http.Client
	if len(proxy) == 0 {
		client = &http.Client{Timeout: cfg.RequestTimeout * time.Second}
	} else {
		proxyAddr, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyAddr),
			},
			Timeout: cfg.RequestTimeout * time.Second,
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("请求GTP出错了，gpt api status code not equals 200,code is %d ,details:  %v ", response.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("response gpt json string : %v", string(body)))

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v.Message.Content
			break
		}
	}
	logger.Info(fmt.Sprintf("gpt response text: %s ", reply))
	return reply, nil
}
