# Altenar – Casino Transaction Management System

> **Stack:** Go 1.24 · Gin · pgx · Kafka-go · PostgreSQL · Docker Compose

---
## 1  Quick start
```bash
# 1. copy example secrets
cp .env.example .env
```
```bash
# 2. build & run the whole stack
make run        # docker-compose up -d --build
```
Stop & clean up:
```bash
make down       # docker-compose down -v
```

---
## 2  Sending test messages to Kafka

```bash
curl -X POST http://localhost:8082/topics/transactions   -H 'Content-Type: application/vnd.kafka.json.v2+json'   -H 'Accept: application/vnd.kafka.v2+json'   -d '{
        "records":[{
          "value":{
            "user_id": 123,
            "transaction_type": "bet",
            "amount": 50.5,
            "created_at": "2024-01-01T15:04:05Z"
          }
        }]
      }'
```
---

## 3  API endpoints
| Method | Path | Query params | Description |
|--------|------|-------------|-------------|
| `GET` | `/transactions` | `user_id`, `transaction_type` | List transactions (optionally filtered). |

Response example:
```json
[
  {
    "id": 1,
    "user_id": 123,
    "transaction_type": "bet",
    "amount": 50.50,
    "created_at": "2024-01-01T15:04:05Z"
  }
]
```

---
## 4  Local development
### Make targets
| Target | Purpose |
|--------|---------|
| `make test` | Unit tests (`go test` + race) + coverage.out |
| `make coverage` | HTML report from **coverage.out** |
| `make integration-test` | Docker-backed integration tests (tag `integration`) |
| `make icov` | HTML report for integration coverage |
| `make tidy` | `go mod tidy` |
| `make run` / `make down` | Start / stop the full compose stack |

### Build tags
* `integration` – runs tests that spin up real Postgres via **dockertest**; not executed by default.

### Directory map
```
cmd/             entrypoints (api, consumer)
internal/
  api/           Gin handlers
  kafka/         consumer logic
  db/            pgx repository, migrations, test-util
  config/        env helpers
```

---

## 5  Database schema
```sql
CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  transaction_type VARCHAR(10) NOT NULL CHECK (transaction_type IN ('bet','win')),
  amount NUMERIC(12,2) NOT NULL CHECK (amount > 0),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```
Initial migration lives in `internal/db/migrations/001_create_transactions.up.sql` and is applied by the **migrate/migrate** docker job when the stack starts.