package repositories

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cloudparallax/parallax/internal/domain/entities"
	"github.com/cloudparallax/parallax/internal/domain/repositories"
)

// memoryBlogRepository implements BlogPostRepository interface using in-memory storage
type memoryBlogRepository struct {
	posts map[string]*entities.BlogPost
	mutex sync.RWMutex
}

// NewMemoryBlogRepository creates a new in-memory blog repository
func NewMemoryBlogRepository() repositories.BlogPostRepository {
	repo := &memoryBlogRepository{
		posts: make(map[string]*entities.BlogPost),
		mutex: sync.RWMutex{},
	}
	
	// Initialize with sample data
	repo.initializeSampleData()
	
	return repo
}

// Create creates a new blog post
func (r *memoryBlogRepository) Create(ctx context.Context, post *entities.BlogPost) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.posts[post.ID]; exists {
		return errors.New("blog post already exists")
	}
	
	// Create a copy to avoid external modifications
	postCopy := *post
	r.posts[post.ID] = &postCopy
	
	return nil
}

// GetByID retrieves a blog post by its ID
func (r *memoryBlogRepository) GetByID(ctx context.Context, id string) (*entities.BlogPost, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("blog post not found")
	}
	
	// Return a copy to avoid external modifications
	postCopy := *post
	return &postCopy, nil
}

// GetBySlug retrieves a blog post by its slug
func (r *memoryBlogRepository) GetBySlug(ctx context.Context, slug string) (*entities.BlogPost, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, post := range r.posts {
		if post.Slug == slug {
			// Return a copy to avoid external modifications
			postCopy := *post
			return &postCopy, nil
		}
	}
	
	return nil, errors.New("blog post not found")
}

// GetAll retrieves all blog posts with optional filters
func (r *memoryBlogRepository) GetAll(ctx context.Context, filters repositories.BlogPostFilters) ([]*entities.BlogPost, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var posts []*entities.BlogPost
	
	for _, post := range r.posts {
		if r.matchesFilters(post, filters) {
			postCopy := *post
			posts = append(posts, &postCopy)
		}
	}
	
	// Sort posts
	r.sortPosts(posts, filters.SortBy, filters.SortOrder)
	
	// Apply pagination
	start := filters.Offset
	end := start + filters.Limit
	
	if start > len(posts) {
		return []*entities.BlogPost{}, nil
	}
	
	if end > len(posts) {
		end = len(posts)
	}
	
	if filters.Limit > 0 {
		posts = posts[start:end]
	}
	
	return posts, nil
}

// Update updates an existing blog post
func (r *memoryBlogRepository) Update(ctx context.Context, post *entities.BlogPost) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.posts[post.ID]; !exists {
		return errors.New("blog post not found")
	}
	
	// Create a copy to avoid external modifications
	postCopy := *post
	r.posts[post.ID] = &postCopy
	
	return nil
}

// Delete deletes a blog post by ID
func (r *memoryBlogRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.posts[id]; !exists {
		return errors.New("blog post not found")
	}
	
	delete(r.posts, id)
	return nil
}

// GetByTags retrieves blog posts by tags
func (r *memoryBlogRepository) GetByTags(ctx context.Context, tags []string) ([]*entities.BlogPost, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var posts []*entities.BlogPost
	
	for _, post := range r.posts {
		if r.hasAnyTag(post, tags) {
			postCopy := *post
			posts = append(posts, &postCopy)
		}
	}
	
	// Sort by creation date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})
	
	return posts, nil
}

// GetPublished retrieves only published blog posts
func (r *memoryBlogRepository) GetPublished(ctx context.Context, limit, offset int) ([]*entities.BlogPost, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var posts []*entities.BlogPost
	
	for _, post := range r.posts {
		if post.IsPublished {
			postCopy := *post
			posts = append(posts, &postCopy)
		}
	}
	
	// Sort by publication date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PublishedAt.After(posts[j].PublishedAt)
	})
	
	// Apply pagination
	start := offset
	end := start + limit
	
	if start > len(posts) {
		return []*entities.BlogPost{}, nil
	}
	
	if end > len(posts) {
		end = len(posts)
	}
	
	return posts[start:end], nil
}

// Count returns the total number of blog posts
func (r *memoryBlogRepository) Count(ctx context.Context, filters repositories.BlogPostFilters) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	count := int64(0)
	
	for _, post := range r.posts {
		if r.matchesFilters(post, filters) {
			count++
		}
	}
	
	return count, nil
}

// Helper methods

func (r *memoryBlogRepository) matchesFilters(post *entities.BlogPost, filters repositories.BlogPostFilters) bool {
	// Filter by author
	if filters.AuthorID != "" && post.Author != filters.AuthorID {
		return false
	}
	
	// Filter by published status
	if filters.IsPublished != nil && post.IsPublished != *filters.IsPublished {
		return false
	}
	
	// Filter by tags
	if len(filters.Tags) > 0 && !r.hasAnyTag(post, filters.Tags) {
		return false
	}
	
	// Filter by date range
	if filters.FromDate != nil {
		if fromDate, err := time.Parse("2006-01-02", *filters.FromDate); err == nil {
			if post.CreatedAt.Before(fromDate) {
				return false
			}
		}
	}
	
	if filters.ToDate != nil {
		if toDate, err := time.Parse("2006-01-02", *filters.ToDate); err == nil {
			if post.CreatedAt.After(toDate.Add(24 * time.Hour)) {
				return false
			}
		}
	}
	
	// Filter by search query
	if filters.SearchQuery != "" {
		query := strings.ToLower(filters.SearchQuery)
		if !strings.Contains(strings.ToLower(post.Title), query) &&
			!strings.Contains(strings.ToLower(post.Summary), query) &&
			!strings.Contains(strings.ToLower(post.Content), query) {
			return false
		}
	}
	
	return true
}

func (r *memoryBlogRepository) hasAnyTag(post *entities.BlogPost, tags []string) bool {
	for _, filterTag := range tags {
		for _, postTag := range post.Tags {
			if strings.EqualFold(postTag, filterTag) {
				return true
			}
		}
	}
	return false
}

func (r *memoryBlogRepository) sortPosts(posts []*entities.BlogPost, sortBy, sortOrder string) {
	if sortBy == "" {
		sortBy = "created_at"
	}
	
	if sortOrder == "" {
		sortOrder = "desc"
	}
	
	sort.Slice(posts, func(i, j int) bool {
		var less bool
		
		switch sortBy {
		case "title":
			less = posts[i].Title < posts[j].Title
		case "updated_at":
			less = posts[i].UpdatedAt.Before(posts[j].UpdatedAt)
		case "published_at":
			less = posts[i].PublishedAt.Before(posts[j].PublishedAt)
		default: // created_at
			less = posts[i].CreatedAt.Before(posts[j].CreatedAt)
		}
		
		if sortOrder == "desc" {
			return !less
		}
		
		return less
	})
}

func (r *memoryBlogRepository) initializeSampleData() {
	samplePosts := []*entities.BlogPost{
		{
			ID:          "post_1",
			Title:       "Getting Started with HTMX and Go",
			Slug:        "getting-started-with-htmx",
			Summary:     "Learn how to build dynamic web applications using HTMX with Go backend. This tutorial covers the basics of creating interactive UIs without JavaScript frameworks.",
			Content:     r.getSampleContent1(),
			Author:      "CloudParallax Team",
			PublishedAt: time.Now().AddDate(0, 0, -7),
			CreatedAt:   time.Now().AddDate(0, 0, -7),
			UpdatedAt:   time.Now().AddDate(0, 0, -7),
			Tags:        []string{"HTMX", "Go", "Web Development"},
			ReadTime:    5,
			IsPublished: true,
		},
		{
			ID:          "post_2",
			Title:       "Modern Web Architecture with Templ and Fiber",
			Slug:        "modern-web-architecture",
			Summary:     "Explore how to build scalable web applications using Go's Templ templating engine with Fiber framework. A comprehensive guide to modern backend development.",
			Content:     r.getSampleContent2(),
			Author:      "CloudParallax Team",
			PublishedAt: time.Now().AddDate(0, 0, -14),
			CreatedAt:   time.Now().AddDate(0, 0, -14),
			UpdatedAt:   time.Now().AddDate(0, 0, -14),
			Tags:        []string{"Go", "Templ", "Fiber", "Architecture"},
			ReadTime:    8,
			IsPublished: true,
		},
		{
			ID:          "post_3",
			Title:       "Cloud Infrastructure Best Practices",
			Slug:        "cloud-infrastructure-best-practices",
			Summary:     "Essential guidelines for building robust and scalable cloud infrastructure. Learn about deployment strategies, monitoring, and security considerations.",
			Content:     r.getSampleContent3(),
			Author:      "CloudParallax Team",
			PublishedAt: time.Now().AddDate(0, 0, -21),
			CreatedAt:   time.Now().AddDate(0, 0, -21),
			UpdatedAt:   time.Now().AddDate(0, 0, -21),
			Tags:        []string{"Cloud", "Infrastructure", "DevOps", "Best Practices"},
			ReadTime:    12,
			IsPublished: true,
		},
	}
	
	for _, post := range samplePosts {
		r.posts[post.ID] = post
	}
}

func (r *memoryBlogRepository) getSampleContent1() string {
	return `# Getting Started with HTMX and Go

HTMX is a powerful library that allows you to access modern browser features directly from HTML, without writing JavaScript. When combined with Go's robust backend capabilities, you can create dynamic, interactive web applications with minimal complexity.

## Why HTMX?

HTMX extends HTML with attributes that enable:
- AJAX requests directly from HTML elements
- WebSocket connections
- Server-sent events
- CSS transitions
- And much more!

## Basic Example

Here's a simple example of an HTMX-powered button:

` + "```html" + `
<button hx-get="/api/hello" hx-target="#result">
    Click Me
</button>
<div id="result"></div>
` + "```" + `

When clicked, this button will make a GET request to /api/hello and replace the content of the #result div with the response.

## Go Backend

On the Go side, you can handle this request easily:

` + "```go" + `
app.Get("/api/hello", func(c *fiber.Ctx) error {
    return c.SendString("<p>Hello from the server!</p>")
})
` + "```" + `

## Benefits

- **Simplicity**: No complex JavaScript frameworks
- **Progressive Enhancement**: Works without JavaScript enabled
- **Server-Centric**: Keep your logic on the server
- **Performance**: Minimal client-side overhead

Start building your next web application with HTMX and Go today!`
}

func (r *memoryBlogRepository) getSampleContent2() string {
	return `# Modern Web Architecture with Templ and Fiber

Building modern web applications requires the right combination of tools and architectural patterns. In this post, we'll explore how Templ and Fiber work together to create fast, maintainable web applications.

## Introduction to Templ

Templ is a language for writing HTML user interfaces in Go. It provides:
- Type safety
- Great IDE support
- Component composition
- No runtime overhead

## Why Fiber?

Fiber is an Express-inspired web framework built on top of Fasthttp:
- **Fast**: Built on Fasthttp, the fastest HTTP engine for Go
- **Flexible**: Modular design with extensive middleware support
- **Developer Friendly**: Express-like API that's easy to learn

## Project Structure

A well-organized project structure is crucial:

` + "```" + `
project/
├── internal/
│   ├── handlers/
│   ├── middleware/
│   └── services/
├── web/
│   ├── templates/
│   │   ├── components/
│   │   ├── layouts/
│   │   └── pages/
│   └── static/
└── cmd/
    └── server/
` + "```" + `

## Template Organization

Organize your templates by type:
- **Layouts**: Base page structures
- **Components**: Reusable UI elements
- **Pages**: Full page templates

This approach promotes reusability and maintainability.

## Best Practices

1. **Separation of Concerns**: Keep business logic in services
2. **Middleware**: Use middleware for cross-cutting concerns
3. **Error Handling**: Implement proper error handling strategies
4. **Testing**: Write tests for your handlers and services

Building with Templ and Fiber gives you the performance of Go with the productivity of modern web development patterns.`
}

func (r *memoryBlogRepository) getSampleContent3() string {
	return `# Cloud Infrastructure Best Practices

As organizations scale their applications, having robust cloud infrastructure becomes critical. This guide covers essential practices for building reliable, scalable, and secure cloud systems.

## Infrastructure as Code

Always define your infrastructure using code:
- **Version Control**: Track infrastructure changes
- **Reproducibility**: Deploy consistent environments
- **Documentation**: Code serves as documentation

## Monitoring and Observability

Implement comprehensive monitoring:

### Metrics to Track
- Application performance metrics
- Infrastructure resource utilization
- Business metrics
- Error rates and response times

### Tools and Techniques
- **Logging**: Structured logging with correlation IDs
- **Metrics**: Time-series data for trending
- **Tracing**: Distributed request tracing
- **Alerting**: Proactive issue detection

## Security Considerations

Security should be built into every layer:

### Network Security
- Use VPCs and private subnets
- Implement proper firewall rules
- Enable encryption in transit

### Application Security
- Regular security updates
- Secrets management
- Authentication and authorization
- Input validation

## Deployment Strategies

Choose the right deployment strategy:

### Blue-Green Deployment
- Zero-downtime deployments
- Easy rollback capability
- Resource overhead considerations

### Rolling Updates
- Gradual deployment
- Lower resource requirements
- Longer deployment times

### Canary Releases
- Test with subset of users
- Risk mitigation
- Gradual traffic shifting

## Disaster Recovery

Plan for the worst-case scenarios:
- **Backups**: Regular, tested backups
- **Recovery Procedures**: Documented recovery steps
- **RTO/RPO**: Define recovery objectives
- **Testing**: Regular disaster recovery drills

## Cost Optimization

Keep costs under control:
- **Right-sizing**: Match resources to actual needs
- **Auto-scaling**: Scale based on demand
- **Reserved Instances**: Long-term cost savings
- **Monitoring**: Track and optimize spending

Following these practices will help you build infrastructure that can scale with your business while maintaining reliability and security.`
}