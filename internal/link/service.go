// service.go - Business logic layer
// ทำหน้าที่เป็นตัวกลางระหว่าง handler และ repository
// จัดการ logic ต่างๆ เช่น การสร้าง short code
package link

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
)

// Service คือ struct ที่เก็บ repository
// ใช้สำหรับทำ business logic
type Service struct {
	repo *Repository
}

// NewService สร้าง Service instance ใหม่
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateLink สร้าง affiliate link ใหม่
// รับ original URL และคืนค่า Link ที่สร้างสำเร็จ
func (s *Service) CreateLink(originalURL string) (*Link, error) {
	// ตรวจสอบว่า URL ไม่ว่าง
	if strings.TrimSpace(originalURL) == "" {
		return nil, errors.New("original_url is required")
	}

	// สร้าง short code ที่ไม่ซ้ำกัน
	shortCode, err := s.generateUniqueShortCode()
	if err != nil {
		return nil, err
	}

	// สร้าง Link object
	link := &Link{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}

	// บันทึกลง database
	id, err := s.repo.Create(link)
	if err != nil {
		return nil, err
	}

	link.ID = id
	return link, nil
}

// generateUniqueShortCode สร้าง short code แบบสุ่ม
// ลองสร้างใหม่ถ้าซ้ำ (สูงสุด 10 ครั้ง)
func (s *Service) generateUniqueShortCode() (string, error) {
	// ลองสร้าง short code สูงสุด 10 ครั้ง
	for i := 0; i < 10; i++ {
		code := generateRandomCode(6)

		// ตรวจสอบว่าซ้ำหรือไม่
		exists, err := s.repo.ShortCodeExists(code)
		if err != nil {
			return "", err
		}

		// ถ้าไม่ซ้ำ ใช้ได้เลย
		if !exists {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique short code")
}

// generateRandomCode สร้างรหัสสุ่มความยาว n ตัวอักษร
// ใช้ crypto/rand เพื่อความปลอดภัย
func generateRandomCode(length int) string {
	// สร้าง random bytes
	// ใช้ขนาดใหญ่กว่า length เพราะ base64 จะขยายขนาด
	bytes := make([]byte, length)
	rand.Read(bytes)

	// แปลงเป็น base64 และตัดให้ได้ความยาวที่ต้องการ
	// ใช้ URL-safe encoding
	encoded := base64.URLEncoding.EncodeToString(bytes)

	// ตัดเอาแค่ความยาวที่ต้องการ และเอาเฉพาะตัวอักษร/ตัวเลข
	// ลบ - และ _ ออก เพื่อให้ URL สวยงาม
	clean := strings.ReplaceAll(encoded, "-", "")
	clean = strings.ReplaceAll(clean, "_", "")
	clean = strings.ReplaceAll(clean, "=", "")

	if len(clean) > length {
		clean = clean[:length]
	}

	return clean
}

// GetLinkByCode ค้นหา link จาก short code
func (s *Service) GetLinkByCode(code string) (*Link, error) {
	return s.repo.FindByShortCode(code)
}

// IncrementClick เพิ่มจำนวนคลิก
func (s *Service) IncrementClick(id int64) error {
	return s.repo.IncrementClickCount(id)
}

// GetAllLinks ดึง links ทั้งหมด
func (s *Service) GetAllLinks() ([]Link, error) {
	return s.repo.FindAll()
}
