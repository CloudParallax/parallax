package entities

import (
	"time"
)

// BlogPost represents a blog post entity in the domain
type BlogPost struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []string  `json:"tags"`
	ReadTime    int       `json:"read_time"`
	IsPublished bool      `json:"is_published"`
}

// NewBlogPost creates a new blog post entity
func NewBlogPost(title, slug, summary, content, author string, tags []string) *BlogPost {
	now := time.Now()
	return &BlogPost{
		Title:       title,
		Slug:        slug,
		Summary:     summary,
		Content:     content,
		Author:      author,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsPublished: false,
		ReadTime:    calculateReadTime(content),
	}
}

// Publish marks the blog post as published
func (bp *BlogPost) Publish() {
	bp.IsPublished = true
	bp.PublishedAt = time.Now()
	bp.UpdatedAt = time.Now()
}

// Update updates the blog post content
func (bp *BlogPost) Update(title, summary, content string, tags []string) {
	bp.Title = title
	bp.Summary = summary
	bp.Content = content
	bp.Tags = tags
	bp.UpdatedAt = time.Now()
	bp.ReadTime = calculateReadTime(content)
}

// calculateReadTime estimates reading time based on word count
func calculateReadTime(content string) int {
	// Average reading speed: 200 words per minute
	wordCount := len(content) / 5 // Rough estimate: 5 characters per word
	readTime := wordCount / 200
	if readTime < 1 {
		readTime = 1
	}
	return readTime
}