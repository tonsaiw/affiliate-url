Affiliate URL Content Service
=============================

This is a Golang backend that turns an affiliate URL into ready-to-post marketing
content. The system follows Clean Architecture / Hexagonal principles so domain
logic stays framework-agnostic and external dependencies stay behind interfaces.

Features
--------
- Accept an affiliate URL
- Resolve redirects
- Scrape product metadata (title, price, images, description)
- Normalize product data into a domain model
- Generate marketing content (hook, body, CTA, hashtags)
- Apply platform-specific templates (TikTok, Facebook)
- Return structured JSON

Architecture Overview
---------------------
- `internal/domain`: core entities and value types
- `internal/port`: interfaces (ports) for external dependencies
- `internal/usecase`: business orchestration
- `internal/adapter`: delivery adapters (REST handlers)
- `internal/infrastructure`: concrete implementations (resolver, scraper, etc.)

Running the API
--------------
```sh
go run ./cmd/api
```

Example request
---------------
```sh
curl -X POST http://localhost:8080/v1/content \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/affiliate","platform":"tiktok"}'
```

Notes
-----
- The scraper in `internal/infrastructure/scraper/generic_scraper.go` is a stub
  that returns an error. Replace it with a real scraper for production use.
- AI-based generation can be added later by implementing `port.ContentGenerator`.


---------------------

internal meaning

port = สิ่งที่ usecase ติดต่อกับภายนอก
usecase = หัวข้อการทำงานหลักในระบบ ต้องพึ่งพา port
adapter = logic ที่ทำงานกับภายนอก (ทำงานกับ database) ผ่าน port ทั้ง input, output
