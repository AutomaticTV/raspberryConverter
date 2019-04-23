package storage

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "errors"
  "golang.org/x/crypto/bcrypt"
)

type User struct {
 Username string `gorm:"primary_key"`
 HashedPassword []byte
}

func GetHashedPasswword(username string) ([]byte, error) {
  // connect to db
  db, err := gorm.Open("sqlite3", "gorm.db")
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

func createDefaultUser() User {
  hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), 14)
  return User{Username: "admin", HashedPassword: hashedPassword}
}
