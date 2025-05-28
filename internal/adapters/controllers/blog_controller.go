package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/cloudparallax/parallax/internal/adapters/http/dto"
	"github.com/cloudparallax/parallax/internal/usecases"
	"github.com/cloudparallax/parallax/pkg/response"
	"github.com/gofiber/fiber/v3"
)

// BlogController handles HTTP requests for blog operations
type BlogController struct {
	blogUseCase usecases.BlogUseCase
}

// NewBlogController creates a new blog controller
func NewBlogController(blogUseCase usecases.BlogUseCase) *BlogController {
	return &BlogController{
		blogUseCase: blogUseCase,
	}
}

// CreatePost handles creating a new blog post
func (c *BlogController) CreatePost(ctx fiber.Ctx) error {
	var req dto.CreateBlogPostRequest
	
	if err := response.ParseJSON(ctx, &req); err != nil {
		return response.Error(ctx, err)
	}

	input := usecases.CreatePostInput{
		Title:   req.Title,
		Summary: req.Summary,
		Content: req.Content,
		Author:  req.Author,
		Tags:    req.Tags,
	}

	post, err := c.blogUseCase.CreatePost(ctx.Context(), input)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Created(ctx, dto.ToBlogPostResponse(post))
}

// GetPost handles getting a blog post by ID
func (c *BlogController) GetPost(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Post ID is required")
	}

	post, err := c.blogUseCase.GetPostByID(ctx.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Blog post not found")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToBlogPostResponse(post))
}

// GetPostBySlug handles getting a blog post by slug
func (c *BlogController) GetPostBySlug(ctx fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return response.BadRequest(ctx, "Post slug is required")
	}

	post, err := c.blogUseCase.GetPostBySlug(ctx.Context(), slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Blog post not found")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToBlogPostResponse(post))
}

// GetPosts handles getting all blog posts with filters
func (c *BlogController) GetPosts(ctx fiber.Ctx) error {
	var req dto.GetBlogPostsRequest
	
	if err := ctx.Bind().Query(&req); err != nil {
		return response.BadRequest(ctx, "Invalid query parameters")
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Parse tags
	var tags []string
	if req.Tags != "" {
		tags = strings.Split(req.Tags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	// Parse dates
	var fromDate, toDate *time.Time
	if req.FromDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.FromDate); err == nil {
			fromDate = &parsed
		}
	}
	if req.ToDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.ToDate); err == nil {
			toDate = &parsed
		}
	}

	input := usecases.GetPostsInput{
		AuthorID:    req.AuthorID,
		Tags:        tags,
		IsPublished: req.IsPublished,
		FromDate:    fromDate,
		ToDate:      toDate,
		SearchQuery: req.SearchQuery,
		Limit:       req.Limit,
		Offset:      (req.Page - 1) * req.Limit,
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
	}

	posts, err := c.blogUseCase.GetAllPosts(ctx.Context(), input)
	if err != nil {
		return response.Error(ctx, err)
	}

	// For simplicity, we'll assume total count is the length of posts
	// In a real implementation, you'd get this from the repository
	total := len(posts)
	meta := response.NewMeta(req.Page, req.Limit, total)

	return response.SuccessWithMeta(ctx, dto.ToBlogPostListResponse(posts, &dto.MetaResponse{
		Page:       meta.Page,
		Limit:      meta.Limit,
		Total:      meta.Total,
		TotalPages: meta.TotalPages,
	}), meta)
}

// UpdatePost handles updating a blog post
func (c *BlogController) UpdatePost(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Post ID is required")
	}

	var req dto.UpdateBlogPostRequest
	
	if err := response.ParseJSON(ctx, &req); err != nil {
		return response.Error(ctx, err)
	}

	input := usecases.UpdatePostInput{
		Title:   req.Title,
		Summary: req.Summary,
		Content: req.Content,
		Tags:    req.Tags,
	}

	post, err := c.blogUseCase.UpdatePost(ctx.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Blog post not found")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToBlogPostResponse(post))
}

// DeletePost handles deleting a blog post
func (c *BlogController) DeletePost(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Post ID is required")
	}

	err := c.blogUseCase.DeletePost(ctx.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Blog post not found")
		}
		return response.Error(ctx, err)
	}

	return response.NoContent(ctx)
}

// PublishPost handles publishing a blog post
func (c *BlogController) PublishPost(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Post ID is required")
	}

	post, err := c.blogUseCase.PublishPost(ctx.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Blog post not found")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToBlogPostResponse(post))
}

// UnpublishPost handles unpublishing a blog post
func (c *BlogController) UnpublishPost(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.BadRequest(ctx, "Post ID is required")
	}

	post, err := c.blogUseCase.UnpublishPost(ctx.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Blog post not found")
		}
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToBlogPostResponse(post))
}

// GetPublishedPosts handles getting only published blog posts
func (c *BlogController) GetPublishedPosts(ctx fiber.Ctx) error {
	page := 1
	limit := 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	posts, err := c.blogUseCase.GetPublishedPosts(ctx.Context(), limit, offset)
	if err != nil {
		return response.Error(ctx, err)
	}

	// For simplicity, we'll assume total count is the length of posts
	// In a real implementation, you'd get this from the repository
	total := len(posts)
	meta := response.NewMeta(page, limit, total)

	return response.SuccessWithMeta(ctx, dto.ToBlogPostListResponse(posts, &dto.MetaResponse{
		Page:       meta.Page,
		Limit:      meta.Limit,
		Total:      meta.Total,
		TotalPages: meta.TotalPages,
	}), meta)
}

// SearchPosts handles searching blog posts
func (c *BlogController) SearchPosts(ctx fiber.Ctx) error {
	var req dto.SearchBlogPostsRequest
	
	if err := ctx.Bind().Query(&req); err != nil {
		return response.BadRequest(ctx, "Invalid query parameters")
	}

	if req.Query == "" {
		return response.BadRequest(ctx, "Search query is required")
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	posts, err := c.blogUseCase.SearchPosts(ctx.Context(), req.Query, req.Limit, offset)
	if err != nil {
		return response.Error(ctx, err)
	}

	// For simplicity, we'll assume total count is the length of posts
	// In a real implementation, you'd get this from the repository
	total := len(posts)
	meta := response.NewMeta(req.Page, req.Limit, total)

	return response.SuccessWithMeta(ctx, dto.ToBlogPostListResponse(posts, &dto.MetaResponse{
		Page:       meta.Page,
		Limit:      meta.Limit,
		Total:      meta.Total,
		TotalPages: meta.TotalPages,
	}), meta)
}

// GetPostsByTags handles getting blog posts by tags
func (c *BlogController) GetPostsByTags(ctx fiber.Ctx) error {
	tagsStr := ctx.Query("tags")
	if tagsStr == "" {
		return response.BadRequest(ctx, "Tags parameter is required")
	}

	tags := strings.Split(tagsStr, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	posts, err := c.blogUseCase.GetPostsByTags(ctx.Context(), tags)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.ToBlogPostListResponse(posts, nil))
}