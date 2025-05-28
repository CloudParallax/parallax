package repositories

import (
	"context"
	"github.com/cloudparallax/parallax/internal/domain/entities"
)

// BlogPostRepository defines the interface for blog post data operations
type BlogPostRepository interface {
	// Create creates a new blog post
	Create(ctx context.Context, post *entities.BlogPost) error
	
	// GetByID retrieves a blog post by its ID
	GetByID(ctx context.Context, id string) (*entities.BlogPost, error)
	
	// GetBySlug retrieves a blog post by its slug
	GetBySlug(ctx context.Context, slug string) (*entities.BlogPost, error)
	
	// GetAll retrieves all blog posts with optional filters
	GetAll(ctx context.Context, filters BlogPostFilters) ([]*entities.BlogPost, error)
	
	// Update updates an existing blog post
	Update(ctx context.Context, post *entities.BlogPost) error
	
	// Delete deletes a blog post by ID
	Delete(ctx context.Context, id string) error
	
	// GetByTags retrieves blog posts by tags
	GetByTags(ctx context.Context, tags []string) ([]*entities.BlogPost, error)
	
	// GetPublished retrieves only published blog posts
	GetPublished(ctx context.Context, limit, offset int) ([]*entities.BlogPost, error)
	
	// Count returns the total number of blog posts
	Count(ctx context.Context, filters BlogPostFilters) (int64, error)
}

// BlogPostFilters represents filters for querying blog posts
type BlogPostFilters struct {
	AuthorID    string
	Tags        []string
	IsPublished *bool
	FromDate    *string
	ToDate      *string
	SearchQuery string
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string
}