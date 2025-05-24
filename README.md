# Parallax

An opinionated, modern Go web starter kit featuring Gin, HTMX, Tailwind CSS v4, and Templ with live reloading.

*Thanks to [RevanchistX](https://github.com/RevanchistX) & [Skopsgo](https://github.com/nikolastojkov/skopsgo) &  for this wonderful rendition.*

## ✨ Features

- **🚀 Go 1.24** - Latest Go version with modern language features
- **🌐 Gin Framework** - High-performance HTTP web framework
- **⚡ HTMX** - Modern web interactivity without complex JavaScript
- **🎨 Tailwind CSS v4** - Latest utility-first CSS framework
- **📄 Templ** - Type-safe Go templating
- **🔄 Live Reload** - Hot reloading with Air for seamless development
- **🐳 Docker Ready** - Complete containerization support
- **📦 pnpm** - Fast, disk space efficient package manager
- **🛠️ Make Integration** - Simplified build and development commands

## 🎯 Why Parallax?

This starter kit addresses the need for a modern, batteries-included Go web development setup that combines the best of server-side rendering with modern frontend tooling. Perfect for developers who want the simplicity of Go with the interactivity of HTMX and the styling power of Tailwind CSS.

## 📋 Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [pnpm](https://pnpm.io/installation) (recommended) or npm
- [Docker](https://www.docker.com/) (optional, for containerized development)

## 🚀 Quick Start

### Method 1: Local Development (Recommended)

1. **Clone and setup**
   ```bash
   git clone https://github.com/cloudparallax/parallax.git
   cd parallax
   ```

2. **Install dependencies and setup environment**
   ```bash
   make setup
   cp .env.example .env
   ```

3. **Start development server**
   ```bash
   make dev
   ```

4. **Open your browser**
   Navigate to `http://localhost:8081` to see the welcome page with the interactive counter demo.

### Method 2: Docker Development

1. **Build the container**
   ```bash
   docker build --rm -t skopsgo .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 skopsgo
   ```

## 🛠️ Available Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make setup` | Complete first-time development setup |
| `make dev` | Start development server with live reload |
| `make build` | Build production binary |
| `make run` | Run the application |
| `make test` | Run test suite |
| `make clean` | Clean build artifacts |
| `make css-build` | Build CSS with Tailwind |
| `make templ-generate` | Generate Go code from Templ templates |

## 📁 Project Structure

```
skopsgo/
├── cmd/parallax/           # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP request handlers
│   └── middleware/        # Custom middleware
├── web/
│   ├── app/              # Main application logic
│   ├── static/           # Static assets (CSS, JS, images)
│   └── templates/        # Templ template files
├── .air.toml             # Air live reload configuration
├── Dockerfile            # Container configuration
├── Makefile              # Build automation
└── tailwind.config.js    # Tailwind CSS configuration
```

## 🎨 Getting Started with Development

After running the setup, you'll see a splash screen with an example HTMX counter to verify everything works correctly.

<img src="https://i.ibb.co/RhHDTRd/splash-Final.jpg" height="480px" />

### Creating Your First Page

1. **Create a new template** in `web/templates/`
2. **Add a handler** in `internal/handlers/`
3. **Register the route** in your handler setup
4. **Style with Tailwind** classes

### Cleaning Up Demo Content

When ready to start your project, you can safely remove:
- `web/templates/splash.templ` - Demo splash page template
- `internal/handlers/counter.go` - Counter demo handlers
- `LoadCounterHandler` function call in `internal/handlers/handlers.go`

## 🔧 Customization

### Renaming the Project

Update the module name throughout your project:
```bash
go mod edit -module your-new-module-name
```

Then update references in:
- `Makefile` (APP_NAME variable)
- `cmd/` directory structure
- Import paths in Go files

### Environment Configuration

Copy `.env.example` to `.env` and configure:
```env
PORT=8080
GIN_MODE=debug
# Add your environment-specific variables
```

## 🧪 Testing

Run the test suite:
```bash
make test
```

## 🚀 Production Deployment

### Building for Production

```bash
make build
```

This creates an optimized binary in `./bin/parallax` with:
- Minified CSS
- Generated templates
- Compressed binary (using `-ldflags="-s -w"`)

### Docker Production Build

The included Dockerfile uses multi-stage builds for optimal production images:

```bash
docker build -t skopsgo:production .
docker run -p 8080:8080 skopsgo:production
```

## 🤝 Contributing

This is a personal starter template, but contributions are welcome! Feel free to:
- Submit bug reports
- Propose new features
- Submit pull requests

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🙏 Acknowledgments

- [Skopsgo](https://github.com/nikolastojkov/skopsgo) Original Package
- [HTMX](https://htmx.org/) for making web development fun again
- [Templ](https://github.com/a-h/templ) for type-safe Go templates
- [Gin](https://github.com/gin-gonic/gin) for the excellent web framework
- [Tailwind CSS](https://tailwindcss.com/) for utility-first styling
- [Air](https://github.com/air-verse/air) for seamless live reloading
