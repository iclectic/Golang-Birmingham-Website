# Golang Birmingham

The website for [Golang Birmingham](https://www.meetup.com/golang-birmingham/), the Go community meetup for Birmingham and the West Midlands.

It's a small, fast, server-rendered Go application with no JavaScript framework, no database and no tracking. Events, speakers, organisers and sponsors are plain YAML files, so updating the site is a pull request.

## Pages

| Path | What's there |
| --- | --- |
| `/` | Hero, next meetup, what to expect, stats, recent events, Go learning resources, sponsors, contact |
| `/events` | Upcoming and past events |
| `/events/{slug}` | Event details, venue/maps link, speakers, RSVP, recording and slides |
| `/events/{slug}/calendar.ics` | "Add to calendar" download |
| `/speaking` | Talk formats, ideas, how we support speakers, past speakers |
| `/sponsors` | Current sponsors and ways to support the group |
| `/organisers` | The organising team |
| `/code-of-conduct` | Code of Conduct and how to report a problem |
| `/privacy` | Privacy notice |
| `/meetup` | Redirects to the Meetup group |
| `/health`, `/robots.txt`, `/sitemap.xml` | Operational endpoints |

Photos live on Meetup, so there's no gallery; old `/gallery` links redirect there.

## Running locally

Requires Go 1.22+.

```bash
make run          # or: go run ./cmd/web
open http://localhost:8080
```

Other useful targets: `make test`, `make test-race`, `make vet`, `make fmt`, `make validate-content`, `make docker-build`, `make docker-run`.

### Configuration

All settings are optional environment variables (see `.env.example`):

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | HTTP port |
| `BASE_URL` | `http://localhost:8080` | Public URL used in canonical links, Open Graph tags, the sitemap and calendar files |
| `ENVIRONMENT` | `development` | `production` enables HSTS |
| `CONTENT_DIR` | `content` | Location of the content files |

## Editing content

Everything lives in `content/`. The loader is strict: unknown fields, duplicate slugs, events that end before they start and references to speakers that don't exist all stop the server from starting (and fail `make validate-content` / CI), so mistakes never reach the live site.

### Site settings: `content/site.json`

Title, tagline, description, Meetup URL, talk submission URL, contact and Code of Conduct emails, homepage stats and social links.

### Events: `content/events/*.yaml`

```yaml
slug: "go-microservices-workshop"        # URL: /events/go-microservices-workshop
title: "Building Microservices with Go"
summary: "One or two sentences shown on cards."
description: |
  Full description. Blank lines separate paragraphs,
  and lines starting with "- " become bullet lists.
start_time: "2026-11-26T18:30:00Z"       # RFC 3339. Use +01:00 for summer (BST) events
end_time: "2026-11-26T21:00:00Z"
venue_name: "Birmingham Tech Hub"
venue_address: "Colmore Row, Birmingham B3 2BJ"
topics: ["microservices", "grpc"]
speakers: ["sarah-chen"]                 # slugs from content/speakers
registration_url: "https://www.meetup.com/golang-birmingham/events/..."
video_url: ""                            # shown once the event is over
slides_url: ""
```

Events move from "Upcoming" to "Past" automatically when they finish. Times are always shown in UK time.

### Speakers: `content/speakers/*.yaml`

```yaml
slug: "sarah-chen"
name: "Sarah Chen"
role: "Senior Software Engineer"
organisation: "TechCorp Solutions"
biography: "..."
image: ""          # e.g. /static/images/speakers/sarah-chen.jpg, initials are shown if empty
linkedin_url: ""
github_url: ""
website_url: ""
```

### Organisers: `content/team/*.yaml`

`name`, `role`, `biography`, `image`, `linkedin_url`, `github_url`, `website_url`. Files are listed in filename order, so prefix them `01-`, `02-` and so on.

### Sponsors: `content/sponsors/*.yaml`

`name`, `category` (e.g. `Venue Partner`, `Food & Drink Sponsor`, `Community Partner`), `logo`, `website_url`, `description`. Sponsors are grouped by category in filename order.

Images go in `web/static/images/{speakers,team,sponsors}/`.

## How it's built

```
cmd/web              entry point, graceful shutdown on SIGINT/SIGTERM
internal/config      environment configuration
internal/content     strict YAML/JSON loading, validation and queries
internal/domain      types
internal/calendar    RFC 5545 iCalendar generation
internal/rendering   html/template renderer and template helpers
internal/handlers    HTTP handlers, SEO metadata and JSON-LD
internal/middleware  security headers, caching, logging, panic recovery
internal/server      routing (Go 1.22 ServeMux patterns) and http.Server setup
web/                 templates and static assets, embedded with go:embed
```

- **Templates and assets are embedded** in the binary and parsed once at startup. Asset URLs carry a content hash, so they're cached as immutable.
- **Security headers** on every response: a strict Content-Security-Policy, `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, `Permissions-Policy`, and HSTS in production.
- **SEO**: canonical URLs, Open Graph/Twitter cards, `Organization` and `Event` structured data, sitemap.
- **Accessibility**: semantic HTML, skip link, visible focus styles, reduced-motion support, and a mobile menu that works without JavaScript.

## Deploying

```bash
docker build -t golang-birmingham-website .
docker run -p 8080:8080 -e BASE_URL=https://golangbirmingham.com golang-birmingham-website
```

The image is a single static binary on distroless, running as a non-root user. Any platform that runs containers (Fly.io, Cloud Run, Render, Railway...) will work; point its health check at `/health`.

## Contributing

Pull requests are welcome, whether that's a new event, a typo fix or a feature. Please run `make fmt vet test` first, and follow the [Code of Conduct](https://golangbirmingham.com/code-of-conduct).

The Go gopher was designed by [Renée French](https://reneefrench.blogspot.com/) and is licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).
