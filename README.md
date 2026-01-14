# Affiliate URL System

ระบบสร้าง Short URL สำหรับ Affiliate Links ด้วย Go และ SQLite

## 📦 Project Structure

```
affiliate-url/
├── cmd/
│   └── api/
│       └── main.go          # จุดเริ่มต้น application
├── internal/
│   ├── db/
│   │   └── sqlite.go        # Database connection
│   └── link/
│       ├── handler.go       # HTTP handlers
│       ├── model.go         # Data models
│       ├── repository.go    # Database operations
│       └── service.go       # Business logic
├── go.mod
└── README.md
```

## 🚀 วิธีการรัน Backend

### 1. ติดตั้ง Dependencies
```bash
go mod tidy
```

### 2. Build และ Run
```bash
# วิธีที่ 1: Run โดยตรง
go run ./cmd/api/main.go

# วิธีที่ 2: Build แล้วค่อย Run
go build -o server ./cmd/api/main.go
./server
```

### 3. Server จะรันที่
```
http://localhost:8080
```

---

## 🧪 API Testing with cURL

### 1. สร้าง Affiliate Link ใหม่

**Request:**
```bash
curl -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://shopee.co.th/product/123?affiliate=abc"}'
```

**Response:**
```json
{
  "id": 1,
  "short_url": "/a/7iJ5d3"
}
```

---

### 2. ดูรายการ Links ทั้งหมด

**Request:**
```bash
curl http://localhost:8080/links
```

**Response:**
```json
[
  {
    "id": 1,
    "original_url": "https://shopee.co.th/product/123?affiliate=abc",
    "short_code": "7iJ5d3",
    "click_count": 0,
    "created_at": "2026-01-14T12:49:55Z"
  }
]
```

---

### 3. Redirect (เปิดใน Browser หรือใช้ curl)

**Request:**
```bash
# ดู redirect header
curl -I http://localhost:8080/a/7iJ5d3

# หรือ follow redirect
curl -L http://localhost:8080/a/7iJ5d3
```

**Response (Header):**
```
HTTP/1.1 302 Found
Location: https://shopee.co.th/product/123?affiliate=abc
```

---

## 📮 Postman Collection

### Import ใน Postman

สร้าง Collection ใหม่ใน Postman แล้วเพิ่ม Requests ดังนี้:

#### Request 1: Create Link
- **Method:** `POST`
- **URL:** `http://localhost:8080/links`
- **Headers:**
  - `Content-Type`: `application/json`
- **Body (raw JSON):**
```json
{
  "original_url": "https://shopee.co.th/product/123?affiliate=abc"
}
```

#### Request 2: List Links
- **Method:** `GET`
- **URL:** `http://localhost:8080/links`

#### Request 3: Redirect
- **Method:** `GET`
- **URL:** `http://localhost:8080/a/{short_code}`
- **Settings:** Disable "Automatically follow redirects" เพื่อดู 302 response

---

## 📝 API Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/links` | สร้าง affiliate link ใหม่ |
| GET | `/links` | ดูรายการ links ทั้งหมด |
| GET | `/a/{code}` | Redirect ไปยัง original URL |

---

## 🗄️ Database

SQLite database file: `./affiliate.db`

### Schema
```sql
CREATE TABLE links (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  original_url TEXT NOT NULL,
  short_code TEXT NOT NULL UNIQUE,
  click_count INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 📋 ตัวอย่าง cURL ที่พร้อมใช้งาน (Copy-Paste)

```bash
# 1. สร้าง Link สำหรับ Shopee
curl -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://shopee.co.th/product/123?affiliate=abc"}'

# 2. สร้าง Link สำหรับ Lazada
curl -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://lazada.co.th/products/456?aff_id=xyz"}'

# 3. สร้าง Link สำหรับ Amazon
curl -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://amazon.com/dp/B09V3KXJPB?tag=myaffiliate-20"}'

# 4. ดู Links ทั้งหมด
curl http://localhost:8080/links

# 5. ดู Links แบบ Pretty Print (ต้องมี jq)
curl http://localhost:8080/links | jq

# 6. Test Redirect
curl -I http://localhost:8080/a/YOUR_SHORT_CODE
```
