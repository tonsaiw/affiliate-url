// repository.go - จัดการการติดต่อกับ database โดยตรง
// ทำหน้าที่ CRUD operations สำหรับ links table
package link

import (
	"database/sql"
)

// Repository คือ struct ที่เก็บ database connection
// ใช้สำหรับทำ database operations ทั้งหมด
type Repository struct {
	db *sql.DB // database connection
}

// NewRepository สร้าง Repository instance ใหม่
// รับ database connection เป็น parameter
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create บันทึก link ใหม่ลงใน database
// คืนค่า ID ของ link ที่ถูกสร้าง
func (r *Repository) Create(link *Link) (int64, error) {
	// SQL สำหรับ insert link ใหม่
	// ใช้ ? เป็น placeholder เพื่อป้องกัน SQL injection
	query := `
		INSERT INTO links (original_url, short_code, click_count)
		VALUES (?, ?, 0)
	`

	// Execute insert statement
	result, err := r.db.Exec(query, link.OriginalURL, link.ShortCode)
	if err != nil {
		return 0, err
	}

	// ดึง ID ของ record ที่เพิ่งสร้าง
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// FindByShortCode ค้นหา link จาก short_code
// ใช้สำหรับ redirect - เมื่อ user เข้า /a/{code}
func (r *Repository) FindByShortCode(code string) (*Link, error) {
	query := `
		SELECT id, original_url, short_code, click_count, created_at
		FROM links
		WHERE short_code = ?
	`

	// QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์แค่ 1 row
	row := r.db.QueryRow(query, code)

	// Scan ค่าจาก row ลงใน struct
	link := &Link{}
	err := row.Scan(
		&link.ID,
		&link.OriginalURL,
		&link.ShortCode,
		&link.ClickCount,
		&link.CreatedAt,
	)

	if err != nil {
		// ถ้าไม่เจอ record จะได้ sql.ErrNoRows
		return nil, err
	}

	return link, nil
}

// IncrementClickCount เพิ่ม click_count ขึ้น 1
// เรียกใช้ทุกครั้งที่มีการ redirect
func (r *Repository) IncrementClickCount(id int64) error {
	query := `UPDATE links SET click_count = click_count + 1 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// FindAll ดึง links ทั้งหมดจาก database
// ใช้สำหรับ GET /links endpoint
func (r *Repository) FindAll() ([]Link, error) {
	query := `
		SELECT id, original_url, short_code, click_count, created_at
		FROM links
		ORDER BY created_at DESC
	`

	// Query ใช้เมื่อคาดว่าจะได้ผลลัพธ์หลาย rows
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	// defer จะทำให้ rows.Close() ถูกเรียกเมื่อ function จบ
	// สำคัญมาก! ป้องกัน memory leak
	defer rows.Close()

	// สร้าง slice เก็บผลลัพธ์
	var links []Link

	// วนอ่านแต่ละ row
	for rows.Next() {
		var link Link
		err := rows.Scan(
			&link.ID,
			&link.OriginalURL,
			&link.ShortCode,
			&link.ClickCount,
			&link.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	// ตรวจสอบ error ที่อาจเกิดระหว่างการวน loop
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

// ShortCodeExists ตรวจสอบว่า short_code นี้มีอยู่แล้วหรือไม่
// ใช้เพื่อป้องกันการสร้าง duplicate short_code
func (r *Repository) ShortCodeExists(code string) (bool, error) {
	query := `SELECT COUNT(*) FROM links WHERE short_code = ?`

	var count int
	err := r.db.QueryRow(query, code).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
