# Parallax

An opinionated, modern Go web starter kit featuring Fiber v3, HTMX, Tailwind CSS v4, and Templ with live reloading.

*Thanks to [RevanchistX](https://github.com/RevanchistX) & [Skopsgo](https://github.com/nikolastojkov/skopsgo) for this wonderful rendition.*

## ✨ Features

- **🚀 Go 1.24** - Latest Go version with modern language features
- **🌐 Fiber Framework v3** - High-performance HTTP web framework built on Fasthttp
- **⚡ HTMX** - Modern web interactivity without complex JavaScript
- **🎨 Tailwind CSS v4** - Latest utility-first CSS framework
- **📄 Templ** - Type-safe Go templating
- **🔄 Live Reload** - Hot reloading with Air for seamless development
- **🐳 Docker Ready** - Complete containerization support
- **📦 pnpm** - Fast, disk space efficient package manager
- **🛠️ Make Integration** - Simplified build and development commands
- **🔒 Built-in Middleware** - Compression, logging, recovery, and CORS support

## 🎯 Why Parallax?

This starter kit addresses the need for a modern, batteries-included Go web development setup that combines the best of server-side rendering with modern frontend tooling. Perfect for developers who want the performance of Fiber with the interactivity of HTMX and the styling power of Tailwind CSS.

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
   Navigate to `http://localhost:8080` to see the welcome page with the interactive counter demo.

### Method 2: Docker Development

1. **Build the container**
   ```bash
   docker build --rm -t parallax .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 parallax
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
| `make install` | Install all dependencies |

## 📁 Project Structure

```
parallax/
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
└── package.json          # Frontend dependencies
```

## 🎨 Getting Started with Development

After running the setup, you'll see a splash screen with an example HTMX counter to verify everything works correctly.

### Creating Your First Page

1. **Create a new template** in `web/templates/`
2. **Add a handler** in `internal/handlers/`
3. **Register the route** in your handler setup
4. **Style with Tailwind** classes

### Cleaning Up Demo Content

When ready to start your project, you can safely remove:
- Demo splash page templates in `web/templates/`
- Counter demo handlers in `internal/handlers/`
- Demo route registrations in your handlers

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
SERVER_PORT=8080

# CORS Configuration
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization,X-Requested-With
CORS_EXPOSE_HEADERS=Content-Length
CORS_ALLOW_CREDENTIALS=false
CORS_MAX_AGE=86400
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
docker build -t parallax:production .
docker run -p 8080:8080 parallax:production
```

## 🏗️ Architecture

This starter kit is built on:

- **Fiber v3**: High-performance web framework built on Fasthttp
- **Templ**: Type-safe HTML templating for Go
- **HTMX**: Modern web interactions without heavy JavaScript
- **Tailwind CSS v4**: Utility-first CSS framework
- **Air**: Live reloading for development

## 🤝 Contributing

This is a personal starter template, but contributions are welcome! Feel free to:
- Submit bug reports
- Propose new features
- Submit pull requests

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🙏 Acknowledgments

- [Skopsgo](https://github.com/nikolastojkov/skopsgo) - Original inspiration
- [Fiber](https://github.com/gofiber/fiber) - High-performance web framework
- [HTMX](https://htmx.org/) - Making web development fun again
- [Templ](https://github.com/a-h/templ) - Type-safe Go templates
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first styling
- [Air](https://github.com/air-verse/air) - Seamless live reloading