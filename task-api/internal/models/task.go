package models

type Task struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Done   bool   `json:"done"`
    UserID int    `json:"userId,omitempty"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
}

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type DetailedErrorResponse struct {
    Error     string            `json:"error"`
    Details   string            `json:"details,omitempty"`
    Errors    []ValidationError `json:"validation_errors,omitempty"`
    RequestID string            `json:"request_id,omitempty"`
}

type SuccessResponse struct {
    Message string `json:"message"`
    ID      int    `json:"id,omitempty"`
    Updated bool   `json:"updated,omitempty"`
    Deleted bool   `json:"deleted,omitempty"`
}

type ExternalTodo struct {
    ID        int    `json:"id"`
    UserID    int    `json:"userId"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

type CreatePostRequest struct {
    Title  string `json:"title"`
    Body   string `json:"body"`
    UserID int    `json:"userId"`
}

type CreatePostResponse struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Body   string `json:"body"`
    UserID int    `json:"userId"`
}