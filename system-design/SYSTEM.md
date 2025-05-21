# 🧠 System Design Document: Access Key Management & Token Info System

## 🎯 Objective

Build two independent but collaborative microservices communicating via Kafka, handling:

- API key creation, expiration, and rate limiting
- Token info retrieval for authorized users
- Asynchronous communication via event streaming

---

## 🧱 High-Level Architecture


```
graph TD

User[🧑 Admin / User]
Gateway[🌐 API Gateway]

User --> Gateway

subgraph Microservices
  AKM[🔐 Access Key Management Service]
  TokenService[💰 Token Info Service]
end

Gateway --> AKM
Gateway --> TokenService

AKM -- Kafka: akm --> TokenService

subgraph Storage
  AKMDB[(🗄️ PostgreSQL - AKM)]
  TokenDB[(🗄️ PostgreSQL - Token)]
  Redis[(⚡ Redis - Rate Limiting)]
end

AKM --> AKMDB
TokenService --> TokenDB
TokenService --> Redis
```

---

## ⚙️ Components Breakdown

### 🧩 1. Access Key Management Service

#### Responsibilities
- Admin operations: create, update, delete keys
- User queries: view key details
- Publish Kafka events for key updates

#### 🔐 Admin Auth
- Uses simple JWT (mocked)
- No full-fledged auth service (per project spec)

#### ✅ API Endpoints

| Endpoint            | Method | Description                      |
|---------------------|--------|----------------------------------|
| `/admin/key`        | POST   | Create new access key            |
| `/admin/key/:id`    | PUT    | Update rate limit/expiration     |
| `/admin/key/:id`    | DELETE | Delete key                       |
| `/admin/keys`       | GET    | List all keys                    |
| `/key/info`         | GET    | Get user’s API key plan info     |

---

### 🧩 2. Token Information Service

#### Responsibilities
- Validate API keys
- Enforce rate limits using Redis (fallback psql DB)
- Provide price info (mock/static)
- Listen to Kafka events to update local cache or DB

#### ✅ API Endpoint

| Endpoint       | Method | Description                                  |
|----------------|--------|----------------------------------------------|
| `/price`       | GET    | Returns token info (requires API key)        |

---

## 🔄 Kafka Communication Design

| Topic Name          | Producer                  | Consumer                 | Payload                                |
|---------------------|---------------------------|--------------------------|----------------------------------------|
| `akm` | Access Key Management Svc | Token Info Service  | JSON (hashKey, rate_limit, expires_at, event_type, ReqID) |

- Ensures **asynchronous**, decoupled communication
- Enables **event-driven updates** without REST/HTTP calls

---

## 🧠 Storage Strategy

| Resource  | Used By       | Purpose                                  |
|-----------|----------------|------------------------------------------|
| PostgreSQL| Both Services  | Stores access key metadata               |
| Redis     | Token Service  | Stores rate limiter state, key TTL, etc. |

- **Service 1** writes canonical data to Postgres
- **Service 2** reads from cache (Redis) with fallback to DB

---

## 🚦 Rate Limiting & Key Validation (Service 2)

- Redis key format: `ratelimit:<key>`
- Token Bucket / Sliding Window per minute
- Expired/disabled key → TTL/Deletion in Redis
- On Redis miss → fallback to PostgreSQL check

---

## 📦 Deployment

- `docker-compose` includes:
  - PostgreSQL for each service
  - Redis
  - Kafka + Zookeeper
  - Both Go microservices

---

## 🧰 Developer Utilities

### `.env` Configuration
- Each service has its own `.env` (DB_URL, KAFKA_BROKER, REDIS_URL, etc.)


---

## 📈 Observability

- Logs contain `request_id` (injected per request/message)
- Logs structured for filtering and tracing

---

---

## ✅ Final Notes

- Microservices are **isolated**, **asynchronous**, and **stateless**
- Data is duplicated for performance (Redis cache)
- Clean architecture and testability were core design goals
