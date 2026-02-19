package external

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "task-api/internal/models"
    "time"
)

type APIClient struct {
    client *http.Client
}

func NewAPIClient() *APIClient {
    return &APIClient{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

//получает задачи с внешнего апи
func (c *APIClient) GetExternalTodos() ([]models.ExternalTodo, error) {
    resp, err := c.client.Get("https://jsonplaceholder.typicode.com/todos")
    if err != nil {
        return nil, fmt.Errorf("failed to fetch todos: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    var todos []models.ExternalTodo
    if err := json.Unmarshal(body, &todos); err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }
    
    return todos, nil
}

// пост на внешнем апи
func (c *APIClient) CreateExternalPost(post models.CreatePostRequest) (*models.CreatePostResponse, error) {
    jsonData, err := json.Marshal(post)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal post: %w", err)
    }
    
    resp, err := c.client.Post(
        "https://jsonplaceholder.typicode.com/posts",
        "application/json",
        bytes.NewReader(jsonData),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create post: %w", err)
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    var createdPost models.CreatePostResponse
    if err := json.Unmarshal(body, &createdPost); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }
    
    return &createdPost, nil
}