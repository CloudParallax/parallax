package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cloudparallax/parallax/web/templates/pages"
)

type BlogService struct {
	contentPath string
}

func NewBlogService(contentPath string) *BlogService {
	return &BlogService{
		contentPath: contentPath,
	}
}

func (bs *BlogService) GetAllPosts() ([]pages.BlogPost, error) {
	posts := []pages.BlogPost{}

	// Create some sample posts if content directory is empty
	if _, err := os.Stat(bs.contentPath); os.IsNotExist(err) {
		return bs.getSamplePosts(), nil
	}

	err := filepath.Walk(bs.contentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".md" {
			post, err := bs.parseMarkdownFile(path)
			if err != nil {
				return err
			}
			posts = append(posts, post)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort posts by publication date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PublishedAt.After(posts[j].PublishedAt)
	})

	return posts, nil
}

func (bs *BlogService) GetPostBySlug(slug string) (*pages.BlogPost, error) {
	posts, err := bs.GetAllPosts()
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		if post.Slug == slug {
			return &post, nil
		}
	}

	return nil, fmt.Errorf("post with slug '%s' not found", slug)
}

func (bs *BlogService) parseMarkdownFile(path string) (pages.BlogPost, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return pages.BlogPost{}, err
	}

	// Simple frontmatter parser
	lines := strings.Split(string(content), "\n")
	var frontmatterEnd int
	var frontmatter map[string]string = make(map[string]string)

	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				frontmatterEnd = i + 1
				break
			}

			parts := strings.SplitN(lines[i], ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				frontmatter[key] = strings.Trim(value, `"`)
			}
		}
	}

	// Extract markdown content
	var markdownContent string
	if frontmatterEnd > 0 {
		markdownContent = strings.Join(lines[frontmatterEnd:], "\n")
	} else {
		markdownContent = string(content)
	}

	// Convert markdown to HTML (simple implementation)
	htmlContent := bs.markdownToHTML(markdownContent)

	// Parse publication date
	publishedAt := time.Now()
	if dateStr, exists := frontmatter["date"]; exists {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			publishedAt = parsed
		}
	}

	// Parse tags
	var tags []string
	if tagsStr, exists := frontmatter["tags"]; exists {
		tags = strings.Split(tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	// Generate slug from filename
	slug := strings.TrimSuffix(filepath.Base(path), ".md")

	return pages.BlogPost{
		ID:          slug,
		Title:       frontmatter["title"],
		Slug:        slug,
		Summary:     frontmatter["summary"],
		Content:     htmlContent,
		Author:      frontmatter["author"],
		PublishedAt: publishedAt,
		Tags:        tags,
		ReadTime:    bs.calculateReadTime(markdownContent),
	}, nil
}

func (bs *BlogService) markdownToHTML(markdown string) string {
	// Simple markdown to HTML conversion
	lines := strings.Split(markdown, "\n")
	var html strings.Builder

	inCodeBlock := false
	inParagraph := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Handle code blocks
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				html.WriteString("</code></pre>\n")
				inCodeBlock = false
			} else {
				if inParagraph {
					html.WriteString("</p>\n")
					inParagraph = false
				}
				language := strings.TrimPrefix(line, "```")
				if language != "" {
					html.WriteString(fmt.Sprintf(`<pre class="bg-gray-100 rounded-lg p-4 overflow-x-auto"><code class="language-%s">`, language))
				} else {
					html.WriteString(`<pre class="bg-gray-100 rounded-lg p-4 overflow-x-auto"><code>`)
				}
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			html.WriteString(line + "\n")
			continue
		}

		// Handle headers
		if strings.HasPrefix(line, "# ") {
			if inParagraph {
				html.WriteString("</p>\n")
				inParagraph = false
			}
			html.WriteString(fmt.Sprintf(`<h1 class="text-3xl font-bold text-gray-900 mb-4">%s</h1>`, strings.TrimPrefix(line, "# ")))
			continue
		}
		if strings.HasPrefix(line, "## ") {
			if inParagraph {
				html.WriteString("</p>\n")
				inParagraph = false
			}
			html.WriteString(fmt.Sprintf(`<h2 class="text-2xl font-bold text-gray-900 mb-3 mt-6">%s</h2>`, strings.TrimPrefix(line, "## ")))
			continue
		}
		if strings.HasPrefix(line, "### ") {
			if inParagraph {
				html.WriteString("</p>\n")
				inParagraph = false
			}
			html.WriteString(fmt.Sprintf(`<h3 class="text-xl font-bold text-gray-900 mb-2 mt-4">%s</h3>`, strings.TrimPrefix(line, "### ")))
			continue
		}

		// Handle empty lines
		if line == "" {
			if inParagraph {
				html.WriteString("</p>\n")
				inParagraph = false
			}
			continue
		}

		// Handle lists
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			if inParagraph {
				html.WriteString("</p>\n")
				inParagraph = false
			}
			content := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
			html.WriteString(fmt.Sprintf(`<ul class="list-disc list-inside mb-4"><li class="text-gray-700">%s</li></ul>`, content))
			continue
		}

		// Handle regular paragraphs
		if !inParagraph {
			html.WriteString(`<p class="text-gray-700 mb-4 leading-relaxed">`)
			inParagraph = true
		} else {
			html.WriteString(" ")
		}

		// Process inline formatting
		content := bs.processInlineFormatting(line)
		html.WriteString(content)
	}

	if inParagraph {
		html.WriteString("</p>\n")
	}

	return html.String()
}

func (bs *BlogService) processInlineFormatting(text string) string {
	// Handle bold text
	text = strings.ReplaceAll(text, "**", `<strong>`)

	// Handle italic text
	text = strings.ReplaceAll(text, "*", `<em>`)

	// Handle inline code
	parts := strings.Split(text, "`")
	for i := 1; i < len(parts); i += 2 {
		if i < len(parts) {
			parts[i] = fmt.Sprintf(`<code class="bg-gray-100 px-1 py-0.5 rounded text-sm">%s</code>`, parts[i])
		}
	}
	text = strings.Join(parts, "")

	return text
}

func (bs *BlogService) calculateReadTime(content string) int {
	words := strings.Fields(content)
	// Average reading speed: 200 words per minute
	readTime := len(words) / 200
	if readTime < 1 {
		readTime = 1
	}
	return readTime
}

func (bs *BlogService) getSamplePosts() []pages.BlogPost {
	return []pages.BlogPost{
		{
			ID:          "getting-started-with-htmx",
			Title:       "Getting Started with HTMX and Go",
			Slug:        "getting-started-with-htmx",
			Summary:     "Learn how to build dynamic web applications using HTMX with Go backend. This tutorial covers the basics of creating interactive UIs without JavaScript frameworks.",
			Content:     bs.markdownToHTML(samplePost1),
			Author:      "CloudParallax Team",
			PublishedAt: time.Now().AddDate(0, 0, -7),
			Tags:        []string{"HTMX", "Go", "Web Development"},
			ReadTime:    5,
		},
		{
			ID:          "modern-web-architecture",
			Title:       "Modern Web Architecture with Templ and Fiber",
			Slug:        "modern-web-architecture",
			Summary:     "Explore how to build scalable web applications using Go's Templ templating engine with Fiber framework. A comprehensive guide to modern backend development.",
			Content:     bs.markdownToHTML(samplePost2),
			Author:      "CloudParallax Team",
			PublishedAt: time.Now().AddDate(0, 0, -14),
			Tags:        []string{"Go", "Templ", "Fiber", "Architecture"},
			ReadTime:    8,
		},
		{
			ID:          "cloud-infrastructure-best-practices",
			Title:       "Cloud Infrastructure Best Practices",
			Slug:        "cloud-infrastructure-best-practices",
			Summary:     "Essential guidelines for building robust and scalable cloud infrastructure. Learn about deployment strategies, monitoring, and security considerations.",
			Content:     bs.markdownToHTML(samplePost3),
			Author:      "CloudParallax Team",
			PublishedAt: time.Now().AddDate(0, 0, -21),
			Tags:        []string{"Cloud", "Infrastructure", "DevOps", "Best Practices"},
			ReadTime:    12,
		},
	}
}

const samplePost1 = `# Getting Started with HTMX and Go

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

const samplePost2 = `# Modern Web Architecture with Templ and Fiber

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

const samplePost3 = `# Cloud Infrastructure Best Practices

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
