package storage

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Network struct {
  gorm.Model
  Mode string
  IP string
  Netmask string
  Gateway string
  DNS1 string
  DNS2 string
}

func GetNetwork () (interface{}, error) {
  // connect to db
  db, err := gorm.Open("sqlite3", "gorm.db")
  var network Network
  defer db.Close()
  if err != nil {
    return &network, err
  }
  // get user with given username
  if dbc := db.Last(&network); dbc.Error != nil {
    return network, dbc.Error
  }
  return &network, nil
}

func SetNetwork(network Network) error {
  db, err := gorm.Open("sqlite3", "gorm.db")
  defer db.Close()
  if err != nil {
    return err
  }
  db.Create(&network)
  return nil
}

func createDefaultNetwork() error {
  db, err := gorm.Open("sqlite3", "gorm.db")
  defer db.Close()
  if err != nil {
    return err
  }
  registeredValues := 0
  db.Find(&Network{}).Count(&registeredValues)
  if registeredValues == 0 {
    db.Create(&Network{
      Mode: "DHCP",
      IP: "",
      Netmask: "",
      Gateway: "",
      DNS1: "",
      DNS2: "",
    })
  }
  return nil
}
