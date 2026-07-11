# 🌤️ Weather Aggregator Microservice

A highly scalable, serverless-ready Weather Aggregator built in idiomatic Go. This project fetches, normalizes, and aggregates weather data from multiple third-party providers in parallel, serving it through a high-performance Redis cache and a responsive, glassmorphism UI.

## 🚀 Features

* **Parallel API Aggregation:** Concurrently fetches data from Open-Meteo and Wttr.in, mathematically aggregating temperatures and conditions.
* **Adapter Pattern:** Implements an anti-corruption layer to standardize proprietary API responses into a unified, strictly-typed Domain Model based on WMO scientific codes.
* **Serverless Redis Caching:** Completely bypasses third-party rate limits and cold-start penalties by seamlessly integrating a stateless Upstash Redis cache via a custom HTTP REST implementation.
* **Clean Architecture:** Strictly adheres to the MVC/Service-Controller pattern, ensuring the HTTP handlers are entirely decoupled from the business logic and external clients.
* **Zero-Downtime Deployment:** Fully optimized for serverless deployments on Vercel via Go's `//go:embed` directive and an isolated `api` handler.
* **Premium UI:** Features a custom CSS design system with glassmorphism, dynamic hourly-forecast modals, and responsive vector icons.

## 🏗️ Architecture

```text
.
├── api/                  # Vercel serverless entrypoint (index.go)
├── client/               # External API integrations
│   ├── openmeteo/        # Open-Meteo client
│   ├── upstash/          # Custom Redis REST wrapper
│   └── wttr/             # Wttr.in client & data normalizer
├── handler/              # HTTP layer and HTML template embedding
├── service/              # Core business logic and aggregation algorithms
├── weather/              # Shared Domain Models and WMO definitions
├── main.go               # Local persistent server entrypoint
└── vercel.json           # Vercel serverless routing configuration
```

## 🛠️ Local Development

### Prerequisites
* Go 1.20+
* A free [Upstash Redis](https://upstash.com/) database.

### Setup

1. Clone the repository and navigate to the root directory.
2. Create a `.env` file with your Upstash credentials:
```env
UPSTASH_REDIS_REST_URL="https://your-database-url.upstash.io"
UPSTASH_REDIS_REST_TOKEN="your_rest_token"
CACHE_EXPIRY_SECONDS="900"
```
3. Spin up the local development server:
```bash
go run ./cmd/server/main.go
```
4. Visit `http://localhost:8080` in your browser.

## ☁️ Deployment (Vercel)

This project is natively configured for Serverless Vercel deployments. 

1. Connect your GitHub repository to Vercel.
2. In the Vercel Dashboard, navigate to your project settings and add your `.env` variables.
3. Vercel will automatically detect the `vercel.json` router, compile the Go binary with the embedded HTML from the `handler` package, and deploy your microservice to a global CDN.

---
*Built as a masterclass in Go concurrency, clean architecture, and serverless design patterns.*