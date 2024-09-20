package users

import (
	"errors"
	"go-api/model"
	"go-api/token"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo Repository
}

func NewService(repo Repository) UserService {
	return UserService{repo: repo}
}

func (s *UserService) GetUsers() ([]model.User, error) {
	return s.repo.GetUsers()
}

func (s *UserService) GetUserByID(id uint) (model.User, error) {
	return s.repo.GetUserByID(id)
}

// ทำ LOGIN
// UserLogin authenticates a user by comparing the provided password with the hashed password from the database.
func (s *UserService) UserLogin(user model.User) (model.User, string, error) {
	// 1. Check username
	userFromdb, err := s.repo.GetUserByUsername(user.Username)

	if err != nil {
		// ถ้าผู้ใช้ไม่พบ ให้ส่งข้อผิดพลาดที่ชัดเจน
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, "", errors.New("username not found")
		}
		return model.User{}, "", err
	}

	// 2. เปรียบเทียบรหัสผ่านที่เข้ารหัส
	if err := bcrypt.CompareHashAndPassword([]byte(userFromdb.Password), []byte(user.Password)); err != nil {
		// ให้ข้อความข้อผิดพลาดที่ชัดเจนเมื่อรหัสผ่านไม่ตรงกัน
		return model.User{}, "", errors.New("incorrect password")
	}

	// 3. สร้าง JWT Token
	token, err := token.GenerateToken(userFromdb.Username, os.Getenv("JWT_SECRET"))
	if err != nil {
		return model.User{}, "", err
	}

	userFromdb.Password = ""

	return userFromdb, token, nil
}

// 1 รับ user ที่ส่งมา
// 2 ค้นหา username ใน db
// 3 เช็ค password ที่ส่งมากับ db ว่าตรงกันไหม
// 4 ถ้าตรงกันให้ส่ง user ไป

func (s *UserService) CreateUser(user *model.User) error {
	// เข้ารหัสรหัสผ่านก่อนที่จะบันทึกลงฐานข้อมูล
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.CreateUser(user)
}

func (s *UserService) UpdateUser(user model.User) error {
	// ตรวจสอบว่ามีการเปลี่ยนแปลงรหัสผ่านหรือไม่
	if user.Password != "" {
		// เข้ารหัสรหัสผ่านใหม่ก่อนที่จะอัปเดต
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(user model.User) error {
	return s.repo.DeleteUser(user)
}
