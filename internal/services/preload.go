package services

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type AssetInfo struct {
	Path         string
	ContentType  string
	Hash         string
	Size         int64
	LastModified time.Time
	Preload      bool
	Critical     bool
}

type PreloadService struct {
	assets     map[string]*AssetInfo
	mu         sync.RWMutex
	staticDir  string
	production bool
}

func NewPreloadService(staticDir string, production bool) *PreloadService {
	service := &PreloadService{
		assets:     make(map[string]*AssetInfo),
		staticDir:  staticDir,
		production: production,
	}

	// Initialize asset discovery
	go service.discoverAssets()

	return service
}

func (p *PreloadService) discoverAssets() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Critical assets that should always be preloaded
	criticalAssets := map[string]string{
		"/static/dist/output.css": "text/css",
		"/static/favicon.ico":     "image/x-icon",
	}

	// Scan static directory
	filepath.Walk(p.staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Convert file path to URL path
		relPath := strings.TrimPrefix(path, p.staticDir)
		relPath = strings.TrimPrefix(relPath, "/")
		relPath = strings.TrimPrefix(relPath, "\\")

		// If relPath already starts with "static", don't add another /static prefix
		var urlPath string
		if strings.HasPrefix(relPath, "static/") || strings.HasPrefix(relPath, "static\\") {
			urlPath = "/" + filepath.ToSlash(relPath)
		} else {
			urlPath = "/static/" + filepath.ToSlash(relPath)
		}

		// Determine content type
		contentType := getContentType(filepath.Ext(path))
		if contentType == "" {
			return nil
		}

		// Check if it's a critical asset
		isCritical := false
		if expectedType, exists := criticalAssets[urlPath]; exists && expectedType == contentType {
			isCritical = true
		}

		// Calculate hash for cache busting
		hash := ""
		if p.production {
			hash = p.calculateHash(path)
		}

		p.assets[urlPath] = &AssetInfo{
			Path:         urlPath,
			ContentType:  contentType,
			Hash:         hash,
			Size:         info.Size(),
			LastModified: info.ModTime(),
			Preload:      shouldPreload(contentType, isCritical),
			Critical:     isCritical,
		}

		return nil
	})
}

func (p *PreloadService) watchAssets() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		p.checkForChanges()
	}
}

func (p *PreloadService) checkForChanges() {
	cssPath := filepath.Join(p.staticDir, "dist", "output.css")
	if info, err := os.Stat(cssPath); err == nil {
		p.mu.Lock()
		if asset, exists := p.assets["/static/dist/output.css"]; exists {
			if !info.ModTime().Equal(asset.LastModified) {
				asset.LastModified = info.ModTime()
				asset.Size = info.Size()
				if p.production {
					asset.Hash = p.calculateHash(cssPath)
				}
			}
		}
		p.mu.Unlock()
	}
}

func (p *PreloadService) GetPreloadHeaders() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var headers []string

	for _, asset := range p.assets {
		if asset.Preload {
			header := fmt.Sprintf("<%s>; rel=preload; as=%s",
				asset.Path, getAsType(asset.ContentType))

			if asset.Critical {
				header += "; crossorigin"
			}

			if p.production && asset.Hash != "" {
				header += fmt.Sprintf("; integrity=\"md5-%s\"", asset.Hash)
			}

			headers = append(headers, header)
		}
	}

	return headers
}

func (p *PreloadService) GetPushTargets() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var targets []string

	for _, asset := range p.assets {
		if asset.Critical && asset.Size < 50*1024 { // Only push small critical assets
			targets = append(targets, asset.Path)
		}
	}

	return targets
}

func (p *PreloadService) GetAsset(path string) (*AssetInfo, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	asset, exists := p.assets[path]
	return asset, exists
}

func (p *PreloadService) calculateHash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func getContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff", ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	default:
		return ""
	}
}

func getAsType(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "text/css"):
		return "style"
	case strings.HasPrefix(contentType, "application/javascript"):
		return "script"
	case strings.HasPrefix(contentType, "image/"):
		return "image"
	case strings.HasPrefix(contentType, "font/"):
		return "font"
	default:
		return "fetch"
	}
}

func shouldPreload(contentType string, isCritical bool) bool {
	if isCritical {
		return true
	}

	// Preload CSS and critical fonts
	switch {
	case strings.HasPrefix(contentType, "text/css"):
		return true
	case strings.HasPrefix(contentType, "font/"):
		return true
	default:
		return false
	}
}
