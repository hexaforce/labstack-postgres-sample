package schemas

// -- Requests and Responses

type CommentRequest struct {
	Comment string `json:"comment"`
	UserID  int64  `json:"user_id"`
}

type CommentResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error"`
	Comment *models.Comment `json:"comment"`
}

type CommentsResponse struct {
	Success  bool              `json:"success"`
	Error    string            `json:"error"`
	Comments []*models.Comment `json:"comments"`
}
