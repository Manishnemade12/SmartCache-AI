# ⚡ SmartCache AI

> **Async AI Processing & Caching Engine** — Go · Valkey · Gemini · React

A scalable backend system that accepts text/URL inputs, processes them asynchronously through a goroutine worker pool, caches results in Valkey, and generates AI-powered summaries via Gemini.

---

## 🏗️ Architecture

```
Client → Go API (Gin) → Valkey (Cache + Queue) → Worker Pool (goroutines) → Gemini AI → Valkey → Client
```

**Request path:**
1. POST `/api/submit` → check Valkey cache → HIT: return instantly | MISS: enqueue job, return job_id
2. GET `/api/status/:job_id` → return job state + result when complete

**Worker path:**
1. BLPOP from `job_queue`
2. Fetch URL content (if URL input)
3. Call Gemini AI for summary + tags
4. Store result in Valkey with TTL
5. Update job status → `completed`

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.21+ · Gin |
| Cache + Queue | Valkey 8 (Redis-compatible) |
| AI | Google Gemini 1.5 Flash |
| Frontend | React 18 · TypeScript · Vite 5 |

---

## 🚀 Quick Start

### 1. Start Valkey

```bash
docker compose up -d
```

### 2. Configure Backend

```bash
cp backend/.env.example backend/.env
# Edit backend/.env and set your GEMINI_API_KEY
```

### 3. Run Backend

```bash
cd backend
go run cmd/server/main.go
```

Backend runs at `http://localhost:8080`

### 4. Run Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs at `http://localhost:5173`

---

## 🔌 API Reference

### `POST /api/submit`
Submit text or a URL for async AI summarization.

```json
// Request
{ "input": "Paste your text or https://example.com/article" }

// Response (cache miss)
{ "job_id": "uuid", "status": "pending", "cached": false }

// Response (cache hit — instant)
{ "job_id": "hash", "status": "completed", "cached": true, "summary": "...", "tags": [] }
```

### `GET /api/status/:job_id`
Poll job status.

```json
{ "job_id": "...", "status": "completed", "summary": "...", "tags": ["AI"], "duration_ms": 1230 }
```

### `GET /api/analytics`
Get system metrics.

```json
{ "total_requests": 42, "cache_hits": 30, "cache_misses": 12, "queue_size": 0, "avg_processing_time_ms": 1240 }
```

### `GET /api/health`
Health check.

---

## 📁 Project Structure

```
smartcache-ai/
├── backend/
│   ├── cmd/server/main.go          # Entry point
│   ├── config/config.go            # Env config
│   ├── internal/
│   │   ├── ai/                     # Gemini client + prompts
│   │   ├── analytics/              # Metrics tracking
│   │   ├── api/handlers/           # HTTP handlers
│   │   ├── cache/                  # Valkey client
│   │   ├── services/               # Processor (URL fetch + AI)
│   │   └── worker/                 # Job struct + goroutine pool
│   ├── .env                        # ← Add your GEMINI_API_KEY here
│   └── go.mod
├── frontend/
│   └── src/
│       ├── components/             # SubmitForm, JobStatus, ResultCard
│       ├── pages/                  # Home, Analytics
│       └── services/api.ts         # Typed API client
└── docker-compose.yml              # Valkey service
```

---

## ⚙️ Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `GEMINI_API_KEY` | *(required)* | Google AI Studio API key |
| `PORT` | `8080` | Backend port |
| `REDIS_URL` | `redis://localhost:6379` | Valkey connection URL |
| `WORKER_COUNT` | `3` | Number of goroutine workers |
| `CACHE_TTL` | `300` | Summary cache TTL in seconds |

---

## 🔑 Valkey Key Design

| Key Pattern | Purpose |
|------------|---------|
| `summary:{hash}` | Cached AI result |
| `job:{id}` | Job state (pending → processing → completed) |
| `job_queue` | Redis list used as FIFO queue |
| `metrics:*` | Analytics counters |
