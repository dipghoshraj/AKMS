# 🔐 Access Key Management & Token Info System – Testing Guide

## 📦 Requirements

Ensure you have the following installed:

- **Docker & Docker Compose**
- **Go 1.20+** (only if you plan to run locally)
- **curl** or **Postman** (for API testing)

---

## 🚀 Quick Setup (via Docker Compose)

### 1. Clone the Repository

```bash
git clone https://github.com/dipghoshraj/AKMS
cd AKMS
```

### 2. Set Up Environment Variables

Create a .env file in the project root with the following content:


```
# Enviroment variable for token verification system .
DB_HOST=servic_db_2
DB_USER=user
DB_PASSWORD=password
DB_NAME=servic_db_2
DB_PORT=5432

KAFKA_BROKER=kafka:9093
KAFKA_TOPIC=akm
KAFKA_GROUP_ID=akm-group


REDIS_HOST=redis
REDIS_PORT=6379
```

```
# Enviroment variable for access key management service.
DB_HOST=servic_db_1
DB_USER=user
DB_PASSWORD=password
DB_NAME=servic_db_1
DB_PORT=5432

KAFKA_URL=kafka:9093
```


### 3. Start the Services
```
docker-compose up --build
```

## This will start
- PostgreSQL
- Redis
- Kafka + Zookeeper
- Service 1 (Access Key Management)
- Service 2 (Token Info)


## 🧪 How to Test the System
### 🔐 Admin – Create Access Key

```
curl --location 'http://localhost:8081/admin/key' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{token}}' \
--data '{
    "rate_limit_per_minute": 1,
    "expires_at": 10
}'
```
Expected: 201 Created and an access key in the response

### 📋 Admin – List Keys

```
curl --location 'http://localhost:8081/admin/key' \
--header 'Authorization: Bearer {{token}}'
```

### 🛠️ Admin – Update Key

```
curl --location --request PUT 'http://localhost:8081/admin/keys/7ba98678-2f23-4a77-8607-f3047775bf57' \
--header 'Authorization: Bearer {{token}}' \
--header 'Content-Type: application/json' \
--data '{
    "rate_limit_per_minute": 10,
    "expires_at": "2025-05-23T15:00:00Z"
}'
```

### ❌ Admin – Delete Key

```
curl --location --request DELETE 'http://localhost:8080/admin/keys/7a1326db-9827-43b9-9b48-1be79d1da2cd' \
--header 'Authorization: Bearer {{token}}'
```


### 📄 User – Get Plan Info

```
curl --location 'http://localhost:8081/key/info' \
--header 'x-api-key: {{token-key}}'
```

### 🚫 User – Disable Token

```
curl --location --request POST 'http://localhost:8081/key/disable' \
--header 'x-api-key: {{token-key}} '
```

### 🧠 User – Get Proce Info (Service 2)

```
curl --location 'http://localhost:8082/price' \
--header 'x-api-key:  {{token-key}}'
```



## 📤 Kafka Flow Test
After creating a key in Service 1:

- A message is pushed to Kafka
- Service 2 listens and updates its Redis cache
- This enables valid token access via /price

✅ This is automatically validated when a new key works on the /token-info endpoint.