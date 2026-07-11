# Golang Birmingham Website

A complete, production-quality website for Golang Birmingham, an independent community group for Go developers in Birmingham, United Kingdom.

## Features

- **Server-Rendered Go Application**: Built with Go's standard library (net/http, html/template)
- **Responsive Design**: Mobile-first design with Go-inspired color palette
- **Accessibility**: WCAG compliant with semantic HTML and ARIA support
- **Progressive Enhancement**: Works without JavaScript, enhanced with minimal JS
- **Event Management**: Calendar integration with ICS file generation
- **Content Management**: YAML-based content system
- **Security**: Security headers, CSRF protection, and input validation
- **SEO Optimized**: Meta tags, structured data, sitemap, and robots.txt
- **Docker Ready**: Multi-stage Dockerfile and docker-compose configuration

## Technology Stack

- **Backend**: Go 1.21+ with standard library only
- **Frontend**: Semantic HTML5, CSS3, minimal vanilla JavaScript
- **Templates**: Go html/template with component-based architecture
- **Styling**: Custom CSS with Go-inspired design system
- **Content**: YAML files for events, speakers, team, sponsors
- **Deployment**: Docker, docker-compose, environment variables

## Quick Start

### Prerequisites

- Go 1.21 or later
- Make (optional, for convenient commands)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/placeholder/golang-birmingham-website.git
   cd golang-birmingham-website
   ```

2. **Install dependencies**
   ```bash
   make install
   # or
   go mod download
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the development server**
   ```bash
   make run
   # or
   go run ./cmd/web
   ```

5. **Visit the website**
   Open http://localhost:8080 in your browser

### Using Docker

1. **Build and run with Docker Compose**
   ```bash
   make docker-build
   make docker-run
   # or
   docker-compose up --build
   ```

2. **View logs**
   ```bash
   make docker-logs
   # or
   docker-compose logs -f
   ```

3. **Stop the services**
   ```bash
   make docker-stop
   # or
   docker-compose down
   ```

## Project Structure

```
golang-birmingham-website/
├── cmd/
│   └── web/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── content/                 # Content loading and parsing
│   ├── calendar/                # ICS calendar generation
│   ├── domain/                  # Domain models and types
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # Security and utility middleware
│   ├── rendering/               # Template rendering
│   └── server/                  # Server setup and routing
├── web/
│   ├── templates/
│   │   ├── layouts/             # Base templates
│   │   ├── components/          # Reusable components
│   │   └── pages/               # Page templates
│   └── static/
│       ├── css/                 # Stylesheets
│       ├── js/                  # JavaScript files
│       └── images/              # Static images
├── content/
│   ├── events/                  # Event content (YAML)
│   ├── speakers/                # Speaker profiles (YAML)
│   ├── team/                    # Team members (YAML)
│   ├── sponsors/                # Sponsor information (YAML)
│   ├── gallery/                 # Gallery items (YAML)
│   └── site.json                # Site configuration
├── tests/                       # Test files
├── Dockerfile                   # Multi-stage Docker configuration
├── docker-compose.yml           # Docker Compose configuration
├── Makefile                     # Build and development commands
├── .env.example                 # Environment variables template
└── README.md                    # This file
```

## Configuration

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
# Server Configuration
PORT=8080
BASE_URL=http://localhost:8080
ENVIRONMENT=development

# Contact Information
CONTACT_EMAIL=your-email@example.com
SPEAKER_FORM_URL=https://forms.example.com/speaker

# Social Media Links
MEETUP_URL=https://meetup.com/golang-birmingham
LINKEDIN_URL=https://linkedin.com/company/golang-birmingham
YOUTUBE_URL=https://youtube.com/@golangbirmingham
GITHUB_URL=https://github.com/golang-birmingham

# Security
CSRF_SECRET=your-random-secret-here
SESSION_SECRET=your-session-secret-here
ENABLE_HTTPS=false
```

### Site Configuration

Edit `content/site.json` to customize:

- Site title and description
- Community statistics
- Social media links
- Contact information

## Content Management

### Adding Events

Create a new YAML file in `content/events/`:

```yaml
id: "event-001"
slug: "your-event-slug"
title: "Event Title"
summary: "Brief description"
description: |
  Full event description in Markdown.
start_time: "2027-07-25T18:30:00+01:00"
end_time: "2027-07-25T21:00:00+01:00"
venue_name: "Venue Name"
venue_address: "Full Address"
image: "/static/images/events/event-image.jpg"
topics: ["go", "web", "microservices"]
speakers:
  - slug: "speaker-slug"
registration_url: "https://meetup.com/event-url"
video_url: "https://youtube.com/watch?v=..."
slides_url: "https://slides.com/..."
featured: true
```

### Adding Speakers

Create a new YAML file in `content/speakers/`:

```yaml
id: "speaker-001"
slug: "speaker-slug"
name: "Speaker Name"
role: "Job Title"
organisation: "Company Name"
biography: "Speaker bio in Markdown"
image: "/static/images/speakers/speaker-photo.jpg"
linkedin_url: "https://linkedin.com/in/speaker"
github_url: "https://github.com/speaker"
website_url: "https://speaker-website.com"
```

### Adding Team Members

Create a new YAML file in `content/team/`:

```yaml
name: "Team Member Name"
role: "Organiser Role"
biography: "Team member bio"
image: "/static/images/team/member-photo.jpg"
linkedin_url: "https://linkedin.com/in/member"
github_url: "https://github.com/member"
website_url: "https://member-website.com"
```

### Adding Sponsors

Create a new YAML file in `content/sponsors/`:

```yaml
name: "Sponsor Name"
category: "Venue Partner"
logo: "/static/images/sponsors/sponsor-logo.png"
website_url: "https://sponsor-website.com"
description: "Brief sponsor description"
```

### Adding Gallery Items

Create a new YAML file in `content/gallery/`:

```yaml
id: "gallery-001"
event_slug: "event-slug"
year: 2027
category: "speakers"
image: "/static/images/gallery/full-image.jpg"
thumbnail: "/static/images/gallery/thumbnails/thumb.jpg"
caption: "Image caption"
alt_text: "Descriptive alt text for accessibility"
```

## Available Commands

### Development Commands

```bash
make run          # Start development server
make build        # Build the application
make test         # Run tests
make test-race    # Run tests with race detection
make vet          # Run go vet
make fmt          # Format code
make lint         # Run linter (requires golangci-lint)
make clean        # Clean build artifacts
```

### Docker Commands

```bash
make docker-build # Build Docker image
make docker-run   # Run with Docker Compose
make docker-stop  # Stop Docker Compose
make docker-logs  # Show Docker logs
```

### Setup Commands

```bash
make install      # Install dependencies
make dev-setup    # Set up development environment
make dev          # Quick development start
```

### Quality Checks

```bash
make pre-commit   # Run pre-commit checks
make check-all    # Run all quality checks
make release-check # Run release checks
```

## Deployment

### Production Deployment

1. **Set production environment variables**
   ```bash
   export ENVIRONMENT=production
   export BASE_URL=https://your-domain.com
   export ENABLE_HTTPS=true
   export CSRF_SECRET=your-production-secret
   export SESSION_SECRET=your-production-session-secret
   ```

2. **Build for production**
   ```bash
   make prod-build
   # or
   CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/web-linux ./cmd/web
   ```

3. **Deploy with Docker**
   ```bash
   docker-compose -f docker-compose.yml --profile production up -d
   ```

### Platform-Specific Deployment

The application can be deployed to any platform that supports Go binaries:

- **Google Cloud Run**: `gcloud run deploy`
- **Fly.io**: `fly deploy`
- **Render**: Connect GitHub repository
- **Railway**: Connect GitHub repository
- **VPS**: Copy binary and run directly

### Nginx Configuration (Optional)

For production, you might want to use Nginx as a reverse proxy:

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Run specific test
go test ./internal/handlers

# Run tests with coverage
go test -cover ./...
```

### Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Security

### Security Features

- **Security Headers**: CSP, HSTS, X-Frame-Options, etc.
- **Input Validation**: Server-side form validation
- **CSRF Protection**: Token-based CSRF protection
- **Rate Limiting**: Basic rate limiting for form submissions
- **Secure Cookies**: HttpOnly, Secure, SameSite cookies
- **Content Security Policy**: Restricts resource loading

### Security Best Practices

1. **Environment Variables**: Never commit secrets to version control
2. **HTTPS**: Always use HTTPS in production
3. **Regular Updates**: Keep Go and dependencies updated
4. **Input Validation**: Validate all user inputs
5. **Error Handling**: Don't expose sensitive information in errors

## Accessibility

### Accessibility Features

- **Semantic HTML**: Proper use of HTML5 semantic elements
- **ARIA Labels**: Screen reader support
- **Keyboard Navigation**: Full keyboard accessibility
- **Focus Management**: Visible focus states
- **Color Contrast**: WCAG AA compliant color ratios
- **Alt Text**: Descriptive alt text for images
- **Skip Links**: Skip to main content link

### Testing Accessibility

- Use keyboard-only navigation
- Test with screen readers
- Check color contrast ratios
- Validate with accessibility tools

## Performance

### Performance Features

- **Server-Side Rendering**: Fast initial page loads
- **Minimal JavaScript**: Only essential enhancements
- **Image Optimization**: Responsive images with lazy loading
- **CSS Optimization**: Minified production CSS
- **Caching Headers**: Appropriate cache headers for static assets
- **Compression**: Gzip/Brotli compression (via hosting)

### Performance Monitoring

- Monitor Core Web Vitals
- Check PageSpeed Insights
- Use browser dev tools for analysis
- Monitor server response times

## Contributing

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and quality checks
5. Submit a pull request

### Code Style

- Use `gofmt` for code formatting
- Follow Go conventions
- Write meaningful commit messages
- Add tests for new features
- Update documentation

### Pull Request Process

1. Ensure all tests pass
2. Run `make pre-commit`
3. Update documentation
4. Describe changes clearly
5. Link to relevant issues

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Attribution

### Go Gopher

The Go Gopher was designed by Renée French and is licensed under the Creative Commons Attribution 4.0 license.

### Disclaimer

Golang Birmingham is an independently organised community group and is not affiliated with or endorsed by Google or the Go project.

## Support

### Getting Help

- **Issues**: Report bugs via GitHub Issues
- **Discussions**: Use GitHub Discussions for questions
- **Email**: Contact the team via the contact form
- **Community**: Join our Meetup group

### Resources

- [Go Documentation](https://golang.org/doc/)
- [Go Blog](https://blog.golang.org/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go by Example](https://gobyexample.com/)

## Acknowledgments

- The Go community for inspiration and support
- All speakers and attendees who make our events possible
- Our sponsors and venue partners
- The broader Birmingham tech community

---

Built with ❤️ using Go and standard library only.
