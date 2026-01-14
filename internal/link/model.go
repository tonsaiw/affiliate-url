// Package link จัดการทุกอย่างเกี่ยวกับ affiliate links
// model.go - กำหนดโครงสร้างข้อมูล (struct) ที่ใช้ในระบบ
package link

import "time"

// Link คือ struct หลักที่แทน affiliate link ในระบบ
// ใช้สำหรับทั้ง database และ JSON response
type Link struct {
	ID          int64     `json:"id"`           // Primary key
	OriginalURL string    `json:"original_url"` // URL ต้นฉบับ (affiliate URL)
	ShortCode   string    `json:"short_code"`   // รหัสสั้น 6-8 ตัวอักษร
	ClickCount  int64     `json:"click_count"`  // จำนวนครั้งที่ถูกคลิก
	CreatedAt   time.Time `json:"created_at"`   // เวลาที่สร้าง
}

// CreateLinkRequest คือ request body สำหรับสร้าง link ใหม่
// ต้องการแค่ original_url เท่านั้น
type CreateLinkRequest struct {
	OriginalURL string `json:"original_url"` // URL ที่ต้องการสร้าง short link
}

// CreateLinkResponse คือ response เมื่อสร้าง link สำเร็จ
// ส่งกลับ short_url ที่พร้อมใช้งาน
type CreateLinkResponse struct {
	ID       int64  `json:"id"`        // ID ของ link ที่สร้าง
	ShortURL string `json:"short_url"` // URL สั้นที่พร้อมใช้งาน เช่น /a/abc123
}
