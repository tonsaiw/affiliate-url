// Package main - จุดเริ่มต้นของ application
// ทำหน้าที่:
// 1. Initialize database
// 2. สร้าง dependencies (repository, service, handler)
// 3. ตั้งค่า routes
// 4. Start HTTP server
package main

import (
	"log"
	"net/http"

	"affiliate-url/internal/db"
	"affiliate-url/internal/link"
)

func main() {
	// ===== 1. Initialize Database =====
	// สร้างและเชื่อมต่อ SQLite database
	// ไฟล์ affiliate.db จะถูกสร้างใน directory ปัจจุบัน
	database, err := db.InitDB("./affiliate.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// defer จะทำให้ database.Close() ถูกเรียกเมื่อ main() จบ
	defer database.Close()

	// ===== 2. สร้าง Dependencies =====
	// ใช้ pattern: Repository -> Service -> Handler
	// - Repository: จัดการ database
	// - Service: จัดการ business logic
	// - Handler: จัดการ HTTP request/response

	// สร้าง repository (data layer)
	repo := link.NewRepository(database)

	// สร้าง service (business logic layer)
	service := link.NewService(repo)

	// สร้าง handler (presentation layer)
	handler := link.NewHandler(service)

	// ===== 3. ตั้งค่า Routes =====
	// ใช้ http.HandleFunc จาก standard library
	// ไม่ใช้ framework เพื่อความเรียบง่าย

	// POST /links - สร้าง link ใหม่
	http.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateLink(w, r)
		case http.MethodGet:
			handler.ListLinks(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /a/{code} - Redirect ไปยัง original URL
	http.HandleFunc("/a/", handler.Redirect)

	// ===== 4. Start Server =====
	port := ":8080"
	log.Printf("🚀 Server starting on http://localhost%s", port)
	log.Println("Available endpoints:")
	log.Println("  POST /links     - Create new affiliate link")
	log.Println("  GET  /links     - List all links")
	log.Println("  GET  /a/{code}  - Redirect to original URL")

	// ListenAndServe จะ block จนกว่าจะมี error
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
