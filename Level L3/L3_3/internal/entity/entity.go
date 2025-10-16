package entity

import (
	"time"
)

// NewComment структура нового комментария
type NewComment struct {
	CommentID   string    `json:"comment_id"`
	ParentID    *string   `json:"parent_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
}

// Comment - структура получаемых комментариев
type Comment struct {
	ParentID    string       `json:"parent_id"`
	AllComments []NewComment `json:"all_comments"`
}

// SearchComment поиск по комментам
type SearchComment struct {
	Text string `json:"text"`
}
