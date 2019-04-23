package storage

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Status struct {
  gorm.Model
  CPU int
  RAM int
  URL string
  Video string
  Status string
}

func GetStatus () (interface{}, error) {
  // connect to db
  db, err := gorm.Open("sqlite3", "gorm.db")
  var status Status
  defer db.Close()
  if err != nil {
    return &status, err
  }
  // get user with given username
  if dbc := db.Last(&status); dbc.Error != nil {
    return status, dbc.Error
  }
  return &status, nil
}

func SetStatus(status Status) error {
  db, err := gorm.Open("sqlite3", "gorm.db")
  defer db.Close()
  if err != nil {
    return err
  }
  db.Create(&status)
  return nil
}

func createDefaultStatus() error {
  db, err := gorm.Open("sqlite3", "gorm.db")
  defer db.Close()
  if err != nil {
    return err
  }
  registeredValues := 0
  db.Find(&Status{}).Count(&registeredValues)
  if registeredValues == 0 {
    db.Create(&Status{
      CPU: 25,
      RAM: 75,
      URL: "https://some.example",
      Video: "1080p59.94",
      Status: "No Signal",
    })
  }
  return nil
}
