package storage

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Player struct {
  gorm.Model
  Video string
  AudioDecoding string
  URL string
  Transport string
  Buffer int
  Username string
  Password string
  Volume int
  Autoplay string
}

func GetPlayer () (interface{}, error) {
  // connect to db
  db, err := gorm.Open("sqlite3", "gorm.db")
  var player Player
  defer db.Close()
  if err != nil {
    return &player, err
  }
  // get user with given username
  if dbc := db.Last(&player); dbc.Error != nil {
    return player, dbc.Error
  }
  return &player, nil
}

func SetPlayer(player Player) error {
  db, err := gorm.Open("sqlite3", "gorm.db")
  defer db.Close()
  if err != nil {
    return err
  }
  db.Create(&player)
  return nil
}

func createDefaultPlayer() error {
  db, err := gorm.Open("sqlite3", "gorm.db")
  defer db.Close()
  if err != nil {
    return err
  }
  registeredValues := 0
  db.Find(&Player{}).Count(&registeredValues)
  if registeredValues == 0 {
    db.Create(&Player{
      Video: "1080p59.94",
      AudioDecoding: "Software",
      URL: "https://some.example",
      Transport: "HTTP",
      Buffer: 300,
      Username: "service_sername",
      Password: "service_password",
      Volume: 0,
      Autoplay: "No",
    })
  }
  return nil
}
