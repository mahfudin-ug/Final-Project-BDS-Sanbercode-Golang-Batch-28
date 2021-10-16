package models

import (
	"api-ecommerce/utils/token"
	"fmt"
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	User struct {
		ID        uint      `json:"id" gorm:"primary_key"`
		Username  string    `json:"username" gorm:"not null; unique"`
		Email     string    `json:"email" gorm:"not null; unique"`
		Password  string    `json:"password" gorm:"not null"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Gender    string    `json:"gender"`
		Phone     string    `json:"phone"`
		PhotoPath string    `json:"photo_path"`
		RoleID    uint      `json:"role_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Address   []Address `json:"-"`
	}
)

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username, password string, db *gorm.DB) (string, error) {
	var err error
	u := User{}
	err = db.Model(User{}).Where("username = ?", username).Take(&u).Error
	fmt.Println("------username found")
	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	fmt.Println("------verfied found")

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	fmt.Println("------token")
	return token, nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	// turn password into hash
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		return &User{}, errPassword
	}

	u.Password = string(hashedPassword)
	// remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	var err error = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}
