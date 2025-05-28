package dto

import (
	"time"

	"github.com/cloudparallax/parallax/internal/domain/entities"
)

// CreateBlogPostRequest represents the request to create a blog post
type CreateBlogPostRequest struct {
	Title   string   `json:"title" validate:"required,min=1,max=200"`
	Summary string   `json:"summary" validate:"required,min=1,max=500"`
	Content string   `json:"content" validate:"required,min=1"`
	Author  string   `json:"author" validate:"required,min=1,max=100"`
	Tags    []string `json:"tags" validate:"dive,min=1,max=50"`
}

// UpdateBlogPostRequest represents the request to update a blog post
type UpdateBlogPostRequest struct {
	Title   string   `json:"title" validate:"required,min=1,max=200"`
	Summary string   `json:"summary" validate:"required,min=1,max=500"`
	Content string   `json:"content" validate:"required,min=1"`
	Tags    []string `json:"tags" validate:"dive,min=1,max=50"`
}

// BlogPostResponse represents the response for a blog post
type BlogPostResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []string  `json:"tags"`
	ReadTime    int       `json:"read_time"`
	IsPublished bool      `json:"is_published"`
}

// BlogPostListResponse represents the response for a list of blog posts
type BlogPostListResponse struct {
	Posts []BlogPostSummaryResponse `json:"posts"`
	Meta  *MetaResponse             `json:"meta,omitempty"`
}

// BlogPostSummaryResponse represents a summary response for blog posts in lists
type BlogPostSummaryResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Summary     string     `json:"summary"`
	Author      string     `json:"author"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	Tags        []string   `json:"tags"`
	ReadTime    int        `json:"read_time"`
	IsPublished bool       `json:"is_published"`
}

// GetBlogPostsRequest represents query parameters for getting blog posts
type GetBlogPostsRequest struct {
	AuthorID    string `query:"author_id"`
	Tags        string `query:"tags"` // comma-separated tags
	IsPublished *bool  `query:"is_published"`
	FromDate    string `query:"from_date" validate:"omitempty,datetime=2006-01-02"`
	ToDate      string `query:"to_date" validate:"omitempty,datetime=2006-01-02"`
	SearchQuery string `query:"q"`
	Page        int    `query:"page" validate:"min=1"`
	Limit       int    `query:"limit" validate:"min=1,max=100"`
	SortBy      string `query:"sort_by" validate:"omitempty,oneof=created_at updated_at published_at title"`
	SortOrder   string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// SearchBlogPostsRequest represents the request for searching blog posts
type SearchBlogPostsRequest struct {
	Query string `query:"q" validate:"required,min=1"`
	Page  int    `query:"page" validate:"min=1"`
	Limit int    `query:"limit" validate:"min=1,max=100"`
}

// MetaResponse represents pagination metadata
type MetaResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ToBlogPostResponse converts a blog post entity to response DTO
func ToBlogPostResponse(post *entities.BlogPost) *BlogPostResponse {
	var publishedAt *time.Time
	if post.IsPublished && !post.PublishedAt.IsZero() {
		publishedAt = &post.PublishedAt
	}

	return &BlogPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Summary:     post.Summary,
		Content:     post.Content,
		Author:      post.Author,
		PublishedAt: publishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Tags:        post.Tags,
		ReadTime:    post.ReadTime,
		IsPublished: post.IsPublished,
	}
}

// ToBlogPostSummaryResponse converts a blog post entity to summary response DTO
func ToBlogPostSummaryResponse(post *entities.BlogPost) *BlogPostSummaryResponse {
	var publishedAt *time.Time
	if post.IsPublished && !post.PublishedAt.IsZero() {
		publishedAt = &post.PublishedAt
	}

	return &BlogPostSummaryResponse{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Summary:     post.Summary,
		Author:      post.Author,
		PublishedAt: publishedAt,
		CreatedAt:   post.CreatedAt,
		Tags:        post.Tags,
		ReadTime:    post.ReadTime,
		IsPublished: post.IsPublished,
	}
}

// ToBlogPostListResponse converts a list of blog post entities to list response DTO
func ToBlogPostListResponse(posts []*entities.BlogPost, meta *MetaResponse) *BlogPostListResponse {
	summaries := make([]BlogPostSummaryResponse, len(posts))
	for i, post := range posts {
		summaries[i] = *ToBlogPostSummaryResponse(post)
	}

	return &BlogPostListResponse{
		Posts: summaries,
		Meta:  meta,
	}
}