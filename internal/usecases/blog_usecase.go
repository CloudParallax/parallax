package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
)

// BlogUseCase defines the interface for blog post business logic
type BlogUseCase interface {
	CreatePost(ctx context.Context, input CreatePostInput) (*entities.BlogPost, error)
	GetPostByID(ctx context.Context, id string) (*entities.BlogPost, error)
	GetPostBySlug(ctx context.Context, slug string) (*entities.BlogPost, error)
	GetAllPosts(ctx context.Context, filters GetPostsInput) ([]*entities.BlogPost, error)
	UpdatePost(ctx context.Context, id string, input UpdatePostInput) (*entities.BlogPost, error)
	DeletePost(ctx context.Context, id string) error
	PublishPost(ctx context.Context, id string) (*entities.BlogPost, error)
	UnpublishPost(ctx context.Context, id string) (*entities.BlogPost, error)
	GetPostsByTags(ctx context.Context, tags []string) ([]*entities.BlogPost, error)
	GetPublishedPosts(ctx context.Context, limit, offset int) ([]*entities.BlogPost, error)
	SearchPosts(ctx context.Context, query string, limit, offset int) ([]*entities.BlogPost, error)
}

// blogUseCase implements BlogUseCase interface
type blogUseCase struct {
	blogRepo repositories.BlogPostRepository
}

// NewBlogUseCase creates a new blog use case
func NewBlogUseCase(blogRepo repositories.BlogPostRepository) BlogUseCase {
	return &blogUseCase{
		blogRepo: blogRepo,
	}
}

// CreatePostInput represents input for creating a blog post
type CreatePostInput struct {
	Title   string   `json:"title" validate:"required,min=1,max=200"`
	Summary string   `json:"summary" validate:"required,min=1,max=500"`
	Content string   `json:"content" validate:"required,min=1"`
	Author  string   `json:"author" validate:"required,min=1,max=100"`
	Tags    []string `json:"tags" validate:"dive,min=1,max=50"`
}

// UpdatePostInput represents input for updating a blog post
type UpdatePostInput struct {
	Title   string   `json:"title" validate:"required,min=1,max=200"`
	Summary string   `json:"summary" validate:"required,min=1,max=500"`
	Content string   `json:"content" validate:"required,min=1"`
	Tags    []string `json:"tags" validate:"dive,min=1,max=50"`
}

// GetPostsInput represents input for getting blog posts with filters
type GetPostsInput struct {
	AuthorID    string
	Tags        []string
	IsPublished *bool
	FromDate    *time.Time
	ToDate      *time.Time
	SearchQuery string
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string
}

// CreatePost creates a new blog post
func (uc *blogUseCase) CreatePost(ctx context.Context, input CreatePostInput) (*entities.BlogPost, error) {
	// Validate input
	if err := uc.validateCreateInput(input); err != nil {
		return nil, err
	}

	// Generate slug from title
	slug := uc.generateSlug(input.Title)

	// Check if slug already exists
	existingPost, _ := uc.blogRepo.GetBySlug(ctx, slug)
	if existingPost != nil {
		slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	}

	// Create new blog post entity
	post := entities.NewBlogPost(
		input.Title,
		slug,
		input.Summary,
		input.Content,
		input.Author,
		input.Tags,
	)

	// Generate ID (in a real app, this might be done by the database)
	post.ID = uc.generateID()

	// Save to repository
	if err := uc.blogRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create blog post: %w", err)
	}

	return post, nil
}

// GetPostByID retrieves a blog post by its ID
func (uc *blogUseCase) GetPostByID(ctx context.Context, id string) (*entities.BlogPost, error) {
	if id == "" {
		return nil, errors.New("post ID is required")
	}

	post, err := uc.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog post: %w", err)
	}

	return post, nil
}

// GetPostBySlug retrieves a blog post by its slug
func (uc *blogUseCase) GetPostBySlug(ctx context.Context, slug string) (*entities.BlogPost, error) {
	if slug == "" {
		return nil, errors.New("post slug is required")
	}

	post, err := uc.blogRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog post: %w", err)
	}

	return post, nil
}

// GetAllPosts retrieves all blog posts with optional filters
func (uc *blogUseCase) GetAllPosts(ctx context.Context, filters GetPostsInput) ([]*entities.BlogPost, error) {
	repoFilters := repositories.BlogPostFilters{
		AuthorID:    filters.AuthorID,
		Tags:        filters.Tags,
		IsPublished: filters.IsPublished,
		SearchQuery: filters.SearchQuery,
		Limit:       filters.Limit,
		Offset:      filters.Offset,
		SortBy:      filters.SortBy,
		SortOrder:   filters.SortOrder,
	}

	if filters.FromDate != nil {
		fromDate := filters.FromDate.Format("2006-01-02")
		repoFilters.FromDate = &fromDate
	}

	if filters.ToDate != nil {
		toDate := filters.ToDate.Format("2006-01-02")
		repoFilters.ToDate = &toDate
	}

	posts, err := uc.blogRepo.GetAll(ctx, repoFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog posts: %w", err)
	}

	return posts, nil
}

// UpdatePost updates an existing blog post
func (uc *blogUseCase) UpdatePost(ctx context.Context, id string, input UpdatePostInput) (*entities.BlogPost, error) {
	if id == "" {
		return nil, errors.New("post ID is required")
	}

	// Validate input
	if err := uc.validateUpdateInput(input); err != nil {
		return nil, err
	}

	// Get existing post
	post, err := uc.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog post: %w", err)
	}

	// Update post
	post.Update(input.Title, input.Summary, input.Content, input.Tags)

	// Save changes
	if err := uc.blogRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update blog post: %w", err)
	}

	return post, nil
}

// DeletePost deletes a blog post
func (uc *blogUseCase) DeletePost(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("post ID is required")
	}

	// Check if post exists
	_, err := uc.blogRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get blog post: %w", err)
	}

	// Delete post
	if err := uc.blogRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete blog post: %w", err)
	}

	return nil
}

// PublishPost publishes a blog post
func (uc *blogUseCase) PublishPost(ctx context.Context, id string) (*entities.BlogPost, error) {
	if id == "" {
		return nil, errors.New("post ID is required")
	}

	// Get post
	post, err := uc.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog post: %w", err)
	}

	// Publish post
	post.Publish()

	// Save changes
	if err := uc.blogRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to publish blog post: %w", err)
	}

	return post, nil
}

// UnpublishPost unpublishes a blog post
func (uc *blogUseCase) UnpublishPost(ctx context.Context, id string) (*entities.BlogPost, error) {
	if id == "" {
		return nil, errors.New("post ID is required")
	}

	// Get post
	post, err := uc.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog post: %w", err)
	}

	// Unpublish post
	post.IsPublished = false
	post.UpdatedAt = time.Now()

	// Save changes
	if err := uc.blogRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to unpublish blog post: %w", err)
	}

	return post, nil
}

// GetPostsByTags retrieves blog posts by tags
func (uc *blogUseCase) GetPostsByTags(ctx context.Context, tags []string) ([]*entities.BlogPost, error) {
	if len(tags) == 0 {
		return nil, errors.New("at least one tag is required")
	}

	posts, err := uc.blogRepo.GetByTags(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to get blog posts by tags: %w", err)
	}

	return posts, nil
}

// GetPublishedPosts retrieves only published blog posts
func (uc *blogUseCase) GetPublishedPosts(ctx context.Context, limit, offset int) ([]*entities.BlogPost, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := uc.blogRepo.GetPublished(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get published blog posts: %w", err)
	}

	return posts, nil
}

// SearchPosts searches for blog posts
func (uc *blogUseCase) SearchPosts(ctx context.Context, query string, limit, offset int) ([]*entities.BlogPost, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	filters := repositories.BlogPostFilters{
		SearchQuery: query,
		Limit:       limit,
		Offset:      offset,
	}

	posts, err := uc.blogRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search blog posts: %w", err)
	}

	return posts, nil
}

// Helper methods

func (uc *blogUseCase) validateCreateInput(input CreatePostInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(input.Summary) == "" {
		return errors.New("summary is required")
	}
	if strings.TrimSpace(input.Content) == "" {
		return errors.New("content is required")
	}
	if strings.TrimSpace(input.Author) == "" {
		return errors.New("author is required")
	}
	return nil
}

func (uc *blogUseCase) validateUpdateInput(input UpdatePostInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(input.Summary) == "" {
		return errors.New("summary is required")
	}
	if strings.TrimSpace(input.Content) == "" {
		return errors.New("content is required")
	}
	return nil
}

func (uc *blogUseCase) generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (simplified)
	allowedChars := "abcdefghijklmnopqrstuvwxyz0123456789-"
	var result strings.Builder
	for _, char := range slug {
		if strings.ContainsRune(allowedChars, char) {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func (uc *blogUseCase) generateID() string {
	return fmt.Sprintf("post_%d", time.Now().UnixNano())
}