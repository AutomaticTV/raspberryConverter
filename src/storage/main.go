package storage

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "fmt"
)

func InitStorage() {
  // DB connection
  db, err := gorm.Open("sqlite3", "gorm.db")
  defer db.Close()
  if err != nil {
    fmt.Println("error connecting to sqlite: ", err)
  }

  // DB table creation / migration
  db.AutoMigrate(&User{})

  // set default values if necesary
  registeredUsers := 0
  db.Find(&User{}).Count(&registeredUsers)
  if registeredUsers == 0 {
    fmt.Println("Seting default user")
    defaultUser := createDefaultUser()
    db.Create(&defaultUser)
  }
}
