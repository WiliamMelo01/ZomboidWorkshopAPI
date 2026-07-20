# 🚀 Zomboid Workshop API

High-performance Go + Redis API microservice for indexing, storing, and serving Project Zomboid Steam Workshop Mod ID to Workshop ID mappings.

---

## 📌 Features

- **Lightning Fast Lookup**: Built with Go 1.22 and Redis in-memory storage.
- **CORS Enabled**: Ready for frontend applications like [Project Zomboid Mod Manager (PZMM)](https://github.com/WiliamMelo01/ProjectZomboidModManager).
- **Bulk Upload Support**: Fast `/mappings/bulk` endpoint for multi-threaded Workshop scanners.
- **Systemd Daemon Ready**: Includes systemd unit file for production Ubuntu / Linux deployment.

---

## 🛠️ Endpoints

### 1. `GET /mappings`
Returns all Mod ID to Workshop ID mappings in JSON.

**Response:**
```json
{
  "Hydrocraft": "514427485",
  "Base.Vehicle": "2948504608",
  "Brita": "2460154811"
}
```

---

### 2. `POST /mappings/bulk`
Ingests a batch of new Mod ID to Workshop ID mappings into Redis.

**Request:**
```json
{
  "mappings": [
    { "modId": "Hydrocraft", "workshopId": "514427485" },
    { "modId": "Brita", "workshopId": "2460154811" }
  ]
}
```

**Response:**
```json
{
  "status": "ok",
  "processed": 2
}
```

---

## 🚀 How to Run Locally

### 1. Requirements
- Go 1.22+
- Redis 6.0+

### 2. Start Service
```bash
go run main.go
```

The service will start on port `8080` (or `PORT` environment variable).

---

## 📄 License
MIT License
