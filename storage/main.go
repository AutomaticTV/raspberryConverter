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
  db.AutoMigrate(&Status{})
  db.AutoMigrate(&Player{})

  // set default values (if needed!)
  createDefaultUser()
  createDefaultStatus()
  createDefaultPlayer()
}
