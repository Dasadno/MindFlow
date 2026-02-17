// Package llm provides an HTTP client for Ollama-compatible LLM APIs.
//
// Клиент инкапсулирует коммуникацию с языковой моделью через Ollama.
// Поддерживает любую модель, загруженную в Ollama (Gemma, LLaMA, Mistral и т.д.).

package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Client — HTTP-клиент для Ollama API.
type Client struct {
	// Semaphore
	sem chan struct{}
	// BaseURL — endpoint Ollama сервера (default: http://localhost:11434).
	BaseURL string

	// Model — имя модели в Ollama (default: gemma3).
	Model string

	// HTTPClient — HTTP-клиент с таймаутами.
	HTTPClient *http.Client

	// Config — параметры по умолчанию для запросов.
	Config ClientConfig
}

// ClientConfig — конфигурация LLM-клиента.
type ClientConfig struct {
	// DefaultTemperature — температура по умолчанию (0.0–2.0).
	DefaultTemperature float64

	// MaxTokens — максимальное количество токенов в ответе.
	MaxTokens int

	// Timeout — таймаут HTTP-запроса.
	Timeout time.Duration
}

// CompletionRequest — запрос на генерацию текста.
type CompletionRequest struct {
	// SystemPrompt — системный промпт (личность агента).
	SystemPrompt string

	// Messages — история сообщений (user/assistant).
	Messages []Message

	// Temperature — переопределение температуры. nil = дефолт.
	Temperature *float64

	// MaxTokens — переопределение лимита токенов. nil = дефолт.
	MaxTokens *int
}

// Message — одно сообщение в контексте разговора (формат Ollama).
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponse — ответ LLM.
type CompletionResponse struct {
	// Content — сгенерированный текст.
	Content string

	// Model — модель, которая ответила.
	Model string

	// Duration — время запроса.
	Duration time.Duration
}

// NewClient создаёт Ollama-клиент из env-переменных.
// OLLAMA_URL — адрес сервера (default: http://localhost:11434)
// OLLAMA_MODEL — модель (default: gemma3)
func NewClient() *Client {
	baseURL := os.Getenv("OLLAMA_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "gemma3:4b"
	}
	return &Client{
		sem:     make(chan struct{}, 4),
		BaseURL: baseURL,
		Model:   model,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		Config: ClientConfig{
			DefaultTemperature: 0.7,
			MaxTokens:          512,
			Timeout:            5 * time.Minute,
		},
	}
}

// ollamaChatRequest — тело запроса к Ollama /api/chat.
type ollamaChatRequest struct {
	Model    string         `json:"model"`
	Messages []Message      `json:"messages"`
	Stream   bool           `json:"stream"`
	Options  map[string]any `json:"options,omitempty"`
}

// ollamaChatResponse — ответ от Ollama /api/chat (stream: false).
type ollamaChatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Model string `json:"model"`
	Done  bool   `json:"done"`
}

// Complete отправляет запрос в Ollama и возвращает ответ.
func (c *Client) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	c.sem <- struct{}{} // with semaphore
	defer func() { <-c.sem }()
	start := time.Now()

	messages := req.Messages
	if req.SystemPrompt != "" {
		messages = append([]Message{{Role: "system", Content: req.SystemPrompt}}, messages...)
	}

	temp := c.Config.DefaultTemperature
	if req.Temperature != nil {
		temp = *req.Temperature
	}

	body := ollamaChatRequest{
		Model:    c.Model,
		Messages: messages,
		Stream:   false,
		Options:  map[string]any{"temperature": temp},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("Complete marshal: %w", err)
	}

	// 1. Создаем запрос с поддержкой контекста
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/chat", bytes.NewReader(data))
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("Complete create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// 2. Отправляем через Do вместо Post
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		// Здесь теперь будут ловиться таймауты: "context deadline exceeded"
		return CompletionResponse{}, fmt.Errorf("Complete http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CompletionResponse{}, fmt.Errorf("Complete status %d", resp.StatusCode)
	}

	var ollResp ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollResp); err != nil {
		return CompletionResponse{}, fmt.Errorf("Complete decode: %w", err)
	}

	return CompletionResponse{
		Content:  ollResp.Message.Content,
		Model:    ollResp.Model,
		Duration: time.Since(start),
	}, nil
}
