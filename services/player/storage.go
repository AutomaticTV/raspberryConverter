package player

import (
	"errors"
	"fmt"
	"raspberryConverter/models"

	"gopkg.in/validator.v2"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //gorm requires sqlite
)

const dbPath = "/var/lib/raspberryConverter/playerConfig.db"

func initStorage() error {
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.AutoMigrate(&models.PlayerConfig{})
	return createDefaultConfig()
}

// GetConfig is a function that retrieves the player config from persistance
func GetConfig() (models.PlayerConfig, error) {
	// connect to db
	db, err := gorm.Open("sqlite3", dbPath)
	var c models.PlayerConfig
	defer db.Close()
	if err != nil {
		return c, err
	}
	// get user with given username
	if dbc := db.Last(&c); dbc.Error != nil {
		return c, dbc.Error
	}
	return c, nil
}

// SetConfig is a function that stores to persistance
func SetConfig(c models.PlayerConfig) error {
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	err = validator.Validate(c)
	if err != nil {
		fmt.Println(err)
		return errors.New("The player configuration is not valid")
	}
	db.Create(&c)
	return nil
}

func createDefaultConfig() error {
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	registeredValues := 0
	db.Find(&models.PlayerConfig{}).Count(&registeredValues)
	if registeredValues == 0 {
		db.Create(&models.PlayerConfig{
			Video:         "1080p59.94",
			AudioDecoding: "Software",
			URL:           "https://some.example",
			Transport:     "HTTP",
			Buffer:        300,
			Username:      "service_username",
			Password:      "service_password",
			Volume:        0,
			Autoplay:      "No",
		})
	}
	return nil
}
