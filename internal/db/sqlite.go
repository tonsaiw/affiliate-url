// Package db จัดการการเชื่อมต่อกับ SQLite database
// ทำหน้าที่สร้าง connection และ initialize schema
package db

import (
	"database/sql"
	"log"

	// Import SQLite driver
	// ใช้ underscore (_) เพราะเราไม่ได้เรียกใช้โดยตรง แต่ให้ register ตัวเองกับ database/sql
	_ "github.com/mattn/go-sqlite3"
)

// InitDB สร้าง connection ไปยัง SQLite และสร้าง table ถ้ายังไม่มี
// คืนค่า *sql.DB ที่พร้อมใช้งาน
func InitDB(dbPath string) (*sql.DB, error) {
	// เปิด connection ไปยัง SQLite file
	// ถ้า file ไม่มี จะสร้างใหม่อัตโนมัติ
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// ทดสอบว่า connection ใช้งานได้จริง
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// สร้าง table links ถ้ายังไม่มี
	// - id: primary key, auto increment
	// - original_url: URL ต้นฉบับที่จะ redirect ไป
	// - short_code: รหัสสั้นสำหรับ URL (unique)
	// - click_count: จำนวนครั้งที่ถูกคลิก
	// - created_at: เวลาที่สร้าง (default เป็นเวลาปัจจุบัน)
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL,
		short_code TEXT NOT NULL UNIQUE,
		click_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	// Execute คำสั่ง SQL สร้าง table
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
