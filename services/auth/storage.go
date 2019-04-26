package auth

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // required by gorm
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username       string
	HashedPassword []byte
}

// Init is a function that initializes the persistance for authentication
func Init() error {
	db, err := gorm.Open("sqlite3", "auth.db")
	defer db.Close()
	if err != nil {
		return err
	}
	db.AutoMigrate(&User{})
	createDefaultUser()
	return nil
}

func getHashedPasswword(username string) ([]byte, error) {
	// connect to db
	db, err := gorm.Open("sqlite3", "auth.db")
	defer db.Close()
	if err != nil {
		return nil, err
	}
	// get user with given username
	var user User
	if dbc := db.Where("username = ?", username).First(&user); dbc.Error != nil {
		return nil, errors.New("Username not in DB")
	}
	// return hashed password
	return user.HashedPassword, nil
}

func updatePassword(username string, oldPass string, newPass string) error {
	// connect to db
	db, err := gorm.Open("sqlite3", "auth.db")
	defer db.Close()
	if err != nil {
		return err
	}
	// get user with given username
	var user User
	dbc := db.Where(
		"username = ?",
		username,
	).First(&user)
	if dbc.Error != nil {
		return dbc.Error
	}
	if bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(oldPass)) != nil {
		return errors.New("Incorrect old password")
	}
	hashedNew, err := bcrypt.GenerateFromPassword([]byte(newPass), 14)
	if err != nil {
		return errors.New("Error hashing password")
	}
	user.HashedPassword = hashedNew
	if dbc := db.Save(user); dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

func createDefaultUser() error {
	db, err := gorm.Open("sqlite3", "auth.db")
	defer db.Close()
	if err != nil {
		return err
	}
	registeredValues := 0
	db.Find(&User{}).Count(&registeredValues)
	if registeredValues == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), 14)
		db.Create(&User{
			Username:       "admin",
			HashedPassword: hashedPassword,
		})
	}
	return nil
}
