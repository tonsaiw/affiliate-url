// handler.go - HTTP handlers สำหรับ API endpoints
// ทำหน้าที่รับ request, ประมวลผล, และส่ง response
package link

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Handler คือ struct ที่เก็บ service
// ใช้สำหรับ handle HTTP requests
type Handler struct {
	service *Service
}

// NewHandler สร้าง Handler instance ใหม่
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateLink - POST /links
// สร้าง affiliate link ใหม่
func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบ HTTP method ต้องเป็น POST เท่านั้น
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body
	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// ถ้า parse ไม่ได้ ส่ง 400 Bad Request
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// เรียก service เพื่อสร้าง link
	link, err := h.service.CreateLink(req.OriginalURL)
	if err != nil {
		log.Printf("Error creating link: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// สร้าง response
	resp := CreateLinkResponse{
		ID:       link.ID,
		ShortURL: "/a/" + link.ShortCode,
	}

	// ส่ง JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(resp)
}

// Redirect - GET /a/{code}
// Redirect ไปยัง original URL และเพิ่ม click count
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	// ดึง short code จาก URL path
	// URL format: /a/{code}
	// ตัด "/a/" ออกเพื่อให้ได้ code
	path := r.URL.Path
	code := strings.TrimPrefix(path, "/a/")

	// ตรวจสอบว่า code ไม่ว่าง
	if code == "" || code == path {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// ค้นหา link จาก short code
	link, err := h.service.GetLinkByCode(code)
	if err != nil {
		// ถ้าไม่เจอ ส่ง 404 Not Found
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	// เพิ่ม click count (ทำ async ก็ได้ แต่เราเลือกทำ sync เพื่อความง่าย)
	if err := h.service.IncrementClick(link.ID); err != nil {
		// Log error แต่ไม่ต้อง fail การ redirect
		log.Printf("Error incrementing click count: %v", err)
	}

	// Redirect ไปยัง original URL
	// ใช้ 302 Found (temporary redirect)
	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

// ListLinks - GET /links
// ดึงรายการ links ทั้งหมด
func (h *Handler) ListLinks(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบ HTTP method ต้องเป็น GET เท่านั้น
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ดึง links ทั้งหมดจาก service
	links, err := h.service.GetAllLinks()
	if err != nil {
		log.Printf("Error fetching links: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ถ้าไม่มี links ส่ง empty array แทน null
	if links == nil {
		links = []Link{}
	}

	// ส่ง JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
